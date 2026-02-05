---
description: Collaborative vibe coding sessions that free you from spec ceremony while delivering sane software
---

# SDD Vibe

## Required Skills

- `spec-driven-development` (state management, lane flow)
- `research`

## Inputs

> [!IMPORTANT]
> Resolve required inputs from the workspace first; only ask the user when resolution is ambiguous or missing (for example, multiple change sets exist).

- **Exploration Context**: Ask the user what they want to prototype, explore, or experiment with.
- **Change Set Name**: Resolve by checking existing change sets. If only one exists, proceed with it; if multiple exist, ask which to continue. If none exist, ask the user for a new name to track this experiment (e.g., `vibe-api-redesign`).

## Instructions

Vibe lane is your partner for getting through SDD without the ceremony. Think collaborative vibe coding: throw ideas, prototype freely, create documents as insights emerge, let good ones stick. Goal: sane software, not perfect specs.

1. **Resolution**: Check existing change sets. If one exists, proceed. If multiple, ask whether new or continuing. When continuing, inspect `state.toml` and `context.md` (if it exists). Ask where to resume if jump point unclear.

2. **Initialize** (if new):

   ```bash
   ae sdd init <name> --lane vibe
   ```

   This creates `state.toml` with lane=vibe, phase=context, status=in_progress.

   Create `changes/<name>/context.md` with intent, initial thoughts, and curiosity.

3. **Explore Together**: Dive in. Prototype freely, create documents when insights emerge, capture thoughts. Collaboratively refine into at least a light proposal leading to research and planning. No wrong turn if you document learning.

4. **Keep Context Alive**: Continuously update state and context:
   - Update `context.md` as intent evolves
   - Update state notes in state.toml to capture decisions and next steps
   - Captures progress for seamless continuation—never lose the thread

5. **Deliver Sanity, Not Ceremony**: Build plans, run validations. Ensure exploration converges on sane software. Inform user specifications captured during `reconcile` if kept—until then, focus on vibe, not paperwork.

## Success Criteria

- Change set initialized with `vibe` lane and collaborative context captured
- User flows through exploration, creating documents and artifacts as insights emerge
- Session maintains continuous state updates for seamless continuation
- Exploration converges on actionable results leading to sane implementation

## Usage Examples

### Do: Embrace the vibe

"What's on your mind? I'll fire up vibe lane `vibe-api-redesign` and we can prototype endpoints, document what works, and shape it into something solid—all without getting bogged down in specs."

### Don't: Use vibe to skip necessary rigor

Avoid vibe lane for major features where architectural alignment impacts production readiness. Suggest full lane instead when the scope demands structured specification.
