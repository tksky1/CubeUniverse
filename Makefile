WORKDIR = $(abspath .)

DEVTMPDIR = $(WORKDIR)/dev-tmp
MAINDIR = $(WORKDIR)/main
BUILDERERDIR = $(WORKDIR)/universeBuilder
DOCKERFILEDIR = $(WORKDIR)/dockerfiles

$(shell mkdir -p $(DEVTMPDIR))

BUILDERSRC = $(shell find $(BUILDERERDIR) -name "*.go")
MAINSRC = $(shell find $(MAINDIR) -name "*.go")
BUILDERDOCKERFILE = $(DOCKERFILEDIR)/universeBuilder_dev.Dockerfile
MAINDOCKERFILE = $(DOCKERFILEDIR)/main_dev.Dockerfile

BUILDERTAR = $(DEVTMPDIR)/builder-dev.tar
MAINTAR = $(DEVTMPDIR)/main-dev.tar

GO = go
DOCKER = docker
SSH = ssh
CAT = cat

$(BUILDERTAR): $(BUILDERSRC)
	$(GO) build -o $(DEVTMPDIR)/main $^
	$(DOCKER) build -t builder-dev -f $(BUILDERDOCKERFILE) $(WORKDIR)
	$(DOCKER) save builder-dev -o $@

$(MAINTAR): $(MAINSRC)
	$(GO) build -o $(DEVTMPDIR)/main $^
	$(DOCKER) build -t main-dev -f $(MAINDOCKERFILE) $(WORKDIR)
	$(DOCKER) save main-dev -o $@

builder: $(BUILDERTAR)
	$(CAT) $^ | $(SSH) 192.168.79.12 'docker load'
	$(CAT) $^ | $(SSH) 192.168.79.13 'docker load'

main: $(MAINTAR)
	$(CAT) $^ | $(SSH) 192.168.79.12 'docker load'
	$(CAT) $^ | $(SSH) 192.168.79.13 'docker load'

.PHONY: builder main
