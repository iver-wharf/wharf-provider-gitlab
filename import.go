package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-api-client-go/pkg/wharfapi"
	log "github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

// runGitLabHandler godoc
// @Summary Import projects from gitlab or refresh existing one
// @Accept  json
// @Produce  json
// @Param import body main.Import _ "import object"
// @Success 201 "Successfully imported"
// @Failure 400 {object} string "Bad request"
// @Failure 401 {object} string "Unauthorized or missing jwt token"
// @Router /gitlab [post]
func runGitLabHandler(c *gin.Context) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	i := Import{}
	err := c.BindJSON(&i)
	if handleError(c, err) {
		return
	}

	importer, err := newGitLabImporter(c.GetHeader("Authorization"), &i)
	if handleError(c, err) {
		return
	}

	if i.whatToImport() == importProject {
		err = importer.importProject(i.Group, i.Project)
	} else if i.whatToImport() == importGroup {
		err = importer.importGroup(i.Group)
	} else if i.whatToImport() == importAllGroups {
		err = importer.importAll()
	} else {
		err = fmt.Errorf("invalid import data")
	}

	if handleError(c, err) {
		return
	}

	c.Status(http.StatusCreated)
}

func handleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if authErr, ok := err.(*wharfapi.AuthError); ok {
		c.Header("WWW-Authenticate", authErr.Realm)
		c.JSON(http.StatusUnauthorized, "Unauthorized")
	} else {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Error: %+v", err))
	}

	log.Errorln(err)
	return true
}

type gitLabImporter struct {
	gitLabClient gitLabFetcher
	wharfClient  wharfClientAPIFetcher
	mapper       mapper
}

func newGitLabImporter(authHeader string, importData *Import) (*gitLabImporter, error) {
	wharfClient := newWharfClient(authHeader)

	token, err := obtainToken(&wharfClient, importData)
	if err != nil {
		return nil, err
	}

	provider, err := obtainProvider(&wharfClient, token.TokenID, importData)
	if err != nil {
		return nil, err
	}

	gitLabClient, err := getGitLabClient(token.Token, provider.URL)
	if err != nil {
		return nil, err
	}

	return &gitLabImporter{
		wharfClient:  wharfClient,
		gitLabClient: gitLabClient,
		mapper:       mapper{token.TokenID, provider.ProviderID},
	}, nil
}

func obtainToken(wharfClient *wharfapi.Client, importData *Import) (wharfapi.Token, error) {
	if importData.TokenID != 0 {
		token, err := wharfClient.GetTokenById(importData.TokenID)
		if err != nil || token.TokenID == 0 {
			log.WithError(err).Errorln("unable to get token")
			return wharfapi.Token{}, err
		}
		return token, nil
	}

	token, err := wharfClient.GetToken(importData.Token, "")
	if authErr, ok := err.(*wharfapi.AuthError); ok {
		return wharfapi.Token{}, authErr
	}

	if err != nil || token.TokenID == 0 {
		token, err = wharfClient.PostToken(wharfapi.Token{Token: importData.Token})
		if err != nil {
			log.WithError(err).Errorln("unable to post token")
			return wharfapi.Token{}, err
		}
		log.WithField("tokenID", token.TokenID).Debugln("Successfully created token")
	}
	return token, nil
}

func obtainProvider(wharfClient *wharfapi.Client, tokenID uint, importData *Import) (wharfapi.Provider, error) {
	if importData.ProviderID != 0 {
		provider, err := wharfClient.GetProviderById(importData.ProviderID)
		if err != nil || provider.ProviderID == 0 {
			log.WithError(err).Errorln("unable to get provider")
			return wharfapi.Provider{}, err
		}
		return provider, nil
	}

	provider, err := wharfClient.GetProvider(ProviderName, importData.URL, "", tokenID)
	if err != nil || provider.ProviderID == 0 {
		if authErr, ok := err.(*wharfapi.AuthError); ok {
			return wharfapi.Provider{}, authErr
		}

		provider, err = wharfClient.PostProvider(wharfapi.Provider{
			Name:    ProviderName,
			URL:     importData.URL,
			TokenID: tokenID})
		if err != nil {
			log.WithError(err).Errorln("unable to post provider")
			return wharfapi.Provider{}, err
		}
	}

	log.WithField("provider", provider).Debugln("provider from db")
	return provider, nil
}

func (importer *gitLabImporter) importProject(groupName string, projectName string) error {
	gitLabProject, err := importer.gitLabClient.getProject(groupName, projectName)
	if err != nil {
		log.WithError(err).Errorln("failed to get project")
		return err
	}

	wharfProject, err := importer.putProject(*gitLabProject)
	if err != nil {
		log.WithError(err).WithField("git lab project", gitLabProject).Errorln("failed to put project")
		return err
	}

	err = importer.importBranches(wharfProject.ProjectID, gitLabProject.ID)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"git lab project ID": gitLabProject.ID,
			"wharf project":      wharfProject,
		}).Errorln("unable to import branches")
		return err
	}

	return nil
}

func (importer *gitLabImporter) importGroup(groupName string) error {
	return importPaginatedProjects(func(page int) ([]*gitlab.Project, gitLabPaging, error) {
		return importer.gitLabClient.listProjectsFromGroup(groupName, page)
	}, importer.importProjects)
}

func (importer *gitLabImporter) importAll() error {
	return importPaginatedProjects(importer.gitLabClient.listProjects, importer.importProjects)
}

func (importer gitLabImporter) importProjects(projects []*gitlab.Project) string {
	var errMessage string
	for _, project := range projects {
		wharfProject, err := importer.putProject(*project)
		if err != nil {
			log.WithError(err).WithField("git lab project", project).Errorln("failed to put project")
			errMessage += fmt.Sprintf("proj: %v err: %+v \n", project, err)
			continue
		}

		err = importer.importBranches(wharfProject.ProjectID, project.ID)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"git lab project ID": project.ID,
				"wharf project":      wharfProject,
			}).Errorln("unable to import branches")
			errMessage += fmt.Sprintf("proj: %v err: %+v \n", wharfProject, err)
		}
	}
	return errMessage
}

func (importer gitLabImporter) putProject(gitLabProject gitlab.Project) (wharfapi.Project, error) {
	buildDef, err := importer.gitLabClient.getBuildDefinitionIfExists(gitLabProject.ID, gitLabProject.DefaultBranch)
	if err != nil {
		return wharfapi.Project{}, err
	}

	wharfProject := importer.mapper.mapProjectToWharfEntity(gitLabProject, buildDef)

	dbProject, err := importer.wharfClient.PutProject(wharfProject)
	if err != nil {
		log.WithError(err).Errorln("unable to put project")
		return wharfapi.Project{}, err
	}

	return dbProject, nil
}

func (importer gitLabImporter) importBranches(wharfProjectID uint, gitLabProjectID int) error {
	errMessage := ""
	page := 0
	for page >= 0 {
		branches, paging, err := importer.gitLabClient.getBranches(gitLabProjectID, page)
		if err != nil {
			log.Errorln(err)
			return err
		}

		if len(branches) > 0 {
			wharfBranches := importer.mapper.mapBranchesToWharfEntity(wharfProjectID, branches)

			_, err = importer.wharfClient.PutBranches(wharfBranches)
			if err != nil {
				log.WithError(err).Errorln("unable to put branches")
				errMessage += err.Error()
			}
		}

		page = paging.next()
	}

	if errMessage != "" {
		return fmt.Errorf(errMessage)
	}

	return nil
}
