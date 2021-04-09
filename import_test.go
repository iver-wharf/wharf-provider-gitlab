package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/iver-wharf/wharf-api-client-go/pkg/wharfapi"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/xanzy/go-gitlab"
	"github.com/iver-wharf/wharf-provider-gitlab/testdoubles"
)

type importTestSuite struct {
	suite.Suite
	data Import
	sut  gitLabImporter
}

func TestImports(t *testing.T) {
	suite.Run(t, new(importTestSuite))
}

func (suite *importTestSuite) SetupSuite() {
	logger := log.New()
	logger.Out = ioutil.Discard

	suite.data = getTestImport()

	allProjects := readProjectsFromFile(suite.T(), "testdata/projects_all.json")
	spGroupGitLabProjects := readProjectsFromFile(suite.T(), "testdata/groups/default_9/super-project_84/projects.json")
	mushroomGroupGitLabProjects := readProjectsFromFile(suite.T(), "testdata/groups/basket_25/mushroom_87/projects.json")
	defaultGroupGitLabProjects := readProjectsFromFile(suite.T(), "testdata/groups/default_9/projects.json")

	gitLabMock := new(gitLabClientMock)
	gitLabMock.On("listProjects", 0).Return(allProjects, getSampleGitLabPaging(len(allProjects)), nil)
	gitLabMock.On("listProjectsFromGroup", spGroupGitLabProjects[0].Namespace.FullPath, 0).
		Return(spGroupGitLabProjects, getSampleGitLabPaging(len(spGroupGitLabProjects)), nil)
	gitLabMock.On("listProjectsFromGroup", mushroomGroupGitLabProjects[0].Namespace.FullPath, 0).
		Return(mushroomGroupGitLabProjects, getSampleGitLabPaging(len(mushroomGroupGitLabProjects)), nil)
	gitLabMock.On("listProjectsFromGroup", defaultGroupGitLabProjects[0].Namespace.FullPath, 0).
		Return(defaultGroupGitLabProjects, getSampleGitLabPaging(len(defaultGroupGitLabProjects)), nil)
	for _, p := range allProjects {
		newProj := *p
		gitLabMock.On("getProject", p.Namespace.FullPath, p.Name).Return(&newProj, nil)
	}

	gitLabMock.
		On("getBuildDefinitionIfExists", mock.AnythingOfType("int"), mock.AnythingOfType("string")).
		Return("", nil)

	gitLabMock.
		On("getBranches", mock.AnythingOfType("int"), 0).
		Return([]*gitlab.Branch{{
			Name:               "master",
			CanPush:            true,
			Default:            true,
			DevelopersCanMerge: true,
			DevelopersCanPush:  true,
			Merged:             true,
			Protected:          true,
			WebURL:             "",
			Commit:             nil,
		}}, getSampleGitLabPaging(1), nil)

	wharfClientMock := new(testdoubles.WharfClientApiFetcherMock)
	for _, p := range allProjects {
		wProj := mapToWharfProj(p, suite.data.TokenID, suite.data.ProviderID)
		wharfClientMock.On("PutProject", mock.MatchedBy(func(proj wharfapi.Project) bool {
			return proj.ProviderID == wProj.ProviderID &&
				proj.TokenID == wProj.TokenID &&
				proj.GroupName == wProj.GroupName &&
				proj.Name == wProj.Name &&
				proj.GitURL == wProj.GitURL &&
				proj.AvatarUrl == wProj.AvatarUrl &&
				proj.Description == wProj.Description &&
				proj.BuildDefinition == wProj.BuildDefinition
		})).Return(wProj, nil)
	}

	wharfClientMock.On("PutBranches", mock.AnythingOfType(reflect.TypeOf([]wharfapi.Branch{}).Name())).Return([]wharfapi.Branch{}, nil)

	suite.sut = gitLabImporter{
		gitLabClient: gitLabMock,
		wharfClient:  wharfClientMock,
		mapper:       mapper{suite.data.TokenID, suite.data.ProviderID},
	}
}

func (suite *importTestSuite) SetupTest() {
	suite.sut.wharfClient.(*testdoubles.WharfClientApiFetcherMock).Mock.Calls = suite.sut.wharfClient.(*testdoubles.WharfClientApiFetcherMock).Mock.Calls[:0]

	suite.data = getTestImportWithoutProjectID()
}

func (suite *importTestSuite) TestImportProject() {
	wantProject := "builder"
	wantGroup := "default/super-project"
	suite.data.Group = wantGroup
	suite.data.Project = wantProject

	err := suite.sut.importProject(wantGroup, wantProject)

	require.Nilf(suite.T(), err, "Import return error: %v", err)

	apiMock := suite.sut.wharfClient.(*testdoubles.WharfClientApiFetcherMock)
	apiMock.AssertNumberOfCalls(suite.T(), "PutProject", 1)
	apiMock.AssertCalled(suite.T(), "PutProject", mock.MatchedBy(func(p wharfapi.Project) bool { return p.Name == "builder" }))

	suite.sut.wharfClient.(*testdoubles.WharfClientApiFetcherMock).AssertNumberOfCalls(suite.T(), "PutBranches", 1)
}

func (suite *importTestSuite) TestImportGroup() {
	want := "default/super-project"
	suite.data.Project = ""
	suite.data.Group = want

	err := suite.sut.importGroup(want)
	require.Nilf(suite.T(), err, "Import return error: %v", err)

	apiMock := suite.sut.wharfClient.(*testdoubles.WharfClientApiFetcherMock)
	apiMock.AssertNumberOfCalls(suite.T(), "PutProject", 3)
	apiMock.AssertCalled(suite.T(), "PutProject", mock.MatchedBy(func(p wharfapi.Project) bool { return p.Name == "web" }))
	apiMock.AssertCalled(suite.T(), "PutProject", mock.MatchedBy(func(p wharfapi.Project) bool { return p.Name == "builder" }))
	apiMock.AssertCalled(suite.T(), "PutProject", mock.MatchedBy(func(p wharfapi.Project) bool { return p.Name == "docs" }))
}

func (suite *importTestSuite) TestImportAll() {
	suite.data.Project = ""
	suite.data.Group = ""

	err := suite.sut.importAll()
	require.Nilf(suite.T(), err, "Import return error: %v", err)

	apiMock := suite.sut.wharfClient.(*testdoubles.WharfClientApiFetcherMock)
	apiMock.AssertNumberOfCalls(suite.T(), "PutProject", 6)
	apiMock.AssertCalled(suite.T(), "PutProject", mock.MatchedBy(func(p wharfapi.Project) bool { return p.Name == "web" }))
	apiMock.AssertCalled(suite.T(), "PutProject", mock.MatchedBy(func(p wharfapi.Project) bool { return p.Name == "builder" }))
	apiMock.AssertCalled(suite.T(), "PutProject", mock.MatchedBy(func(p wharfapi.Project) bool { return p.Name == "docs" }))
	apiMock.AssertCalled(suite.T(), "PutProject", mock.MatchedBy(func(p wharfapi.Project) bool { return p.Name == "main_test-proj" }))
	apiMock.AssertCalled(suite.T(), "PutProject", mock.MatchedBy(func(p wharfapi.Project) bool { return p.Name == "super-project-messages" }))
	apiMock.AssertCalled(suite.T(), "PutProject", mock.MatchedBy(func(p wharfapi.Project) bool { return p.Name == "Boletus" }))
}

func readProjectsFromFile(t *testing.T, fName string) []*gitlab.Project {
	content, err := ioutil.ReadFile(fName)
	if err != nil {
		t.Fatalf("unable to open file %v, err: %v", fName, err)
	}

	var projects []*gitlab.Project
	err = json.Unmarshal(content, &projects)
	if err != nil {
		t.Fatalf("unable to unmarshal projects: %v, err: %v", fName, err)
	}

	return projects
}

func mapToWharfProj(proj *gitlab.Project, tokenID uint, providerID uint) wharfapi.Project {
	return wharfapi.Project{
		ProjectID:       uint(proj.ID),
		GroupName:       proj.Namespace.FullPath,
		Name:            proj.Name,
		ProviderID:      providerID,
		TokenID:         tokenID,
		Description:     proj.Description,
		BuildDefinition: "",
		AvatarUrl:       proj.AvatarURL,
		GitURL:          proj.SSHURLToRepo,
	}
}

func getSampleGitLabPaging(count int) gitLabPaging {
	return gitLabPaging{
		totalItems:   count,
		totalPages:   1,
		itemsPerPage: 20,
		currentPage:  1,
		nextPage:     1,
		previousPage: 1,
	}
}
