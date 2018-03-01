RUN_CONTAINER := neo-agent-tmp

.PHONY: check
check: lint vet test

.PHONY: test
test: templates
	go test ./...

.PHONY: vet
vet: templates
	# Only consider it a failure if issues are in non-test files
	! go vet ./... 2>&1 | tee /dev/tty | grep '.go' | grep -v '_test.go'

.PHONY: vetall
vetall: templates
	go vet ./...

.PHONY: lint
lint:
	golint -set_exit_status ./cmd/... ./internal/...

templates:
	scripts/make-templates

.PHONY: image
image:
	./scripts/build

.PHONY: vendor
vendor:
	dep ensure

signalfx-agent: templates
	CGO_ENABLED=0 go build \
		-ldflags "-X main.Version=$(AGENT_VERSION) -X main.BuiltTime=$$(date +%FT%T%z)" \
		-o signalfx-agent \
		github.com/signalfx/signalfx-agent/cmd/agent

.PHONY: bundle
bundle:
	BUILD_BUNDLE=true scripts/build

.PHONY: attach-image
run-shell:
# Attach to the running container kicked off by `make run-image`.
	docker exec -it $(RUN_CONTAINER) bash

.PHONY: dev-image
dev-image:
	bash -ec "source scripts/common.sh && do_docker_build signalfx-agent-dev latest dev-extras"

.PHONY: debug
debug:
	dlv debug ./cmd/agent

.PHONY: run-dev-image
run-dev-image:
	docker exec -it signalfx-agent-dev bash 2>/dev/null || docker run --rm -it \
		--privileged \
		--net host \
		-p 6060:6060 \
		--name signalfx-agent-dev \
		-v $(PWD)/local-etc:/etc/signalfx \
		-v /:/bundle/hostfs:ro \
		-v /var/run/docker.sock:/var/run/docker.sock:ro \
		-v $(PWD):/go/src/github.com/signalfx/signalfx-agent:cached \
		-v $(PWD)/collectd:/usr/src/collectd:cached \
		-v /tmp/scratch:/tmp/scratch \
		signalfx-agent-dev /bin/bash

.PHONY: docs
docs:
	scripts/docs/make-monitor-docs
