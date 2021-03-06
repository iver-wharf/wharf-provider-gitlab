package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

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
	project, resp, err := client.Projects.GetProject(fmt.Sprintf("%v/%v", groupName, projectName), nil, nil)
	ev := log.Error().
		WithString("group", groupName).
		WithString("project", projectName)

	if err != nil {
		ev.WithError(err).Message("Failed to get project.")
		return nil, err
	}
	if resp == nil {
		err := errors.New("nil response")
		ev.WithError(err).Message("Failed to get project.")
		return nil, err
	}

	ev = ev.WithString("status", resp.Status)

	if resp.StatusCode == http.StatusNotFound {
		ev.Message("Project not found.")
		projects, _, err := client.Search.Projects(projectName, &gitlab.SearchOptions{})
		if err != nil {
			ev.WithError(err).Message("Failed searching by project name as fallback.")
			return nil, err
		}
		for _, proj := range projects {
			if strings.EqualFold(proj.Namespace.FullPath, groupName) {
				return proj, nil
			}
		}

		return nil, fmt.Errorf("no project found matching: %s/%s", groupName, projectName)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non-2xx status code: %s", resp.Status)
	}

	return project, nil
}

func (client *gitLabClient) listProjectsFromGroup(groupName string, page int) ([]*gitlab.Project, gitLabPaging, error) {
	opt := gitlab.ListGroupProjectsOptions{
		OrderBy: gitlab.String("id"),
	}
	if page != 0 {
		opt.Page = page
	}

	log.Debug().
		WithString("groupName", groupName).
		WithInt("page", page).
		Message("Listing projects for group.")

	projects, resp, err := client.Groups.ListGroupProjects(groupName, &opt, nil)
	if err != nil {
		log.Error().
			WithError(err).
			WithString("URL", resp.Request.URL.String()).
			WithString("status", resp.Status).
			WithString("group", groupName).
			WithInt("page", page).
			Message("Failed to list projects for group.")
		return nil, mapToPaging(resp), err
	}

	log.Debug().
		WithString("groupName", groupName).
		WithInt("page", page).
		Message("Successfully listed projects for group.")

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
			Messagef("Unable to get %s file.", BuildDefinitionFileName)
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
