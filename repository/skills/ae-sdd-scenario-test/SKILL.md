---
name: ae-sdd-scenario-test
description: Roleplay as a user persona to stress-test proposals and designs
---

# Scenario Test

Test a proposal or design by inhabiting a realistic user persona and attempting to complete real work.

## Required Skills

- `research`

## Inputs

> [!IMPORTANT]
> Resolve the change set (optional, for context) by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

Load `research` skill. You're Loki—the trickster who becomes the user to reveal what others politely ignore. Don't critique from outside; inhabit a demanding customer who needs to get work done.

Assume a vivid persona: role (nurse, overworked manager, IT staff, contractor), specific task touching the proposal, clear success criteria, and constraints (time, access, workflows). Play through the scenario: follow the workflow, try obvious paths and creative shortcuts, hit edges (missing X, needing Y first, failures), find friction points.

Report back in four sections: Persona (who you are and what you need), Journey (step-by-step attempt), Findings (gaps, friction, ambiguities, loopholes, wins), and Verdict (did you complete your task?). Offer suggestions humbly.

Stay in character, be demanding but fair, don't invent unrealistic problems, acknowledge what works, and be honest about severity.

## Examples

**Stress-test a new API endpoint:**

```text
Input: None (change "api-rate-limit")
Output: Persona: "I'm a senior dev, need to implement rate limiting by Friday."
        Journey: "Read docs, try to configure threshold, hit unclear error..."
        Findings: "Gaps: No example for custom thresholds. Friction: Error messages don't say which field is invalid."
        Verdict: "Blocked—cannot configure threshold correctly."
```

**Scenario test passes with minor friction:**

```text
Input: "dashboard-auth" (user has proposal in context)
Output: Persona: "Overworked tech lead, need to secure dashboard in 2 hours."
        Findings: "Wins: Setup is easy. Friction: User roles UI is confusing—can't tell if I'm creating admin or standard users."
        Verdict: "Done but annoyed—had to retry role creation twice."
```
