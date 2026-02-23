# AI coding agent extensions: a standardization feasibility analysis

**The agent-skills specification has already achieved partial standardization for rules/skills, but commands, hooks, and subagents present significantly greater challenges.** Six of nine tools now support the SKILL.md format with YAML frontmatter, creating a viable foundation. Commands show the most promise for expanded standardization—all tools support markdown-based prompts, though argument handling varies substantially. Hooks and subagents remain highly tool-specific, with only 4-5 tools offering comparable systems.

This analysis examines extension systems across Claude Code, Codex, OpenCode, Augment, Cursor, Windsurf, Cline, Kilo Code, and Droid to assess feasibility of a unified "agent extensions" repository.

## Support matrix reveals clear standardization opportunities

The following matrix shows current support levels across all tools and extension types:

| Tool | Skills/Rules | Commands | Subagents | Hooks | MCP/Plugins |
|------|--------------|----------|-----------|-------|-------------|
| **Claude Code** | ✅ Full (SKILL.md, CLAUDE.md) | ✅ Full ($ARGUMENTS) | ✅ Full (markdown+YAML) | ✅ Full (10 events) | ✅ Full |
| **Codex** | ✅ Full (SKILL.md, AGENTS.md) | ✅ Full (prompts/*.md) | ⚠️ Partial (profiles only) | ⚠️ Limited (notify only) | ✅ Full |
| **OpenCode** | ✅ Full (SKILL.md, AGENTS.md) | ✅ Full (shell injection) | ✅ Full (primary/subagent) | ✅ Full (25+ events via plugins) | ✅ Full |
| **Augment** | ✅ Full (rules/, AGENTS.md) | ✅ Full ($ARGUMENTS) | ✅ Full (markdown+YAML) | ❌ No (tool permissions only) | ✅ Full |
| **Cursor** | ✅ Full (.mdc rules) | ✅ Full (commands/*.md) | ❌ Removed (was supported) | ✅ Beta (6 events) | ✅ Full |
| **Windsurf** | ✅ Full (rules/, SKILL.md) | ✅ Full (workflows/) | ⚠️ Partial (via rules) | ✅ Full (11 events) | ✅ Full |
| **Cline** | ✅ Full (.clinerules, SKILL.md) | ✅ Full (workflows/) | ⚠️ Partial (Plan/Act only) | ✅ Full (6 events) | ✅ Full |
| **Kilo Code** | ✅ Full (rules/, SKILL.md) | ✅ Full (workflows/) | ✅ Full (custom modes) | ❌ No | ✅ Full |
| **Droid** | ✅ Full (AGENTS.md, skills/) | ✅ Full ($ARGUMENTS) | ✅ Full (droids/) | ✅ Full (9 events) | ✅ Full |

**Skills/Rules show 100% coverage** with only minor format variations. **Commands achieve 100% support** but with divergent argument systems. **Hooks are supported by 6/9 tools** with significant architectural differences. **MCP integration is universal**, making it the strongest standardization candidate after skills.

## Skills standardization is already working via agentskills.io

The Vercel agent-skills pattern demonstrates successful cross-tool deployment using a single specification:

**Universal SKILL.md format:**

```yaml
---
name: skill-name                    # Required, 1-64 chars
description: What it does           # Required, 1-1024 chars
license: Apache-2.0                 # Optional
compatibility: Requires git         # Optional
metadata:                           # Optional key-value pairs
  author: example-org
---
# Skill instructions in markdown...
```

**Installation paths by tool (add-skill CLI handles mapping):**

| Tool | Project Path | Global Path |
|------|--------------|-------------|
| Claude Code | `.claude/skills/<name>/` | `~/.claude/skills/<name>/` |
| Codex | `.codex/skills/<name>/` | `~/.codex/skills/<name>/` |
| OpenCode | `.opencode/skills/<name>/` | `~/.config/opencode/skills/<name>/` |
| Cursor | `.cursor/skills/<name>/` | `~/.cursor/skills/<name>/` |
| Windsurf | `.windsurf/skills/<name>/` | `~/.codeium/windsurf/skills/<name>/` |
| Cline | `.clinerules/skills/<name>/` | `~/Documents/Cline/Skills/<name>/` |
| Kilo Code | `.kilocode/skills/<name>/` | `~/.kilocode/skills/<name>/` |
| Droid | `.factory/skills/<name>/` | `~/.factory/skills/<name>/` |

The pattern works because **content is identical; only destination paths differ**. A unified repository can include a path-mapping configuration file that the installer reads.

**Rules (AGENTS.md/CLAUDE.md) achieve similar portability.** OpenCode, Augment, Cline, Windsurf, and Droid all support `AGENTS.md` as a fallback, providing cross-tool compatibility. Claude Code uses `CLAUDE.md` natively, but most tools will discover and use either filename.

## Commands can be standardized with "ask-first" pattern replacing arguments

All nine tools support markdown-based custom commands, but argument handling diverges significantly:

| Tool | Argument Syntax | Shell Injection | Model Override |
|------|-----------------|-----------------|----------------|
| Claude Code | `$ARGUMENTS` only | Via Bash tool | No |
| Codex | `$1-$9`, `$ARGUMENTS`, `KEY=value` | No | No |
| OpenCode | `$1-$9`, `$ARGUMENTS`, `$NAMED` | ✅ `!`backtick`` syntax | ✅ `model:` frontmatter |
| Augment | `$ARGUMENTS` implied | No | ✅ `model:` frontmatter |
| Cursor | Informal placeholders | No | No |
| Windsurf | `[placeholder]` syntax | No | No |
| Cline | Via `ask_followup_question` | No | No |
| Kilo Code | `$1-$9` | No | ✅ `model:` frontmatter |
| Droid | `$ARGUMENTS`, `$1-$9` | ✅ Executable scripts | No |

**The lowest common denominator for commands:**

```yaml
---
description: Brief description shown in slash menu
argument-hint: [optional-args]      # User guidance only
---
# Command instructions in markdown...
Use `$ARGUMENTS` for any user input.
```

**Recommendation: Standardize on "ask questions" pattern.** Rather than relying on argument interpolation (which varies), commands should explicitly instruct the agent to ask for required information:

```markdown
---
description: Create a new component with tests
---
Ask the user for:
1. Component name
2. Framework (React/Vue/Svelte)
3. Whether to include tests

Then create the component following project conventions...
```

This approach works universally because **every tool's agent can ask follow-up questions**, whereas argument parsing differs. OpenCode's shell injection (`!`cmd``) and Droid's executable scripts would need to be dropped for cross-tool compatibility.

**Features to drop for command standardization:**

- Shell injection syntax (OpenCode)
- Executable script commands (Droid)
- Model override frontmatter (OpenCode, Augment, Kilo Code)
- Positional arguments `$1-$9` (use `$ARGUMENTS` only)
- Named parameter syntax `$KEY` (Codex, OpenCode)

## Subagents show least standardization potential

Only **5 of 9 tools** have full subagent support with comparable architectures:

**Tools with file-based subagent definitions:**

| Tool | Location | Format | Key Features |
|------|----------|--------|--------------|
| Claude Code | `.claude/agents/<name>.md` | Markdown + YAML | `tools`, `model`, `permissionMode`, `skills` |
| OpenCode | `.opencode/agents/<name>.md` | Markdown + YAML | `mode`, `tools`, `permission`, `temperature` |
| Augment | `.augment/agents/<name>.md` | Markdown + YAML | `model`, `color` |
| Kilo Code | `.kilocodemodes` (YAML) | YAML | `groups`, `fileRegex`, `roleDefinition` |
| Droid | `.factory/droids/<name>.md` | Markdown + YAML | `tools`, `model`, `reasoningEffort` |

**Minimal viable subagent format:**

```yaml
---
name: code-reviewer
description: Reviews code for quality issues
model: inherit                      # Or specific model ID
tools: [Read, Grep, Glob]          # Tool restrictions
---
You are a code reviewer. Focus on...
```

**Blocking issues for subagent standardization:**

1. **Tool naming differs radically** — Claude Code uses `Bash`, OpenCode uses `execute`, Kilo Code uses `command`
2. **Permission models vary** — Some tools allow regex patterns, others use categories
3. **Model identifiers aren't portable** — `sonnet` vs `anthropic/claude-sonnet-4` vs provider-specific IDs
4. **4 tools lack comparable systems** — Codex (profiles only), Cursor (removed), Windsurf (via rules), Cline (Plan/Act only)

**Recommendation:** Subagents are **not feasible for cross-tool standardization** in the near term. The architectural differences are too significant. Instead, document tool-specific configurations and focus standardization efforts on skills and commands.

## Hooks require tool-specific implementations

Hooks represent the **most divergent** extension type:

| Tool | Config Format | Events | Can Block Actions |
|------|--------------|--------|-------------------|
| Claude Code | JSON in settings.json | 10 (PreToolUse, PostToolUse, Stop, etc.) | ✅ Yes |
| OpenCode | JS/TS plugins | 25+ (via plugin events) | ✅ Yes |
| Windsurf | JSON in hooks.json | 11 (pre/post for read/write/run) | ✅ Yes |
| Cursor | JSON in hooks.json | 6 (beforeShell, afterFileEdit, etc.) | ✅ Yes |
| Cline | Executable scripts | 6 (TaskStart, PreToolUse, etc.) | ✅ Yes |
| Droid | JSON in settings.json | 9 (mirrors Claude Code) | ✅ Yes |

**Common hook events across supporting tools:**

- Pre-tool execution (validate/block)
- Post-tool execution (format/log)
- Session start/end
- Task completion

**Standardization approach for hooks:**

A unified hooks spec could define **portable hook events** mapped to tool-specific implementations:

```json
{
  "hooks": {
    "before_file_write": {
      "command": "./hooks/format.sh",
      "timeout": 30
    },
    "after_task_complete": {
      "command": "./hooks/notify.sh"
    }
  }
}
```

The installer would transform these to tool-specific configurations:

- Claude Code → `PostToolUse` with `Write|Edit` matcher
- Windsurf → `post_write_code`
- Cursor → `afterFileEdit`
- Cline → `PostToolUse` script

**Features that cannot be standardized:**

- OpenCode's plugin-based hooks (require JS/TS)
- Event-specific matchers and regex patterns
- Input modification capabilities (only some tools support)
- Tool-specific environment variables

## MCP configuration is nearly universal but already fragmented

All 9 tools support MCP with JSON configuration, but **config file locations and formats already differ**:

| Tool | Config Location | Format Notes |
|------|-----------------|--------------|
| Claude Code | `.mcp.json`, `~/.claude.json` | Standard mcpServers object |
| Codex | `~/.codex/config.toml` | TOML format, unique |
| OpenCode | `opencode.json` | Under `mcp` key |
| Augment | Settings panel / JSON import | Standard format |
| Cursor | `.cursor/mcp.json` | Standard mcpServers object |
| Windsurf | `~/.codeium/windsurf/mcp_config.json` | Standard format |
| Cline | `mcp_settings.json` (in VS Code storage) | Adds `alwaysAllow`, `disabled` |
| Kilo Code | `.kilocode/mcp.json` | Adds `disabledTools` |
| Droid | `.factory/mcp.json` | Standard format |

**MCP standardization is trivial** — the config format is already identical across most tools. A unified repository simply needs installers that place the same `mcpServers` configuration in tool-specific locations.

## Recommended standardization tiers

Based on this analysis, extension types fall into three feasibility tiers:

### Tier 1: Ready for standardization now

- **Skills (SKILL.md)** — Already working via agentskills.io spec
- **Rules (AGENTS.md)** — Widely supported with fallback compatibility
- **MCP configurations** — Identical format, just different file locations

### Tier 2: Feasible with constraints

- **Commands** — Standardize on markdown + `$ARGUMENTS` only, use "ask questions" pattern instead of complex argument passing
- **Hooks** — Define portable event names mapped to tool-specific implementations; accept that some tools won't support all events

### Tier 3: Not recommended for standardization

- **Subagents** — Too architecturally different; maintain tool-specific configurations
- **Plugins (beyond MCP)** — Claude Code plugins, OpenCode plugins are completely different systems

## Repository architecture recommendation

A unified agent-extensions repository should mirror the Vercel agent-skills pattern:

```text
agent-extensions/
├── skills/                     # Tier 1 - already standardized
│   └── <skill-name>/SKILL.md
├── rules/                      # Tier 1 - portable rules
│   └── <rule-name>.md
├── commands/                   # Tier 2 - constrained standardization
│   └── <command-name>.md
├── hooks/                      # Tier 2 - portable event definitions
│   └── hooks.json
├── mcp/                        # Tier 1 - server configurations
│   └── servers.json
├── tool-configs/              # Tool-specific overrides
│   ├── claude-code/
│   ├── cursor/
│   └── ...
└── add-extension              # CLI installer (extends add-skill)
```

The installer CLI would:

1. Read extension files from unified format
2. Transform to tool-specific formats where needed
3. Copy to correct locations per tool
4. Support `--tool` flag for targeted installation

**The Vercel add-skill pattern proves this works.** The same content deploys to 25+ tools via simple path mapping. Extending this to commands and hooks requires additional transformation logic but follows the same principle: **universal source format, tool-specific destinations**.

## Key tradeoffs for cross-tool compatibility

| Feature | Keep for standardization | Drop for standardization |
|---------|-------------------------|--------------------------|
| **Skills** | YAML frontmatter (name, description), markdown body, supporting files | `allowed-tools` (experimental), tool-specific metadata |
| **Commands** | `$ARGUMENTS`, description frontmatter, markdown prompts | Shell injection, executable scripts, model override, positional args |
| **Hooks** | Common events (pre/post write, task complete), shell command execution | Input modification, tool-specific matchers, plugin-based hooks |
| **Subagents** | N/A — keep tool-specific | All standardization attempts |

The goal is maximizing reach across tools while accepting that power users needing tool-specific features can maintain separate configurations. A 90% solution that works everywhere beats a 100% solution that only works in one tool.
