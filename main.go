package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-api-client-go/pkg/wharfapi"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/iver-wharf/wharf-provider-gitlab/docs"
)

// BUILD_DEFINITION_FILE_NAME is a name of the file that should exists in project repository if project will have ability to be built by Wharf.
const BUILD_DEFINITION_FILE_NAME = ".wharf-ci.yml"

// PROVIDER_NAME is a provider name that is used in whole wharf system for GitLab.
const PROVIDER_NAME = "gitlab"

// @title Swagger import API
// @version 1.0
// @description Wharf import server.

// @Host
// @BasePath /import
func main() {
	initLogger(log.TraceLevel)

	r := gin.Default()

	allowCors, ok := os.LookupEnv("ALLOW_CORS")
	if ok && allowCors == "YES" {
		log.Infof("Allowing CORS\n")
		r.Use(cors.Default())
	}

	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	r.POST("/import/gitlab", RunGitLabHandler)
	r.POST("/import/gitlab/trigger", RunGitLabTrigger)
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
