#!/usr/bin/env node
"use strict";

// Thin launcher for the tty-clock Go binary.
//
// The actual binaries ship in per-platform packages (tty-clock-darwin-arm64,
// tty-clock-linux-x64, ...) declared as optionalDependencies and pinned with
// os/cpu, so npm installs only the one matching this host. We resolve that
// package's binary and exec it with our argv/stdio inherited (so the TUI gets
// a real TTY and signals flow through).

const { spawnSync } = require("node:child_process");

function resolveBinary() {
  const platform = process.platform; // 'darwin' | 'linux' | 'win32' | ...
  const arch = process.arch; // 'arm64' | 'x64' | ...
  const ext = platform === "win32" ? ".exe" : "";
  const pkg = `tty-clock-${platform}-${arch}`;
  try {
    return require.resolve(`${pkg}/bin/tty-clock${ext}`);
  } catch {
    return null;
  }
}

const bin = resolveBinary();
if (!bin) {
  console.error(
    `tty-clock: no prebuilt binary for ${process.platform}-${process.arch}.\n` +
      "Supported: darwin/linux/win32 on x64/arm64.\n" +
      "Build from source instead:\n" +
      "  go install github.com/michi-1221/tty-clock@latest",
  );
  process.exit(1);
}

const result = spawnSync(bin, process.argv.slice(2), { stdio: "inherit" });

if (result.error) {
  console.error(`tty-clock: failed to launch binary: ${result.error.message}`);
  process.exit(1);
}

// Re-raise the terminating signal so callers see the right exit status; fall
// back to the numeric code (or 1 if neither is set).
if (result.signal) {
  process.kill(process.pid, result.signal);
}
process.exit(result.status === null ? 1 : result.status);
