package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/pkg/ginutil"
	log "github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

// REF is a default branch name. Used when default not set.
const REF = "master"

type gitLabClient struct {
	*gitlab.Client
	repositoryFiles gitLabRepoFilesReader
	branches        gitLabBranchesReader
	projects        gitLabProjectsReader
}

func getGitLabClientWritesProblem(c *gin.Context, token string, url string) (*gitLabClient, bool) {
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		ginutil.WriteInvalidBindError(c, err,
			"Creating the GitLab client failed because of an invalid URL. Please double check the Upload URL.")
		log.WithError(err).Fatalf("Failed to create client.")
		return nil, false
	}

	return &gitLabClient{git, git.RepositoryFiles, git.Branches, git.Projects}, true
}

func (client *gitLabClient) listProjects(page int) ([]*gitlab.Project, gitLabPaging, error) {
	opt := gitlab.ListProjectsOptions{OrderBy: gitlab.String("id")}
	if page != 0 {
		opt.Page = page
	}

	projects, resp, err := client.Projects.ListProjects(&opt)
	if err != nil {
		log.WithError(err).
			WithField("page", page).
			Errorln("failed to list projects")
		return nil, mapToPaging(resp), err
	}

	return projects, mapToPaging(resp), nil
}

func (client *gitLabClient) getProject(groupName string, projectName string) (*gitlab.Project, error) {
	project, _, err := client.Projects.GetProject(fmt.Sprintf("%v/%v", groupName, projectName), nil)
	if err != nil {
		log.WithError(err).Errorln("failed to get project")
		return nil, fmt.Errorf("unable to get project %s/%s: %w", groupName, projectName, err)
	}

	return project, nil
}

func (client *gitLabClient) listProjectsFromGroup(groupName string, page int) ([]*gitlab.Project, gitLabPaging, error) {
	opt := gitlab.ListProjectsOptions{
		OrderBy:          gitlab.String("id"),
		SearchNamespaces: gitlab.Bool(true),
		Search:           gitlab.String(groupName),
	}
	if page != 0 {
		opt.Page = page
	}

	projects, resp, err := client.Projects.ListProjects(&opt)
	if err != nil {
		log.WithError(err).
			WithFields(log.Fields{"group name": groupName, "page": page}).
			Errorln("failed to list projects")
		return nil, mapToPaging(resp), err
	}

	return projects, mapToPaging(resp), nil
}

func (client *gitLabClient) getBuildDefinitionIfExists(projectID int, defaultBranch string) (string, error) {
	if defaultBranch == "" {
		log.Debugln("default branch name cannot be empty")
		defaultBranch = REF
	}

	ref := defaultBranch
	opts := &gitlab.GetRawFileOptions{Ref: &ref}
	bytes, resp, err := client.repositoryFiles.GetRawFile(projectID, BuildDefinitionFileName, opts)
	if resp == nil {
		return "", err
	}

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	if resp.StatusCode != http.StatusOK {
		log.WithError(err).
			WithFields(log.Fields{
				"project ID":          projectID,
				"default branch name": defaultBranch,
				"file name":           BuildDefinitionFileName,
				"status code":         resp.StatusCode}).
			Errorln("unable to get build definition")
		return "", err
	}

	return string(bytes), err
}

func (client *gitLabClient) getBranches(gitLabProjectID int, page int) ([]*gitlab.Branch, gitLabPaging, error) {
	opt := gitlab.ListBranchesOptions{}
	if page != 0 {
		opt.Page = page
	}

	branches, resp, err := client.branches.ListBranches(gitLabProjectID, &opt)
	if err != nil {
		log.WithError(err).
			WithFields(log.Fields{"project id": gitLabProjectID, "page": page}).
			Errorln("unable to list branches")
		return nil, mapToPaging(resp), err
	}

	return branches, mapToPaging(resp), nil
}
