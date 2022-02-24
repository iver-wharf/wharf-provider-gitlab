package main

import (
	"strconv"

	"github.com/iver-wharf/wharf-api-client-go/v2/pkg/model/request"
	"github.com/xanzy/go-gitlab"
)

type mapper struct {
	tokenID    uint
	providerID uint
}

func (m *mapper) mapProjectToWharfEntity(proj gitlab.Project, buildDef string) request.Project {
	groupName := ""

	if proj.Namespace != nil {
		groupName = proj.Namespace.FullPath
	}

	return request.Project{
		Name:            proj.Name,
		BuildDefinition: buildDef,
		Description:     proj.Description,
		AvatarURL:       proj.AvatarURL,
		GitURL:          proj.SSHURLToRepo,
		TokenID:         m.tokenID,
		ProviderID:      m.providerID,
		GroupName:       groupName,
		RemoteProjectID: strconv.Itoa(proj.ID),
	}
}

func (m *mapper) mapBranchToWharfEntity(branch gitlab.Branch) request.Branch {
	return request.Branch{
		Name:    branch.Name,
		Default: branch.Default,
	}
}
