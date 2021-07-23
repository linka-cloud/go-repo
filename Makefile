PROJECT = go-repo

IMAGE = linkacloud/$(PROJECT)

VERSION = $(shell git describe --tags `git rev-list --tags --max-count=1` 2> /dev/null || echo "v0.0.0-`git rev-parse --short HEAD`")

show-version:
	@echo $(VERSION)

docker: docker-build docker-push

docker-build:
	@docker image build -t $(IMAGE):$(VERSION) -t $(IMAGE):latest .

docker-push:
	@docker image push --all-tags $(IMAGE)
