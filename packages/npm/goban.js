#!/usr/bin/env node
"use strict";

const { spawnSync } = require("child_process");
const https = require("https");
const fs = require("fs");
const path = require("path");
const zlib = require("zlib");

const REPO = "icl00ud/goban";
const VERSION = require("./package.json").version;
const BIN_DIR = path.join(__dirname, "bin");
const IS_WIN = process.platform === "win32";
const BIN_PATH = path.join(BIN_DIR, IS_WIN ? "goban.exe" : "goban");

function run() {
  const result = spawnSync(BIN_PATH, process.argv.slice(2), { stdio: "inherit" });
  process.exit(result.status ?? 0);
}

// Binary already installed — just run it
if (fs.existsSync(BIN_PATH)) {
  run();
  process.exit(0);
}

// First run: download the binary
const PLATFORM_MAP = { linux: "linux", darwin: "darwin", win32: "windows" };
const ARCH_MAP = { x64: "amd64", arm64: "arm64" };

const platform = PLATFORM_MAP[process.platform];
const arch = ARCH_MAP[process.arch];

if (!platform || !arch) {
  console.error(`[goban] Unsupported platform: ${process.platform}/${process.arch}`);
  console.error("[goban] Download manually from: https://github.com/icl00ud/goban/releases/latest");
  process.exit(1);
}

const ext = IS_WIN ? "zip" : "tar.gz";
const filename = `goban_${VERSION}_${platform}_${arch}.${ext}`;
const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${filename}`;

console.error(`[goban] First run — downloading binary for ${platform}/${arch}...`);
fs.mkdirSync(BIN_DIR, { recursive: true });

follow(url, (err, res) => {
  if (err) {
    console.error("[goban] Download failed:", err.message);
    console.error("[goban] Download manually from: https://github.com/icl00ud/goban/releases/latest");
    process.exit(1);
  }

  if (IS_WIN) {
    const chunks = [];
    res.on("data", (c) => chunks.push(c));
    res.on("end", () => {
      const zipPath = path.join(BIN_DIR, "goban.zip");
      fs.writeFileSync(zipPath, Buffer.concat(chunks));
      const { execFileSync } = require("child_process");
      execFileSync("powershell", ["-Command", `Expand-Archive -Force '${zipPath}' '${BIN_DIR}'`]);
      fs.unlinkSync(zipPath);
      console.error("[goban] Ready.");
      run();
    });
  } else {
    extractTarGz(res, BIN_PATH, () => {
      fs.chmodSync(BIN_PATH, 0o755);
      console.error("[goban] Ready.");
      run();
    });
  }
});

function follow(url, cb) {
  https.get(url, (res) => {
    if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
      follow(res.headers.location, cb);
    } else if (res.statusCode !== 200) {
      cb(new Error(`HTTP ${res.statusCode}`));
    } else {
      cb(null, res);
    }
  }).on("error", cb);
}

function extractTarGz(stream, destFile, done) {
  const gunzip = zlib.createGunzip();
  stream.pipe(gunzip);
  let buf = Buffer.alloc(0);
  gunzip.on("data", (chunk) => { buf = Buffer.concat([buf, chunk]); });
  gunzip.on("end", () => {
    let offset = 0;
    while (offset + 512 <= buf.length) {
      const header = buf.slice(offset, offset + 512);
      const name = header.slice(0, 100).toString("utf8").replace(/\0/g, "");
      if (!name) break;
      const size = parseInt(header.slice(124, 136).toString("utf8").trim(), 8) || 0;
      offset += 512;
      if (path.basename(name) === "goban" || path.basename(name) === "goban.exe") {
        fs.writeFileSync(destFile, buf.slice(offset, offset + size));
        done();
        return;
      }
      offset += Math.ceil(size / 512) * 512;
    }
    console.error("[goban] Binary not found in archive.");
    process.exit(1);
  });
  gunzip.on("error", (e) => { console.error("[goban] Decompression error:", e.message); process.exit(1); });
}
