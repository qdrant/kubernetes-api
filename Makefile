## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
ENVTEST ?= $(LOCALBIN)/setup-envtest
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.30.0
ENVTEST_VERSION ?= release-0.19

CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
CRD_REF_DOCS ?= $(LOCALBIN)/crd-ref-docs
GCI ?= $(LOCALBIN)/gci

CONTROLLER_TOOLS_VERSION ?= v0.18.0
CRD_REF_DOCS_VERSION ?= v0.1.0
CHART_DIR ?= charts/qdrant-kubernetes-api
CRDS_DIR ?= crds

lint:
	bash -c 'files=$$(gofmt -l .) && echo $$files && [ -z "$$files" ]'
	helm lint $(CHART_DIR)
	golangci-lint run

.PHONY: gen
gen: manifests generate format vet ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.

.PHONY: manifests
manifests: controller-gen ## Generate CustomResourceDefinition objects.
	rm $(CHART_DIR)/templates/management-crds/*.yaml
	rm $(CHART_DIR)/templates/region-crds/*.yaml
	$(CONTROLLER_GEN) crd paths="./..." output:crd:artifacts:config=$(CRDS_DIR)
	mv $(CRDS_DIR)/qdrant.io_qdrantreleases.yaml $(CHART_DIR)/templates/management-crds/
	cp $(CRDS_DIR)/qdrant*.yaml $(CHART_DIR)/templates/region-crds/
	for file in $(CHART_DIR)/templates/management-crds/*.yaml; do \
		echo "{{ if .Values.includeManagementCRDs }}" | cat - $$file > temp && mv temp $$file; \
		echo "{{ end }}" >> $$file; \
	done
	for file in $(CHART_DIR)/templates/region-crds/*.yaml; do \
		echo "{{ if .Values.includeRegionCRDs }}" | cat - $$file > temp && mv temp $$file; \
		echo "{{ end }}" >> $$file; \
	done
	helm lint $(CHART_DIR)

.PHONY: generate
generate: controller-gen crd-ref-docs ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
	$(CRD_REF_DOCS) --config .crd-ref-docs.yaml --renderer markdown --output-path docs/api.md

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary. If wrong version is installed, it will be overwritten.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen && $(LOCALBIN)/controller-gen --version | grep -q $(CONTROLLER_TOOLS_VERSION) || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: crd-ref-doc
crd-ref-docs: $(CRD_REF_DOCS) ## Download crd-ref-doc locally if necessary. If wrong version is installed, it will be overwritten.
$(CRD_REF_DOCS): $(LOCALBIN)
	test -s $(LOCALBIN)/crd-ref-docs || \
	GOBIN=$(LOCALBIN) go install github.com/elastic/crd-ref-docs@$(CRD_REF_DOCS_VERSION)

.PHONY: go_fmt
go_fmt:
	gofmt -w -s .

.PHONY: fmt_imports
fmt_imports: $(GCI)
	$(GCI) write ./ --skip-generated -s standard -s default -s 'prefix(github.com/qdrant)'

.PHONY: fmt
format: go_fmt fmt_imports

fmt: format

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

$(GCI): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/daixiang0/gci@latest


.PHONY: envtest
envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

test: manifests generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)"  go test $(shell go list ./... | grep -v /test/) -coverprofile cover.out


# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef