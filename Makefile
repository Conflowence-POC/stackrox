ROX_PROJECT=apollo
TESTFLAGS=-race -p 4
BASE_DIR=$(CURDIR)
TAG=$(shell git describe --tags --abbrev=10 --dirty)

RELEASE_GOTAGS := release
ifdef CI
ifneq ($(CIRCLE_TAG),)
GOTAGS := $(RELEASE_GOTAGS)
endif
endif

null :=
space := $(null) $(null)
comma := ,

FORMATTING_FILES=$(shell find . -name vendor -prune -o -name generated -prune -o -name mocks -prune -o -name '*_easyjson.go' -prune -o -name '*.go' -print)

.PHONY: all
all: deps style test image

###########
## Style ##
###########
.PHONY: style
style: fmt imports lint vet blanks validateimports no-large-files only-store-storage-protos storage-protos-compatible no-unchecked-errors ui-lint qa-tests-style

.PHONY: qa-tests-style
qa-tests-style:
	@echo "+ $@"
	make -C qa-tests-backend/ style

.PHONY: ui-lint
ui-lint:
	@echo "+ $@"
	make -C ui lint

.PHONY: fmt
fmt:
	@echo "+ $@"
ifdef CI
		@echo "The environment indicates we are in CI; checking gofmt."
		@echo 'If this fails, run `make style`.'
		@$(eval FMT=`echo $(FORMATTING_FILES) | xargs gofmt -s -l`)
		@echo "gofmt problems in the following files, if any:"
		@echo $(FMT)
		@test -z "$(FMT)"
endif
	@echo $(FORMATTING_FILES) | xargs gofmt -s -l -w

.PHONY: imports
imports: deps volatile-generated-srcs
	@echo "+ $@"
ifdef CI
		@echo "The environment indicates we are in CI; checking goimports."
		@echo 'If this fails, run `make style`.'
		@$(eval IMPORTS=`echo $(FORMATTING_FILES) | xargs goimports -l`)
		@echo "goimports problems in the following files, if any:"
		@echo $(IMPORTS)
		@test -z "$(IMPORTS)"
endif
	@echo $(FORMATTING_FILES) | xargs goimports -w

.PHONY: validateimports
validateimports:
	@echo "+ $@"
	@go run $(BASE_DIR)/tools/validateimports/verify.go $(shell go list -e ./...)

.PHONY: no-large-files
no-large-files:
	@echo "+ $@"
	@$(BASE_DIR)/tools/large-git-files/find.sh

.PHONY: only-store-storage-protos
only-store-storage-protos:
	@echo "+ $@"
	@go run $(BASE_DIR)/tools/storedprotos/verify.go $(shell go list github.com/stackrox/rox/central/...)


.PHONY: no-unchecked-errors
no-unchecked-errors:
	@echo "+ $@"
	@go run $(BASE_DIR)/tools/uncheckederrors/cmd/main.go $(shell go list -e ./... | grep -v -e 'stackrox/rox/image')


PROTOLOCK_BIN := $(GOPATH)/bin/protolock
$(PROTOLOCK_BIN):
	@echo "+ $@"
	$(BASE_PATH)/scripts/go-get-version.sh github.com/viswajithiii/protolock 43bb8a9ba4e8de043a5ffacc64b1c38d95419e1d --skip-install
	mkdir -p $(GOPATH)/src/github.com/nilslice
	mv $(GOPATH)/src/github.com/viswajithiii/protolock $(GOPATH)/src/github.com/nilslice/protolock
	go install github.com/nilslice/protolock/...

.PHONY: storage-protos-compatible
storage-protos-compatible: $(PROTOLOCK_BIN)
	@echo "+ $@"
	@protolock status -lockdir=$(BASE_DIR)/proto/storage -protoroot=$(BASE_DIR)/proto/storage

.PHONY: lint
lint:
	@echo "+ $@"
	@set -e; git -C $(CURDIR) ls-files '*.go' | xargs -n 1 dirname | sort | uniq | while IFS='' read -r dir || [[ -n "$$dir" ]]; do golint -set_exit_status "$$dir"/*.go ; done

.PHONY: vet
vet:
	@echo "+ $@"
	@$(BASE_DIR)/tools/go-vet.sh -tags "$(subst $(comma),$(space),$(GOTAGS))" $(shell go list -e ./... | grep -v generated | grep -v vendor)
ifdef CI
	@echo "+ $@ ($(RELEASE_GOTAGS))"
	@$(BASE_DIR)/tools/go-vet.sh -tags "$(subst $(comma),$(space),$(RELEASE_GOTAGS))" $(shell go list -e ./... | grep -v generated | grep -v vendor)
endif

.PHONY: blanks
blanks:
	@echo "+ $@"
	@find . \( \( -name vendor -o -name generated \) -type d -prune \) -o \( -name \*.go -print0 \) | xargs -0 $(BASE_PATH)/tools/import_validate.py

.PHONY: dev
dev:
	@echo "+ $@"
	@go get -u golang.org/x/lint/golint
	@go get -u golang.org/x/tools/cmd/goimports
	@go get -u github.com/jstemmer/go-junit-report
	@go get -u github.com/golang/dep/cmd/dep


#####################################
## Generated Code and Dependencies ##
#####################################

PROTO_GENERATED_SRCS = $(GENERATED_PB_SRCS) $(GENERATED_API_GW_SRCS)

include make/protogen.mk

STRINGER_BIN := $(GOPATH)/bin/stringer
$(STRINGER_BIN):
	@echo "+ $@"
	@go get golang.org/x/tools/cmd/stringer

MOCKGEN_BIN := $(GOPATH)/bin/mockgen
$(MOCKGEN_BIN):
	@echo "+ $@"
	@$(BASE_PATH)/scripts/go-get-version.sh golang.org/x/tools e21233ffa6c386fc230b4328493f77af54ff9372 --skip-install
	@$(BASE_PATH)/scripts/go-get-version.sh github.com/golang/mock/mockgen 8a44ef6e8be577e050008c7886f24fc705d709fb

GENNY_BIN := $(GOPATH)/bin/genny
$(GENNY_BIN):
	@echo "+ $@"
	@$(BASE_PATH)/scripts/go-get-version.sh github.com/mauricelam/genny e937528877485c089aa62cfa9f60968749d650f1

PACKR_BIN := $(GOPATH)/bin/packr
$(PACKR_BIN):
	@echo "+ $@"
	@$(BASE_PATH)/scripts/go-get-version.sh github.com/gobuffalo/packr/packr 899fe0e4176fca9bca81763c810d74af82548c78

.PHONY: go-packr-srcs
go-packr-srcs: $(PACKR_BIN)
	@echo "+ $@"
	@packr

# For some reasons, a `packr clean` is much slower than the `find`. It also does not work.
.PHONY: clean-packr-srcs
clean-packr-srcs:
	@echo "+ $@"
	@find . -name '*-packr.go' -exec rm {} \;

EASYJSON_BIN := $(GOPATH)/bin/easyjson
$(EASYJSON_BIN):
	@echo "+ $@"
	@$(BASE_PATH)/scripts/go-get-version.sh github.com/mailru/easyjson/easyjson 60711f1a8329503b04e1c88535f419d0bb440bff

.PHONY: go-easyjson-srcs
go-easyjson-srcs: $(EASYJSON_BIN)
	@echo "+ $@"
	@easyjson -pkg pkg/docker/types.go

.PHONY: clean-easyjson-srcs
clean-easyjson-srcs:
	@echo "+ $@"
	@find . -name '*_easyjson.go' -exec rm {} \;

.PHONY: go-generated-srcs
go-generated-srcs: go-easyjson-srcs $(MOCKGEN_BIN) $(STRINGER_BIN) $(GENNY_BIN)
	@echo "+ $@"
	PATH=$(PATH):$(BASE_DIR)/tools/generate-helpers go generate ./...

proto-generated-srcs: $(PROTO_GENERATED_SRCS)
	@echo "+ $@"
	@touch proto-generated-srcs

# volatile-generated-srcs are all generated sources that are NOT committed
.PHONY: volatile-generated-srcs
volatile-generated-srcs: proto-generated-srcs go-packr-srcs

.PHONY: generated-srcs
generated-srcs: volatile-generated-srcs go-generated-srcs

.PHONY: clean-generated-srcs
clean-generated-srcs: clean-packr-srcs clean-easyjson-srcs
	@echo "+ $@"
	git clean -xdf generated

deps: Gopkg.toml Gopkg.lock proto-generated-srcs
	@echo "+ $@"
ifdef CI
	@# `dep check` exits with a nonzero code if there is a toml->lock mismatch.
	dep check -skip-vendor
endif
	@# `dep ensure` can be flaky sometimes, so try rerunning it if it fails.
	dep ensure || (rm -rf vendor .vendor-new && dep ensure)
	@touch deps

.PHONY: clean-deps
clean-deps:
	@echo "+ $@"
	@rm -f deps

###########
## Build ##
###########
PURE := --features=pure
RACE := --features=race
LINUX_AMD64 := --cpu=k8
VARIABLE_STAMPS := --workspace_status_command=$(BASE_DIR)/status.sh
BAZEL_OS=linux
ifeq ($(UNAME_S),Darwin)
    BAZEL_OS=darwin
endif
PLATFORMS := --platforms=@io_bazel_rules_go//go/toolchain:$(BAZEL_OS)_amd64

BAZEL_FLAGS := $(PURE) $(LINUX_AMD64) $(VARIABLE_STAMPS) --define gotags=$(GOTAGS)
cleanup:
	@echo "Total BUILD.bazel files deleted: "
	@git status --ignored --untracked-files=all --porcelain | grep '^\(!!\|??\) ' | cut -d' ' -f 2- | grep '\(/\|^\)BUILD\.bazel$$' | xargs rm -v | wc -l

.PHONY: gazelle
gazelle: deps volatile-generated-srcs cleanup
	bazel run //:gazelle -- -build_tags=$(GOTAGS)

cli: gazelle
	bazel build $(BAZEL_FLAGS) --platforms=@io_bazel_rules_go//go/toolchain:darwin_amd64 -- //roxctl
	bazel build $(BAZEL_FLAGS) --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 -- //roxctl
	bazel build $(BAZEL_FLAGS) --platforms=@io_bazel_rules_go//go/toolchain:windows_amd64 -- //roxctl

	# Copy the user's specific OS into gopath
	cp bazel-bin/roxctl/$(BAZEL_OS)_amd64_pure_stripped/roxctl $(GOPATH)/bin/roxctl
	chmod u+w $(GOPATH)/bin/roxctl

.PHONY: main-build
main-build: gazelle
	@echo "+ $@"
	bazel build $(BAZEL_FLAGS) \
		//central \
		//migrator \
		//sensor/kubernetes \
		//compliance/collection

.PHONY: scale-build
scale-build: gazelle
	@echo "+ $@"
	bazel build $(BAZEL_FLAGS) \
		//scale/mocksensor \
		//scale/mockcollector \
		//scale/profiler

.PHONY: genericserver-build
genericserver-build: gazelle
	@echo "+ $@"
	bazel build $(BAZEL_FLAGS) \
		//genericserver

.PHONY: mock-grpc-server-build
mock-grpc-server-build: gazelle
	@echo "+ $@"
	bazel build $(BAZEL_FLAGS) \
		//integration-tests/mock-grpc-server

.PHONY: gendocs
gendocs: $(GENERATED_API_DOCS)
	@echo "+ $@"

# We don't need to do anything here, because the $(MERGED_API_SWAGGER_SPEC) target already performs validation.
.PHONY: swagger-docs
swagger-docs: $(MERGED_API_SWAGGER_SPEC)
	@echo "+ $@"

.PHONY: bazel-test
bazel-test: gazelle
	-rm vendor/github.com/coreos/pkg/BUILD
	-rm vendor/github.com/cloudflare/cfssl/script/BUILD
	-rm vendor/github.com/grpc-ecosystem/grpc-gateway/BUILD
	@# Be careful if you add action_env arguments; their values can invalidate cached
	@# test results. See https://github.com/bazelbuild/bazel/issues/2574#issuecomment-320006871.
	bazel coverage $(BAZEL_FLAGS) $(RACE) \
	    --test_output=errors \
	    -- \
	    //... -proto/... -tests/... -vendor/...

.PHONY: ui-test
ui-test:
	@# UI tests don't work in Bazel yet.
	make -C ui test

.PHONY: test
test: bazel-test ui-test collector-tag

.PHONY: integration-unit-tests
integration-unit-tests: gazelle
	 go test -tags=integration $(shell go list ./... | grep  "registries\|scanners\|notifiers")

upload-coverage:
	@# 'mode: set' is repeated in each coverage file, but Coveralls only wants it
	@# exactly once at the head of the file.
	@# We might be able to use Coveralls parallel builds to resolve this:
	@#     https://docs.coveralls.io/parallel-build-webhook

	@echo 'mode: set' > combined_coverage.dat
	@find ./bazel-testlogs/ -name 'coverage.dat' | xargs -I {} cat "{}" | grep -v 'mode: set' | grep -v vendor >> combined_coverage.dat
	goveralls -coverprofile="combined_coverage.dat" -ignore 'central/graphql/resolvers/generated.go,generated/storage/*,generated/*/*/*' -service=circle-ci -repotoken="$$COVERALLS_REPO_TOKEN"

.PHONY: coverage
coverage:
	@echo "+ $@"
	@go test -cover -coverprofile coverage.out $(shell go list -e ./... | grep -v /tests)
	@go tool cover -html=coverage.out -o coverage.html

###########
## Image ##
###########

# Exists for compatibility reasons. Please consider migrating to using `make main-image`.
.PHONY: image
image: main-image monitoring-image

.PHONY: monitoring-image
monitoring-image:
	docker build -t stackrox/monitoring:$(TAG) monitoring

.PHONY: main-image
main-image: cli main-build clean-image $(MERGED_API_SWAGGER_SPEC)
	make -C ui build
	make docker-build-main-image

# This target copies compiled artifacts into the expected locations and
# runs the docker build.
# Please DO NOT invoke this target directly unless you know what you're doing;
# you probably want to run `make main-image`. This target is only in Make for convenience;
# it assumes the caller has taken care of the dependencies, and does not
# declare its dependencies explicitly.
.PHONY: docker-build-main-image
docker-build-main-image:
	cp -r ui/build image/ui/
	cp bazel-bin/central/linux_amd64_pure_stripped/central image/bin/central
	cp bazel-bin/roxctl/linux_amd64_pure_stripped/roxctl image/bin/roxctl-linux
	cp bazel-bin/roxctl/darwin_amd64_pure_stripped/roxctl image/bin/roxctl-darwin
	cp bazel-bin/roxctl/windows_amd64_pure_stripped/roxctl.exe image/bin/roxctl-windows.exe
	cp bazel-bin/migrator/linux_amd64_pure_stripped/migrator image/bin/migrator
	cp bazel-bin/sensor/kubernetes/linux_amd64_pure_stripped/kubernetes image/bin/kubernetes-sensor
	cp bazel-bin/compliance/collection/linux_amd64_pure_stripped/collection image/bin/compliance

ifdef CI
	@[ -f image/NOTICE.txt ] || { echo "image/NOTICE.txt file not found! It is required for CI-built images."; exit 1; }
else
	@[ -f image/NOTICE.txt ] || touch image/NOTICE.txt
endif
	@[ -d image/docs ] || { echo "Generated docs not found in image/docs. They are required for build."; exit 1; }
	docker build -t stackrox/main:$(TAG) image/
	@echo "Built main image with tag: $(TAG)"
	@echo "You may wish to:       export MAIN_IMAGE_TAG=$(TAG)"

.PHONY: scale-image
scale-image: scale-build clean-image
	cp bazel-bin/scale/mocksensor/linux_amd64_pure_stripped/mocksensor scale/image/bin/mocksensor
	cp bazel-bin/scale/mockcollector/linux_amd64_pure_stripped/mockcollector scale/image/bin/mockcollector
	cp bazel-bin/scale/profiler/linux_amd64_pure_stripped/profiler scale/image/bin/profiler
	chmod +w scale/image/bin/*
	docker build -t stackrox/scale:$(TAG) -f scale/image/Dockerfile scale

genericserver-image: genericserver-build
	mkdir genericserver/bin
	cp bazel-bin/genericserver/linux_amd64_pure_stripped/genericserver genericserver/bin/genericserver
	chmod +w genericserver/bin/genericserver
	docker build -t stackrox/genericserver:latest -f genericserver/Dockerfile genericserver

.PHONY: mock-grpc-server-image
mock-grpc-server-image: mock-grpc-server-build clean-image
	cp bazel-bin/integration-tests/mock-grpc-server/linux_amd64_pure_stripped/mock-grpc-server integration-tests/mock-grpc-server/image/bin/mock-grpc-server
	docker build -t stackrox/grpc-server:$(TAG) integration-tests/mock-grpc-server/image

###########
## Clean ##
###########
.PHONY: clean
clean: clean-image clean-generated-srcs
	@echo "+ $@"

.PHONY: clean-image
clean-image:
	@echo "+ $@"
	git clean -xf image/bin
	git clean -xdf image/ui image/docs
	git clean -xf integration-tests/mock-grpc-server/image/bin/mock-grpc-server

.PHONY: tag
tag:
ifdef COMMIT
	@git describe $(COMMIT) --tags --abbrev=10
else
	@echo $(TAG)
endif

ossls-audit:
	ossls -audit

ossls-notice:
	ossls -notice | tee image/NOTICE.txt

.PHONY: collector-tag
collector-tag:
	@cat COLLECTOR_VERSION
