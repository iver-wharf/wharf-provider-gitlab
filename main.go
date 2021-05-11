package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-api-client-go/pkg/wharfapi"
	"github.com/iver-wharf/wharf-provider-gitlab/docs"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// BuildDefinitionFileName is a name of the file that should exists in project
// repository if project will have ability to be built by Wharf.
const BuildDefinitionFileName = ".wharf-ci.yml"

// ProviderName is a provider name that is used in whole wharf system for GitLab.
const ProviderName = "gitlab"

// @title Wharf provider API for GitLab
// @description Wharf backend API for integrating GitLab repositories with
// @description the Wharf main API.
// @license.name MIT
// @license.url https://github.com/iver-wharf/wharf-provider-gitlab/blob/master/LICENSE
// @contact.name Iver Wharf GitLab provider API support
// @contact.url https://github.com/iver-wharf/wharf-provider-gitlab/issues
// @contact.email wharf@iver.se
// @basePath /import
func main() {
	docs.SwaggerInfo.Version = ApiVersion.Version

	initLogger(log.TraceLevel)

	r := gin.Default()

	allowCors, ok := os.LookupEnv("ALLOW_CORS")
	if ok && allowCors == "YES" {
		log.Infof("Allowing CORS\n")
		r.Use(cors.Default())
	}

	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	r.POST("/import/gitlab", runGitLabHandler)
	r.POST("/import/gitlab/trigger", runGitLabTriggerHandler)
	r.GET("/import/gitlab/version", getVersionHandler)
	r.GET("/import/gitlab/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run()
	if err != nil {
		log.Infof("unable to run gin, error: %+v\n", err)
	}
}

func initLogger(level log.Level) {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	})
	log.SetLevel(level)
}

func newWharfClient(authHeader string) wharfapi.Client {
	return wharfapi.Client{
		ApiUrl:     os.Getenv("WHARF_API_URL"),
		AuthHeader: authHeader,
	}
}
