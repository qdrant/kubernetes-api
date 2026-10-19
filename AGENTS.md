# kubernetes-api

Go types and generated CRDs for Qdrant's Kubernetes resources. This repository defines the
contract only; it contains no controllers.

## Layout

Three API groups, each with its own Helm chart toggle:

| Package | API group | Chart directory (value) |
|---|---|---|
| `api/v1` | `qdrant.io/v1` | `region-crds` (`includeRegionCRDs`); `QdrantRelease` goes to `management-crds` (`includeManagementCRDs`) |
| `api/auth/v1alpha1` | `auth.qdrant.io/v1alpha1` | `auth-crds` (`includeAuthCRDs`) |
| `api/routing/v1alpha1` | `routing.qdrant.io/v1alpha1` | `routing-crds` (`includeRoutingCRDs`) |

The `manifests` target in the `Makefile` decides which chart directory each CRD lands in;
see `docs/api.md` for all resources.

- `api/**/*_types.go`: the types and kubebuilder markers. In `api/`, edit only these, the
  tests and `groupversion_info.go` (for a new group) by hand.
- Generated, never edit by hand: `api/**/zz_generated.deepcopy.go`, `crds/`,
  `charts/qdrant-kubernetes-api/templates/*-crds/`, `docs/api.md`.
- A new resource in `qdrant.io` lands in `region-crds` unless the `manifests` target moves it
  elsewhere.

## Changing types

- After any change under `api/`, run `make gen` and commit all generated files.
- New API group: create `charts/qdrant-kubernetes-api/templates/<group>-crds/`, add the `rm`,
  `mv` and template loop lines to the `manifests` target, add `include<Group>CRDs` to
  `values.yaml`, add the helm lint lines to both the `lint` target and
  `.github/workflows/pr-workflow.yaml`, and update the Layout table above. Otherwise its CRD
  stays in `crds/` and never reaches the chart.
- New optional fields: pointer type, `+optional`, `omitempty`, and no `+kubebuilder:default`
  unless explicitly required. A default changes every existing CR.
- Fields under `spec.config` that map directly to a Qdrant config key use Qdrant's snake_case
  name (`max_log_files`). Kubernetes references such as secrets use camelCase. `maxCollections`
  is an existing exception; don't copy it. Fields outside `spec.config` use camelCase.
- Restrict enums with `+kubebuilder:validation:Enum=`, using Qdrant's exact values.
- If a field needs a minimum Qdrant version, say so in the field comment ("Requires Qdrant
  v1.x.y or later").
- Every field except `Spec`, `Status` and `Items` gets a doc comment; it becomes the CRD
  description, and non-status comments also end up in `docs/api.md`.
- Additions only. Removing or renaming a field breaks existing CRs. To retire a field, keep it
  and add a `// Deprecated: Use X instead` comment.
- Don't reorder imports by hand; `make gen` runs `gci` with a `github.com/qdrant` group.

## Validate

- `make gen` leaves no uncommitted changes.
- `make lint` (gofmt, golangci-lint, helm lint for the chart flag combinations). CI also runs
  `kubeconform` on `helm template` output, which `make lint` does not.
- `make test` (unit tests in `api/v1` and `api/routing/v1alpha1`, envtest integration tests
  in `api/v1`). It runs `manifests`, `generate`, `fmt` and `vet` first, and downloads envtest
  assets, so it needs network access. Add a case to the table tests in `*_types_test.go` for
  new validation logic.

## Branches and release

- Name branches `feature/...` (or `feat/...`), `fix/...` or `chore/...`. The PR label is set from the branch
  name, and release notes are grouped by that label.
- A release is a push to a `releases/vX.Y.Z` branch; CI creates the tag and GitHub release.
