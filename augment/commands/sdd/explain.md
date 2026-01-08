---
description: Explain SDD concepts and workflow
argument-hint: [topic]
---

# Explain SDD

Explain SDD concepts, workflow, or specific phases to help users understand the system.

## Arguments

- `$ARGUMENTS` - Topic to explain (optional)

## Instructions

### FIRST: Check for Topic Argument

**If `$ARGUMENTS` is provided**: Skip directly to "Specific Topic" section below. Do NOT output full overview.

**If `$ARGUMENTS` is empty/not provided**: Output full overview below.

---

### Full Overview (only when no argument provided)

Provide a comprehensive explanation of SDD:

---

**Spec-Driven Development (SDD)** is a disciplined approach where specifications drive implementation. The core principle: *define what you're building before building it*.

#### Why SDD?

- **Clarity**: Specs force you to think through requirements before coding
- **Traceability**: Every implementation decision traces back to a spec
- **Verification**: Reconciliation ensures code matches intent
- **Collaboration**: Specs are a shared contract between human and AI

#### The Change Set

Everything in SDD revolves around a **change set** - a named collection of files tracking one logical change:

```
changes/
  add-user-auth/           # Change set directory
    state.md               # Lane, phase, phase status, notes, pending
    proposal.md            # What we're building and why
    specs/                 # Change-set specs (full lane only; `kind: new|delta`)
      api.md
      components.md
    tasks.md               # Ordered implementation tasks
    plans/                 # Per-task implementation plans
      01.md
      02.md
```

#### Three Lanes

| Lane | When to Use | Flow |
|------|-------------|------|
| **Full** | New features, architectural changes | Specs first, then implement |
| **Vibe** | Prototypes, experiments, quick enhancements | Implement first, specs later (if keeping) |
| **Bug** | Defect fixes | Fix first, spec impact assessed at reconcile |

The key insight: **Full lane** writes specs before implementation. **Vibe/Bug lanes** invert this - implement first, capture specs from diff at reconcile (if work is worth keeping).

#### Phase Progression

**Full Lane** (specs drive implementation):
```
init -> proposal -> specs -> discovery -> tasks -> plan -> implement -> reconcile -> finish
```

**Vibe Lane** (freedom to explore):
```
/sdd:fast:vibe <context> -> plan -> implement -> [reconcile -> finish]
```

**Bug Lane** (fix with triage):
```
/sdd:fast:bug <context> -> plan -> implement -> [reconcile -> finish]
```

Vibe/Bug lanes have optional completion - if throwing away the work, stop after implement. If keeping, reconcile captures specs from what was built.

#### Key Commands

| Command | Purpose |
|---------|---------|
| `/sdd:init <name>` | Start a new change set (full lane) |
| `/sdd:fast:vibe <context>` | Start vibe lane - exploration mode |
| `/sdd:fast:bug <context>` | Start bug lane - fix with triage |
| `/sdd:brainstorm` | Explore problem space before proposing |
| `/sdd:proposal` | Draft/refine proposal (full lane) |
| `/sdd:specs` | Write change-set specifications (full lane) |
| `/sdd:discovery` | Verify specs fit repo architecture |
| `/sdd:tasks` | Generate implementation tasks (full lane) |
| `/sdd:plan` | Create implementation plan |
| `/sdd:implement` | Execute current plan |
| `/sdd:reconcile` | Verify implementation, capture specs if needed |
| `/sdd:finish` | Close change set |
| `/sdd:status` | Show current state |

#### Example: Vibe Lane Flow

```
User: /sdd:fast:vibe "I want to try a new caching approach for the API layer"
      [Creates context.md with freeform exploration context]

User: /sdd:plan
      [Creates single plan.md combining research + implementation steps]

User: /sdd:implement
      [Executes plan, makes changes, iterates freely]

# If throwing away: DONE - no further steps needed

# If keeping the work:
User: /sdd:reconcile
      [Reviews diff, captures change-set specs from implementation]

User: /sdd:finish
      [Change set closed; `kind: new` specs moved, `kind: delta` specs merged to canonical]
```

#### Example: Bug Lane Flow

```
User: /sdd:fast:bug "Users are seeing 500 errors when submitting forms with special characters"
      [Triages: actual bug or behavioral change?]
      [If behavioral change -> redirects to full lane]
      [If bug -> researches, creates context.md with root cause]

User: /sdd:plan
      [Creates single plan.md with fix approach]

User: /sdd:implement
      [Applies fix, validates]

# If one-off fix: DONE

# If keeping for spec update:
User: /sdd:reconcile
      [Assesses spec impact, updates if warranted]

User: /sdd:finish
```

#### Example: Full Lane Flow

```
User: /sdd:init add-plugin-system
User: /sdd:proposal
      "Add a plugin architecture to allow extending functionality"
      [Forge helps draft proposal.md, sets lane: full]

User: /sdd:specs
      [Forge writes change-set specs defining interfaces, behaviors]

User: /sdd:discovery
      [Uses architecture-fit-check to verify specs fit existing architecture]

User: /sdd:tasks
      [Generates tasks from change-set specs]

User: /sdd:plan -> /sdd:implement (for each task)

User: /sdd:reconcile
      [Verifies implementation satisfies all spec requirements]

User: /sdd:finish
```

---

### Specific Topic (when argument provided)

When user provides a topic in `$ARGUMENTS`, explain **only that topic** in depth. Do not output full overview above.

Use this format for the topic:
1. **What it is** - Clear definition
2. **Why it exists** - The problem it solves
3. **How it works** - Mechanics and structure
4. **Example** - Concrete illustration
5. **Common questions** - Anticipate confusion points

#### Known Topics

| Topic | What to Explain |
|-------|-----------------|
| `phases` | Each phase in detail, what happens, gates between them |
| `lanes` | Full vs vibe vs bug - when to use each, spec relationship |
| `specs` | Change-set spec format (`kind: new|delta`), EARS syntax, Before/After blocks |
| `tasks` | Task format, checkboxes, requirements sections (full lane) |
| `plans` | Plan structure - per-task for full lane, single plan.md for vibe/bug |
| `reconcile` | What reconciliation checks, spec capture for vibe/bug |
| `state` | state.md structure (lane, phase, phase status, notes), flexible progression |
| `commands` | All available commands and when to use them |
| `tools` | `/sdd:tools:*` commands for critique, scenario testing, taxonomy |

If topic doesn't match a known topic exactly, interpret user's intent and explain relevant concept.

### Research if Needed

If user asks about something implementation-specific (e.g., "how does reconciliation actually compare specs to code?"), delegate to `@librarian` to investigate codebase for accurate details rather than guessing.
