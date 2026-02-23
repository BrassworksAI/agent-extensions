# Script Runtime (Node-only)

These skills use Node.js for running scripts with the Node standard library.

## Requirements

- Node.js must be available: `node --version`
- Use `.mjs` files
- Use Node standard library only (`node:*` imports)
- Avoid third-party dependencies by default

## Invocation

```bash
node scripts/<script>.mjs -- <args>
```

## Syntax validation

Validate scripts parse before relying on them:

```bash
node -c scripts/<script>.mjs
```
