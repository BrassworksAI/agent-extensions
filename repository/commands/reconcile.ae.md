---
description: Ensure bi-directional match between specs and implementation
---

# Reconcile

Perform bi-directional reconciliation between change set specs and implementation. The final guard to ensure spec-driven development was truly followed. Verify both: (1) all specs are fully implemented, and (2) all implementation changes have corresponding spec coverage for business logic requirements.

## Required Skills

- `spec-driven-development` (state management, completion workflow)
- `spec-format`

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

1. Load `spec-driven-development` and `spec-format` skills. Read state and apply state entry check per skill guidelines.

2. Read tasks from `changes/<name>/tasks.toml` if full lane.

3. Get the implementation diff for the change set. Execute bi-directional reconciliation:

### Forward Reconciliation (Specs → Implementation)

For each spec in specs/:

- Verify the spec describes behavior that is actually present in the implementation
- Confirm all conditions, flows, and capabilities specified are implemented
- Identify specs without corresponding implementation changes
- Check for partial implementations (spec promises more than code delivers)

### Backward Reconciliation (Implementation → Specs)

For the implementation diff:

- Scan all code changes for new business logic not captured in specs
- Identify control flow changes (if statements, switches, loops) that introduce conditional logic
- Look for new functions, API endpoints, or data structures that represent capabilities
- Check for added error handling, validation, or edge cases that specify requirements
- Validate that behavioral changes (what the code does, not how) have spec coverage
- Note: refactors, test changes, and infrastructure updates may not need specs

Examples of unspecced implementation changes worth capturing:

- Added if statement introduces conditional logic (e.g., user role checks)
- New function exposes capability not previously available
- Added validation logic implies new requirements
- Changed error handling indicates new failure modes to specify
- New API parameters or return values need specification

Examples that typically don't require specs:

- Code refactoring without behavior change
- Test additions or modifications
- Build configuration updates
- Comment-only changes
- Whitespace/formatting adjustments

1. Present findings to the user:

**Forward reconciliation results:**

- Specs fully implemented: [list]
- Specs partially implemented: [list with gaps]
- Specs without implementation: [list]

**Backward reconciliation results:**

- Implementation captured in specs: [confirmation]
- Unspecced business logic found: [list with details]
- No new business logic: [statement]

1. Get user input:

- Add missing specs for implementation changes
- Remove specs not applicable to this change
- Modify specs to match actual implementation
- Mark items as out-of-scope (if applicable)

1. If specs/ doesn't exist, perform backward reconciliation and analyze whether implementation adds/removes business logic worth specifying. Present analysis and ask user to capture specs. If yes, create specs/ and write change-set specs (kind: new and/or delta). If no, document that specs were not created.

2. Document bi-directional reconciliation findings in `changes/<name>/reconciliation.md` including:

- Forward reconciliation audit results
- Backward reconciliation audit results
- User decisions on each discrepancy
- Final determination: reconciled, requires work, or trivial

Do not update phase status in this command. When the user wants to proceed, suggest `next.ae <name>`.

## Examples

**Bi-directional reconciliation, specs match implementation:**

```text
Input: "password-reset" (specs exist)
Forward reconciling...
- Spec "password-reset-flow": MATCHES implementation
- Spec "email-notification": MATCHES implementation
Backward reconciling...
- All implementation changes captured in specs
No discrepancies found. Reconciliation complete.
```

**Backward reconciliation finds unspecced logic:**

```text
Input: "user-preferences" (specs exist)
Forward reconciling...
- All specs implemented: ✓
Backward reconciling...
- Found unspecced conditional logic in settings.service:
  Line 42: if (user.isPremium) { ... }
  This introduces premium-only behavior not captured in specs
Action: Capture as new spec or confirm out-of-scope?
```

**No specs, but implementation adds business logic:**

```text
Input: "checkout-flow" (no specs)
Backward reconciling...
- Added function validateShippingAddress(): new validation logic
- Added if statement for international shipping: conditional logic
- Added new error handling for invalid addresses
Found 3 business logic changes worth specifying. Create specs?
```

**No specs, changes are purely refactoring:**

```text
Input: "code-cleanup" (user has implementation in context)
Backward reconciling...
- Changes are code refactoring (renamed variables, extracted functions)
- No behavioral changes detected
- No new business logic added
Not spec-worthy. Documenting as refactoring-only in reconciliation.md.
```
