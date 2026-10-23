# Qdrant Kubernetes API

This repository holds the API definitions and documentation for the Qdrant Kubernetes operator.

## Development

To generate the API code and docs, run:

```bash
make gen
```

## Releasing a new version

Releases are built from a `releases/vX.Y.Z` branch. The version must be in the format `vX.Y.Z`, prefixed with a `v`, because this repository is a Go module.

Only members of `cloud-admins` can create `v*` tags (SEC-41 tag ruleset), so the tag is pushed by a cloud-admin before the release branch:

1. Create the release branch locally from the commit to release:

   ```bash
   git checkout -b releases/v1.2.3
   ```

2. Ask a cloud-admin to tag that commit and push the tag. If you are not a cloud-admin, send them the commit SHA:

   ```bash
   git tag -a v1.2.3 <commit-sha> -m v1.2.3
   git push origin refs/tags/v1.2.3
   ```

3. Push the release branch:

   ```bash
   git push origin releases/v1.2.3
   ```

The release workflow finds the tag on the branch head and skips tag creation, then publishes the release notes and the Helm chart. If the branch was pushed before the tag, the run fails at the tag step; push the tag and re-run the failed job. A tag that exists on a different commit fails the release.

Deployment to different environments is coupled with operator release, please see [qdrant/operator: Release process](https://github.com/qdrant/operator?tab=readme-ov-file#release-process) for details.
