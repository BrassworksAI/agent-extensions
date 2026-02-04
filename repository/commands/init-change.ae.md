---
description: Initialize a new SDD change set
---

# Initialize Change Set

Create a new SDD change set to track progress.

## Required Skills

You MUST load `spec-driven-development` before proceeding. Stop and ask if unavailable.

## Workflow

1. **Validate**: Ensure name is kebab-case (lowercase, no spaces).
2. **Conflict Check**: Verify `changes/<name>/` doesn't exist.
3. **Scaffold**:
    - `changes/<name>/state.md`: Initialize with lane (pending), phase (init), and status (complete).
4. **Confirm**: Suggest next step (draft proposal).

## Usage Examples

- **Do**: `/init-change my-new-feature`
- **Don't**: `/init-change "My New Feature"` (no spaces or uppercase)

## Success Criteria

- [ ] Directory `changes/<name>/` exists.
- [ ] `state.md` and `proposal.md` are initialized correctly.
- [ ] User is informed of the next steps.

## Deliverable

List the created files to confirm success.
