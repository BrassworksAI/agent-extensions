---
name: sdd/fast/bug
description: Bug investigation and fix - triage, research, plan, fix
agent: sdd/plan
---

# Required Skills (Must Load)

You MUST load and follow these skills before doing anything else:

- `sdd-state-management`
- `research`

If any required skill content is missing or not available in context, you MUST stop and ask the user to re-run the command or otherwise provide the missing skill content. Do NOT proceed without them.

# Bug Fix

Investigate and fix a bug. This command triages the issue, researches the codebase, and gets you ready to plan and fix.

## Arguments

- `$ARGUMENTS` - Bug context (what's wrong, error messages, reproduction steps, etc.)

## Instructions

### Gather Context

If no context provided, ask user:
- What's happening?
- What did you expect to happen?
- Any error messages?

### Triage: Bug or Behavioral Change?

Before investigating, determine what type of issue this is:

**Actual Bug** (proceed with bug lane):
- Runtime error, crash, exception
- Code doesn't do what it's supposed to do
- Implementation doesn't match existing specs
- Regression from previous working behavior

**Behavioral Change** (redirect to full lane):
- User wants different behavior than what's specified
- "It works, but I want it to work differently"
- New capability or feature request disguised as bug

**If it's a behavioral change:**
> This looks like a behavioral change rather than a bug. The current behavior may be working as specified.
>
> I recommend using the full SDD lane to properly spec out the new behavior:
> - `/sdd/init <name>` → `/sdd/proposal` → `/sdd/specs`
>
> This ensures the change is captured in specs before implementation.

**If it's an actual bug, continue.**

### Derive Change Set Name

From the bug context, derive a descriptive kebab-case name:
- "Login fails on Firefox" → `fix-login-firefox`
- "Null pointer when adding user" → `fix-add-user-null`
- "Export downloads wrong file" → `fix-export-file`

### Research the Issue

Use the `research` skill to investigate:

1. **Locate the problem**: Where is this happening in the codebase?
2. **Understand the flow**: Trace the execution path
3. **Identify root cause**: Why is this happening?
4. **Check specs**: Are there specs that define expected behavior?

### Initialize Change Set

Create `changes/<name>/`:

**state.md:**
```markdown
# SDD State: <name>

## Lane

bug

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
# Bug Fix: <name>

## Problem

<What's happening>

## Expected Behavior

<What should happen>

## Root Cause

<Findings from research>

## Fix Approach

<High-level approach>

## Spec Impact

<None expected / May affect specs - will verify at reconcile>
```

### Spec Impact Assessment

Based on your research:

- **No spec impact**: Fix is purely implementation - proceed normally
- **Possible spec impact**: Note it in context.md, reconcile will catch it
- **Obvious spec impact**: Flag to user:
  > This fix will change documented behavior. Consider running `/sdd/specs` first to update specs, or proceed and capture changes at reconcile.

Most bugs are implementation issues where specs remain aligned. Trust reconcile to catch edge cases.

### Next Steps

Tell the user:
- Change set created
- Root cause identified
- Run `/sdd/plan <name>` to plan and implement the fix

### The Bug Flow

```
/sdd/fast/bug <context>  →  /sdd/plan  →  /sdd/implement
                                              ↓
                              [if specs affected]
                                              ↓
                          /sdd/reconcile  →  /sdd/finish
```

Reconcile and finish are optional if the fix doesn't affect specs. If behavior changed, reconcile captures the spec updates.
