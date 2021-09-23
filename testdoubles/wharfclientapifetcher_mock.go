package testdoubles

import (
	"github.com/iver-wharf/wharf-api-client-go/pkg/wharfapi"
	"github.com/stretchr/testify/mock"
)

// WharfClientAPIFetcherMock is a mock variant of the wharfapi.Client with the
// help of github.com/stretchr/testify/mock.
type WharfClientAPIFetcherMock struct {
	mock.Mock
}

// GetTokenByID returns a token as identified by its ID.
func (m *WharfClientAPIFetcherMock) GetTokenByID(tokenID uint) (wharfapi.Token, error) {
	args := m.Called(tokenID)
	return args.Get(0).(wharfapi.Token), args.Error(1)
}

// GetToken returns a token as identified by its token and username strings.
func (m *WharfClientAPIFetcherMock) GetToken(token, userName string) (wharfapi.Token, error) {
	args := m.Called(token, userName)
	return args.Get(0).(wharfapi.Token), args.Error(1)
}

// PostToken creates a new token and returns the created token,
// populated with its newly assigned ID.
func (m *WharfClientAPIFetcherMock) PostToken(token wharfapi.Token) (wharfapi.Token, error) {
	args := m.Called(token)
	return args.Get(0).(wharfapi.Token), args.Error(1)
}

// GetProviderByID returns a provider as identified by its ID.
func (m *WharfClientAPIFetcherMock) GetProviderByID(providerID uint) (wharfapi.Provider, error) {
	args := m.Called(providerID)
	return args.Get(0).(wharfapi.Provider), args.Error(1)
}

// GetProvider returns a provider as identified by its name, URL, and upload
// URL strings, as well as its token ID reference.
func (m *WharfClientAPIFetcherMock) GetProvider(providerName, urlStr string, tokenID uint) (wharfapi.Provider, error) {
	args := m.Called(providerName, urlStr, tokenID)
	return args.Get(0).(wharfapi.Provider), args.Error(1)
}

// PostProvider creates a new provider and returns the created provider,
// populated with its newly assigned ID.
func (m *WharfClientAPIFetcherMock) PostProvider(provider wharfapi.Provider) (wharfapi.Provider, error) {
	args := m.Called(provider)
	return args.Get(0).(wharfapi.Provider), args.Error(1)
}

// PutProject creates or updates a project, based on wether the project ID value
// is non-zero, and returns the created or updated project,
// populated with its possibly newly assigned ID.
func (m *WharfClientAPIFetcherMock) PutProject(project wharfapi.Project) (wharfapi.Project, error) {
	args := m.Called(project)
	return args.Get(0).(wharfapi.Project), args.Error(1)
}

// PutBranches replaces the list of branches for a project.
func (m *WharfClientAPIFetcherMock) PutBranches(branches []wharfapi.Branch) ([]wharfapi.Branch, error) {
	args := m.Called(branches)
	return args.Get(0).([]wharfapi.Branch), args.Error(1)
}
