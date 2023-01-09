PROJECT := patrol/patrol-api

REPO := appfactory.arfa.wise-paas.com/

.PHONY: gomod, build, build-docker, build-exe

VERSION := ${VERSION}

# Build dockr image
build-docker: gomod
	@echo "build docker"
	docker build -t   ${REPO}${PROJECT}:${VERSION} .
	docker push ${REPO}${PROJECT}:${VERSION}


# Build linux binary
build: gomod
	@echo "building"
	go build -o ./bin/partol-server  ./cmd/partol-server

# Build windows binary
build-exe:
	@echo "building exe"
	GOOS=windows go build -o ./bin/partol-server.exe  ./cmd/partol-server

# Download go modules
gomod:
	go mod download

# Clean binary files and temporary files
clean:
	$(RM) -rf bin
