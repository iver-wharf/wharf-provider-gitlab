package main

import "github.com/xanzy/go-gitlab"

type gitLabFetcher interface {
	listProjects(page int) ([]*gitlab.Project, gitLabPaging, error)
	listProjectsFromGroup(groupName string, page int) ([]*gitlab.Project, gitLabPaging, error)
	getProject(groupName string, projectName string) (*gitlab.Project, error)
	getBuildDefinitionIfExists(projectID int, defaultBranch string) (string, error)
	getBranches(gitLabProjectID int, page int) ([]*gitlab.Branch, gitLabPaging, error)
}

type gitLabRepoFilesReader interface {
	GetRawFile(pid any, fileName string, opt *gitlab.GetRawFileOptions, options ...gitlab.RequestOptionFunc) ([]byte, *gitlab.Response, error)
}

type gitLabBranchesReader interface {
	ListBranches(pid any, opts *gitlab.ListBranchesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Branch, *gitlab.Response, error)
}

type gitLabProjectsReader interface {
	ListProjects(opt *gitlab.ListProjectsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error)
}
