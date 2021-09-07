package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func runGitLabTriggerHandler(c *gin.Context) {
	log.Debug().Message("GitLab triggered.")

	var event Event
	if err := c.ShouldBindBodyWith(&event, binding.JSON); err != nil {
		log.Panic().WithError(err).Message("Could not bind event.")
	}
	log.Info().WithString("event", event.Name).Message("Successfully binded event.")

	log.Debug().Message("GitLab trigger finished.")
}
