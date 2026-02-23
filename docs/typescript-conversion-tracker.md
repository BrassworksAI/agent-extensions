# TypeScript Scripts Requiring Conversion to .mjs

This document catalogs all `.ts` scripts in the repository that should be converted to `.mjs` (ES Modules) for Node.js compatibility.

## Rationale

The project has migrated from Bun to Node.js as the primary runtime. All scripts should use:

- `.mjs` extension (ES Modules)
- Node.js standard library only (`node:*` imports)
- Node.js shebang: `#!/usr/bin/env node`

## Scripts Awaiting Conversion

### Priority 1: Used by Production Commands

These scripts are actively used by commands in the repository and should be converted next.

1. **repository/skills/design-case-study-generator/scripts/copy-version.ts**
   - Purpose: Copies a version directory for design system iterations
   - Used by: `design-case-study-generator` skill
   - Status: Needs conversion to .mjs
   - Complexity: Moderate (file system operations, path handling)

2. **repository/hooks/opencode/spec-validate.ts**
   - Purpose: Validates OpenCode specification files
   - Used by: OpenCode hooks
   - Status: Needs conversion to .mjs
   - Complexity: Low (file validation)

3. **repository/hooks/opencode/markdownlint.mjs**
   - Purpose: Runs markdownlint on files
   - Used by: OpenCode hooks
   - Status: Converted to .mjs
   - Complexity: Low (wrapper around markdownlint)

### Priority 2: Already Have .mjs Equivalents

These scripts have been converted to .mjs. The original .ts files should be removed after confirming .mjs versions work correctly.

1. **repository/skills/merge-change-specs/scripts/merge-change-specs.ts**
   - Purpose: Merges change-set specs into canonical specs
   - Converted to: `merge-change-specs.mjs`
   - Status: .mjs version created, .ts can be removed after validation

2. **repository/skills/spec-format/scripts/validate-change-spec.ts**
   - Purpose: Validates change-set spec markdown format
   - Converted to: `validate-change-spec.mjs`
   - Status: .mjs version created, .ts can be removed after validation

## Conversion Checklist

When converting a .ts script to .mjs:

1. **Rename file** from `.ts` to `.mjs`
2. **Update shebang** to `#!/usr/bin/env node`
3. **Remove TypeScript types** (all type annotations, interfaces, type aliases)
4. **Update imports** (ensure `.mjs` extensions for local imports)
5. **Update all references** in:
   - Skill files (`SKILL.md`)
   - Command files (`.md` in `repository/commands/`)
   - Any documentation files that reference script paths
6. **Test with node**: `node -c <script>.mjs` to verify syntax
7. **Test execution**: Run the script with real data to ensure it works

## Removal Process

After confirming .mjs versions work correctly:

1. Run all relevant commands/skills to verify .mjs version works
2. Check for any remaining references to the .ts file using `grep -r "<filename>.ts"`
3. Remove the .ts file: `rm repository/<path>/<filename>.ts`
4. Update any documentation that still references the .ts version
5. Verify no build/deployment scripts expect the .ts file

## Notes

- Legacy uninstall scripts (`uninstall.legacy.sh` and `uninstall.legacy.ps1`) reference old .ts files installed in user environments. These should not be modified as they're removing old installations.
- The `bun-shell-commands` skill is specifically about Bun Shell for OpenCode and remains as-is since it documents OpenCode's shell syntax.
