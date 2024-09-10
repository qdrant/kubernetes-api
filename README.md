# Qdrant Kubernetes API

This repository holds the API definitions and documentation for the Qdrant Kubernetes operator.

## Development

To generate the API code and docs, run:

```bash
make gen
```

## Releasing a new version

Create and push a new release branch like this:

```bash
git checkout -b releases/v1.2.3
git push origin releases/v1.2.3
```

A GitHub Action will create and tag a new release with the given version number. Important: the version number must be in the format `vX.Y.Z`, and prefixed with a `v`.
