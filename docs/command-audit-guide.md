# Command Audit Guide

This guide provides criteria for auditing commands in the repository based on best practices from Claude, Cursor, Windsurf, and other AI coding platforms.

## Audit Checklist

### 1. Purpose & Clarity

- [ ] Command has a single, clear purpose (one thing well)
- [ ] Title is action-oriented and descriptive
- [ ] Description provides a one-sentence summary
- [ ] Command solves a specific, recurring problem

### 2. Structure & Organization

- [ ] Frontmatter includes ONLY: description (other fields break multi-agent support)
- [ ] Content follows logical flow: title (H1) → context → instructions → examples
- [ ] Uses clear headers (`##`) for sections
- [ ] Word count within bounds: 150-400 words total

### 3. Content Quality

- [ ] Core instructions are 60-120 words (concise)
- [ ] Focuses on "why" not just "what" (decisions/constraints)
- [ ] Specifies behavior over implementation details
- [ ] Includes edge cases in instructions
- [ ] Provides clear error handling guidance

### 4. Examples & Usage

- [ ] Includes concrete usage examples
- [ ] Examples show "do X, not Y" patterns
- [ ] Success criteria are explicit and verifiable
- [ ] Examples are 20-60 words each

### 5. Context Awareness

- [ ] Command knows its scope (file, directory, project)
- [ ] Respects environment conventions
- [ ] Adapts to project structure
- [ ] Uses @mentions/context flags appropriately

### 6. Progressive Disclosure

- **Layer 1 (Essential)**: Always visible
  - [ ] Command name and one-line purpose
  - [ ] Key parameters with defaults
  - [ ] One usage example
  - [ ] Success criteria

- **Layer 2 (Key Details)**: Expandable
  - [ ] Step-by-step instructions
  - [ ] Project-specific conventions
  - [ ] Common edge cases

- **Layer 3 (Advanced)**: On-demand
  - [ ] Detailed rationale
  - [ ] Alternative approaches
  - [ ] Integration notes (if complex)

### 7. Anti-Patterns Check

- [ ] No vague instructions (e.g., "improve the code")
- [ ] No missing constraints (e.g., "fix tests" without boundaries)
- [ ] Success criteria are present
- [ ] Not overly prescriptive (dictating implementation)
- [ ] Assumes defined conventions

### 8. Consistency

- [ ] Follows project coding style conventions
- [ ] Uses consistent terminology with other commands
- [ ] Matches existing command structure patterns

### 9. Completeness

- [ ] All necessary steps are included
- [ ] No "magic" steps or undefined references
- [ ] Prerequisites are clear if needed
- [ ] Dependencies on other commands are documented

### 10. Testability

- [ ] Command can be executed successfully
- [ ] Expected outcomes are clear
- [ ] Verification steps are provided

### 11. Interaction & State

- [ ] **Mandatory Questions**: Command must ask the user for any required inputs (paths, names, etc.) instead of assuming CLI arguments are supported.
- [ ] **Input Formatting**: Sections requiring user input follow the standard callout format:

  ```markdown
  ## Inputs

  > [!IMPORTANT]
  > You must ask the user for the following information; do not assume CLI arguments are provided.

  - **Input Name**: Description
  ```

- [ ] **SDD State Management**: For all SDD commands, the agent must:
  - Check for existing change set folders.
  - If only one exists, proceed with it.
  - If multiple exist, ask the user to specify which change set name to use.

## Quality Thresholds

### Red Flags (Needs Major Revision)

- Exceeds 400 words total
- No clear purpose or tries to do multiple things
- Missing success criteria
- Vague or ambiguous instructions
- No examples

### Yellow Flags (Needs Minor Revision)

- Between 300-400 words borderline
- Examples could be more concrete
- Edge cases not addressed
- Inconsistent with other commands

### Green Flags (Ready)

- 150-300 words total
- Single clear purpose
- Concrete examples with success criteria
- Follows structure conventions
- Progressive disclosure applied

## Audit Process

1. **Read Command**: Read the full command file
2. **Apply Checklist**: Go through each checklist item
3. **Identify Issues**: Note red/yellow flags
4. **Propose Changes**: List specific improvements needed
5. **Rate Quality**: Red (rewrite), Yellow (update), Green (approve)

## Notes

- Commands exceeding 400 words should be split into multiple focused commands
- Commands with 3-5 distinct steps may need chaining logic
- Optional branches should be separate command variants
- Focus on reducing cognitive load while enabling expert outcomes
