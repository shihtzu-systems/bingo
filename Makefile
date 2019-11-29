# touch app.version
# touch app.datestamp
# touch app.timestamp

version := $(shell cat app.version)
datestamp := $(shell cat app.datestamp)
timestamp := $(shell cat app.timestamp)

githubProject := github.com/shihtzu-systems
imageRepo := shihtzu
app := bingo

stamp:
	printf `/bin/date "+%Y%m%d"` > app.datestamp
	printf `/bin/date "+%H%M%S"` > app.timestamp
	printf "$(version)" > app.version

build: fmt vet test build-binary build-version-container archive

build-binary:
	GOOS=linux  GOARCH=amd64 go build -o bin/linux_amd64/$(app)  main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64/$(app) main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

build-version-container:
	docker build . \
		-t local/$(app) \
		-t $(imageRepo)/$(app):$(version)-on.$(datestamp).at.$(timestamp) \
		-t $(imageRepo)/$(app):$(version) \
		-t $(imageRepo)/$(app):latest

git-master-branch:
	git checkout master

git-commit:
	git add --all
	git commit

git-tag:
	git tag v$(version)

git-push-tags:
	git push origin --tags

git-push:
	git push origin

push-version-container: build-version-container
	docker push $(imageRepo)/$(app):$(version)-on.$(datestamp).at.$(timestamp)
	docker push $(imageRepo)/$(app):$(version)
	docker push $(imageRepo)/$(app):latest