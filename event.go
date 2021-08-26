package main

// RepositoryUpdateEvent is the event type name for a GitLab repository update
// event.
const RepositoryUpdateEvent = "repository_update"

// RepositoryUpdate is a type of event regarding an update to a GitLab
// repository.
type RepositoryUpdate struct {
	Name    string `json:"event_name"`
	Project struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		SSHURL      string `json:"ssh_url"`
		Namespace   string `json:"namespace"`
	} `json:"project"`
	Changes []struct {
		Before string `json:"before"`
		After  string `json:"after"`
		Ref    string `json:"ref"`
	} `json:"changes"`
	Refs []string `json:"refs"`
}

// Event is the base event type, used to figure out what event type the message
// holds.
type Event struct {
	Name string `json:"event_name"`
}
