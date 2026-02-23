---
name: research
description: Research and discover codebase information by breaking tasks into focused steps. Use when you need to understand existing patterns, find relevant files, locate integration points, or trace how something works before making changes.
---

# Research - Codebase Discovery

Guidance for effective codebase research.

## When to Research

Research when you need to:

- **Understand patterns** before proposing changes
- **Find relevant files** for a task
- **Locate integration points** for new functionality
- **Trace flows** before modifying them
- **Verify assumptions** about structure
- **Find similar implementations** to follow

## Research Strategy

### Break Into Focused Steps

Decompose research into small, specific questions:

```markdown
❌ "How does authentication work?"
✅ "Where are login routes defined?"
✅ "How are passwords hashed?"
✅ "Where are JWT tokens validated?"
```

Each question should:

- Have a clear, verifiable answer
- Focus on one aspect at a time
- Be answerable through code inspection

### Use Available Tools

Search the codebase with available tools:

1. **Start with semantic search** - Ask the context engine for concepts
2. **Find files by pattern** - Use glob for known naming conventions
3. **Search for strings** - Use grep for specific terms
4. **Read key files** - Examine implementation directly

### Delegate Background Tasks (When Supported)

If your tools support parallel/background execution:

1. **Identify parallel research paths** - What can be investigated independently?
2. **Launch focused searches** - Each search answers one specific question
3. **Synthesize findings** - Combine results into coherent understanding

Example:

```markdown
Search 1: "Where are user models defined?"
Search 2: "Where is authentication middleware?"  
Search 3: "How are session tokens stored?"
```

## Research Patterns

### Finding Files/Locations

```markdown
"Where are [X] handlers defined?"
"Find files related to [feature]"
"Locate [pattern] implementations"
```

### Understanding Patterns

```markdown
"How does this codebase handle [X]?"
"Show me error handling patterns"
"Find examples of [technique]"
```

### Tracing Flows

```markdown
"Trace: API route → controller → service → database"
"How does a request become a response?"
"Where is [X] transformed to [Y]?"
```

### Finding Integration Points

```markdown
"Where should [new feature] hook in?"
"What files would [change] affect?"
```

## Documenting Findings

Capture discoveries concisely:

```markdown
## Finding: Authentication Flow

- Entry: `src/routes/auth.js:15`
- Validation: `src/middleware/auth.js:42`
- Storage: `src/models/user.js:88`
```

## When NOT to Research

Skip research when:

- File paths are already known
- Change is isolated and understood
- Following an existing researched plan
- Task is purely about artifacts (state.md, plan.md, etc.)

## Research First, Then Act

For non-trivial changes:

1. **Decompose** - Break into specific questions
2. **Execute** - Search and investigate
3. **Synthesize** - Combine findings into understanding
4. **Document** - Capture paths and patterns
5. **Implement** - With confidence
