---
description: Explain SDD concepts and workflow
---

# SDD Explain

Explain Spec-Driven Development (SDD) concepts, workflows, or specific phases to help users navigate the system. This command serves as a technical mentor for SDD.

## Required Skills

- `spec-format`
- `spec-driven-development`
- `research`

## Instructions

1. **Identify Intent**: Ask if they want a general overview or a specific deep-dive.
2. **Why SDD?**: SDD ensures quality through **Clarity** (thinking before coding), **Traceability** (code maps to intent), and **Confidence** (reconciliation).
3. **Workflow Logic**: Explain the lanes and step purposes:
   - **Full**: High-confidence features. `init` (isolate) > `specs` (contract) > `discovery` (validate) > `tasks` (roadmap) > `plan` (strategy) > `implement`.
   - **Vibe/Bug**: Rapid fixes/prototypes. `implement` first to explore, then `reconcile` to capture specs once the solution is proven.
4. **The Change Set**: Explain the directory structure:

   ```text
   changes/<name>/
     state.md      # Lane/phase tracker
     specs/        # Change-set 'Contract'
     thoughts/     # Capture insights for planning
     tasks.md      # Implementation 'Roadmap'
     plans/        # Per-task 'Strategy'
   ```

5. **Command Reference**:

   | Command | Purpose | Why use it? |
   |---------|---------|-------------|
   | `/sdd/init` | Start change set | Isolate work & track progress |
   | `/sdd/fast/vibe` | Prototyping | Explore solutions without early specs |
   | `/sdd/fast/bug` | Defect repair | Targeted fixes with triage |
   | `/sdd/specs` | Define specs | Create the feature's source of truth |
   | `/sdd/discovery` | Verify fit | Ensure specs align with architecture |
   | `/sdd/tasks` | Generate tasks | Implementation checklist |
   | `/sdd/plan` | Create plans | Strategy for a specific task |
   | `/sdd/implement` | Execute plans | Turn strategy into code |
   | `/sdd/reconcile` | Verify code | Audit diff against specs |
   | `/sdd/finish` | Close & merge | Finalize specs and cleanup |

6. **Accuracy**: Use `research` for repo-specific implementation details.

## Success Criteria

- User understands the logic and purpose of each phase.
- Folder structure and command relationships are clearly articulated.
- Explanations prioritize rationale ("Why") over structure ("What").

## Usage Examples

### Do: Explain the Vibe lane

"Vibe lane is for 'implement-first' workflows. It's useful for prototyping. You reconcile later to capture what you built."

### Don't: Be overly prescriptive

Avoid dictating a specific lane. Explain trade-offs so the user can choose.

## Followup Question

> [!IMPORTANT]
> Ask the user: "Is there a specific SDD concept, command, or workflow you would like to dive into next?"
