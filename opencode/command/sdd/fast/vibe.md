---
name: sdd/fast/vibe
description: Quick prototyping and exploration - skip specs, get to building
agent: sdd/plan
---

<skill>sdd-state-management</skill>

# Vibe

Freedom to explore. Quick prototypes, architectural experiments, "what if" explorations.

Skip the spec ceremony - get to building fast.

## Arguments

- `$ARGUMENTS` - What you want to build/try (freeform context)

## Instructions

### Gather Context

If no context provided, ask user what they want to explore:
- What are you trying to build or prototype?
- What are you curious about?

Keep it loose. This isn't a formal proposal.

### Initialize Change Set

Derive a kebab-case name from the context. Create `changes/<name>/`:

**state.md:**
```markdown
# SDD State: <name>

## Lane

vibe

## Phase

plan

## Phase Status

in_progress

## Pending

- None

## Notes
```

**context.md:**
```markdown
# Vibe: <name>

## What We're Exploring

<Capture user's intent in their words>

## Initial Thoughts

<Any immediate observations or directions>
```

### Next Steps

Tell the user:
- Change set created
- Run `/sdd/plan <name>` to research, plan, and start building

### The Vibe Flow

```
/sdd/fast/vibe <context>  →  /sdd/plan  →  /sdd/implement
                                              ↓
                              [if keeping it]
                                              ↓
                          /sdd/reconcile  →  /sdd/finish
```

Reconcile and finish are optional. If you're throwing it away, stop after implement. If you're keeping it, reconcile captures specs from your implementation.
