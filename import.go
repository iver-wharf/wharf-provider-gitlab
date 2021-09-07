package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-api-client-go/pkg/wharfapi"
	"github.com/iver-wharf/wharf-core/pkg/ginutil"
	"github.com/xanzy/go-gitlab"
)

type importModule struct {
	config *Config
}

func (m importModule) register(r gin.IRouter) {
	r.POST("/import/gitlab", m.runGitLabHandler)
}

// runGitLabHandler godoc
// @Summary Import projects from gitlab or refresh existing one
// @Accept  json
// @Produce  json
// @Param import body main.Import _ "import object"
// @Success 201 "Successfully imported"
// @Failure 400 {object} problem.Response "Bad request"
// @Failure 401 {object} problem.Response "Unauthorized or missing jwt token"
// @Failure 502 {object} problem.Response "Bad gateway"
// @Router /gitlab [post]
func (m importModule) runGitLabHandler(c *gin.Context) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	i := Import{}
	err := c.ShouldBindJSON(&i)
	if err != nil {
		ginutil.WriteInvalidBindError(c, err,
			"One or more parameters failed to parse when reading the request body for GitHub projects import/refresh")
		return
	}

	wharfClient := wharfapi.Client{
		AuthHeader: c.GetHeader("Authorization"),
		APIURL:     m.config.API.URL,
	}

	importer, ok := newGitLabImporterWritesProblem(c, wharfClient, &i)
	if !ok {
		return
	}

	var detail string
	switch i.whatToImport() {
	case importProject:
		err = importer.importProject(i.Group, i.Project)
		detail = fmt.Sprintf("Unable to import GitLab project %q", i.Project)
	case importGroup:
		err = importer.importGroup(i.Group)
		detail = fmt.Sprintf("Unable to import GitLab group %q", i.Group)
	case importAllGroups:
		err = importer.importAll()
		detail = "Unable to import GitLab groups"
	default:
		err = fmt.Errorf("invalid import data")
		detail = fmt.Sprintf("You need to specify either group, group and project, or neither. "+
			"Specifying only project is invalid. "+
			"Group=%q, Project=%q", i.Group, i.Project)
		ginutil.WriteInvalidParamError(c, err, "Group or Project", detail)
		return
	}

	if err != nil {
		ginutil.WriteAPIClientWriteError(c, err, detail)
		return
	}

	c.Status(http.StatusCreated)
}

type gitLabImporter struct {
	gitLabClient gitLabFetcher
	wharfClient  wharfClientAPIFetcher
	mapper       mapper
}

func newGitLabImporterWritesProblem(c *gin.Context, wharfClient wharfapi.Client, importData *Import) (*gitLabImporter, bool) {
	token, ok := obtainTokenWritesProblem(c, wharfClient, importData)
	if !ok {
		return nil, false
	}

	provider, ok := obtainProviderWritesProblem(c, wharfClient, token.TokenID, importData)
	if !ok {
		return nil, false
	}

	gitLabClient, ok := getGitLabClientWritesProblem(c, token.Token, provider.URL)
	if !ok {
		return nil, false
	}

	return &gitLabImporter{
		wharfClient:  wharfClient,
		gitLabClient: gitLabClient,
		mapper:       mapper{token.TokenID, provider.ProviderID},
	}, true
}

func obtainTokenWritesProblem(c *gin.Context, wharfClient wharfapi.Client, importData *Import) (wharfapi.Token, bool) {
	if importData.TokenID != 0 {
		token, err := wharfClient.GetTokenByID(importData.TokenID)
		if err != nil {
			ginutil.WriteAPIClientReadError(c, err,
				fmt.Sprintf(
					"Unable to get token by ID %d. Likely because of a failed request or malformed response.",
					importData.TokenID))
		} else if token.TokenID == 0 {
			err = fmt.Errorf("token with ID %d not found", importData.TokenID)
			ginutil.WriteAPIClientReadError(c, err,
				fmt.Sprintf("Token with ID %d not found.", importData.TokenID))
		}

		if err != nil {
			return wharfapi.Token{}, false
		}
		return token, true
	}

	token, err := wharfClient.GetToken(importData.Token, "")
	if authErr, ok := err.(*wharfapi.AuthError); ok {
		c.Header("WWW-Authenticate", authErr.Realm)
		ginutil.WriteUnauthorizedError(c, authErr,
			"You are not allowed to use this functionality. Please make sure your token is correct.")
		return wharfapi.Token{}, false
	}

	if err != nil || token.TokenID == 0 {
		token, err = wharfClient.PostToken(wharfapi.Token{Token: importData.Token})
		if err != nil {
			ginutil.WriteAPIClientWriteError(c, err,
				"Unable to post token to the API. This issue might be temporary. Please try again later.")
			log.Error().WithError(err).Message("Unable to post token.")
			return wharfapi.Token{}, false
		}
		log.Debug().WithUint("tokenId", token.TokenID).Message("Successfully created token.")
	}
	return token, true
}

func obtainProviderWritesProblem(c *gin.Context, wharfClient wharfapi.Client, tokenID uint, importData *Import) (wharfapi.Provider, bool) {
	if importData.ProviderID != 0 {
		provider, err := wharfClient.GetProviderByID(importData.ProviderID)
		if err != nil || provider.ProviderID == 0 {
			ginutil.WriteAPIClientReadError(c, err,
				fmt.Sprintf("Unable to get provider by ID %d.", importData.ProjectID))
			log.Error().WithError(err).Message("Unable to get provider.")
			return wharfapi.Provider{}, false
		}
		return provider, true
	}

	provider, err := wharfClient.GetProvider(ProviderName, importData.URL, "", tokenID)
	if err != nil || provider.ProviderID == 0 {
		if authErr, ok := err.(*wharfapi.AuthError); ok {
			c.Header("WWW-Authenticate", authErr.Realm)
			ginutil.WriteUnauthorizedError(c, authErr,
				"You are not allowed to get a provider. Please make sure your token is correct.")
			return wharfapi.Provider{}, false
		}

		provider, err = wharfClient.PostProvider(wharfapi.Provider{
			Name:    ProviderName,
			URL:     importData.URL,
			TokenID: tokenID})
		if err != nil {
			ginutil.WriteAPIClientWriteError(c, err,
				"Unable to post provider to the API. This issue might be temporary. Please try again later.")
			log.Error().WithError(err).Message("Unable to create provider.")
			return wharfapi.Provider{}, false
		}
	}

	log.Debug().
		WithString("provider", provider.Name).
		WithString("providerUrl", provider.URL).
		Message("Provider from DB.")
	return provider, true
}

func (importer *gitLabImporter) importProject(groupName string, projectName string) error {
	gitLabProject, err := importer.gitLabClient.getProject(groupName, projectName)
	if err != nil {
		log.Error().WithError(err).Message("Failed to get project.")
		return err
	}

	wharfProject, err := importer.putProject(*gitLabProject)
	if err != nil {
		log.Error().
			WithError(err).
			WithString("gitLabProject", gitLabProject.NameWithNamespace).
			Message("Failed to create project.")
		return err
	}

	err = importer.importBranches(wharfProject.ProjectID, gitLabProject.ID)
	if err != nil {
		log.Error().
			WithString("gitLabProject", gitLabProject.NameWithNamespace).
			WithStringf("wharfProject", "%s/%s", wharfProject.GroupName, wharfProject.Name).
			Message("Unable to import branches.")
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
			log.Error().
				WithError(err).
				WithString("gitLabProject", project.NameWithNamespace).
				Message("Failed to create project")
			errMessage += fmt.Sprintf("proj: %v err: %+v \n", project, err)
			continue
		}

		err = importer.importBranches(wharfProject.ProjectID, project.ID)
		if err != nil {
			log.Error().
				WithString("gitLabProject", project.NameWithNamespace).
				WithStringf("wharfProject", "%s/%s", wharfProject.GroupName, wharfProject.Name).
				Message("Unable to import branches.")
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
		log.Error().WithError(err).Message("Unable to create project.")
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
			log.Error().WithError(err).Message("Failed to get branches.")
			return err
		}

		if len(branches) > 0 {
			wharfBranches := importer.mapper.mapBranchesToWharfEntity(wharfProjectID, branches)

			_, err = importer.wharfClient.PutBranches(wharfBranches)
			if err != nil {
				log.Error().WithError(err).Message("Failed to reset branches.")
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
