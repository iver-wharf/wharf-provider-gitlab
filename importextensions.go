package main

// Import is the data that is required by the import endpoint.
type Import struct {
	// used in refresh only
	TokenID uint   `json:"tokenId" example:"0"`
	Token   string `json:"token" example:"sample token"`
	User    string `json:"user" example:"sample user name"`
	URL     string `json:"url" example:"https://gitlab.local"`
	// used in refresh only
	ProviderID uint `json:"providerId" example:"0"`
	// used in refresh only
	ProjectID uint   `json:"projectId" example:"0"`
	Project   string `json:"project" example:"sample project name"`
	Group     string `json:"group" example:"default"`
}

type operationType int

const (
	getOperation operationType = iota
	putOperation
	invalidOperation
)

type entity int

const (
	tokenEntity entity = iota
	projectEntity
)

type importType int

const (
	importProject importType = iota
	importGroup
	importAllGroups
	invalidImport
)

func (i Import) getOperationType(entity entity) operationType {
	if entity == tokenEntity {
		if i.TokenID == 0 {
			return putOperation
		}
		return getOperation
	} else if entity == projectEntity {
		if i.ProjectID == 0 {
			return putOperation
		}
		return getOperation
	}
	return invalidOperation
}

func (i Import) tokenOperation() operationType {
	return i.getOperationType(tokenEntity)
}

func (i Import) whatToImport() importType {
	if i.Project != "" && i.Group != "" {
		return importProject
	} else if i.Project == "" && i.Group != "" {
		return importGroup
	} else if i.Project == "" && i.Group == "" {
		return importAllGroups
	}
	return invalidImport
}
