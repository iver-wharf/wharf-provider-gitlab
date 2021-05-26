package main

type Change struct {
	Before string `json:"before"`
	After  string `json:"after"`
	Ref    string `json:"ref"`
}

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	SSHURL      string `json:"ssh_url"`
	Namespace   string `json:"namespace"`
}

const RepositoryUpdateEvent = "repository_update"

type RepositoryUpdate struct {
	Name    string   `json:"event_name"`
	Project Project  `json:"project"`
	Changes []Change `json:"changes"`
	Refs    []string `json:"refs"`
}

type Event struct {
	Name string `json:"event_name"`
}
