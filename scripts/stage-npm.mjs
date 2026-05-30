// Stage npm packages from GoReleaser output.
//
// Reads dist/metadata.json (version) and dist/artifacts.json (built binaries),
// then writes a publishable tree under dist/npm/:
//
//   dist/npm/tty-clock/                  main package (bin shim + optionalDeps)
//   dist/npm/tty-clock-<platform>-<arch>/  one per binary, with os/cpu pins
//
// The main package's bin shim (npm/tty-clock/bin/tty-clock.js) resolves the
// platform package matching the host and execs the binary inside it. Because
// the platform packages are optionalDependencies pinned with os/cpu, npm only
// installs the one that matches — so `npx tty-clock` fetches a single binary.
//
// Run after `goreleaser release` / `goreleaser build`:  node scripts/stage-npm.mjs

import {
  readFileSync,
  writeFileSync,
  mkdirSync,
  copyFileSync,
  chmodSync,
  rmSync,
  cpSync,
} from "node:fs";
import { join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(fileURLToPath(import.meta.url), "../..");
const dist = join(root, "dist");
const outDir = join(dist, "npm");
const mainSrc = join(root, "npm", "tty-clock");

// GoReleaser GOOS/GOARCH -> Node process.platform/process.arch (so the shim's
// `@myaruran/tty-clock-${process.platform}-${process.arch}` lookup matches names).
const GOOS_TO_NODE = { darwin: "darwin", linux: "linux", windows: "win32" };
const GOARCH_TO_NODE = { amd64: "x64", arm64: "arm64" };

const repository = {
  type: "git",
  url: "git+https://github.com/michi-1221/tty-clock.git",
};

function readJSON(path) {
  return JSON.parse(readFileSync(path, "utf8"));
}

const { version } = readJSON(join(dist, "metadata.json"));
const artifacts = readJSON(join(dist, "artifacts.json"));
const binaries = artifacts.filter((a) => a.type === "Binary");

if (binaries.length === 0) {
  throw new Error("no Binary artifacts found in dist/artifacts.json");
}

rmSync(outDir, { recursive: true, force: true });
mkdirSync(outDir, { recursive: true });

const optionalDependencies = {};

for (const b of binaries) {
  const platform = GOOS_TO_NODE[b.goos];
  const arch = GOARCH_TO_NODE[b.goarch];
  if (!platform || !arch) {
    console.warn(`skip unmapped target ${b.goos}/${b.goarch}`);
    continue;
  }

  // Flat folder name (drives the `tty-clock-*` publish glob); scoped npm name
  // (@myaruran/…) so per-platform packages don't trip npm's spam detection.
  const dirName = `tty-clock-${platform}-${arch}`;
  const pkgName = `@myaruran/${dirName}`;
  const ext = platform === "win32" ? ".exe" : "";
  const binDir = join(outDir, dirName, "bin");
  mkdirSync(binDir, { recursive: true });

  const dest = join(binDir, `tty-clock${ext}`);
  copyFileSync(join(root, b.path), dest);
  chmodSync(dest, 0o755); // npm preserves file mode in the tarball

  writeFileSync(
    join(outDir, dirName, "package.json"),
    JSON.stringify(
      {
        name: pkgName,
        version,
        description: `Prebuilt tty-clock binary for ${platform}-${arch}.`,
        repository,
        license: "MIT",
        os: [platform],
        cpu: [arch],
        files: ["bin/"],
      },
      null,
      2,
    ) + "\n",
  );

  optionalDependencies[pkgName] = version;
}

// Main package: copy the committed template, then inject the version and the
// optionalDependencies map derived from the binaries we actually staged.
const mainDest = join(outDir, "tty-clock");
cpSync(mainSrc, mainDest, { recursive: true });

const mainPkgPath = join(mainDest, "package.json");
const mainPkg = readJSON(mainPkgPath);
mainPkg.version = version;
mainPkg.optionalDependencies = Object.fromEntries(
  Object.entries(optionalDependencies).sort(([a], [b]) => a.localeCompare(b)),
);
writeFileSync(mainPkgPath, JSON.stringify(mainPkg, null, 2) + "\n");

console.log(`Staged npm packages for v${version} under dist/npm/:`);
console.log(`  main  tty-clock`);
for (const name of Object.keys(mainPkg.optionalDependencies)) {
  console.log(`  dep   ${name}`);
}
