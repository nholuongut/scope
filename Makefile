.PHONY: all cri deps static clean realclean client-lint client-test client-sync backend frontend shell lint ui-upload

# If you can use Docker without being root, you can `make SUDO= <target>`
SUDO=$(shell docker info >/dev/null 2>&1 || echo "sudo -E")
DOCKERHUB_USER=nholuongutworks
SCOPE_EXE=prog/scope
SCOPE_EXPORT=scope.tar
CLOUD_AGENT_EXPORT=cloud-agent.tar
SCOPE_UI_BUILD_IMAGE=$(DOCKERHUB_USER)/scope-ui-build
SCOPE_UI_BUILD_UPTODATE=.scope_ui_build.uptodate
SCOPE_BACKEND_BUILD_IMAGE=$(DOCKERHUB_USER)/scope-backend-build
SCOPE_BACKEND_BUILD_UPTODATE=.scope_backend_build.uptodate
SCOPE_VERSION=$(shell git rev-parse --short HEAD)
GIT_REVISION=$(shell git rev-parse HEAD)
nholuongutNET_VERSION=2.1.3
RUNSVINIT=vendor/runsvinit/runsvinit
CODECGEN_DIR=vendor/github.com/ugorji/go/codec/codecgen
CODECGEN_EXE=$(CODECGEN_DIR)/bin/codecgen_$(shell go env GOHOSTOS)_$(shell go env GOHOSTARCH)
CODECGEN_UID=0
GET_CODECGEN_DEPS=$(shell find $(1) -maxdepth 1 -type f -name '*.go' -not -name '*_test.go' -not -name '*.codecgen.go' -not -name '*.generated.go')
CODECGEN_TARGETS=report/report.codecgen.go render/detailed/detailed.codecgen.go
RM=--rm
RUN_FLAGS=-ti
BUILD_IN_CONTAINER=true
GO_ENV=GOGC=off
GO_BUILD_INSTALL_DEPS=-i
GO_BUILD_TAGS='netgo unsafe'
GO_BUILD_FLAGS=$(GO_BUILD_INSTALL_DEPS) -ldflags "-extldflags \"-static\" -X main.version=$(SCOPE_VERSION) -s -w" -tags $(GO_BUILD_TAGS)
GOOS=$(shell go tool dist env | grep GOOS | sed -e 's/GOOS="\(.*\)"/\1/')

ifeq ($(GOOS),linux)
GO_ENV+=CGO_ENABLED=1
endif

ifeq ($(GOARCH),arm)
ARM_CC=CC=/usr/bin/arm-linux-gnueabihf-gcc
endif

GO=env $(GO_ENV) $(ARM_CC) go

NO_CROSS_COMP=unset GOOS GOARCH
GO_HOST=$(NO_CROSS_COMP); env $(GO_ENV) go
WITH_GO_HOST_ENV=$(NO_CROSS_COMP); $(GO_ENV)
IMAGE_TAG=$(shell ./tools/image-tag)

all: $(SCOPE_EXPORT)

update-cri:
	curl https://raw.githubusercontent.com/kubernetes/kubernetes/master/pkg/kubelet/apis/cri/runtime/v1alpha2/api.proto > cri/runtime/api.proto

protoc-gen-gofast:
	@go get -u -v github.com/gogo/protobuf/protoc-gen-gofast

# Use cri target to download latest cri proto files and regenerate CRI runtime files. 
cri: update-cri protoc-gen-gofast
	@cd $(GOPATH)/src;protoc --proto_path=$(GOPATH)/src --gofast_out=plugins=grpc:. github.com/nholuongut/scope/cri/runtime/api.proto

docker/nholuongut:
	curl -L https://github.com/nholuongutworks/nholuongut/releases/download/v$(nholuongutNET_VERSION)/nholuongut -o docker/nholuongut
	chmod u+x docker/nholuongut

docker/nholuongututil:
	$(SUDO) docker run --rm  --entrypoint=cat nholuongutworks/nholuongutexec:$(nholuongutNET_VERSION) /usr/bin/nholuongututil > $@
	chmod +x $@

docker/%: %
	cp $* docker/

%.tar: docker/Dockerfile.%
	$(SUDO) docker build --build-arg=revision=$(GIT_REVISION) -t $(DOCKERHUB_USER)/$* -f $< docker/
	$(SUDO) docker tag $(DOCKERHUB_USER)/$* $(DOCKERHUB_USER)/$*:$(IMAGE_TAG)
	$(SUDO) docker save $(DOCKERHUB_USER)/$*:latest > $@

$(CLOUD_AGENT_EXPORT): docker/Dockerfile.cloud-agent docker/$(SCOPE_EXE) docker/nholuongut docker/nholuongututil

$(SCOPE_EXPORT): docker/Dockerfile.scope $(CLOUD_AGENT_EXPORT) docker/$(RUNSVINIT) docker/demo.json docker/run-app docker/run-probe docker/entrypoint.sh

$(RUNSVINIT): vendor/runsvinit/*.go

$(SCOPE_EXE): $(shell find ./ -path ./vendor -prune -o -type f -name '*.go') prog/staticui/staticui.go prog/externalui/externalui.go $(CODECGEN_TARGETS)

report/report.codecgen.go: $(call GET_CODECGEN_DEPS,report/)
render/detailed/detailed.codecgen.go: $(call GET_CODECGEN_DEPS,render/detailed/)
static: prog/staticui/staticui.go prog/externalui/externalui.go
prog/staticui/staticui.go: client/build/index.html
prog/externalui/externalui.go: client/build-external/index.html

ifeq ($(BUILD_IN_CONTAINER),true)

$(SCOPE_EXE) $(RUNSVINIT) lint tests shell prog/staticui/staticui.go prog/externalui/externalui.go: $(SCOPE_BACKEND_BUILD_UPTODATE)
	@mkdir -p $(shell pwd)/.pkg
	$(SUDO) docker run $(RM) $(RUN_FLAGS) \
		-v $(shell pwd):/go/src/github.com/nholuongut/scope \
		-v $(shell pwd)/.pkg:/go/pkg \
		--net=host \
		-e GOARCH -e GOOS -e CIRCLECI -e CIRCLE_BUILD_NUM -e CIRCLE_NODE_TOTAL \
		-e CIRCLE_NODE_INDEX -e COVERDIR -e SLOW -e TESTDIRS \
		$(SCOPE_BACKEND_BUILD_IMAGE) SCOPE_VERSION=$(SCOPE_VERSION) GO_BUILD_INSTALL_DEPS=$(GO_BUILD_INSTALL_DEPS) CODECGEN_UID=$(CODECGEN_UID) $@

else

$(SCOPE_EXE):
	time $(GO) build $(GO_BUILD_FLAGS) -o $@ ./$(@D)
	@strings $@ | grep cgo_stub\\\.go >/dev/null || { \
	        rm $@; \
	        echo "\nYour go standard library was built without the 'netgo' build tag."; \
	        echo "To fix that, run"; \
	        echo "    sudo go clean -i net"; \
	        echo "    sudo go install -tags netgo std"; \
	        false; \
	    }

%.codecgen.go: $(CODECGEN_EXE)
	rm -f $@; $(GO_HOST) build $(GO_BUILD_FLAGS) ./$(@D) # workaround for https://github.com/ugorji/go/issues/145
	cd $(@D) && $(WITH_GO_HOST_ENV) $(shell pwd)/$(CODECGEN_EXE) -d $(CODECGEN_UID) -rt $(GO_BUILD_TAGS) -u -o $(@F) $(notdir $(call GET_CODECGEN_DEPS,$(@D)))

$(CODECGEN_EXE): $(CODECGEN_DIR)/*.go
	mkdir -p $(@D)
	$(GO_HOST) build $(GO_BUILD_FLAGS) -o $@ ./$(CODECGEN_DIR)

$(RUNSVINIT):
	time $(GO) build $(GO_BUILD_FLAGS) -o $@ ./$(@D)

shell:
	/bin/bash

tests: $(CODECGEN_TARGETS) prog/staticui/staticui.go prog/externalui/externalui.go
	./tools/test -no-go-get -tags $(GO_BUILD_TAGS)

lint:
	./tools/lint
	./tools/shell-lint tools

prog/staticui/staticui.go:
	mkdir -p prog/staticui
	esc -o $@ -pkg staticui -prefix client/build client/build

prog/externalui/externalui.go:
	mkdir -p prog/externalui
	esc -o $@ -pkg externalui -prefix client/build-external -include '\.html$$' client/build-external

endif

ifeq ($(BUILD_IN_CONTAINER),true)

SCOPE_UI_TOOLCHAIN=.cache/build_node_modules
SCOPE_UI_TOOLCHAIN_UPTODATE=$(SCOPE_UI_TOOLCHAIN)/.uptodate

$(SCOPE_UI_TOOLCHAIN_UPTODATE): client/yarn.lock $(SCOPE_UI_BUILD_UPTODATE)
	mkdir -p $(SCOPE_UI_TOOLCHAIN) client/node_modules
	if test "true" != "$(SCOPE_SKIP_UI_ASSETS)"; then \
		$(SUDO) docker run $(RM) $(RUN_FLAGS) \
			-v $(shell pwd)/.cache:/home/nholuongut/scope/.cache \
			-v $(shell pwd)/client:/home/nholuongut/scope/client \
			-v $(shell pwd)/$(SCOPE_UI_TOOLCHAIN):/home/nholuongut/scope/client/node_modules \
			-w /home/nholuongut/scope/client \
			$(SCOPE_UI_BUILD_IMAGE) yarn install; \
	fi
	touch $(SCOPE_UI_TOOLCHAIN_UPTODATE)

client/build/index.html: $(shell find client/app -type f) $(SCOPE_UI_TOOLCHAIN_UPTODATE)
	mkdir -p client/build
	if test "true" != "$(SCOPE_SKIP_UI_ASSETS)"; then \
		$(SUDO) docker run $(RM) $(RUN_FLAGS) \
			-v $(shell pwd)/.cache:/home/nholuongut/scope/.cache \
			-v $(shell pwd)/client:/home/nholuongut/scope/client \
			-v $(shell pwd)/$(SCOPE_UI_TOOLCHAIN):/home/nholuongut/scope/client/node_modules \
			-w /home/nholuongut/scope/client \
			$(SCOPE_UI_BUILD_IMAGE) yarn run build; \
	fi

client/build-external/index.html: $(shell find client/app -type f) $(SCOPE_UI_TOOLCHAIN_UPTODATE)
	mkdir -p client/build-external
	if test "true" != "$(SCOPE_SKIP_UI_ASSETS)"; then \
		$(SUDO) docker run $(RM) $(RUN_FLAGS) \
			-v $(shell pwd)/.cache:/home/nholuongut/scope/.cache \
			-v $(shell pwd)/client:/home/nholuongut/scope/client \
			-v $(shell pwd)/$(SCOPE_UI_TOOLCHAIN):/home/nholuongut/scope/client/node_modules \
			-w /home/nholuongut/scope/client \
			$(SCOPE_UI_BUILD_IMAGE) yarn run build-external; \
	fi

client-test: $(shell find client/app/scripts -type f) $(SCOPE_UI_TOOLCHAIN_UPTODATE)
	$(SUDO) docker run $(RM) $(RUN_FLAGS) \
		-v $(shell pwd)/.cache:/home/nholuongut/scope/.cache \
		-v $(shell pwd)/client:/home/nholuongut/scope/client \
		-v $(shell pwd)/$(SCOPE_UI_TOOLCHAIN):/home/nholuongut/scope/client/node_modules \
		-w /home/nholuongut/scope/client \
		$(SCOPE_UI_BUILD_IMAGE) yarn test

client-lint: $(SCOPE_UI_TOOLCHAIN_UPTODATE)
	$(SUDO) docker run $(RM) $(RUN_FLAGS) \
		-v $(shell pwd)/.cache:/home/nholuongut/scope/.cache \
		-v $(shell pwd)/client:/home/nholuongut/scope/client \
		-v $(shell pwd)/$(SCOPE_UI_TOOLCHAIN):/home/nholuongut/scope/client/node_modules \
		-w /home/nholuongut/scope/client \
		$(SCOPE_UI_BUILD_IMAGE) yarn run lint

client-start: $(SCOPE_UI_TOOLCHAIN_UPTODATE)
	$(SUDO) docker run $(RM) $(RUN_FLAGS) --net=host \
		-v $(shell pwd)/.cache:/home/nholuongut/scope/.cache \
		-v $(shell pwd)/client:/home/nholuongut/scope/client \
		-v $(shell pwd)/$(SCOPE_UI_TOOLCHAIN):/home/nholuongut/scope/client/node_modules \
		-e WEBPACK_SERVER_HOST \
		-w /home/nholuongut/scope/client \
		$(SCOPE_UI_BUILD_IMAGE) yarn start

client/bundle/nholuongut-scope.tgz: $(shell find client/app -type f) $(SCOPE_UI_TOOLCHAIN_UPTODATE)
	$(sudo) docker run $(RUN_FLAGS) \
		-v $(shell pwd)/.cache:/home/nholuongut/scope/.cache \
		-v $(shell pwd)/client:/home/nholuongut/scope/client \
		-v $(shell pwd)/$(SCOPE_UI_TOOLCHAIN):/home/nholuongut/scope/client/node_modules \
		-v $(shell pwd)/tmp:/home/nholuongut/tmp \
		-w /home/nholuongut/scope/client \
		$(SCOPE_UI_BUILD_IMAGE) yarn run bundle

else

SCOPE_UI_TOOLCHAIN=client/node_modules
SCOPE_UI_TOOLCHAIN_UPTODATE=$(SCOPE_UI_TOOLCHAIN)/.uptodate

$(SCOPE_UI_TOOLCHAIN_UPTODATE): client/yarn.lock
	if test "true" = "$(SCOPE_SKIP_UI_ASSETS)"; then mkdir -p $(SCOPE_UI_TOOLCHAIN); else cd client && yarn install; fi
	touch $(SCOPE_UI_TOOLCHAIN_UPTODATE)

client/build/index.html: $(SCOPE_UI_TOOLCHAIN_UPTODATE)
	mkdir -p client/build
	if test "true" != "$(SCOPE_SKIP_UI_ASSETS)"; then cd client && yarn run build; fi

client/build-external/index.html: $(SCOPE_UI_TOOLCHAIN_UPTODATE)
	mkdir -p client/build-external
	if test "true" != "$(SCOPE_SKIP_UI_ASSETS)"; then cd client && yarn run build-external; fi

endif

$(SCOPE_UI_BUILD_UPTODATE): client/Dockerfile client/package.json client/webpack.local.config.js client/webpack.production.config.js client/server.js client/.eslintrc
	$(SUDO) docker build -t $(SCOPE_UI_BUILD_IMAGE) client
	$(SUDO) docker tag $(SCOPE_UI_BUILD_IMAGE) $(SCOPE_UI_BUILD_IMAGE):$(IMAGE_TAG)
	touch $@

$(SCOPE_BACKEND_BUILD_UPTODATE): backend/*
	$(SUDO) docker build -t $(SCOPE_BACKEND_BUILD_IMAGE) backend
	$(SUDO) docker tag $(SCOPE_BACKEND_BUILD_IMAGE) $(SCOPE_BACKEND_BUILD_IMAGE):$(IMAGE_TAG)
	touch $@

ui-upload: client/build-external/index.html
	AWS_ACCESS_KEY_ID=$$UI_BUCKET_KEY_ID \
	AWS_SECRET_ACCESS_KEY=$$UI_BUCKET_KEY_SECRET \
	aws s3 cp client/build-external/ s3://static.nholuongut.works/scope-ui/ --recursive --exclude '*.html'

ui-pkg-upload: client/bundle/nholuongut-scope.tgz
	AWS_ACCESS_KEY_ID=$$UI_BUCKET_KEY_ID \
	AWS_SECRET_ACCESS_KEY=$$UI_BUCKET_KEY_SECRET \
	aws s3 cp client/bundle/nholuongut-scope.tgz s3://nholuongutworks-js-modules/nholuongut-scope/$(shell echo $(SCOPE_VERSION))/nholuongut-scope.tgz

# We don't rmi images here; rm'ing the .uptodate files is enough to
# get the build images rebuilt, and rm'ing the scope exe is enough to
# get the main images rebuilt.
#
# rmi'ng images is desirable sometimes. Invoke `realclean` for that.
clean:
	$(GO) clean ./...
	rm -rf $(SCOPE_EXPORT) $(SCOPE_UI_BUILD_UPTODATE) $(SCOPE_UI_TOOLCHAIN_UPTODATE) $(SCOPE_BACKEND_BUILD_UPTODATE) \
		$(SCOPE_EXE) $(RUNSVINIT) prog/staticui/staticui.go prog/externalui/externalui.go client/build/*.js client/build-external/*.js docker/nholuongut .pkg \
		$(CODECGEN_TARGETS) $(CODECGEN_DIR)/bin

clean-codecgen:
	rm -rf $(CODECGEN_TARGETS) $(CODECGEN_DIR)/bin

# clean + rmi
#
# Removal of the main images ensures that a subsequent build rebuilds
# all their layers, in particular layers installing packages.
# Crucially, we also remove the *base* images, so their latest
# versions will be pulled.
#
# Doing this is important for release builds.
realclean: clean
	rm -rf $(SCOPE_UI_TOOLCHAIN)
	$(SUDO) docker rmi -f $(SCOPE_UI_BUILD_IMAGE) $(SCOPE_BACKEND_BUILD_IMAGE) \
		$(DOCKERHUB_USER)/scope $(DOCKERHUB_USER)/cloud-agent \
		$(DOCKERHUB_USER)/scope:$(IMAGE_TAG) $(DOCKERHUB_USER)/cloud-agent:$(IMAGE_TAG) \
		nholuongutworks/nholuongutexec:$(nholuongutNET_VERSION) \
		ubuntu:yakkety alpine:3.5 node:6.9.0 2>/dev/null || true

# Dependencies are intentionally build without enforcing any tags
# since they are build on the host
deps:
	$(GO) get -u -f \
		github.com/FiloSottile/gvt \
		github.com/mattn/goveralls \
		github.com/nholuongutworks/github-release \
		github.com/2opremio/trifles/wscat

# This target is only intended for use in Netlify CI environment for generating preview pages on feature branches and pull requests.
# We need to obtain website repository (checked out under `site-build`) and place `site` directory into the context (`site-build/_nholuongut_net_docs`).
# We then run make in `site-build` and Netlify will publish the output (`site-build/_site`).
netlify-site-preview:
	@mkdir -p site-build
	@curl --user $(WEBSITE_GITHUB_USER) --silent 'https://codeload.github.com/nholuongutworks/website-next/tar.gz/$(WEBSITE_BRANCH)' \
	  | tar --extract --gunzip --directory site-build --strip 1
	@cp -r site site-build/_nholuongut_scope_docs
	@$(MAKE) -C site-build netlify_ensure_install
	@$(MAKE) -C site-build BUILD_ENV=netlify
