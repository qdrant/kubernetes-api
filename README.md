# Qdrant Kubernetes API

This repository holds the API definitions and documentation for the Qdrant Kubernetes operator.

## Development

To generate the API code and docs, run:

```bash
make gen
```

## Releasing a new version

A release is built from a `vX.Y.Z` tag on a commit in `main`. The version must be in the format `vX.Y.Z`, prefixed with a `v`, because this repository is a Go module.

Only members of `cloud-admins` can create `v*` tags (SEC-41 tag ruleset). If you are not a cloud-admin, send one the commit SHA to release:

```bash
git tag -a v1.2.3 <commit-sha> -m v1.2.3
git push origin refs/tags/v1.2.3
```

The tag push starts the release workflow in the `release` environment, which needs an approval from a reviewer. The workflow fails if the tagged commit is not in the history of `main`. After the approval, it publishes the release notes and the Helm chart.

Deployment to different environments is coupled with operator release, please see [qdrant/operator: Release process](https://github.com/qdrant/operator?tab=readme-ov-file#release-process) for details.
