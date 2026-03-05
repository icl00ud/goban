#!/usr/bin/env node
"use strict";

const { spawnSync } = require("child_process");
const path = require("path");
const fs = require("fs");

const bin = path.join(
  __dirname,
  "bin",
  process.platform === "win32" ? "goban.exe" : "goban"
);

if (!fs.existsSync(bin)) {
  console.error("[goban] Binary not found. Try reinstalling: npm install -g goban-cli");
  process.exit(1);
}

const result = spawnSync(bin, process.argv.slice(2), { stdio: "inherit" });
process.exit(result.status ?? 0);
