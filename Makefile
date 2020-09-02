MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

rwildcard = $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))

docker/service/service: $(call rwildcard,src/,*)
	(cd src; GOOS=linux GOARCH=amd64 go build -trimpath -o $(abspath $@) -ldflags="-s -w")

.PHONY: up
up: docker/service/service
	docker-compose up --build

.PHONY: psql
psql:
	docker-compose exec postgresql psql -U admin main

.PHONY: test
test:
	scripts/example.sh