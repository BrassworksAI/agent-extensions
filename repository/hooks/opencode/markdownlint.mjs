#!/usr/bin/env node
import { existsSync } from "node:fs";

const MARKDOWN_EXTENSIONS = [".md", ".markdown", ".mdown", ".mkd", ".mkdn"];
const MARKDOWNLINT_CONFIG_FILES = [
  ".markdownlint.json",
  ".markdownlint.yaml",
  ".markdownlint.yml",
  ".markdownlint.js",
  ".markdownlint.cjs",
  ".markdownlintrc",
  ".markdownlintrc.json",
  ".markdownlintrc.yaml",
  ".markdownlintrc.yml",
  ".markdownlintrc.js",
  ".markdownlintrc.cjs",
];

function isMarkdownFilePath(filePath) {
  const lower = filePath.toLowerCase();
  return MARKDOWN_EXTENSIONS.some((ext) => lower.endsWith(ext));
}

function findMarkdownlintConfig(worktree) {
  const root = worktree.replace(/[\\/]+$/, "");

  for (const config of MARKDOWNLINT_CONFIG_FILES) {
    const configPath = `${root}/${config}`;
    if (existsSync(configPath)) return configPath;
  }

  return null;
}

export const MarkdownlintPlugin = async ({ $, worktree }) => {
  const configPath = findMarkdownlintConfig(worktree);

  async function validateFile(filePath) {
    if (!isMarkdownFilePath(filePath)) return;

    const result = configPath
      ? await $`markdownlint -c ${configPath} -f ${filePath}`.quiet().nothrow()
      : await $`markdownlint -f ${filePath}`.quiet().nothrow();

    if (result.exitCode !== 0) {
      const errorMsg = [
        `markdownlint failed for: ${filePath}`,
        `exit code: ${result.exitCode}`,
      ];

      if (result.stdout.toString().trim()) {
        errorMsg.push(`stdout: ${result.stdout.toString().trim()}`);
      }
      if (result.stderr.toString().trim()) {
        errorMsg.push(`stderr: ${result.stderr.toString().trim()}`);
      }

      errorMsg.push(`You MUST fix the above issues before proceeding.`);

      throw new Error(errorMsg.join("\n"));
    }
  }

  return {
    "tool.execute.after": async (input, output) => {
      if (input.tool !== "edit" && input.tool !== "write") return;

      await validateFile(output.title);
    },
  };
};
