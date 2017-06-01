
vet: bootstrap
	@echo ">> vetting code"
	go vet `glide novendor`

lint: bootstrap
	@echo ">> linting code"
	golint `glide novendor`

watch: get-deps
	@echo ">> starting goconvey"
	goconvey

get-deps: bootstrap
	@echo ">> installing dependencies"
	glide install

$(GOPATH)/bin/glide:
	@echo ">> installing glide"
	go get -u github.com/Masterminds/glide

$(GOPATH)/src/github.com/golang/lint/golint:
	@echo ">> installing golint"
	go get -u github.com/golang/lint/golint

$(GOPATH)/bin/goconvey:
	@echo ">> installing goconvey"
	go get -u github.com/smartystreets/goconvey

bootstrap: $(GOPATH)/bin/glide $(GOPATH)/src/github.com/golang/lint/golint $(GOPATH)/bin/goconvey

test: get-deps lint vet
	go test -v -cover `glide novendor` -ldflags '$(TESTLDFLAGS)'

.PHONY: test vet lint watch get-deps bootstrap
