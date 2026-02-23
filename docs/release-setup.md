# Release Setup Checklist

## npm Publishing (Trusted Publishing with OIDC)

### First-Time Setup

1. **Manual first publish** (creates the package on npm):

   ```sh
   cd npm
   npm publish --access public
   ```

2. **Configure Trusted Publisher on npmjs.com**:
   - Go to: <https://www.npmjs.com/package/@shanepadgett/agent-extensions/access>
   - Navigate to: Settings → Trusted Publisher
   - Select: **GitHub Actions**
   - Enter:
     - Organization/user: `shanepadgett`
     - Repository: `agent-extensions`
     - Workflow filename: `release.yml`
   - (Optional) Enable "Require two-factor authentication and disallow tokens" for extra security

After setup, all future releases via `v*` tags will auto-publish to npm using OIDC (no tokens).

## Creating a Release

```sh
git checkout main
git pull origin main
git tag v0.1.0
git push origin v0.1.0
```

This triggers:

1. GoReleaser builds binaries for all platforms
2. Creates GitHub Release with assets
3. Publishes to npm with matching version

## Deferred Tasks

### Homebrew Tap (Not Yet Configured)

To enable `brew install shanepadgett/tap/ae`:

1. Create GitHub repo: `shanepadgett/homebrew-tap`
2. Add empty `Casks/` directory
3. Add `homebrew_casks` section back to `.goreleaser.yaml`:

   ```yaml
   homebrew_casks:
     - name: ae
       repository:
         owner: shanepadgett
         name: homebrew-tap
       directory: Casks
       homepage: https://github.com/shanepadgett/agent-extensions
       description: CLI for managing AI coding agent extensions
       binaries:
         - ae
   ```

4. Version tags must follow semver: `v1.0.0`, `v1.2.3`, etc.

### Hooks Support

Commands and skills are implemented. Hooks require JSON merging per tool and were deferred as complex.
