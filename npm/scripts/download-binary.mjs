#!/usr/bin/env node
import { createWriteStream, mkdirSync, chmodSync, existsSync, unlinkSync } from "fs";
import { join, dirname } from "path";
import { fileURLToPath } from "url";
import { execSync } from "child_process";
import https from "https";

const __dirname = dirname(fileURLToPath(import.meta.url));
const packageRoot = join(__dirname, "..");
const binDir = join(packageRoot, "bin");
const binPath = join(binDir, process.platform === "win32" ? "ae.exe" : "ae");

const VERSION = process.env.npm_package_version || "0.1.0";
const REPO = "shanepadgett/agent-extensions";

function getPlatform() {
  const platform = process.platform;
  const arch = process.arch;

  const platformMap = {
    darwin: "darwin",
    linux: "linux",
    win32: "windows",
  };

  const archMap = {
    x64: "amd64",
    arm64: "arm64",
  };

  const os = platformMap[platform];
  const cpu = archMap[arch];

  if (!os || !cpu) {
    throw new Error(`Unsupported platform: ${platform}-${arch}`);
  }

  return { os, arch: cpu };
}

function getDownloadUrl(version, os, arch) {
  const ext = os === "windows" ? "zip" : "tar.gz";
  const filename = `ae_${version}_${os}_${arch}.${ext}`;
  return `https://github.com/${REPO}/releases/download/v${version}/${filename}`;
}

function download(url) {
  return new Promise((resolve, reject) => {
    const request = (url) => {
      https.get(url, (res) => {
        if (res.statusCode === 302 || res.statusCode === 301) {
          request(res.headers.location);
          return;
        }
        if (res.statusCode !== 200) {
          reject(new Error(`Download failed: ${res.statusCode} for ${url}`));
          return;
        }
        const chunks = [];
        res.on("data", (chunk) => chunks.push(chunk));
        res.on("end", () => resolve(Buffer.concat(chunks)));
        res.on("error", reject);
      }).on("error", reject);
    };
    request(url);
  });
}

async function extract(buffer, os) {
  const tmpFile = join(binDir, os === "windows" ? "tmp.zip" : "tmp.tar.gz");
  
  mkdirSync(binDir, { recursive: true });
  
  const fs = await import("fs/promises");
  await fs.writeFile(tmpFile, buffer);

  try {
    if (os === "windows") {
      execSync(`powershell -Command "Expand-Archive -Path '${tmpFile}' -DestinationPath '${binDir}' -Force"`, {
        stdio: "inherit",
      });
    } else {
      execSync(`tar -xzf "${tmpFile}" -C "${binDir}"`, { stdio: "inherit" });
    }
  } finally {
    if (existsSync(tmpFile)) unlinkSync(tmpFile);
  }

  if (process.platform !== "win32") {
    chmodSync(binPath, 0o755);
  }
}

async function main() {
  if (existsSync(binPath)) {
    console.log("ae binary already exists, skipping download");
    return;
  }

  const { os, arch } = getPlatform();
  const url = getDownloadUrl(VERSION, os, arch);

  console.log(`Downloading ae v${VERSION} for ${os}-${arch}...`);
  console.log(`URL: ${url}`);

  try {
    const buffer = await download(url);
    await extract(buffer, os);
    console.log("ae installed successfully!");
  } catch (err) {
    console.error(`Failed to install ae: ${err.message}`);
    console.error("You can manually download from:", `https://github.com/${REPO}/releases`);
    process.exit(1);
  }
}

main();
