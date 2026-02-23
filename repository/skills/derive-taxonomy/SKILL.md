---
name: derive-taxonomy
description: Identify product-domain and capability boundaries for change intents, defining how a taxonomy should be structured.
---

# Derive Taxonomy

## Purpose

Translate change intents into a capability taxonomy by determining product system fit and clear domain/subdomain boundaries.

## Inputs I need

- The change intent(s) or proposal text.
- Relevant context about existing domains and capabilities (review available taxonomy context).

## Workflow

1. Load the `research` skill for taxonomy context and review available domain/capability materials.
2. For each change intent, determine product system fit: which domain/capability boundaries it aligns with and where it does not.
3. Identify all impacted capabilities in the change set, including new capabilities, modifications, and removals across domains.
4. Classify each capability change as one of: entry points (UIs/APIs/CLIs), cross-cutting mechanisms (caching/auth/ranking), core domain model (data/validation), or global invariants (IDs/resolution rules).
5. Choose the shallowest cohesive group structure for each domain and capability; use deeper paths only when boundaries demand it.
6. Produce the mapping output with boundary decisions, grouping, and dependencies.

## Capability framing principles

1. **Capability-first, implementation-agnostic**: boundaries describe observable behavior and constraints, not technology choices.
2. **Centralize invariants**: global rules belong in one capability boundary that other capabilities depend on.
3. **Split surfaces vs mechanisms**: surfaces are reader/author-facing endpoints; mechanisms are cross-cutting behaviors that affect multiple surfaces.
4. **Prefer small capabilities that compose**: split when a capability contains multiple independent concerns.
5. **Make determinism explicit**: ordering, tie-breakers, and generation rules must be stated in the owning capability.
6. **Continuously dedupe**: remove shadow requirements by picking one capability to own a rule.
7. **Decide once, reference everywhere**: decisions live in a single owning capability boundary.

## Output format

Provide:

- Proposal-to-taxonomy mapping (existing capabilities, new capabilities, refactors, removals)
- Boundary decisions (in/out of scope per domain)
- Group structure (flat or grouped, with path depth)
- Dependencies (capability-to-capability or cross-domain)

## Quick heuristics

- If a rule must be identical across multiple surfaces, centralize it and reference it.
- If a capability has more than one entry point, split by entry point.
- If a requirement spans domains, pick a single owning domain and have others reference it.
- If a requirement is not testable as written, rewrite it until it is.

## Guardrails

- Base decisions on domain/capability fit, not convenience.
- Expect change sets to span multiple capabilities and domains.
- Do not collapse or blur product system boundaries.
- Taxonomy depth is unbounded, but keep paths as shallow as possible for cohesive boundaries.
