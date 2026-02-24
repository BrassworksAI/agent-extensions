---
description: Explain SDD concepts, workflow, and CLI-first usage
---

# SDD Explain

Explain Spec-Driven Development (SDD) with a CLI-first mental model. Use this command to teach lane selection, phase flow, and how `.ae` commands and `ae sdd` commands work together.

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
    - Progress is tracked with `ae sdd status [name]` and `next.ae [name]`.
    - Full-lane implementation uses `ae sdd task list|start|complete|complete --next`.
    - Use `ae sdd spec list` to discover canonical spec paths before reading or editing canonical specs.
    - `.ae` commands (`init.ae`, `proposal.ae`, `specs.ae`, `tasks.ae`, etc.) are workflow assistants; the CLI is the source of truth for state transitions.

4. **Deterministic Full-Lane Loop (Plan <-> Implement)**:
   - Explain this exact loop until all tasks are complete:
     1) `plan.ae` prepares the current task plan (task is selected via CLI task commands)
     2) `next.ae` advances `plan -> implement`
     3) `implement.ae` executes and completes the current task via `ae sdd task complete`
     4) `next.ae` decides transition deterministically using CLI guards:
        - remaining tasks: loop `implement -> plan`
        - all tasks complete: advance `implement -> reconcile`
   - Emphasize that workflow commands do not manually force phase transitions.

5. **Lane Flows**:
   - **Full**: `proposal -> specs -> discovery -> tasks -> plan -> implement -> reconcile -> finish`
   - **Vibe**: `context -> plan -> implement -> [reconcile -> finish]`
   - **Bug**: `triage -> plan -> implement -> [reconcile -> finish]`

6. **Change Set Structure**:

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

7. **Command Map**:

   | Command | Use For |
   |---|---|
   | `ae sdd init` | Create a new change set |
   | `ae sdd status` | See current lane/phase/tasks |
   | `next.ae` | Verify current-phase artifacts, then advance phase |
   | `ae sdd task start` | Start next full-lane task |
   | `ae sdd task current` | Show active in-progress task |
   | `ae sdd task next` | Show next pending task |
   | `ae sdd task complete` | Complete current in-progress task |
   | `ae sdd task complete --next` | Finish and chain tasks |
   | `ae sdd phase complete --next` | Underlying guarded transition used by `next.ae` |
   | `ae sdd phase set` | Manual override only by explicit user request |
   | `ae sdd notes set` | Persist resume context/decisions |
   | `ae sdd pending add/clear` | Track unresolved blockers |
   | `ae sdd spec list` | List canonical specs |
   | `init.ae` | Derive and approve a new change-set name, then initialize |
   | `continue.ae` | Resume an existing change set from CLI status |
   | `vibe.ae` | Start or continue vibe-lane exploratory work |
   | `bug.ae` | Triage and initialize bug-lane fixes |
   | `proposal.ae` | Draft and refine proposal |
   | `specs.ae` | Create/update specs |
   | `discovery.ae` | Validate architecture and risks |
   | `tasks.ae` | Build `tasks.toml` from specs |
   | `plan.ae` | Create execution plans |
   | `implement.ae` | Execute planned work |
   | `critique.ae` | Review artifacts for quality and gaps |
   | `scenario-test.ae` | Validate behavior using realistic scenarios |
   | `commit.ae` | Craft commit(s) aligned to SDD progress |
   | `reconcile.ae` | Verify implementation vs specs |
   | `finish.ae` | Close out change set |
   | `explain.ae` | Teach SDD concepts and workflow usage |

8. **How To Guide Users**:
   - Recommend lane choice based on risk and ambiguity.
   - Explain next command, why it is next, and what artifact it should produce.
   - Use `research` for repo-specific details when asked how current code behaves.

## Success Criteria

- User understands when to use full vs vibe vs bug lane.
- User understands that `ae sdd` CLI drives state and phase progression.
- User can name the next command and expected artifact in their current phase.

## Usage Examples

### Explain CLI + command relationship

"Use `specs.ae` to draft the contract, then run `next.ae` when you want to advance phases."

## Followup Question

> [!IMPORTANT]
> End by asking one focused follow-up: "Do you want a lane recommendation for your current task, or a step-by-step walkthrough of your current phase?"
