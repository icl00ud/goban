#!/usr/bin/env node
// Downloads the correct goban binary from GitHub Releases on install.
"use strict";

const https = require("https");
const fs = require("fs");
const path = require("path");
const zlib = require("zlib");

const REPO = "icl00ud/goban";
const VERSION = require("./package.json").version;
const BIN_DIR = path.join(__dirname, "bin");
const BIN_PATH = path.join(BIN_DIR, process.platform === "win32" ? "goban.exe" : "goban");

const PLATFORM_MAP = { linux: "linux", darwin: "darwin", win32: "windows" };
const ARCH_MAP = { x64: "amd64", arm64: "arm64" };

const platform = PLATFORM_MAP[process.platform];
const arch = ARCH_MAP[process.arch];

if (!platform || !arch) {
  console.warn(`[goban] Unsupported platform: ${process.platform}/${process.arch}`);
  console.warn("[goban] Download the binary manually from https://github.com/icl00ud/goban/releases");
  process.exit(0);
}

const ext = platform === "windows" ? "zip" : "tar.gz";
const filename = `goban_${VERSION}_${platform}_${arch}.${ext}`;
const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${filename}`;

fs.mkdirSync(BIN_DIR, { recursive: true });

console.log(`[goban] Downloading ${filename}...`);

function follow(url, cb) {
  https.get(url, (res) => {
    if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
      follow(res.headers.location, cb);
    } else if (res.statusCode !== 200) {
      cb(new Error(`HTTP ${res.statusCode} for ${url}`));
    } else {
      cb(null, res);
    }
  }).on("error", cb);
}

follow(url, (err, res) => {
  if (err) {
    console.error("[goban] Download failed:", err.message);
    console.error(`[goban] Try manually: ${url}`);
    process.exit(1);
  }

  if (platform === "windows") {
    // Buffer the zip and extract via PowerShell
    const chunks = [];
    res.on("data", (c) => chunks.push(c));
    res.on("end", () => {
      const zipPath = path.join(BIN_DIR, "goban.zip");
      fs.writeFileSync(zipPath, Buffer.concat(chunks));
      const { execFileSync } = require("child_process");
      try {
        execFileSync("powershell", [
          "-Command",
          `Expand-Archive -Force '${zipPath}' '${BIN_DIR}'`,
        ]);
        fs.unlinkSync(zipPath);
        console.log("[goban] Installed successfully.");
      } catch (e) {
        console.error("[goban] Extraction failed:", e.message);
        process.exit(1);
      }
    });
  } else {
    // Stream: gzip decompress → tar extract a single file
    extractTarGz(res, BIN_PATH, () => {
      fs.chmodSync(BIN_PATH, 0o755);
      console.log("[goban] Installed successfully.");
    });
  }
});

function extractTarGz(stream, destFile, done) {
  // Minimal tar.gz extractor: find the "goban" file entry and write it.
  const gunzip = zlib.createGunzip();
  stream.pipe(gunzip);

  let buf = Buffer.alloc(0);
  gunzip.on("data", (chunk) => { buf = Buffer.concat([buf, chunk]); });
  gunzip.on("end", () => {
    // Parse POSIX tar: each block is 512 bytes, header + data blocks
    let offset = 0;
    while (offset + 512 <= buf.length) {
      const header = buf.slice(offset, offset + 512);
      const name = header.slice(0, 100).toString("utf8").replace(/\0/g, "");
      if (!name) break; // end-of-archive

      const sizeStr = header.slice(124, 136).toString("utf8").replace(/\0/g, "").trim();
      const size = parseInt(sizeStr, 8) || 0;
      offset += 512;

      const baseName = path.basename(name);
      if (baseName === "goban" || baseName === "goban.exe") {
        fs.writeFileSync(destFile, buf.slice(offset, offset + size));
        done();
        return;
      }
      // Skip data blocks (padded to 512-byte boundary)
      offset += Math.ceil(size / 512) * 512;
    }
    console.error("[goban] Binary not found in archive.");
    process.exit(1);
  });
  gunzip.on("error", (e) => {
    console.error("[goban] Decompression error:", e.message);
    process.exit(1);
  });
}
