package main

import "github.com/iver-wharf/wharf-api-client-go/pkg/wharfapi"

type wharfClientAPIFetcher interface {
	GetTokenByID(tokenID uint) (wharfapi.Token, error)
	GetToken(token string, userName string) (wharfapi.Token, error)
	PostToken(token wharfapi.Token) (wharfapi.Token, error)
	GetProviderByID(providerID uint) (wharfapi.Provider, error)
	GetProvider(providerName string, urlStr string, uploadURLStr string, tokenID uint) (wharfapi.Provider, error)
	PostProvider(provider wharfapi.Provider) (wharfapi.Provider, error)
	PutProject(project wharfapi.Project) (wharfapi.Project, error)
	PutBranches(branches []wharfapi.Branch) ([]wharfapi.Branch, error)
}
