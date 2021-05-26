package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
)

func runGitLabTriggerHandler(c *gin.Context) {
	log.Debug("Gitlab triggered")

	var event Event
	if err := c.ShouldBindBodyWith(&event, binding.JSON); err != nil {
		log.Fatalf("could not bind event: %v", err)
	}
	log.Infof("Got event %v", event.Name)

	if event.Name == RepositoryUpdateEvent {
		_ = RunRepositoryUpdateTrigger(c)
	}

	log.Debug("Gitlab trigger finished")
}

func RunRepositoryUpdateTrigger(c *gin.Context) error {

	var repo RepositoryUpdate
	if err := c.ShouldBindBodyWith(&repo, binding.JSON); err != nil {
		log.Errorf("Error binding RepositoryUpdate: %v", err)
		return err
	}

	log.Infof("Repo %v updated", repo.Project.Name)

	client := newWharfClient(c.GetHeader("Authorization"))

	log.Infof("Got client %v", client)

	return nil
}
