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

// GetBuildArtifactList filters artifacts based on the parameters by invoking the HTTP
// request:
//  GET /api/build/{buildId}/artifact
func (m *WharfClientAPIFetcherMock) GetBuildArtifactList(params wharfapi.ArtifactSearch, buildID uint) (response.PaginatedArtifacts, error) {
	args := m.Called(params, buildID)
	return args.Get(0).(response.PaginatedArtifacts), args.Error(1)
}

// GetBuildArtifact gets an artifact by invoking the HTTP request:
//  GET /api/build/{buildId}/artifact/{artifactId}
func (m *WharfClientAPIFetcherMock) GetBuildArtifact(buildID, artifactID uint) (response.Artifact, error) {
	args := m.Called(buildID, artifactID)
	return args.Get(0).(response.Artifact), args.Error(1)
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

// GetProjectBranchList gets the branches for a project by invoking the HTTP
// request:
//  GET /api/project/{projectId}/branch
func (m *WharfClientAPIFetcherMock) GetProjectBranchList(projectID uint) ([]response.Branch, error) {
	args := m.Called(projectID)
	return args.Get(0).([]response.Branch), args.Error(1)
}

// GetBuildList filters builds based on the parameters by invoking the HTTP
// request:
//  GET /api/build
func (m *WharfClientAPIFetcherMock) GetBuildList(params wharfapi.BuildSearch) (response.PaginatedBuilds, error) {
	args := m.Called(params)
	return args.Get(0).(response.PaginatedBuilds), args.Error(1)
}

// GetBuild gets a build by invoking the HTTP request:
//  GET /api/build/{buildId}
func (m *WharfClientAPIFetcherMock) GetBuild(buildID uint) (response.Build, error) {
	args := m.Called(buildID)
	return args.Get(0).(response.Build), args.Error(1)
}

// UpdateBuildStatus updates a build by invoking the HTTP request:
//  PUT /api/build/{buildId}/status
func (m *WharfClientAPIFetcherMock) UpdateBuildStatus(buildID uint, status request.LogOrStatusUpdate) (response.Build, error) {
	args := m.Called(buildID, status)
	return args.Get(0).(response.Build), args.Error(1)
}

// CreateBuildLog adds a new log to a build by invoking the HTTP request:
//  POST /api/build/{buildId}/log
func (m *WharfClientAPIFetcherMock) CreateBuildLog(buildID uint, buildLog request.LogOrStatusUpdate) error {
	args := m.Called(buildID, buildLog)
	return args.Error(0)
}

// GetBuildLogList gets the logs for a build by invoking the HTTP request:
//  GET /api/build/{buildId}/log
func (m *WharfClientAPIFetcherMock) GetBuildLogList(buildID uint) ([]response.Log, error) {
	args := m.Called(buildID)
	return args.Get(0).([]response.Log), args.Error(1)
}

// StartProjectBuild starts a new build by invoking the HTTP request:
//  POST /api/project/{projectID}/build
func (m *WharfClientAPIFetcherMock) StartProjectBuild(projectID uint, params wharfapi.ProjectStartBuild, inputs request.BuildInputs) (response.BuildReferenceWrapper, error) {
	args := m.Called(projectID, params, inputs)
	return args.Get(0).(response.BuildReferenceWrapper), args.Error(1)
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

// UpdateProvider updates the provider with the specified ID by invoking the
// HTTP request:
//  PUT /api/provider/{providerID}
func (m *WharfClientAPIFetcherMock) UpdateProvider(providerID uint, provider request.ProviderUpdate) (response.Provider, error) {
	args := m.Called(providerID, provider)
	return args.Get(0).(response.Provider), args.Error(1)
}

// CreateProvider creates a new provider by invoking the HTTP request:
//  POST /api/provider
func (m *WharfClientAPIFetcherMock) CreateProvider(provider request.Provider) (response.Provider, error) {
	args := m.Called(provider)
	return args.Get(0).(response.Provider), args.Error(1)
}

// GetBuildAllTestResultDetailList fetches all the test result
// details for the specified build by invoking the HTTP request:
//  GET /api/build/{buildId}/test-result/detail
func (m *WharfClientAPIFetcherMock) GetBuildAllTestResultDetailList(buildID uint) (response.PaginatedTestResultDetails, error) {
	args := m.Called(buildID)
	return args.Get(0).(response.PaginatedTestResultDetails), args.Error(1)
}

// GetBuildAllTestResultSummaryList fetches all the test result
// summaries for the specified build by invoking the HTTP request:
//  GET /api/build/{buildId}/test-result/summary
func (m *WharfClientAPIFetcherMock) GetBuildAllTestResultSummaryList(buildID uint) (response.PaginatedTestResultSummaries, error) {
	args := m.Called(buildID)
	return args.Get(0).(response.PaginatedTestResultSummaries), args.Error(1)
}

// GetBuildTestResultSummary fetches a test result summary by ID by
// invoking the HTTP request:
//  GET /api/build/{buildId}/test-result/summary/{artifactId}
func (m *WharfClientAPIFetcherMock) GetBuildTestResultSummary(buildID, artifactID uint) (response.TestResultSummary, error) {
	args := m.Called(buildID, artifactID)
	return args.Get(0).(response.TestResultSummary), args.Error(1)
}

// GetBuildTestResultDetailList fetches all test result details for the specified
// test result summary by invoking the HTTP request:
//  GET /api/build/{buildId}/test-result/summary/{artifactId}/detail
func (m *WharfClientAPIFetcherMock) GetBuildTestResultDetailList(buildID, artifactID uint) (response.PaginatedTestResultDetails, error) {
	args := m.Called(buildID, artifactID)
	return args.Get(0).(response.PaginatedTestResultDetails), args.Error(1)
}

// GetBuildAllTestResultListSummary fetches the test result list summary of all tests for
// the specified build.
//  GET /api/build/{buildId}/test-result/list-summary
func (m *WharfClientAPIFetcherMock) GetBuildAllTestResultListSummary(buildID uint) (response.TestResultListSummary, error) {
	args := m.Called(buildID)
	return args.Get(0).(response.TestResultListSummary), args.Error(1)
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

// UpdateToken updates the token with the specified ID by invoking the HTTP request:
//  PUT /api/token
func (m *WharfClientAPIFetcherMock) UpdateToken(tokenID uint, token request.TokenUpdate) (response.Token, error) {
	args := m.Called(tokenID, token)
	return args.Get(0).(response.Token), args.Error(1)
}

// CreateToken adds a new a token by invoking the HTTP request:
//  POST /api/token
func (m *WharfClientAPIFetcherMock) CreateToken(token request.Token) (response.Token, error) {
	args := m.Called(token)
	return args.Get(0).(response.Token), args.Error(1)
}
