---
name: derive-taxonomy
description: Derive product-domain taxonomy from change intents and existing specs, producing nested domain/context/capability boundaries.
---

# Derive Taxonomy

## Purpose

Translate change intents into a product-domain taxonomy by:

- understanding existing domain and capability structure,
- evaluating boundary quality,
- proposing nested domain/context/capability paths,
- and producing final capability spec filenames.

Taxonomy derivation must explicitly surface all needed scope across the tree:

- new capabilities under existing domains,
- existing capabilities that likely need modification,
- entirely new domains/contexts when current taxonomy does not fit,
- with nesting depth allowed to grow to any number of levels based on product boundaries.

The taxonomy must stay implementation-agnostic and product-boundary first.

## Inputs I need

- Change intent(s) or proposal text.
- Relevant context about existing domains and capabilities.
- Canonical/change-set spec structure when available.

## Workflow

1. Load the `research` skill and inspect existing specs/taxonomy artifacts first.
2. Build an **Existing Taxonomy Snapshot**:
   - list discovered domains/subdomains/capabilities,
   - identify strong boundaries,
   - flag ambiguous or missing boundaries.
3. Decompose the change intent into product concerns:
   - actors,
   - workflows,
   - business rules/invariants,
   - capability outcomes.
4. Propose a **nested taxonomy tree** using:
   - `Domain -> Subdomain/Context -> Capability Leaf`
   - nesting depth is unbounded and may go to any number of levels when each level introduces a distinct product boundary.
5. Convert the proposed tree into concrete taxonomy artifacts:
    - folder paths,
    - capability leaf names,
    - and spec file targets that should exist under `changes/<name>/specs/...`.
    Explicitly label which paths are:
    - existing paths reused as-is,
    - existing capability files that need updates,
    - and brand-new domain/context/capability paths to create.
6. For resource-oriented API domains, use operation capabilities under resource ownership:
   - `Domain -> Resource -> Operation Capability`
   - examples include `create/read/list/update/delete` and non-CRUD operations such as `search/archive/reconcile` when behavior requires it.
7. Present options, tradeoffs, and a recommended taxonomy. Pause for user confirmation before locking structure.
8. Produce the final taxonomy tree and spec-path plan.

## Capability framing principles

1. **Product-boundary first**: define boundaries by user/business outcomes and rule ownership.
2. **Implementation-agnostic**: do not use frameworks, storage, transport, or infrastructure as primary taxonomy boundaries.
3. **Capability leaves are spec-ready**: each leaf represents one testable capability outcome.
4. **Centralize invariants**: shared rules have one owning capability; dependent capabilities reference it.
5. **Separate independent outcomes**: if one node contains multiple independently testable outcomes, create separate leaves.
6. **Keep inseparable outcomes together**: if two outcomes cannot be specified or validated independently, keep them in one leaf.
7. **Depth follows boundary clarity**: use as many levels as needed to represent distinct product contexts; stop only when the next level would not add a new product-boundary distinction.

## Output format

Provide:

1. **Existing Taxonomy Snapshot**
   - discovered domain/capability paths
   - quality notes (clear, ambiguous, missing)
2. **Intent Decomposition**
   - actors/workflows/rules/outcomes that drive structure
3. **Proposed Nested Taxonomy**
   - tree: `Domain -> Subdomain/Context -> Capability Leaf`
4. **Taxonomy Artifact Plan**
    - folder tree to create/reuse
    - capability files (final spec filenames)
    - per file/path status: `reuse`, `modify`, or `create`
5. **Ownership Decisions**
    - single-owner decisions for cross-domain rules
6. **Spec File Plan**
   - proposed files under `changes/<name>/specs/...`
   - dependency notes where useful
7. **Open Decisions for User Confirmation**
   - options, recommendation, and consequences

## Quick heuristics

- If a boundary is justified by technology rather than product behavior, reframe it.
- If a rule must be identical across multiple contexts, centralize ownership and reference it.
- If a capability node has multiple independent verbs/outcomes, split it into separate leaves.
- If a requirement spans domains, assign one owning domain/capability and reference from others.
- If a proposed leaf is not testable as written, refine until it is.
- If existing taxonomy cannot cleanly contain intent, introduce new domain/context levels rather than forcing fit.

## Guardrails

- Do not derive boundaries from implementation architecture.
- Do not collapse distinct product contexts into one generic bucket.
- Do not over-split capabilities into implementation-level details.
- Keep taxonomy collaborative: propose, explain, and confirm before finalizing.

## Example shapes

- Product example:
  - `gmail/inbox-organization/search-and-filters/full-text-discovery.md`
- Resource-operation example:
  - `account-management/users/create`
  - `account-management/users/read`
  - `account-management/users/list`
  - `account-management/users/update`
  - `account-management/users/delete`
  - `account-management/users/search`
