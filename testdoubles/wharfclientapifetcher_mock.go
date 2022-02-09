package testdoubles

import (
	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/model/request"
	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/model/response"
	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/wharfapi"
	"github.com/stretchr/testify/mock"
)

// WharfClientAPIFetcherMock is a mock variant of the wharfapi.Client with the
// help of github.com/stretchr/testify/mock.
type WharfClientAPIFetcherMock struct {
	mock.Mock
}

// CreateProjectBranch adds a branch to the project with the matching
// project ID by invoking the HTTP request:
//  POST /api/project/{projectId}/branch
func (m *WharfClientAPIFetcherMock) CreateProjectBranch(projectID uint, branch request.Branch) (response.Branch, error) {
	args := m.Called(projectID, branch)
	return args.Get(0).(response.Branch), args.Error(1)
}

// UpdateProjectBranchList resets the default branch and list of branches for a project
// using the project ID from the first branch in the provided list by invoking
// the HTTP request:
//  PUT /api/project/{projectId}/branch
func (m *WharfClientAPIFetcherMock) UpdateProjectBranchList(projectID uint, branches []request.Branch) ([]response.Branch, error) {
	args := m.Called(projectID, branches)
	return args.Get(0).([]response.Branch), args.Error(1)
}

// CreateProject adds a new project to the database by invoking the
// HTTP request:
//  POST /api/project
func (m *WharfClientAPIFetcherMock) CreateProject(project request.Project) (response.Project, error) {
	args := m.Called(project)
	return args.Get(0).(response.Project), args.Error(1)
}

// GetProject fetches a project by ID by invoking the HTTP request:
//  GET /api/project/{projectID}
func (m *WharfClientAPIFetcherMock) GetProject(projectID uint) (response.Project, error) {
	args := m.Called(projectID)
	return args.Get(0).(response.Project), args.Error(1)
}

// GetProjectList filters projects based on the parameters by invoking the HTTP
// request:
//  GET /api/project
func (m *WharfClientAPIFetcherMock) GetProjectList(params wharfapi.ProjectSearch) (response.PaginatedProjects, error) {
	args := m.Called(params)
	return args.Get(0).(response.PaginatedProjects), args.Error(1)
}

// UpdateProject updates a project by ID by invoking the HTTP request:
//  PUT /api/project/{projectID}
func (m *WharfClientAPIFetcherMock) UpdateProject(projectID uint, project request.ProjectUpdate) (response.Project, error) {
	args := m.Called(projectID, project)
	return args.Get(0).(response.Project), args.Error(1)
}

// GetProvider fetches a provider by ID by invoking the HTTP request:
//  GET /api/provider/{providerID}
func (m *WharfClientAPIFetcherMock) GetProvider(providerID uint) (response.Provider, error) {
	args := m.Called(providerID)
	return args.Get(0).(response.Provider), args.Error(1)
}

// GetProviderList filters providers based on the parameters by invoking the HTTP
// request:
//  GET /api/provider
func (m *WharfClientAPIFetcherMock) GetProviderList(params wharfapi.ProviderSearch) (response.PaginatedProviders, error) {
	args := m.Called(params)
	return args.Get(0).(response.PaginatedProviders), args.Error(1)
}

// CreateProvider creates a new provider by invoking the HTTP request:
//  POST /api/provider
func (m *WharfClientAPIFetcherMock) CreateProvider(provider request.Provider) (response.Provider, error) {
	args := m.Called(provider)
	return args.Get(0).(response.Provider), args.Error(1)
}

// GetToken fetches a token by ID by invoking the HTTP request:
//  GET /api/token/{tokenID}
func (m *WharfClientAPIFetcherMock) GetToken(tokenID uint) (response.Token, error) {
	args := m.Called(tokenID)
	return args.Get(0).(response.Token), args.Error(1)
}

// GetTokenList filters tokens based on the parameters by invoking the HTTP
// request:
//  GET /api/token
func (m *WharfClientAPIFetcherMock) GetTokenList(params wharfapi.TokenSearch) (response.PaginatedTokens, error) {
	args := m.Called(params)
	return args.Get(0).(response.PaginatedTokens), args.Error(1)
}

// CreateToken adds a new a token by invoking the HTTP request:
//  POST /api/token
func (m *WharfClientAPIFetcherMock) CreateToken(token request.Token) (response.Token, error) {
	args := m.Called(token)
	return args.Get(0).(response.Token), args.Error(1)
}
