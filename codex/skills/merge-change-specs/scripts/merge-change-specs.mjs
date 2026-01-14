import { existsSync } from "node:fs";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import * as path from "node:path";

import { validateChangeSpecMarkdown } from "../../spec-format/scripts/validate-change-spec.mjs";

function die(message) {
  process.stderr.write(`${message}\n`);
  process.exit(1);
}

function assertRepoRelative(p) {
  if (p.startsWith("/") || p.includes("..")) {
    die(`Refusing unsafe path: ${p}`);
  }
}

function usage() {
  return [
    "Usage:",
    "  node codex/skills/merge-change-specs/scripts/merge-change-specs.mjs --change <name> [--dry-run]",
    "",
    "Examples:",
    "  node codex/skills/merge-change-specs/scripts/merge-change-specs.mjs --change auth-refresh",
    "  node codex/skills/merge-change-specs/scripts/merge-change-specs.mjs --change auth-refresh --dry-run",
  ].join("\n");
}

function parseArgs(argv) {
  let change;
  let dryRun = false;

  for (let index = 0; index < argv.length; index++) {
    const arg = argv[index];
    if (arg === "-h" || arg === "--help") {
      process.stdout.write(usage() + "\n");
      process.exit(0);
    }
    if (arg === "--change" || arg === "-c") {
      change = argv[index + 1];
      index++;
      continue;
    }
    if (arg === "--dry-run") {
      dryRun = true;
      continue;
    }
  }

  if (!change) {
    process.stderr.write(usage() + "\n");
    process.exit(1);
  }

  assertRepoRelative(change);

  return { change, dryRun };
}

function splitFrontmatter(markdown) {
  if (!markdown.startsWith("---\n")) {
    return { frontmatterRaw: null, body: markdown };
  }

  const endIndex = markdown.indexOf("\n---\n", "---\n".length);
  if (endIndex === -1) {
    die("Found starting frontmatter marker but no closing ---");
  }

  const frontmatterRaw = markdown.slice(4, endIndex + 1);
  const body = markdown.slice(endIndex + "\n---\n".length);
  return { frontmatterRaw, body };
}

function parseKind(frontmatterRaw, fileLabel) {
  if (!frontmatterRaw) {
    die(`Missing YAML frontmatter in change spec: ${fileLabel}`);
  }

  const kindMatch = frontmatterRaw.match(/^kind:\s*(new|delta)\s*$/m);
  if (!kindMatch) {
    die(`Missing required 'kind: new|delta' in frontmatter: ${fileLabel}`);
  }

  return kindMatch[1];
}

async function listChangeSpecs(changeSpecsDirAbs) {
  const { readdir } = await import("node:fs/promises");

  const results = [];

  async function walkFs(dirAbs) {
    const dirents = await readdir(dirAbs, { withFileTypes: true });
    const sorted = dirents.slice().sort((a, b) => a.name.localeCompare(b.name));

    for (const dirent of sorted) {
      const nextAbs = path.join(dirAbs, dirent.name);
      if (dirent.isDirectory()) {
        await walkFs(nextAbs);
        continue;
      }

      if (dirent.isFile() && dirent.name.endsWith(".md")) {
        results.push(nextAbs);
      }
    }
  }

  await walkFs(changeSpecsDirAbs);

  results.sort((a, b) => a.localeCompare(b));
  return results;
}

function computeCanonicalRelPath(changeName, changeSpecAbs, repoRootAbs) {
  const changeSpecsDirAbs = path.join(repoRootAbs, "changes", changeName, "specs");
  const relFromChangeSpecs = path.relative(changeSpecsDirAbs, changeSpecAbs);
  if (relFromChangeSpecs.startsWith("..")) {
    die(`Change spec is outside expected directory: ${changeSpecAbs}`);
  }

  const canonicalRelPath = path.posix.join(
    "specs",
    relFromChangeSpecs.split(path.sep).join(path.posix.sep),
  );
  return canonicalRelPath;
}

function normalizeNewSpecBody(changeSpecBody) {
  return changeSpecBody.replace(/^\s+/, "");
}

function extractDeltaBuckets(deltaBody, fileLabel) {
  const requirementsHeader = "## Requirements";
  const requirementsIndex = deltaBody.indexOf(requirementsHeader);
  if (requirementsIndex === -1) {
    die(`Delta spec missing '## Requirements': ${fileLabel}`);
  }

  const nextSectionIndex = deltaBody.indexOf("\n## ", requirementsIndex + requirementsHeader.length);
  const requirementsSection = deltaBody.slice(
    requirementsIndex,
    nextSectionIndex === -1 ? deltaBody.length : nextSectionIndex + 1,
  );

  const addedIndex = requirementsSection.indexOf("### ADDED");
  const modifiedIndex = requirementsSection.indexOf("### MODIFIED");
  const removedIndex = requirementsSection.indexOf("### REMOVED");

  if (modifiedIndex === -1 && addedIndex === -1 && removedIndex === -1) {
    die(`Delta spec has no ADDED/MODIFIED/REMOVED buckets under Requirements: ${fileLabel}`);
  }

  function sliceBucket(start, endCandidates) {
    if (start === -1) return "";
    const end = Math.min(
      ...endCandidates
        .filter((n) => n !== -1 && n > start)
        .concat([requirementsSection.length]),
    );
    return requirementsSection.slice(start, end).trimEnd() + "\n";
  }

  return {
    added: sliceBucket(addedIndex, [modifiedIndex, removedIndex]),
    modified: sliceBucket(modifiedIndex, [addedIndex, removedIndex]),
    removed: sliceBucket(removedIndex, [addedIndex, modifiedIndex]),
  };
}

function patchCanonicalWithDelta(canonicalBody, deltaBody, fileLabel) {
  const buckets = extractDeltaBuckets(deltaBody, fileLabel);

  let next = canonicalBody;

  function removeExactBlocks(bucketContent) {
    const touchedTopics = new Set();
    if (!bucketContent) return touchedTopics;

    const body = bucketContent.replace(/^###\s+REMOVED\s*\n?/m, "").trim();

    const topicChunks = body
      .split(/^####\s+/m)
      .map((s) => s.trim())
      .filter(Boolean);

    for (const chunk of topicChunks) {
      const lines = chunk.split("\n");
      const topic = (lines.shift() ?? "").trim();
      if (topic) touchedTopics.add(topic);

      const remainder = lines.join("\n").trim();
      if (!remainder) continue;

      const reasonIndex = remainder.indexOf("**Reason:**");
      const block = (reasonIndex === -1 ? remainder : remainder.slice(0, reasonIndex)).trim();
      if (!block) continue;

      const normalizedBlock = block.replace(/\r\n/g, "\n");
      const pos = next.indexOf(normalizedBlock);
      if (pos === -1) {
        die(`REMOVED block not found in canonical spec for: ${fileLabel}`);
      }

      next = next.slice(0, pos) + next.slice(pos + normalizedBlock.length);
    }

    return touchedTopics;
  }

  function cleanupEmptyTopics(topics) {
    for (const topic of topics) {
      const topicHeader = `### ${topic}`;
      const headerIndex = next.indexOf(topicHeader);
      if (headerIndex === -1) continue;

      const sectionStart = headerIndex;
      const contentStart = headerIndex + topicHeader.length;
      const followingIndex = next.indexOf("\n### ", contentStart);
      const sectionEnd = followingIndex === -1 ? next.length : followingIndex + 1;

      const sectionBody = next.slice(contentStart, sectionEnd);
      const hasBullet = /^\s*-\s+/m.test(sectionBody);
      if (hasBullet) continue;

      const before = next.slice(0, sectionStart).replace(/\n{3,}$/m, "\n\n");
      const after = next.slice(sectionEnd).replace(/^\n{3,}/m, "\n\n");
      next = before + after;
    }
  }

  function insertAddedBlocks(bucketContent) {
    if (!bucketContent) return;

    const body = bucketContent.replace(/^###\s+ADDED\s*\n?/m, "").trim();

    const topicChunks = body
      .split(/^####\s+/m)
      .map((s) => s.trim())
      .filter(Boolean);

    for (const chunk of topicChunks) {
      const lines = chunk.split("\n");
      const topic = (lines.shift() ?? "").trim();
      const insertion = lines.join("\n").trim();
      if (!topic || !insertion) continue;

      const topicHeader = `### ${topic}`;
      const topicIndex = next.indexOf(topicHeader);
      if (topicIndex === -1) {
        die(`ADDED topic not found in canonical spec (${topicHeader}) for: ${fileLabel}`);
      }

      const searchStart = topicIndex + topicHeader.length;
      const nextTopicIndex = next.indexOf("\n### ", searchStart);
      const insertAt = nextTopicIndex === -1 ? next.length : nextTopicIndex + 1;

      const toInsert = insertion.replace(/\r\n/g, "\n").trimEnd();
      next = next.slice(0, insertAt) + `\n${toInsert}\n` + next.slice(insertAt);
    }
  }

  const removedTopics = removeExactBlocks(buckets.removed);

  if (buckets.modified) {
    const modifiedBody = buckets.modified.replace(/^###\s+MODIFIED\s*\n?/m, "").trim();

    const topicChunks = modifiedBody
      .split(/^####\s+/m)
      .map((s) => s.trim())
      .filter(Boolean);

    for (const chunk of topicChunks) {
      const beforeLabelIndex = chunk.indexOf("**Before:**");
      const afterLabelIndex = chunk.indexOf("**After:**");
      if (beforeLabelIndex === -1 || afterLabelIndex === -1 || afterLabelIndex < beforeLabelIndex) {
        die(`Invalid MODIFIED topic; expected **Before:** then **After:** in: ${fileLabel}`);
      }

      const beforeBlock = chunk.slice(beforeLabelIndex + "**Before:**".length, afterLabelIndex).trim();
      const afterBlock = chunk.slice(afterLabelIndex + "**After:**".length).trim();

      if (!beforeBlock || !afterBlock) {
        die(`Empty Before/After block in MODIFIED topic: ${fileLabel}`);
      }

      const normalizedBefore = beforeBlock.replace(/\r\n/g, "\n");
      const normalizedAfter = afterBlock.replace(/\r\n/g, "\n");

      const beforePos = next.indexOf(normalizedBefore);
      if (beforePos === -1) {
        die(`Before block not found in canonical spec for: ${fileLabel}`);
      }

      next = next.slice(0, beforePos) + normalizedAfter + next.slice(beforePos + normalizedBefore.length);
    }
  }

  insertAddedBlocks(buckets.added);
  cleanupEmptyTopics(removedTopics);

  return next;
}

async function ensureParentDir(fileAbs) {
  await mkdir(path.dirname(fileAbs), { recursive: true });
}

async function main() {
  const { change, dryRun } = parseArgs(process.argv.slice(2));
  const repoRootAbs = process.cwd();

  const changeDirAbs = path.join(repoRootAbs, "changes", change);
  const changeSpecsDirAbs = path.join(changeDirAbs, "specs");
  if (!existsSync(changeSpecsDirAbs)) {
    die(`Missing change specs directory: ${path.relative(repoRootAbs, changeSpecsDirAbs)}`);
  }

  const changeSpecAbsPaths = await listChangeSpecs(changeSpecsDirAbs);

  const plan = [];
  for (const changeSpecAbs of changeSpecAbsPaths) {
    const canonicalRelPath = computeCanonicalRelPath(change, changeSpecAbs, repoRootAbs);
    plan.push({
      canonicalRelPath,
      changeSpecRelPath: path.relative(repoRootAbs, changeSpecAbs),
      kind: "new",
    });
  }

  for (const item of plan) {
    const changeSpecAbs = path.join(repoRootAbs, item.changeSpecRelPath);
    const changeMd = await readFile(changeSpecAbs, "utf8");

    const validation = validateChangeSpecMarkdown(changeMd);
    if (!validation.ok) {
      die(
        `Invalid change-set spec: ${item.changeSpecRelPath}\n` +
          validation.issues.map((i) => `- [${i.code}] ${i.message}`).join("\n"),
      );
    }

    const { frontmatterRaw } = splitFrontmatter(changeMd);
    item.kind = parseKind(frontmatterRaw, item.changeSpecRelPath);
  }

  plan.sort((a, b) => a.canonicalRelPath.localeCompare(b.canonicalRelPath));

  const summary = { created: [], modified: [], skipped: [] };

  for (const item of plan) {
    const changeSpecAbs = path.join(repoRootAbs, item.changeSpecRelPath);
    const canonicalAbs = path.join(repoRootAbs, item.canonicalRelPath);

    const changeMd = await readFile(changeSpecAbs, "utf8");
    const { body } = splitFrontmatter(changeMd);

    if (item.kind === "new") {
      const nextBody = normalizeNewSpecBody(body);
      const existed = existsSync(canonicalAbs);

      if (!dryRun) {
        await ensureParentDir(canonicalAbs);
        await writeFile(canonicalAbs, nextBody, "utf8");
      }

      if (existed) summary.modified.push(item.canonicalRelPath);
      else summary.created.push(item.canonicalRelPath);

      continue;
    }

    if (!existsSync(canonicalAbs)) {
      die(`Delta targets missing canonical spec: ${item.canonicalRelPath}`);
    }

    const canonicalMd = await readFile(canonicalAbs, "utf8");
    const patched = patchCanonicalWithDelta(canonicalMd, body, item.changeSpecRelPath);

    if (!dryRun) {
      await writeFile(canonicalAbs, patched, "utf8");
    }
    summary.modified.push(item.canonicalRelPath);
  }

  summary.created.sort();
  summary.modified.sort();
  summary.skipped.sort();

  process.stdout.write(
    JSON.stringify(
      {
        change,
        dryRun,
        counts: {
          created: summary.created.length,
          modified: summary.modified.length,
          skipped: summary.skipped.length,
        },
        created: summary.created,
        modified: summary.modified,
        skipped: summary.skipped,
      },
      null,
      2,
    ) + "\n",
  );
}

void main();
