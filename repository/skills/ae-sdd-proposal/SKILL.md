---
name: ae-sdd-proposal
description: Draft and refine a product proposal document through collaborative dialogue
---

# Product Proposal

Help the user articulate a complete SDD proposal for the active change set.

## Required Skills

- `spec-driven-development` (load first; use it to resolve/init change sets and manage phase/task state)
- `research`

## Inputs

> [!IMPORTANT]
> Resolve required inputs from the workspace first; only ask the user when resolution is ambiguous or missing.

- **Proposal Context**: What problem or capability should this change set address?
- **Change Set Name**: Resolve from existing `changes/` directories. If none exist, create a new full-lane change set using `ae sdd init <name> --lane full`.

## Instructions

The proposal is the source of truth for how the system should behave. Use SDD workflow conventions and CLI commands, not manual scaffolding.

1. **Load SDD Context**: Load `spec-driven-development`, run `ae sdd status [name]`, and confirm the active change set and current phase.
2. **Research First**: Use the `research` skill to identify existing patterns, related features, and integration points.
3. **Collaborative Dialogue**: Ask targeted questions to close behavior gaps instead of guessing.
4. **Write Proposal**: Create or refine `changes/<name>/proposal.md` with concrete narrative behavior, failure/recovery paths, and clear boundaries.
5. **Move Workflow Forward**: Do not update phase state in this command. If the user is ready to proceed, suggest `ae-sdd-next <name>`.

## Success Criteria

- Active change set is resolved through SDD workflow.
- Proposal includes a clear problem statement (1-2 sentences).
- System behavior is documented in specific narrative prose.
- Error handling and recovery paths are explicit.
- Proposal is saved at `changes/<name>/proposal.md`.

## Usage Examples

### Do

- "Based on the `auth` module, should we handle session expiration like this...?"
- "What is the expected behavior if a user tries to delete a locked record?"

### Don't

- Skip loading `spec-driven-development` before deciding where to write.
- Create ad-hoc files outside `changes/<name>/` for proposal artifacts.
