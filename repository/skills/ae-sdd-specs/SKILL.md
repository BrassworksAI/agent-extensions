---
name: ae-sdd-specs
description: Write change-set specifications for change
---

# Specs

Write change-set specifications for the change set.

## Required Skills

- `spec-driven-development` (state management, phase gates)
- `spec-format`
- `research`

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

1. Load `spec-driven-development`, `spec-format`, and `research` skills. Read state from `changes/<name>/state.toml` and `proposal.md`. Apply state entry check per skill guidelines.

2. **Lane check**: If lane is not `full`, redirect to appropriate command per skill guidelines.

3. Spec writing is collaborative. Don't start writing immediately. First, summarize the goal, actors, workflows, and constraints in your own words. Explain how the change maps to the existing capability hierarchy. Present concrete statements and decisions. Offer options for open decisions. Pause and wait for user confirmation.

4. After confirmation, run `ae sdd spec list` (no flag by default) to discover canonical spec paths, then use `research` skill to understand related capabilities, naming conventions, and build context for spec writing. If the canonical specs directory does not exist yet, treat this as greenfield and proceed.

5. With approval, create specs in `changes/<name>/specs/` using `spec-format` skill:
   - Specs may be nested by domain
   - Use EARS syntax for requirements
   - Review specs for atomic, testable, implementation-agnostic requirements

6. Suggest `ae-sdd-critique specs` when complete for review.

7. Do not update phase status in this command. When the user wants to proceed, suggest `ae-sdd-next <name>`.

## Examples

**Writing specs with collaborative discovery:**

```text
Input: None (change: "password-reset")
Output: "Your goal is to allow users to reset passwords via email. 
        Actors: user, admin. Workflow: request email → click link → set new password. Correct?"
       User confirms.
       Output: "This maps to auth/password spec boundary. 
                Writing delta spec at changes/password-reset/specs/auth/password.md..."
```

**Lane check redirects:**

```text
Input: "vibe-lane" (not full lane)
Output: "Vibe lane should skip specs—redirecting to ae-sdd-plan."
```
