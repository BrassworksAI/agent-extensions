---
name: ae-sdd-commit
description: Create high-quality git commits with clear Conventional Commit messages
---

# Commit

Analyze the current git state and create appropriate commits. By default, attempt to commit everything in a single commit. Only break into multiple commits if there is clear semantic drift between file changes.

## Required Information

Run these commands to gather necessary context:

1. **Current git status** - Run: `git status --short`
2. **Recent commit history for context** - Run: `git log --oneline -10`
3. **Staged changes** - Run: `git diff --cached`
4. **Unstaged changes** - Run: `git diff`

**IMPORTANT:** Do not perform any file searches, glob searches, or grep searches. Use only the summary information from the commands above. If the summaries are truly ambiguous for grouping or message intent, ask for specific file diffs; otherwise proceed with the best possible commit(s).

## Workflow

Based on the information gathered above:

1. **Inspect**: Check state with `git status` and `git diff`.

2. **Boundary**: Group changes into logical units.
   - If all changes are semantically related (e.g., all part of the same feature, bug fix, or refactoring), create a SINGLE commit.
   - ONLY if there is clear semantic drift (e.g., mixing unrelated features, bug fixes with refactors, documentation with code changes that aren't related), break into multiple commits and group related files together.
   - Default behavior: ONE commit unless you can clearly articulate why the changes are semantically unrelated.

3. **Stage**:
   - If no staged changes exist but unstaged changes do, stage them first with `git add -A` or `git add -p` for mixed files.
   - Review with `git diff --cached` to ensure no secrets or debug logs.

4. **Verify**: Run the fastest relevant check (tests/lint) before committing.

5. **Commit**:
   - Use conventional commit format: `type(scope): description`
   - Follow these message rules:
     - Use conventional commits (feat, fix, refactor, docs, test, chore, etc.)
     - Keep the summary line short (under 50 characters)
     - Focus on what the change ENABLES, not just what was changed
     - DO NOT mention specific filenames, file counts, or line changes
     - DO NOT use vague phrases like "updated files" or "made changes"
     - Include body (rationale) and footer (breaking changes) when needed
   - Examples of good messages:
     - `feat(auth): add OAuth2 login flow`
     - `fix(api): handle null response from cache`
     - `refactor(db): simplify connection pooling`
     - `docs(readme): clarify setup instructions`
     - `feat(auth): add login rate limiting to prevent brute force`

6. **Agent behavior (non-interactive by default)**:
   - **Default**: execute this workflow end-to-end without asking the user what to do.
   - Do not ask for confirmation - perform the commit(s) immediately based on your analysis.
   - Only ask a question if grouping or message intent is truly unclear from the summaries. If you must ask, request specific file diffs by path.

### Mandatory guardrails

- Always run `git diff --cached` immediately before committing.
- If you detect secrets/credentials/keys, **stop** and warn the user.
- If checks fail, fix them and rerun checks; do not commit with failing validators unless the user explicitly approves skipping.

## Usage Examples

- **Do**: `feat(auth): add login rate limiting to prevent brute force`
- **Don't**: `update login.js` or `fix some bugs`
- **Do**: Split a risky database migration from a cosmetic UI fix.
- **Don't**: Mix 15 unrelated file changes into a single "cleanup" commit.

## Success Criteria

- [ ] Working tree is clean.
- [ ] Commit history is linear and logical.
- [ ] Messages follow Conventional Commits.
- [ ] Verification steps passed.

## Deliverable

Run `git log --oneline -5` to verify and prove work.
