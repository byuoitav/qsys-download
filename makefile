NAME := clevertouch-control
OWNER := byuoitav
REPO := clevertouch-control
PKG := github.com/$(OWNER)/$(REPO)
DOCKER_URL := docker.pkg.github.com

# version:
# use the git tag, if this commit
# doesn't have a tag, use the git hash
COMMIT_HASH := $(shell git rev-parse --short HEAD)
TAG := $(shell git rev-parse --short HEAD)
ifneq ($(shell git describe --exact-match --tags HEAD 2> /dev/null),)
	TAG := $(shell git describe --exact-match --tags HEAD)
endif

PRD_TAG_REGEX := "v[0-9]+\.[0-9]+\.[0-9]+"
DEV_TAG_REGEX := "v[0-9]+\.[0-9]+\.[0-9]+-.+"

# go stuff

PKG_LIST := $(shell go list $(PKG)/...)

.PHONY: all deps build test test-cov clean

all: build clean

test:
	@go test -v $(PKG_LIST)

test-cov:
	@go test -coverprofile=coverage.txt -covermode=atomic $(PKG_LIST)

lint:
	@golangci-lint run --tests=false

deps:
	@echo Downloading backend dependencies
	@go mod download

build: deps
	@mkdir -p dist

	@echo
	@echo Building for linux-amd64...
	@cd cmd/ && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../dist/$(NAME)-linux-amd64

	@echo
	@echo Building for linux-arm...
	@cd cmd/ && env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ../dist/$(NAME)-linux-arm

	@echo
	@echo Build output is located in ./dist/.

docker: clean build
ifeq ($(COMMIT_HASH),$(TAG))
	@echo Building dev container with tag $(COMMIT_HASH)

	@echo Building container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-dev:$(COMMIT_HASH)
	@docker build -f dockerfile --build-arg NAME=$(NAME)-linux-amd64 -t $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-dev:$(COMMIT_HASH) dist

	@echo Building container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm-dev:$(COMMIT_HASH)
	@docker build -f dockerfile --build-arg NAME=$(NAME)-linux-arm -t $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm-dev:$(COMMIT_HASH) dist
else ifneq ($(shell echo $(TAG) | egrep $(PRD_TAG_REGEX)),)
	@echo Building prod container with tag $(TAG)

	@echo Building container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME):$(TAG)
	@docker build -f dockerfile --build-arg NAME=$(NAME)-linux-amd64 -t $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME):$(TAG) dist

	@echo Building container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm:$(TAG)
	@docker build -f dockerfile --build-arg NAME=$(NAME)-linux-arm -t $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm:$(TAG) dist
else ifneq ($(shell echo $(TAG) | egrep $(DEV_TAG_REGEX)),)
	@echo Building dev container with tag $(TAG)

	@echo Building container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-dev:$(TAG)
	@docker build -f dockerfile --build-arg NAME=$(NAME)-linux-amd64 -t $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-dev:$(TAG) dist

	@echo Building container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm-dev:$(TAG)
	@docker build -f dockerfile --build-arg NAME=$(NAME)-linux-arm -t $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm-dev:$(TAG) dist

endif

deploy: docker
	@echo Logging into Github Package Registry
	@docker login $(DOCKER_URL) -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD)

ifeq ($(COMMIT_HASH),$(TAG))
	@echo Pushing dev container with tag $(COMMIT_HASH)

	@echo Pushing container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-dev:$(COMMIT_HASH)
	@docker push $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-dev:$(COMMIT_HASH)

	@echo Pushing container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm-dev:$(COMMIT_HASH)
	@docker push $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm-dev:$(COMMIT_HASH)

else ifneq ($(shell echo $(TAG) | egrep $(DEV_TAG_REGEX)),)
	@echo Pushing dev container with tag $(TAG)

	@echo Pushing container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-dev:$(TAG)
	@docker push $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-dev:$(TAG)

	@echo Pushing container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm-dev:$(TAG)
	@docker push $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm-dev:$(TAG)

else ifneq ($(shell echo $(TAG) | egrep $(PRD_TAG_REGEX)),)
	@echo Pushing prod container with tag $(TAG)

	@echo Pushing container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-dev:$(TAG)
	@docker push $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-dev:$(TAG)

	@echo Pushing container $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm-dev:$(TAG)
	@docker push $(DOCKER_URL)/$(OWNER)/$(REPO)/$(NAME)-arm-dev:$(TAG)

endif

clean:
	@go clean
	@rm -rf dist/