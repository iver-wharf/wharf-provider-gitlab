package main

import (
	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/model/request"
	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/model/response"
	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/wharfapi"
)

type wharfClientAPIFetcher interface {
	GetBuildArtifactList(params wharfapi.ArtifactSearch, buildID uint) (response.PaginatedArtifacts, error)
	GetBuildArtifact(buildID, artifactID uint) (response.Artifact, error)
	CreateProjectBranch(projectID uint, branch request.Branch) (response.Branch, error)
	UpdateProjectBranchList(projectID uint, branches []request.Branch) ([]response.Branch, error)
	GetProjectBranchList(projectID uint) ([]response.Branch, error)
	GetBuildList(params wharfapi.BuildSearch) (response.PaginatedBuilds, error)
	GetBuild(buildID uint) (response.Build, error)
	UpdateBuildStatus(buildID uint, status request.LogOrStatusUpdate) (response.Build, error)
	CreateBuildLog(buildID uint, buildLog request.LogOrStatusUpdate) error
	GetBuildLogList(buildID uint) ([]response.Log, error)
	StartProjectBuild(projectID uint, params wharfapi.ProjectStartBuild, inputs request.BuildInputs) (response.BuildReferenceWrapper, error)
	CreateProject(project request.Project) (response.Project, error)
	GetProject(projectID uint) (response.Project, error)
	GetProjectList(params wharfapi.ProjectSearch) (response.PaginatedProjects, error)
	UpdateProject(projectID uint, project request.ProjectUpdate) (response.Project, error)
	GetProvider(providerID uint) (response.Provider, error)
	GetProviderList(params wharfapi.ProviderSearch) (response.PaginatedProviders, error)
	UpdateProvider(providerID uint, provider request.ProviderUpdate) (response.Provider, error)
	CreateProvider(provider request.Provider) (response.Provider, error)
	GetBuildAllTestResultDetailList(buildID uint) (response.PaginatedTestResultDetails, error)
	GetBuildAllTestResultSummaryList(buildID uint) (response.PaginatedTestResultSummaries, error)
	GetBuildTestResultSummary(buildID, artifactID uint) (response.TestResultSummary, error)
	GetBuildTestResultDetailList(buildID, artifactID uint) (response.PaginatedTestResultDetails, error)
	GetBuildAllTestResultListSummary(buildID uint) (response.TestResultListSummary, error)
	GetToken(tokenID uint) (response.Token, error)
	GetTokenList(params wharfapi.TokenSearch) (response.PaginatedTokens, error)
	UpdateToken(tokenID uint, token request.TokenUpdate) (response.Token, error)
	CreateToken(token request.Token) (response.Token, error)
}
