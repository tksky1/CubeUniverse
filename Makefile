WORKDIR = $(abspath .)

DEVTMPDIR = $(WORKDIR)/dev-tmp
MAINDIR = $(WORKDIR)/main
BUILDERDIR = $(WORKDIR)/universeBuilder
DOCKERFILEDIR = $(WORKDIR)/dockerfiles
OPERATORDIR = $(WORKDIR)/universeOperator
BACKENDDIR = $(WORKDIR)/control-backend/main

$(shell mkdir -p $(DEVTMPDIR))
export GO111MODULE = on
export GOPROXY = https://goproxy.cn
export GOPATH = /usr/local/go

BUILDERSRC = $(shell find $(BUILDERDIR) -name "*.go")
MAINSRC = $(shell find $(MAINDIR) -name "*.go")
OPERATORSRC = $(shell find $(OPERATORDIR) -name "*.go")
BACKENDSRC= $(shell find $(BACKENDDIR) -name "*.go")


DEVDOCKERFILE = $(DOCKERFILEDIR)/dev-debug.Dockerfile

BUILDERTAR = $(DEVTMPDIR)/builder-dev.tar
MAINTAR = $(DEVTMPDIR)/main-dev.tar
OPERATORTAR = $(DEVTMPDIR)/operator-dev.tar
BACKENDTAR = $(DEVTMPDIR)/backend-dev.tar

GO = go
DOCKER = docker
SSH = ssh
CAT = cat


NODE = 192.168.79.12 192.168.79.13

$(BUILDERTAR): $(BUILDERSRC)
	$(GO) mod download
	$(GO) build -o $(DEVTMPDIR)/main $^
	$(DOCKER) build -t builder-dev -f $(DEVDOCKERFILE) $(WORKDIR)
	$(DOCKER) save builder-dev -o $@

$(MAINTAR): $(MAINSRC)
	$(GO) mod download
	$(GO) build -o $(DEVTMPDIR)/main $^
	$(DOCKER) build -t main-dev -f $(DEVDOCKERFILE) $(WORKDIR)
	$(DOCKER) save main-dev -o $@

$(OPERATORTAR): $(OPERATORSRC)
	$(GO) mod download
	$(GO) build -o $(DEVTMPDIR)/main $^
	$(DOCKER) build -t operator-dev -f $(DEVDOCKERFILE) $(WORKDIR)
	$(DOCKER) save operator-dev -o $@

$(BACKENDTAR): $(BACKENDSRC)
	cd control-backend/main && $(GO) mod download
	cd control-backend/main && $(GO) build -o $(DEVTMPDIR)/main $(BACKENDDIR)/control.go $(BACKENDDIR)/routes.go $(BACKENDDIR)/test.go
	$(DOCKER) build -t backend-dev -f $(DEVDOCKERFILE) $(WORKDIR)
	$(DOCKER) save backend-dev -o $@

builder: $(BUILDERTAR)
	$(foreach node, $(NODE), $(CAT) $^ | $(SSH) $(node) 'docker load';)
	rm -f dev-tmp/*.ta

main: $(MAINTAR)
	$(foreach node, $(NODE), $(CAT) $^ | $(SSH) $(node) 'docker load';)
	rm -f dev-tmp/*.tar

operator: $(OPERATORTAR)
	$(foreach node, $(NODE), $(CAT) $^ | $(SSH) $(node) 'docker load';)
	rm -f dev-tmp/*.tar

backend: $(BACKENDTAR)
	$(foreach node, $(NODE), $(CAT) $^ | $(SSH) $(node) 'docker load';)
	rm -f dev-tmp/*.tar

KUBEADM = kubeadm
RM = rm
XARGS = xargs
TAIL = tail
TEE = tee
SCP = scp
KUBECTL = kubectl

COMMAND = echo "echo 1 > /proc/sys/net/ipv4/ip_forward" >> /etc/rc.d/rc.local; \
echo 1 > /proc/sys/net/ipv4/ip_forward; \
chmod +x /etc/rc.d/rc.local; \
rm -rf /var/lib/rook;

COMMAND2 = echo "export KUBECONFIG=/etc/kubernetes/admin.conf" >> ~/.bash_profile && export KUBECONFIG=/etc/kubernetes/admin.conf

reset:
	$(KUBEADM) reset -f
	$(RM) -rf $(HOME)/.kube
	$(RM) -rf /var/lib/rook
	$(foreach node, $(NODE), $(SSH) $(node) '$(KUBEADM) reset -f';)
	$(foreach node, $(NODE), $(SSH) $(node) '$(COMMAND)';)
	$(KUBEADM) init --config $(HOME)/kubeadm.yaml --upload-certs | $(TEE) $(shell tty) | $(TAIL) -n2 > /tmp/kubeinit.sh
	$(foreach node, $(NODE), $(SSH) $(node) < /tmp/kubeinit.sh;)
	$(SCP) /etc/kubernetes/admin.conf 192.168.79.12:/etc/kubernetes/admin.conf
	$(SCP) /etc/kubernetes/admin.conf 192.168.79.13:/etc/kubernetes/admin.conf
	$(foreach node, $(NODE), $(SSH) $(node) '$(COMMAND2)';)
	$(RM) -f /tmp/kubeinit
	$(KUBECTL) taint nodes master node-role.kubernetes.io/master-
	$(KUBECTL) create -f /home/master/kube-flannel.yml

.PHONY: builder main reset operator backend
