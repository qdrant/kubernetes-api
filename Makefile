## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
CRD_REF_DOCS ?= $(LOCALBIN)/crd-ref-docs
GCI ?= $(LOCALBIN)/gci

CONTROLLER_TOOLS_VERSION ?= v0.16.1
CRD_REF_DOCS_VERSION ?= v0.0.12

lint:
	bash -c 'files=$$(gofmt -l .) && echo $$files && [ -z "$$files" ]'
	helm lint charts/qdrant-operator-crds
	golangci-lint run

.PHONY: gen
gen: manifests generate format vet ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.

.PHONY: manifests
manifests: controller-gen ## Generate CustomResourceDefinition objects.
	rm charts/qdrant-operator-crds/templates/management-crds/*.yaml
	rm charts/qdrant-operator-crds/templates/region-crds/*.yaml
	$(CONTROLLER_GEN) crd paths="./..." output:crd:artifacts:config=charts/qdrant-operator-crds/templates
	mv charts/qdrant-operator-crds/templates/qdrant.io_qdrantreleases.yaml charts/qdrant-operator-crds/templates/management-crds/
	mv charts/qdrant-operator-crds/templates/qdrant*.yaml charts/qdrant-operator-crds/templates/region-crds/
	for file in charts/qdrant-operator-crds/templates/management-crds/*.yaml; do \
		echo "{{ if .Values.includeManagementCRDs }}" | cat - $$file > temp && mv temp $$file; \
		echo "{{ end }}" >> $$file; \
	done
	for file in charts/qdrant-operator-crds/templates/region-crds/*.yaml; do \
		echo "{{ if .Values.includeRegionCRDs }}" | cat - $$file > temp && mv temp $$file; \
		echo "{{ end }}" >> $$file; \
	done
	helm lint charts/qdrant-operator-crds

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