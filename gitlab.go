package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/pkg/ginutil"
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
		log.Panic().WithError(err).Message("Failed to create client.")
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
		log.Error().
			WithError(err).
			WithInt("page", page).
			Message("Failed to list projects.")
		return nil, mapToPaging(resp), err
	}

	return projects, mapToPaging(resp), nil
}

func (client *gitLabClient) getProject(groupName string, projectName string) (*gitlab.Project, error) {
	projects, _, err := client.listProjectsFromGroup(fmt.Sprintf("%v/%v", groupName, projectName), 0)
	if err != nil {
		log.Error().
			WithError(err).
			WithStringf("project", "%s/%s", groupName, projectName).
			Message("Failed to list projects for project name and group.")
		return nil, err
	}

	if len(projects) == 1 {
		return projects[0], nil
	}

	log.Info().WithInt("projectCount", len(projects)).
		WithStringf("project", "%s/%s", groupName, projectName).
		Message("Invalid projects count.")

	return nil, fmt.Errorf("unable to get project %v/%v", groupName, projectName)
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
		log.Error().
			WithError(err).
			WithString("group", groupName).
			WithInt("page", page).
			Message("Failed to list projects for group.")
		return nil, mapToPaging(resp), err
	}

	return projects, mapToPaging(resp), nil
}

func (client *gitLabClient) getBuildDefinitionIfExists(projectID int, defaultBranch string) (string, error) {
	if defaultBranch == "" {
		log.Debug().Message("Default branch name cannot be empty.")
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
		log.Error().
			WithError(err).
			WithInt("projectId", projectID).
			WithString("branch", defaultBranch).
			WithString("status", resp.Status).
			Message("Unable to get .wharf-ci.yml file.")
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
		log.Error().
			WithError(err).
			WithInt("gitLabProjectId", gitLabProjectID).
			WithInt("page", page).
			Message("Failed to list branches.")
		return nil, mapToPaging(resp), err
	}

	return branches, mapToPaging(resp), nil
}
