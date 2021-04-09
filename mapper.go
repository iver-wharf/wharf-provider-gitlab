package main

import (
	"github.com/iver-wharf/wharf-api-client-go/pkg/wharfapi"
	"github.com/xanzy/go-gitlab"
)

type mapper struct {
	tokenID    uint
	providerID uint
}

func (m *mapper) mapProjectToWharfEntity(proj gitlab.Project, buildDef string) wharfapi.Project {
	groupName := ""

	if proj.Namespace != nil {
		groupName = proj.Namespace.FullPath
	}

	return wharfapi.Project{
		Name:            proj.Name,
		BuildDefinition: buildDef,
		Description:     proj.Description,
		AvatarUrl:       proj.AvatarURL,
		GitURL:          proj.SSHURLToRepo,
		TokenID:         m.tokenID,
		ProviderID:      m.providerID,
		GroupName:       groupName,
	}
}

func (m *mapper) mapBranchToWharfEntity(projID uint, branch gitlab.Branch) wharfapi.Branch {
	return wharfapi.Branch{
		ProjectID: projID,
		Name:      branch.Name,
		Default:   branch.Default,
		TokenID:   m.tokenID,
	}
}

func (m *mapper) mapBranchesToWharfEntity(projectID uint, branches []*gitlab.Branch) []wharfapi.Branch {
	var mappedBranches []wharfapi.Branch
	for _, branch := range branches {
		mappedBranches = append(mappedBranches, m.mapBranchToWharfEntity(projectID, *branch))
	}
	return mappedBranches
}
