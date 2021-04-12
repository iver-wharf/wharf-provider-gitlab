package main

import (
	"github.com/stretchr/testify/mock"
	"github.com/xanzy/go-gitlab"
)

type gitLabClientMock struct {
	mock.Mock
}

func (m *gitLabClientMock) listProjects(page int) ([]*gitlab.Project, gitLabPaging, error) {
	args := m.Called(page)
	return args.Get(0).([]*gitlab.Project), args.Get(1).(gitLabPaging), args.Error(2)
}

func (m *gitLabClientMock) listProjectsFromGroup(groupName string, page int) ([]*gitlab.Project, gitLabPaging, error) {
	args := m.Called(groupName, page)
	return args.Get(0).([]*gitlab.Project), args.Get(1).(gitLabPaging), args.Error(2)
}

func (m *gitLabClientMock) getProject(groupName string, projectName string) (*gitlab.Project, error) {
	args := m.Called(groupName, projectName)
	return args.Get(0).(*gitlab.Project), args.Error(1)
}

func (m *gitLabClientMock) getBuildDefinitionIfExists(projectID int, defaultBranch string) (string, error) {
	args := m.Called(projectID, defaultBranch)
	return args.String(0), args.Error(1)
}

func (m *gitLabClientMock) getBranches(gitLabProjectID int, page int) ([]*gitlab.Branch, gitLabPaging, error) {
	args := m.Called(gitLabProjectID, page)
	return args.Get(0).([]*gitlab.Branch), args.Get(1).(gitLabPaging), args.Error(2)
}
