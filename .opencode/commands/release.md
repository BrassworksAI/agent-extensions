---
description: Trigger the Release PR workflow
---

# Release

Current manifest version:
!`node -p "require('./.release-please-manifest.json')['.']"`

Latest tag:
!`git describe --tags --abbrev=0 2>/dev/null || printf 'none'`

Commits since latest tag:
!`LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || true); if [ -n "$LAST_TAG" ]; then git log "$LAST_TAG"..HEAD --oneline; else git log -20 --oneline; fi`

Run the release flow now:

1. Trigger the workflow with `gh workflow run release-pr.yml --ref main`.
2. Show the newest run entry from `gh run list --workflow release-pr.yml --limit 1`.
3. Return the run URL and current status.

If the workflow fails to trigger, report the error and next action.
