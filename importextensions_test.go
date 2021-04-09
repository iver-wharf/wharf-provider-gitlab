package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOperationType(t *testing.T) {
	type testCase struct {
		name   string
		entity entity
		data   Import
		want   operationType
	}

	tests := []testCase{
		{name: "Get operation for Token", entity: tokenEntity, data: getTestImport(), want: getOperation},
		{name: "Put operation for Token", entity: tokenEntity, data: getTestImportWithoutTokenID(), want: putOperation},
		{name: "Get operation for Project", entity: projectEntity, data: getTestImport(), want: getOperation},
		{name: "Put operation for Project", entity: projectEntity, data: getTestImportWithoutProjectID(), want: putOperation},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.data.getOperationType(tc.entity))
		})
	}
}

func TestWhatToImport(t *testing.T) {
	type testCase struct {
		name string
		data Import
		want importType
	}
	tests := []testCase{
		{name: "Import project when group and project set", data: getTestImport(), want: importProject},
		{name: "Import group when group set without project", data: getTestImportWithoutProject(), want: importGroup},
		{name: "Import all groups when no group and no project set",
			data: getTestImportWithoutGroupAndProject(), want: importAllGroups},
		{name: "Invalid import when project set without group", data: getInvalidImport(), want: invalidImport},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.data.whatToImport())
		})
	}
}

func getTestImportWithoutTokenID() Import {
	importData := getTestImport()
	importData.TokenID = 0
	return importData
}

func getTestImportWithoutProjectID() Import {
	importData := getTestImport()
	importData.ProjectID = 0
	return importData
}

func getTestImportWithoutGroupAndProject() Import {
	importData := getTestImport()
	importData.Group = ""
	importData.Project = ""
	return importData
}

func getInvalidImport() Import {
	importData := getTestImport()
	importData.Group = ""
	return importData
}

func getTestImportWithoutProject() Import {
	importData := getTestImport()
	importData.Project = ""
	return importData
}

func getTestImport() Import {
	importSampleData := Import{
		TokenID:    256,
		Token:      "sample token",
		User:       "sample user name",
		URL:        "https://sth.com",
		UploadURL:  "",
		ProviderID: 1,
		ProjectID:  88,
		Project:    "sample project name",
		Group:      "sample group name",
	}
	return importSampleData
}
