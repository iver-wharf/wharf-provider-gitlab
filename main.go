package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/pkg/ginutil"
	"github.com/iver-wharf/wharf-core/pkg/logger"
	"github.com/iver-wharf/wharf-core/pkg/logger/consolepretty"
	"github.com/iver-wharf/wharf-provider-gitlab/docs"
	"github.com/iver-wharf/wharf-provider-gitlab/internal/httputils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// BuildDefinitionFileName is a name of the file that should exists in project
// repository if project will have ability to be built by Wharf.
const BuildDefinitionFileName = ".wharf-ci.yml"

// ProviderName is a provider name that is used in whole wharf system for GitLab.
const ProviderName = "gitlab"

var log = logger.NewScoped("WHARF-PROVIDER-GITLAB")

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
	logger.AddOutput(logger.LevelDebug, consolepretty.Default)

	var (
		config Config
		err    error
	)
	if err = loadEmbeddedVersionFile(); err != nil {
		log.Error().WithError(err).Message("Failed to read embedded version.yaml.")
		os.Exit(1)
	}
	if config, err = loadConfig(); err != nil {
		log.Error().WithError(err).Message("Failed to read config.")
		os.Exit(1)
	}

	docs.SwaggerInfo.Version = AppVersion.Version

	if config.CA.CertsFile != "" {
		client, err := httputils.NewClientWithCerts(config.CA.CertsFile)
		if err != nil {
			log.Error().WithError(err).Message("Failed to get net/http.Client with certs.")
			os.Exit(1)
		}
		http.DefaultClient = client
	}

	gin.DefaultWriter = ginutil.DefaultLoggerWriter
	gin.DefaultErrorWriter = ginutil.DefaultLoggerWriter

	r := gin.New()
	r.Use(
		ginutil.DefaultLoggerHandler,
		ginutil.RecoverProblem,
	)

	if config.HTTP.CORS.AllowAllOrigins {
		log.Info().Message("Allowing all origins in CORS.")
		r.Use(cors.Default())
	}

	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	r.POST("/import/gitlab/trigger", runGitLabTriggerHandler)
	r.GET("/import/gitlab/version", getVersionHandler)
	r.GET("/import/gitlab/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	importModule{&config}.register(r)

	if err := r.Run(config.HTTP.BindAddress); err != nil {
		log.Error().
			WithError(err).
			WithString("address", config.HTTP.BindAddress).
			Message("Failed to start web server.")
		os.Exit(2)
	}
}
