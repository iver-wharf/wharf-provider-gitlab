commit = $(shell git rev-parse HEAD)
version = latest

build: swag
	go build .
	@echo "Built binary found at ./wharf-provider-gitlab or ./wharf-provider-gitlab.exe"

test: swag
	go test ./

docker:
	docker build . \
		--pull \
		-t "quay.io/iver-wharf/wharf-provider-gitlab:latest" \
		-t "quay.io/iver-wharf/wharf-provider-gitlab:$(version)" \
		--build-arg BUILD_VERSION="$(version)" \
		--build-arg BUILD_GIT_COMMIT="$(commit)" \
		--build-arg BUILD_DATE="$(shell date --iso-8601=seconds)"
	@echo ""
	@echo "Push the image by running:"
	@echo "docker push quay.io/iver-wharf/wharf-provider-gitlab:latest"
ifneq "$(version)" "latest"
	@echo "docker push quay.io/iver-wharf/wharf-provider-gitlab:$(version)"
endif

docker-run:
	docker run --rm -it quay.io/iver-wharf/wharf-provider-gitlab:$(version)

serve: swag
	go run .

swag:
	swag init --parseDependency --parseDepth 2

deps:
	cd .. && go get -u github.com/swaggo/swag/cmd/swag@v1.7.1
	go mod download
