#!/usr/bin/env node

import { existsSync } from "fs";
import { join, dirname } from "path";
import { fileURLToPath } from "url";
import { spawn } from "child_process";

const __dirname = dirname(fileURLToPath(import.meta.url));
const packageRoot = join(__dirname, "..");
const executable = process.platform === "win32" ? "ae.exe" : "ae";
const binaryPath = join(packageRoot, "bin", executable);

if (!existsSync(binaryPath)) {
  console.error("ae binary not found in package installation");
  console.error("Try reinstalling: npm rebuild @shanepadgett/agent-extensions");
  process.exit(1);
}

const child = spawn(binaryPath, process.argv.slice(2), { stdio: "inherit" });

child.on("exit", (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 1);
});
