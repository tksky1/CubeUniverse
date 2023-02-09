WORKDIR = $(abspath .)

DEVTMPDIR = $(WORKDIR)/dev-tmp
MAINDIR = $(WORKDIR)/main
BUILDERDIR = $(WORKDIR)/universeBuilder
DOCKERFILEDIR = $(WORKDIR)/dockerfiles

$(shell mkdir -p $(DEVTMPDIR))

BUILDSRC = $(shell find $(BUILDERDIR) -name "*.go")
MAINSRC = $(shell find $(MAINDIR) -name "*.go")
BUILDDOCKERFILE = $(DOCKERFILEDIR)/universeBuilder_dev.Dockerfile
MAINDOCKERFILE = $(DOCKERFILEDIR)/main_dev.Dockerfile

BUILDTAR = $(DEVTMPDIR)/builder-dev.tar
MAINTAR = $(DEVTMPDIR)/main-dev.tar

GO = go
DOCKER = docker
SCP = scp

build: $(BUILDSRC)
	$(GO) build $^ -o $(DEVTMPDIR)/main
	$(DOCKER) build -t builder-dev -f $(BUILDDOCKERFILE) $(WORKDIR)
	$(DOCKER) save builder-dev -o $(BUILDTAR)

main: $(MAINSRC)
	$(GO) build $^ -o $(DEVTMPDIR)/main
	$(DOCKER) build -t main-dev -f $(MAINDOCKERFILE) $(WORKDIR)
	$(DOCKER) save main-dev -o $(MAINTAR)

bscp: $(BUILDTAR)
	$(SCP) $^ 192.168.79.12:/home/node
	$(SCP) $^ 192.168.79.13:/home/node

mscp: $(MAINTAR)
	$(SCP) $^ 192.168.79.12:/home/node
	$(SCP) $^ 192.168.79.13:/home/node

.PHONY: build main bscp mscp
