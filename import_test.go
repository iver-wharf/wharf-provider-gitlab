package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"strconv"
	"testing"

	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/model/request"
	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/model/response"
	"github.com/iver-wharf/wharf-provider-gitlab/testdoubles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/xanzy/go-gitlab"
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
	suite.data = getTestImport()
	nonExistentProjectID := getTestImportWithNonExistentProjectID().ProjectID

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
		}, {
			Name:               "not-master",
			CanPush:            true,
			Default:            false,
			DevelopersCanMerge: true,
			DevelopersCanPush:  true,
			Merged:             true,
			Protected:          true,
			WebURL:             "",
			Commit:             nil,
		}}, getSampleGitLabPaging(2), nil)

	wharfClientMock := new(testdoubles.WharfClientAPIFetcherMock)
	for _, p := range allProjects {
		wProj := mapToWharfProj(p, suite.data.TokenID, suite.data.ProviderID)
		respProj := response.Project{
			ProjectID:       uint(p.ID),
			TokenID:         wProj.TokenID,
			GroupName:       wProj.GroupName,
			Name:            wProj.Name,
			GitURL:          wProj.GitURL,
			AvatarURL:       wProj.AvatarURL,
			Description:     wProj.Description,
			BuildDefinition: wProj.BuildDefinition,
			RemoteProjectID: wProj.RemoteProjectID,
		}
		wharfClientMock.On("CreateProject", mock.MatchedBy(func(proj request.Project) bool {
			return proj.ProviderID == wProj.ProviderID &&
				proj.TokenID == wProj.TokenID &&
				proj.GroupName == wProj.GroupName &&
				proj.Name == wProj.Name &&
				proj.GitURL == wProj.GitURL &&
				proj.AvatarURL == wProj.AvatarURL &&
				proj.Description == wProj.Description &&
				proj.BuildDefinition == wProj.BuildDefinition &&
				proj.RemoteProjectID == wProj.RemoteProjectID
		})).Return(respProj, nil)
		wharfClientMock.On("GetProject", uint(p.ID)).Return(respProj, nil)
		wharfClientMock.
			On("UpdateProject", uint(p.ID), anyOfType(request.ProjectUpdate{})).
			Return(response.Project{}, nil)
		wharfClientMock.
			On("UpdateProjectBranchList", uint(p.ID), anyOfType([]request.Branch{})).
			Return([]response.Branch{}, nil)
	}

	wharfClientMock.On("GetProject", nonExistentProjectID).Return(response.Project{}, errors.New("project with matching ID not found"))
	wharfClientMock.
		On("UpdateProject", nonExistentProjectID, anyOfType(request.ProjectUpdate{})).
		Return(request.ProjectUpdate{}, errors.New("unable to update, project with matching ID not found"))
	wharfClientMock.
		On("UpdateProjectBranchList", nonExistentProjectID, anyOfType([]request.Branch{})).
		Return([]response.Branch{}, errors.New("unable to update project branch list, project with matching ID not found"))
	wharfClientMock.
		On("CreateProjectBranch", mock.AnythingOfType("uint"), anyOfType(request.Branch{})).
		Return(response.Branch{}, nil)

	suite.sut = gitLabImporter{
		gitLabClient: gitLabMock,
		wharfClient:  wharfClientMock,
		mapper:       mapper{suite.data.TokenID, suite.data.ProviderID},
	}
}

func (suite *importTestSuite) SetupTest() {
	suite.sut.wharfClient.(*testdoubles.WharfClientAPIFetcherMock).Mock.Calls = suite.sut.wharfClient.(*testdoubles.WharfClientAPIFetcherMock).Mock.Calls[:0]
	suite.sut.gitLabClient.(*gitLabClientMock).Mock.Calls = suite.sut.gitLabClient.(*gitLabClientMock).Mock.Calls[:0]

	suite.data = getTestImportWithoutProjectID()
}

func (suite *importTestSuite) TestImportProject() {
	wantProject := "builder"
	wantGroup := "default/super-project"
	suite.data.Group = wantGroup
	suite.data.Project = wantProject

	err := suite.sut.importProject(wantGroup, wantProject)

	require.Nilf(suite.T(), err, "Import return error: %v", err)

	apiMock := suite.sut.wharfClient.(*testdoubles.WharfClientAPIFetcherMock)
	apiMock.AssertNumberOfCalls(suite.T(), "CreateProject", 1)
	apiMock.AssertCalled(suite.T(), "CreateProject", mock.MatchedBy(func(p request.Project) bool { return p.Name == "builder" }))

	suite.sut.wharfClient.(*testdoubles.WharfClientAPIFetcherMock).AssertNumberOfCalls(suite.T(), "CreateProjectBranch", 2)
}

func (suite *importTestSuite) TestImportGroup() {
	want := "default/super-project"
	suite.data.Project = ""
	suite.data.Group = want

	err := suite.sut.importGroup(want)
	require.Nilf(suite.T(), err, "Import return error: %v", err)

	apiMock := suite.sut.wharfClient.(*testdoubles.WharfClientAPIFetcherMock)
	apiMock.AssertNumberOfCalls(suite.T(), "CreateProject", 3)
	apiMock.AssertCalled(suite.T(), "CreateProject", mock.MatchedBy(func(p request.Project) bool { return p.Name == "web" }))
	apiMock.AssertCalled(suite.T(), "CreateProject", mock.MatchedBy(func(p request.Project) bool { return p.Name == "builder" }))
	apiMock.AssertCalled(suite.T(), "CreateProject", mock.MatchedBy(func(p request.Project) bool { return p.Name == "docs" }))
}

func (suite *importTestSuite) TestImportAll() {
	suite.data.Project = ""
	suite.data.Group = ""

	err := suite.sut.importAll()
	require.Nilf(suite.T(), err, "Import return error: %v", err)

	apiMock := suite.sut.wharfClient.(*testdoubles.WharfClientAPIFetcherMock)
	apiMock.AssertNumberOfCalls(suite.T(), "CreateProject", 6)
	apiMock.AssertCalled(suite.T(), "CreateProject", mock.MatchedBy(func(p request.Project) bool { return p.Name == "web" }))
	apiMock.AssertCalled(suite.T(), "CreateProject", mock.MatchedBy(func(p request.Project) bool { return p.Name == "builder" }))
	apiMock.AssertCalled(suite.T(), "CreateProject", mock.MatchedBy(func(p request.Project) bool { return p.Name == "docs" }))
	apiMock.AssertCalled(suite.T(), "CreateProject", mock.MatchedBy(func(p request.Project) bool { return p.Name == "main_test-proj" }))
	apiMock.AssertCalled(suite.T(), "CreateProject", mock.MatchedBy(func(p request.Project) bool { return p.Name == "super-project-messages" }))
	apiMock.AssertCalled(suite.T(), "CreateProject", mock.MatchedBy(func(p request.Project) bool { return p.Name == "Boletus" }))
}

func (suite *importTestSuite) TestRefreshProjectSuccess() {
	suite.data = getTestImport()

	err := suite.sut.refreshProject(suite.data.TokenID, suite.data.ProviderID, suite.data.ProjectID)
	require.Nilf(suite.T(), err, "Refresh return error: %v", err)

	apiMock := suite.sut.wharfClient.(*testdoubles.WharfClientAPIFetcherMock)
	apiMock.AssertNumberOfCalls(suite.T(), "GetProject", 1)
	apiMock.AssertNumberOfCalls(suite.T(), "UpdateProject", 1)
	apiMock.AssertNumberOfCalls(suite.T(), "UpdateProjectBranchList", 1)

	gitlabMock := suite.sut.gitLabClient.(*gitLabClientMock)
	gitlabMock.AssertNumberOfCalls(suite.T(), "getProject", 1)
	gitlabMock.AssertNumberOfCalls(suite.T(), "getBuildDefinitionIfExists", 1)
	gitlabMock.AssertNumberOfCalls(suite.T(), "getBranches", 1)
}

func (suite *importTestSuite) TestRefreshProjectFail() {
	suite.data = getTestImportWithNonExistentProjectID()

	err := suite.sut.refreshProject(suite.data.TokenID, suite.data.ProviderID, suite.data.ProjectID)
	require.Errorf(suite.T(), err, "Refresh return error: %v", err)

	apiMock := suite.sut.wharfClient.(*testdoubles.WharfClientAPIFetcherMock)
	apiMock.AssertNumberOfCalls(suite.T(), "GetProject", 1)
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

func mapToWharfProj(proj *gitlab.Project, tokenID uint, providerID uint) request.Project {
	return request.Project{
		GroupName:       proj.Namespace.FullPath,
		Name:            proj.Name,
		ProviderID:      providerID,
		TokenID:         tokenID,
		Description:     proj.Description,
		BuildDefinition: "",
		AvatarURL:       proj.AvatarURL,
		GitURL:          proj.SSHURLToRepo,
		RemoteProjectID: strconv.Itoa(proj.ID),
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

func anyOfType(obj interface{}) mock.AnythingOfTypeArgument {
	return mock.AnythingOfType(reflect.TypeOf(obj).Name())
}

func TestFindTokenByTokenString(t *testing.T) {
	tokens := []response.Token{
		{Token: "abcdef"},
		{Token: "ghijkl"},
	}
	testCases := []struct {
		Name      string
		WantToken response.Token
		WantBool  bool
	}{
		{
			Name: "Found",
			WantToken: response.Token{
				Token: "ghijkl",
			},
			WantBool: true,
		},
		{
			Name:      "Not found",
			WantToken: response.Token{},
			WantBool:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			gotToken, gotBool := findTokenByTokenString(tokens, tc.WantToken.Token)
			assert.Equal(t, tc.WantToken.Token, gotToken.Token, "got the wrong token")
			assert.Equal(t, tc.WantBool, gotBool, "got the wrong bool value")
		})
	}
}

func TestFindProviderByTokenID(t *testing.T) {
	providers := []response.Provider{
		{
			ProviderID: 1,
			TokenID:    1,
		},
		{
			ProviderID: 2,
			TokenID:    2,
		},
	}
	testCases := []struct {
		Name         string
		WantProvider response.Provider
		WantBool     bool
	}{
		{
			Name: "Found",
			WantProvider: response.Provider{
				ProviderID: 2,
				TokenID:    2,
			},
			WantBool: true,
		},
		{
			Name:         "Not found",
			WantProvider: response.Provider{},
			WantBool:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			gotProvider, gotBool := findProviderByTokenID(providers, tc.WantProvider.TokenID)
			assert.Equal(t, tc.WantProvider.ProviderID, gotProvider.ProviderID, "got the wrong token")
			assert.Equal(t, tc.WantBool, gotBool, "got the wrong bool value")
		})
	}
}
