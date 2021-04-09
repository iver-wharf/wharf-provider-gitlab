package main

import (
	"encoding/json"
	"testing"
)

func TestDeserializeRepoUpdate(t *testing.T) {
	example_update := `
{
  "event_name": "repository_update",
  "user_id": 1,
  "user_name": "John Smith",
  "user_email": "admin@example.com",
  "user_avatar": "https://s.gravatar.com/avatar/d4c74594d841139328695756648b6bd6?s=8://s.gravatar.com/avatar/d4c74594d841139328695756648b6bd6?s=80",
  "project_id": 1,
  "project": {
    "name":"Example",
    "description":"",
    "web_url":"http://example.com/jsmith/example",
    "avatar_url":null,
    "git_ssh_url":"git@example.com:jsmith/example.git",
    "git_http_url":"http://example.com/jsmith/example.git",
    "namespace":"Jsmith",
    "visibility_level":0,
    "path_with_namespace":"jsmith/example",
    "default_branch":"master",
    "homepage":"http://example.com/jsmith/example",
    "url":"git@example.com:jsmith/example.git",
    "ssh_url":"git@example.com:jsmith/example.git",
    "http_url":"http://example.com/jsmith/example.git"
  },
  "changes": [
    {
      "before":"8205ea8d81ce0c6b90fbe8280d118cc9fdad6130",
      "after":"4045ea7a3df38697b3730a20fb73c8bed8a3e69e",
      "ref":"refs/heads/master"
    }
  ],
  "refs":["refs/heads/master"]
}
`

	var got RepositoryUpdate

	bytes := []byte(example_update)
	err := json.Unmarshal(bytes, &got)

	if err != nil {
		t.Errorf("Error deserializing: %v", err)
	}

	if got.Name != "repository_update" {
		t.Errorf("Expected repository_update, got %v", got.Name)
	}

	if got.Project.Name != "Example" {
		t.Errorf("Expected name to be Example, got: %v", got.Project.Name)
	}
}
