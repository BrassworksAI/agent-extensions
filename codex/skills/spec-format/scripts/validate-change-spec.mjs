import { readFile } from "node:fs/promises";
import * as path from "node:path";
import { pathToFileURL } from "node:url";

function issue(code, message) {
  return { code, message };
}

function splitFrontmatter(markdown) {
  if (typeof markdown !== "string") {
    return { frontmatterRaw: null, body: "" };
  }

  if (!markdown.startsWith("---\n")) {
    return { frontmatterRaw: null, body: markdown };
  }

  const endIndex = markdown.indexOf("\n---\n", "---\n".length);
  if (endIndex === -1) {
    return { frontmatterRaw: null, body: markdown };
  }

  const frontmatterRaw = markdown.slice(4, endIndex + 1);
  const body = markdown.slice(endIndex + "\n---\n".length);
  return { frontmatterRaw, body };
}

function parseKind(frontmatterRaw) {
  if (!frontmatterRaw) return null;
  const kindMatch = frontmatterRaw.match(/^kind:\s*(new|delta)\s*$/m);
  if (!kindMatch) return null;
  return kindMatch[1];
}

function hasHeadingPrefix(body, prefix) {
  const escaped = prefix.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  return new RegExp(`^${escaped}\\s+.+$`, "m").test(body);
}

function validateNew(body) {
  const issues = [];

  if (!hasHeadingPrefix(body, "#")) {
    issues.push(issue("new.missing_title", "Missing top-level '# <Title>' heading"));
  }
  if (!hasHeadingPrefix(body, "## Overview")) {
    issues.push(issue("new.missing_overview", "Missing '## Overview' section"));
  }
  if (!hasHeadingPrefix(body, "## Requirements")) {
    issues.push(issue("new.missing_requirements", "Missing '## Requirements' section"));
  }

  if (/^###\s+(ADDED|MODIFIED|REMOVED)\s*$/m.test(body)) {
    issues.push(
      issue(
        "new.delta_buckets",
        "`kind: new` spec must not include delta buckets (### ADDED/MODIFIED/REMOVED)",
      ),
    );
  }

  return issues;
}

function validateDelta(body) {
  const issues = [];

  if (!hasHeadingPrefix(body, "#")) {
    issues.push(issue("delta.missing_title", "Missing top-level '# <Title>' heading"));
  }
  if (!hasHeadingPrefix(body, "## Requirements")) {
    issues.push(issue("delta.missing_requirements", "Missing '## Requirements' section"));
  }

  const requirementsIndex = body.indexOf("## Requirements");
  if (requirementsIndex !== -1) {
    const nextSectionIndex = body.indexOf("\n## ", requirementsIndex + "## Requirements".length);
    const requirementsSection = body.slice(
      requirementsIndex,
      nextSectionIndex === -1 ? body.length : nextSectionIndex + 1,
    );

    const hasBucket = /^(###\s+(ADDED|MODIFIED|REMOVED)\s*)$/m.test(requirementsSection);
    if (!hasBucket) {
      issues.push(
        issue(
          "delta.missing_buckets",
          "Delta spec must include at least one bucket under '## Requirements'",
        ),
      );
    }

    if (/^###\s+MODIFIED\s*$/m.test(requirementsSection)) {
      const modifiedStart = requirementsSection.indexOf("### MODIFIED");
      const nextBucketIndex = [
        requirementsSection.indexOf("### ADDED", modifiedStart + 1),
        requirementsSection.indexOf("### REMOVED", modifiedStart + 1),
      ]
        .filter((n) => n !== -1)
        .sort((a, b) => a - b)[0];

      const modifiedBlock = requirementsSection.slice(
        modifiedStart,
        nextBucketIndex === undefined ? requirementsSection.length : nextBucketIndex,
      );

      const topics = modifiedBlock
        .replace(/^###\s+MODIFIED\s*\n?/m, "")
        .trim()
        .split(/^####\s+/m)
        .map((s) => s.trim())
        .filter(Boolean);

      for (const topic of topics) {
        const beforeIndex = topic.indexOf("**Before:**");
        const afterIndex = topic.indexOf("**After:**");
        if (beforeIndex === -1 || afterIndex === -1 || afterIndex < beforeIndex) {
          issues.push(
            issue(
              "delta.modified_malformed",
              "Each MODIFIED topic must contain '**Before:**' followed by '**After:**'",
            ),
          );
          break;
        }
      }
    }
  }

  return issues;
}

export function validateChangeSpecMarkdown(markdown) {
  const { frontmatterRaw, body } = splitFrontmatter(markdown);
  const kind = parseKind(frontmatterRaw);

  const issues = [];

  if (!frontmatterRaw) {
    issues.push(issue("fm.missing", "Missing YAML frontmatter (--- ... ---)"));
  }
  if (!kind) {
    issues.push(issue("fm.kind_missing", "Missing required 'kind: new|delta' in frontmatter"));
  }

  if (kind === "new") issues.push(...validateNew(body));
  if (kind === "delta") issues.push(...validateDelta(body));

  return { ok: issues.length === 0, kind, issues };
}

function formatIssues(rel, issues) {
  return `Spec validation failed: ${rel}\n${issues.map((i) => `- [${i.code}] ${i.message}`).join("\n")}`;
}

function usage() {
  return [
    "Usage:",
    "  node codex/skills/spec-format/scripts/validate-change-spec.mjs <change-spec-path>",
    "",
    "Examples:",
    "  node codex/skills/spec-format/scripts/validate-change-spec.mjs changes/auth-refresh/specs/auth/login.md",
  ].join("\n");
}

async function main() {
  const arg = process.argv[2];
  if (!arg || arg === "-h" || arg === "--help") {
    process.stdout.write(usage() + "\n");
    process.exit(arg ? 0 : 1);
  }

  if (path.isAbsolute(arg)) {
    process.stderr.write("Refusing absolute path; pass a repo-relative path.\n");
    process.exit(1);
  }
  if (arg.includes("..")) {
    process.stderr.write("Refusing path traversal; '..' is not allowed.\n");
    process.exit(1);
  }

  const repoRootAbs = process.cwd();
  const fileAbs = path.join(repoRootAbs, arg);

  const markdown = await readFile(fileAbs, "utf8");
  const result = validateChangeSpecMarkdown(markdown);
  if (result.ok) {
    process.stdout.write("OK\n");
    return;
  }

  process.stderr.write(formatIssues(arg, result.issues) + "\n");
  process.exit(1);
}

if (import.meta.url === pathToFileURL(process.argv[1])?.href) {
  void main();
}
