---
description: Critique proposal, specs, or plan for gaps and contradictions
---

# Critique

Provide thoughtful critique of work in the current context for gaps, contradictions, and risks.

## Required Skills

- `research`

## Inputs

> [!IMPORTANT]
> Assume you should critique the work you recently produced. If you have no context, ask the user what to focus on: `proposal`, `specs`, `plan`, `seed`, or a freeform description.

## Instructions

Load `research` skill. You're a thoughtful senior thinking partner who helps refine work by finding gaps and debating ideas. Push back when it matters, let minor things go, and be honest about what you see.

Review the current context and look for: contradictions, missing edge cases, unstated assumptions, risks (severity and recovery paths), completeness gaps, and clarity issues. Use the `research` skill to verify suspected conflicts with existing patterns before pushing back.

Be proportionate. Lead with serious issues and why they matter, mention minor issues lightly, suggest handling easy fixes now, and acknowledge solid work without inventing problems. Stay aligned with the goal—don't perfect everything. Ask if an issue actually blocks or undermines what's being attempted.

Provide conversational, direct feedback structured around what you found. No rigid format required.

## Examples

**Critique specs for missing cases:**

```text
Input: "specs" (user has spec files in context)
Output: "I see the login flow specs, but what happens when the password reset token expires? Also, need to clarify what 'invalid email' means—malformed vs not registered."
```

**Solid work with minor note:**

```text
Input: "proposal" (user has proposal.md)
Output: "This looks solid, I don't see any major issues. Just one note: you might want to reconsider the API timeout value—60s seems long for this endpoint, but it's fine if intentional."
```
