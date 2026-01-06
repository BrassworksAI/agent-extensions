# SDD State Management

This skill covers how SDD tracks change set state and enforces phase gates.

## State File

Every change set has a state file at `changes/<name>/state.md`.

Phase transitions happen only after explicit user approval—questions and feedback don't count as approval. Task status in `changes/<name>/tasks.md` can move to `[o]` (In Progress) as soon as work (planning or implementation) starts.

```markdown
# SDD State: <name>

## Phase

<current-phase>

## Lane

<full|vibe|bug>

## Pending

- <any blocked items or decisions needed>
```

## Lanes Overview

| Lane | Purpose | Flow |
|------|---------|------|
| **Full** | Enterprise features, architectural changes | Proposal → Specs → Discovery → Tasks → Plan → Implement → Reconcile → Finish |
| **Vibe** | Rapid prototypes, exploration | Context → Plan → Implement → [Reconcile → Finish] |
| **Bug** | Fix defects | Triage + Research → Plan → Implement → [Reconcile → Finish] |

## Full Lane

The complete SDD experience. Specs are written before implementation.

### Phases

```
ideation -> proposal -> specs -> discovery -> tasks -> plan -> implement -> reconcile -> finish
```

| Phase | Purpose | Artifacts |
|-------|---------|-----------|
| `ideation` | Explore problem space | `seed.md` |
| `proposal` | Define what we're building | `proposal.md` |
| `specs` | Write detailed specifications | `specs/**/*.md` (change-set specs: `kind: new|delta`) |
| `discovery` | Understand high-level architecture needs | `thoughts/*.md` |
| `tasks` | Break specs into tasks | `tasks.md` |
| `plan` | Plan current task | `plans/NN.md` |
| `implement` | Execute the plan | Code changes |
| `reconcile` | Verify implementation matches specs | Reconciliation report |
| `finish` | Close and sync specs | `kind: new` moved; `kind: delta` merged into canonical |

### Phase Gates

| From | To | Gate |
|------|----|------|
| ideation | proposal | Seed reviewed |
| proposal | specs | Proposal approved |
| specs | discovery | Change-set specs written (`kind: new|delta`) |
| discovery | tasks | Architecture review complete |
| tasks | plan | Tasks defined |
| plan | implement | Plan approved |
| implement | reconcile | All tasks complete |
| reconcile | finish | Implementation matches specs |

## Vibe Lane

Freedom to explore. Skip specs, get to building fast.

### Flow

```
/sdd/fast/vibe <context>  →  /sdd/plan  →  /sdd/implement  →  [/sdd/reconcile  →  /sdd/finish]
```

### Artifacts

| File | Purpose |
|------|---------|
| `state.md` | Phase and lane tracking |
| `context.md` | Loose capture of what we're exploring |
| `plan.md` | Combined research + planning (single file) |

### Optional Completion

Reconcile and finish are optional:
- **Throwing it away**: Stop after implement
- **Keeping it**: Reconcile captures specs from implementation, finish syncs to canonical

## Bug Lane

Hunt and fix. Triages the issue first.

### Flow

```
/sdd/fast/bug <context>  →  /sdd/plan  →  /sdd/implement  →  [/sdd/reconcile  →  /sdd/finish]
```

### Triage

The bug command determines if this is:
- **Actual bug**: Implementation doesn't match intent → proceed with bug lane
- **Behavioral change**: User wants different behavior → redirect to full lane with specs

### Optional Completion

Reconcile and finish are optional:
- **No spec impact**: Stop after implement (most bug fixes)
- **Specs affected**: Reconcile captures changes, finish syncs

## Thoughts Directory

For complex investigations, use `thoughts/` as a free-form workspace:

```
changes/<name>/
  thoughts/
    investigation.md
    options.md
```

## Pending Semantics

`## Pending` is a ledger of **unresolved** items:

- Contains only unresolved items (blockers, decisions needed)
- Remove items when resolved (don't strike through)
- Leave it blank when nothing is pending
- Deferred ideas go elsewhere (e.g., `docs/future-capabilities.md`)

`## Pending` is NOT an approval log. Do not record approvals there.

## Conversational Phase Advancement

Phase transitions happen through natural conversation, not automatic advancement.

### How It Works

1. **Present work**: Show what was produced (proposals, specs, plans, etc.)
2. **Discuss**: Answer questions, incorporate feedback, refine collaboratively
3. **Wait for explicit approval**: User must clearly signal they're satisfied
4. **Log the decision**: Document what was approved in state.md
5. **Update phase**: Only then advance to the next phase

### Explicit Approval

Look for clear signals that the user approves the work:
- "Looks good", "that works", "approved", "let's move on"
- User runs the next command without raising issues
- Clear statement that the work is complete

**These do NOT count as approval:**
- User asking clarifying questions
- User providing additional context
- Acknowledging without agreement ("interesting", "I see")
- Answering questions you asked

### Logging Decisions

When user explicitly approves a gate:

- Update `changes/<name>/state.md` `## Phase` to the next phase immediately.
- Remove any now-resolved items from `## Pending`.
- Leave `## Pending` blank if nothing is pending.

Do NOT write approval summaries into `## Pending`.

If an approval record is needed, capture it in a separate artifact (e.g., `changes/<name>/thoughts/decisions.md`).

### Natural Flow

Keep it conversational. Don't force "ready to advance?" checkpoints. Work through the step with the user like collaborators solving a problem together. You'll naturally know when they're satisfied and ready to move on.
