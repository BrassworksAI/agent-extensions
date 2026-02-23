---
description: Continue an SDD change set in a fresh chat using CLI status only
---

# SDD Continue

Resume work from current SDD state in a new chat. Rebuild context from CLI status, then recommend the next command.

## Required Skills

- `spec-driven-development` (load first; source of truth for lane/phase/task flow)

## Inputs

> [!IMPORTANT]
> Resolve required inputs from the workspace first; only ask the user when resolution is ambiguous or missing (for example, multiple change sets exist).

- **Change Set Name**: Resolve active change set; if multiple exist, ask which one to continue.

## Hard Constraints

- MUST NOT read any files directly.
- Do not use Read/cat/manual file inspection for state, tasks, plans, or specs.
- Use `ae sdd` CLI commands only.

## Instructions

1. **Load SDD Context**: Load `spec-driven-development` and run `ae sdd status [name]`.

2. **Interpret Status**:
   - Identify lane, current phase, phase status, and the `Next:` recommendation from status output.
   - Use status output as the primary recommendation source for next command.

3. **Task Continuation Rule**:
   - If status includes a task progress line and percentage is not `100%`, run:
     - `ae sdd task current [name]`
   - Report the current task details so the user can continue immediately.
   - If no current task exists while progress is <100%, recommend:
     - `ae sdd task start [name]`

4. **Recommend Next Steps**:
   - Provide a short continuation brief:
     - where they are now,
     - what to do next,
     - exact next command to run.
   - Keep recommendations aligned to lane and phase.

5. **Output Format**:
   - Current lane + phase + completion state
   - Task status (if applicable)
   - Recommended next action
   - Exact command to run next

## Success Criteria

- Session context is restored from `ae sdd status`.
- No direct file reads are performed.
- If task progress is <100%, current task is surfaced via `ae sdd task current`.
- User gets one clear next command to continue.

## Usage Examples

### Incomplete tasks (<100%)

```text
Ran: ae sdd status checkout-fix
Current: full lane, phase=implement, status=in_progress
Tasks: 50% (2/4)
Ran: ae sdd task current checkout-fix
Current task: payment-retry
Next action: Continue current task, then complete and advance.
Next command: ae sdd task complete --next checkout-fix
```

### Tasks complete (100%)

```text
Ran: ae sdd status checkout-fix
Current: full lane, phase=implement, status=in_progress
Tasks: 100% (4/4)
Next action: Phase work appears complete.
Next command: ae sdd phase complete --next checkout-fix
```
