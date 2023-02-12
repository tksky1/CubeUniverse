WORKDIR = $(abspath .)

DEVTMPDIR = $(WORKDIR)/dev-tmp
MAINDIR = $(WORKDIR)/main
BUILDERERDIR = $(WORKDIR)/universeBuilder
DOCKERFILEDIR = $(WORKDIR)/dockerfiles

$(shell mkdir -p $(DEVTMPDIR))
$(shell export GO111MODULE="on")
$(shell export GOPROXY="https://goproxy.cn")
$(shell export GOPATH="")

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
	$(foreach node, $(NODE), $(CAT) $^ | $(SSH) $(node) 'docker load';)

main: $(MAINTAR)
	$(foreach node, $(NODE), $(CAT) $^ | $(SSH) $(node) 'docker load';)

KUBEADM = kubeadm
RM = rm
XARGS = xargs
TAIL = tail
TEE = tee

COMMAND = echo "echo 1 > /proc/sys/net/ipv4/ip_forward" >> /etc/rc.d/rc.local; \
echo 1 > /proc/sys/net/ipv4/ip_forward; \
chmod +x /etc/rc.d/rc.local

reset:
	$(KUBEADM) reset -f
	$(RM) -rf $(HOME)/.kube
	$(foreach node, $(NODE), $(SSH) $(node) '$(KUBEADM) reset -f';)
	$(foreach node, $(NODE), $(SSH) $(node) '$(COMMAND)';)
	$(KUBEADM) init --config $(HOME)/kubeadm.yaml --upload-certs | $(TEE) | $(TAIL) -n2 > /tmp/kubeinit
	$(foreach node, $(NODE), $(CAT) /tmp/kubeinit | $(XARGS) $(SSH) $(node);)
	$(RM) -f /tmp/kubeinit

.PHONY: builder main reset
