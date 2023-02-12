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

NODE = 192.168.79.12 192.168.79.13

$(BUILDERTAR): $(BUILDERSRC)
	$(GO) build -o $(DEVTMPDIR)/main $^
	$(DOCKER) build -t builder-dev -f $(BUILDERDOCKERFILE) $(WORKDIR)
	$(DOCKER) save builder-dev -o $@

$(MAINTAR): $(MAINSRC)
	$(GO) build -o $(DEVTMPDIR)/main $^
	$(DOCKER) build -t main-dev -f $(MAINDOCKERFILE) $(WORKDIR)
	$(DOCKER) save main-dev -o $@

builder: $(BUILDERTAR)
	$(foreach node, $(NODE), $(CAT) $^ | $(SSH) $(node) 'docker load;')

main: $(MAINTAR)
	$(foreach node, $(NODE), $(CAT) $^ | $(SSH) $(node) 'docker load;')

KUBEADM = kubeadm
RM = rm

COMMAND = echo "echo 1 > /proc/sys/net/ipv4/ip_forward" >> /etc/rc.d/rc.local; \
echo 1 > /proc/sys/net/ipv4/ip_forward; \
chmod +x /etc/rc.d/rc.local

reset:
	$(KUBEADM) reset +y
	$(RM) -rf $(HOME)/.kube
	$(foreach node, $(NODE), $(SSH) $(node) '$(KUBEADM) reset +y;')
	$(foreach node, $(NODE), $(SSH) $(node) '$(COMMAND);')
	$(foreach node, $(NODE), $(SSH) $(node) '$(shell $(KUBEADM) init --config $(HOME)/kubeadm.yaml --upload-certs | tail -n2);') 

.PHONY: builder main reset
