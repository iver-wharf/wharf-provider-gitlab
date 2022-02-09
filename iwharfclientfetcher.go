package main

import (
	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/model/request"
	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/model/response"
	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/wharfapi"
)

type wharfClientAPIFetcher interface {
	CreateProject(project request.Project) (response.Project, error)
	CreateProjectBranch(projectID uint, branch request.Branch) (response.Branch, error)
	CreateProvider(provider request.Provider) (response.Provider, error)
	CreateToken(token request.Token) (response.Token, error)
	GetProject(projectID uint) (response.Project, error)
	GetProjectList(params wharfapi.ProjectSearch) (response.PaginatedProjects, error)
	GetProvider(providerID uint) (response.Provider, error)
	GetProviderList(params wharfapi.ProviderSearch) (response.PaginatedProviders, error)
	GetToken(tokenID uint) (response.Token, error)
	GetTokenList(params wharfapi.TokenSearch) (response.PaginatedTokens, error)
	UpdateProject(projectID uint, project request.ProjectUpdate) (response.Project, error)
	UpdateProjectBranchList(projectID uint, branches []request.Branch) ([]response.Branch, error)
}
