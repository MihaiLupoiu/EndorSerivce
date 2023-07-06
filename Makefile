# Makefile

# Variables
UNAME		:= $(shell uname -s)

.EXPORT_ALL_VARIABLES:

# this is godly
# https://news.ycombinator.com/item?id=11939200
.PHONY: help
help:	### this screen. Keep it first target to be default
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@# this is not tested, but prepared in advance for you, Mac drivers
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

# Targets
#
.PHONY: debug
debug:	### Debug Makefile itself
	@echo $(UNAME)

.PHONY: build
build: ### Build the service
	CGO_ENABLED=0 go build -a -o ./target/bin/endorService ./cmd/main.go

.PHONY: test
test: ### Run tests with coverage
	go test ./... -cover

.PHONY: run
run: build ### Run the service
	./target/bin/endorService

.PHONY: clean
clean: ### Clean binary artifacts from build
	rm -fr ./target

.PHONY: docker
docker: build ### Build docker image with latest code change
	docker build -t endor-service:latest ./

.PHONY: publish-docker-image
publish-docker-image: build docker ### Publish docker image to docker hub
	docker tag endor-service:latest myhay/endor-service && docker push myhay/endor-service
