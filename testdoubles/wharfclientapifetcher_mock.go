package testdoubles

import (
	"github.com/iver-wharf/wharf-api-client-go/pkg/wharfapi"
	"github.com/stretchr/testify/mock"
)

type WharfClientApiFetcherMock struct {
	mock.Mock
}

func (m *WharfClientApiFetcherMock) GetTokenById(tokenID uint) (wharfapi.Token, error) {
	args := m.Called(tokenID)
	return args.Get(0).(wharfapi.Token), args.Error(1)
}

func (m *WharfClientApiFetcherMock) GetToken(token string, userName string) (wharfapi.Token, error) {
	args := m.Called(token, userName)
	return args.Get(0).(wharfapi.Token), args.Error(1)
}

func (m *WharfClientApiFetcherMock) PostToken(token wharfapi.Token) (wharfapi.Token, error) {
	args := m.Called(token)
	return args.Get(0).(wharfapi.Token), args.Error(1)
}

func (m *WharfClientApiFetcherMock) GetProviderById(providerID uint) (wharfapi.Provider, error) {
	args := m.Called(providerID)
	return args.Get(0).(wharfapi.Provider), args.Error(1)
}

func (m *WharfClientApiFetcherMock) GetProvider(providerName string, URLStr string, uploadURLStr string, tokenID uint) (wharfapi.Provider, error) {
	args := m.Called(providerName, URLStr, uploadURLStr, tokenID)
	return args.Get(0).(wharfapi.Provider), args.Error(1)
}

func (m *WharfClientApiFetcherMock) PostProvider(provider wharfapi.Provider) (wharfapi.Provider, error) {
	args := m.Called(provider)
	return args.Get(0).(wharfapi.Provider), args.Error(1)
}

func (m *WharfClientApiFetcherMock) PutProject(project wharfapi.Project) (wharfapi.Project, error) {
	args := m.Called(project)
	return args.Get(0).(wharfapi.Project), args.Error(1)
}
func (m *WharfClientApiFetcherMock) PutBranches(branches []wharfapi.Branch) ([]wharfapi.Branch, error) {
	args := m.Called(branches)
	return args.Get(0).([]wharfapi.Branch), args.Error(1)
}
