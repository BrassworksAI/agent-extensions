# Script Runtime (Node-first)

These skills are authored to run scripts with Node by default. The goal is a predictable baseline with minimal runtime branching.

## Detection

Check whats available:

```bash
node --version
bun --version
```

## Decision

- If `node` is available: use Node.
- Else if `bun` is available: use Bun as a fallback runner.
- Else: do not bundle scripts; use manual steps in `SKILL.md`.

## File type and dependencies

- Use `.mjs`.
- Use Node standard library only (`node:*` imports).
- Avoid third-party dependencies by default.

## Invocation

Prefer:

```bash
node scripts/<script>.mjs -- <args>
```

Fallback:

```bash
bun scripts/<script>.mjs -- <args>
```

## Syntax/parse validation

Before relying on a new or modified script, validate that it parses:

- Node: `node -c scripts/<script>.mjs`
- Bun: `bun --check scripts/<script>.mjs`
