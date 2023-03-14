WORKDIR = $(abspath .)

DEVTMPDIR = $(WORKDIR)/dev-tmp
MAINDIR = $(WORKDIR)/main
BUILDERDIR = $(WORKDIR)/universeBuilder
DOCKERFILEDIR = $(WORKDIR)/dockerfiles
OPERATORDIR = $(WORKDIR)/universeOperator

$(shell mkdir -p $(DEVTMPDIR))
export GO111MODULE = on
export GOPROXY = https://goproxy.cn
export GOPATH ?=

BUILDERSRC = $(shell find $(BUILDERDIR) -name "*.go")
MAINSRC = $(shell find $(MAINDIR) -name "*.go")
OPERATORSRC = $(shell find $(OPERATORDIR) -name "*.go")


DEVDOCKERFILE = $(DOCKERFILEDIR)/dev-debug.Dockerfile

BUILDERTAR = $(DEVTMPDIR)/builder-dev.tar
MAINTAR = $(DEVTMPDIR)/main-dev.tar
OPERATORTAR = $(DEVTMPDIR)/operator-dev.tar

GO = go
DOCKER = docker
SSH = ssh
CAT = cat


NODE = 192.168.79.12 192.168.79.13

$(BUILDERTAR): $(BUILDERSRC)
	$(GO) build -o $(DEVTMPDIR)/main $^
	$(DOCKER) build -t builder-dev -f $(DEVDOCKERFILE) $(WORKDIR)
	$(DOCKER) save builder-dev -o $@

$(MAINTAR): $(MAINSRC)
	$(GO) build -o $(DEVTMPDIR)/main $^
	$(DOCKER) build -t main-dev -f $(DEVDOCKERFILE) $(WORKDIR)
	$(DOCKER) save main-dev -o $@

$(OPERATORTAR): $(OPERATORSRC)
	$(GO) build -o $(DEVTMPDIR)/main $^
	$(DOCKER) build -t operator-dev -f $(DEVDOCKERFILE) $(WORKDIR)
	$(DOCKER) save operator-dev -o $@

builder: $(BUILDERTAR)
	$(foreach node, $(NODE), $(CAT) $^ | $(SSH) $(node) 'docker load';)

main: $(MAINTAR)
	$(foreach node, $(NODE), $(CAT) $^ | $(SSH) $(node) 'docker load';)

operator: $(OPERATORTAR)
	$(foreach node, $(NODE), $(CAT) $^ | $(SSH) $(node) 'docker load';)

KUBEADM = kubeadm
RM = rm
XARGS = xargs
TAIL = tail
TEE = tee

COMMAND = echo "echo 1 > /proc/sys/net/ipv4/ip_forward" >> /etc/rc.d/rc.local; \
echo 1 > /proc/sys/net/ipv4/ip_forward; \
chmod +x /etc/rc.d/rc.local; \
rm -rf /var/lib/rook

reset:
	$(KUBEADM) reset -f
	$(RM) -rf $(HOME)/.kube
	$(RM) -rf /var/lib/rook
	$(foreach node, $(NODE), $(SSH) $(node) '$(KUBEADM) reset -f';)
	$(foreach node, $(NODE), $(SSH) $(node) '$(COMMAND)';)
	$(KUBEADM) init --config $(HOME)/kubeadm.yaml --upload-certs | $(TEE) $(shell tty) | $(TAIL) -n2 > /tmp/kubeinit.sh
	$(foreach node, $(NODE), $(SSH) $(node) < /tmp/kubeinit.sh;)
	$(RM) -f /tmp/kubeinit
	kubectl taint nodes master node-role.kubernetes.io/master-
	kubectl create -f /home/master/kube-flannel.yml

.PHONY: builder main reset operator
