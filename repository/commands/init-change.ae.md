---
description: Initialize a new SDD change set
---

# Initialize Change Set

Create a new SDD change set to track progress.

## Required Skills

- `spec-driven-development` (state management, lane setup)

## Workflow

1. **Validate**: Ensure name is kebab-case (lowercase, no spaces).

2. **Conflict Check**: Verify `changes/<name>/` doesn't exist.

3. **Scaffold**:

   ```bash
   ae sdd init <name> --lane full
   ```

   This creates `changes/<name>/state.toml` with:
   - lane: full
   - phase: proposal
   - status: in_progress
   - created_at: current timestamp

4. Create `changes/<name>/proposal.md` with initial structure (problem statement placeholder).

5. **Confirm**: List created files and suggest next step (draft proposal with `/sdd/proposal`).

## Usage Examples

- **Do**: `/init-change my-new-feature`
- **Don't**: `/init-change "My New Feature"` (no spaces or uppercase)

## Success Criteria

- [ ] Directory `changes/<name>/` exists.
- [ ] `state.toml` initialized correctly via `ae sdd init`.
- [ ] `proposal.md` created with initial structure.
- [ ] User informed of next steps.

## Deliverable

List the created files to confirm success.
