# Qdrant operator API

This repository holds the API definitions and documentation for the Qdrant Kubernetes operator.

## Development

To generate the API code and docs, run:

```bash
make gen
```

## Releasing a new version

Create and push a new release branch like this:

```bash
git checkout -b releases/x.y.z
git push origin releases/x.y.z
```

A GitHub Action will create and tag a new release with the given version number.
