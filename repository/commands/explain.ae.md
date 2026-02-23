---
description: Explain SDD concepts, workflow, and CLI-first usage
---

# SDD Explain

Explain Spec-Driven Development (SDD) with a CLI-first mental model. Use this command to teach lane selection, phase flow, and how slash commands and `ae sdd` commands work together.

## Required Skills

- `spec-driven-development`
- `research`

## Instructions

1. **Identify Intent**: Determine whether the user wants a quick overview, lane guidance, or a deep dive on one phase/command.

2. **Why SDD**: Explain the three outcomes:
   - **Clarity**: behavior is defined before implementation drift
   - **Traceability**: tasks and code map back to specs
   - **Confidence**: reconcile/validation checks prove intent matches implementation

3. **CLI-First Operating Model**:
    - Change sets are initialized with `ae sdd init <name> --lane <full|vibe|bug>`.
    - Progress is tracked with `ae sdd status [name]`, `ae sdd phase complete`, and `ae sdd phase next`.
    - Full-lane implementation uses `ae sdd task list|start|complete|complete --next`.
    - Slash commands (`/sdd/init`, `/sdd/proposal`, `/sdd/specs`, `/sdd/tasks`, etc.) are workflow assistants; the CLI is the source of truth for state transitions.

4. **Lane Flows**:
   - **Full**: `proposal -> specs -> discovery -> tasks -> plan -> implement -> reconcile -> finish`
   - **Vibe**: `context -> plan -> implement -> [reconcile -> finish]`
   - **Bug**: `triage -> plan -> implement -> [reconcile -> finish]`

5. **Change Set Structure**:

   ```text
   changes/<name>/
     state.toml      # lane, phase, notes, pending
     proposal.md     # full lane proposal
     context.md      # vibe/bug context when needed
     specs/          # behavior contracts
     thoughts/       # discovery notes
     tasks.toml      # ordered tasks with spec_requirements
     plans/          # implementation plans
   ```

6. **Command Map**:

   | Command | Use For |
   |---|---|
    | `ae sdd init` | Create a new change set |
    | `ae sdd status` | See current lane/phase/tasks |
    | `ae sdd phase complete [--next]` | Mark phase complete and optionally advance |
    | `ae sdd phase next` | Advance only after phase is complete |
    | `ae sdd task start` | Start next full-lane task |
    | `ae sdd task complete --next` | Finish and chain tasks |
    | `/sdd/init` | Derive and approve a new change-set name, then initialize |
    | `/sdd/continue` | Resume an existing change set from CLI status |
    | `/sdd/vibe` | Start or continue vibe-lane exploratory work |
    | `/sdd/bug` | Triage and initialize bug-lane fixes |
    | `/sdd/proposal` | Draft and refine proposal |
    | `/sdd/specs` | Create/update specs |
    | `/sdd/discovery` | Validate architecture and risks |
    | `/sdd/tasks` | Build `tasks.toml` from specs |
    | `/sdd/plan` | Create execution plans |
    | `/sdd/implement` | Execute planned work |
    | `/sdd/critique` | Review artifacts for quality and gaps |
    | `/sdd/scenario-test` | Validate behavior using realistic scenarios |
    | `/sdd/commit` | Craft commit(s) aligned to SDD progress |
    | `/sdd/reconcile` | Verify implementation vs specs |
    | `/sdd/finish` | Close out change set |
    | `/sdd/explain` | Teach SDD concepts and workflow usage |

7. **How To Guide Users**:
   - Recommend lane choice based on risk and ambiguity.
   - Explain next command, why it is next, and what artifact it should produce.
   - Use `research` for repo-specific details when asked how current code behaves.

## Success Criteria

- User understands when to use full vs vibe vs bug lane.
- User understands that `ae sdd` CLI drives state and phase progression.
- User can name the next command and expected artifact in their current phase.

## Usage Examples

### Explain CLI + command relationship

"Use `/sdd/specs` to draft the contract, then run `ae sdd phase complete --next` when the phase is done."

## Followup Question

> [!IMPORTANT]
> End by asking one focused follow-up: "Do you want a lane recommendation for your current task, or a step-by-step walkthrough of your current phase?"
