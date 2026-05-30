# tty-clock

A terminal clock built with **Go + [Bubble Tea](https://github.com/charmbracelet/bubbletea)**.
Configuration-driven, theme-aware, keyboard-controlled. Aims for the "always nice to
look at" feel of k9s / lazygit.

> **This README is the spec.** The JSON configuration schema and the project
> spec below are kept in sync with the implementation — every feature that lands
> is documented here.

## Status

- **Phase 1 (done):** configurable, themeable **digital** clock — giant block
  digits, ASCII fallback, 7 themes, keyboard controls, JSON config, window
  scaling.
- **Phase 2 (done):** **analog** braille dial (`m` toggle, auto-falls back to
  digital when too small or on a non-UTF-8 terminal) and live config reload
  (`r`).
- **Later:** `seg7` font, multi-color hands, fsnotify auto-reload.

## Install

**No toolchain (npm / npx).** A prebuilt binary for your platform is published
to npm; `npx` fetches and runs it:

```bash
npx tty-clock            # run without installing
npm install -g tty-clock # or install the command globally
```

The `tty-clock` npm package is a thin launcher: the binaries ship in
per-platform packages (`tty-clock-darwin-arm64`, `tty-clock-linux-x64`, …)
declared as `optionalDependencies` with `os`/`cpu` pins, so npm downloads only
the one matching your machine. Prebuilt for macOS / Linux / Windows on
x64 / arm64.

**With the Go toolchain.** Run it without installing anything:

```bash
go run github.com/michi-1221/tty-clock@latest
```

Or install the `tty-clock` binary into `$(go env GOPATH)/bin`:

```bash
go install github.com/michi-1221/tty-clock@latest
tty-clock
```

## Run from source

```bash
go run .                                  # XDG config or defaults
go run . --config path/to/config.json     # explicit config
go run . --version                        # print version and exit
go build -o tty-clock .                   # build a binary
```

Flags: `--config <path>` (see resolution order below), `--version` (prints
`tty-clock <version>`; the version is stamped in at release time via the
linker, and reads `dev` for local builds).

## Keybindings

| Key                 | Action                                   |
| ------------------- | ---------------------------------------- |
| `s`                 | show / hide seconds                      |
| `t`                 | cycle theme (clean presets only)         |
| `?`                 | toggle the help line (hidden → clock fills the screen) |
| `q` / `ctrl+c` / `esc` | quit                                  |
| `m`                 | switch mode: digital ⇄ analog            |
| `r`                 | reload the config file                   |

Runtime toggles (`s`, `t`) are **ephemeral** — they are *not* written back to
the config file. Restart or a future reload returns to the file's values.

## Configuration

### Resolution order

1. `--config <path>` — explicit. The file **must exist** (missing is a fatal error).
2. Otherwise — XDG search for `tty-clock/config.json` (e.g.
   `$XDG_CONFIG_HOME/tty-clock/config.json`, or `~/.config/tty-clock/config.json`).
   If not found, **built-in defaults** are used silently.

Loading rules: defaults are the baseline; the file is overlaid on top, so any
omitted key (including nested `format` keys) keeps its default. Unknown keys are
rejected. Syntax and type errors report `line:col`.

### JSON schema

**Top level**

| Field         | Type          | Default        | Allowed / notes                                   |
| ------------- | ------------- | -------------- | ------------------------------------------------- |
| `mode`        | string        | `"digital"`    | `"digital"`, `"analog"` (also toggled live with `m`) |
| `theme`       | string        | `"tokyo-night"`| a preset name (see **Themes**)                    |
| `customTheme` | object \| null| `null`         | partial palette override (see below)              |
| `granularity` | string        | `"seconds"`    | `"seconds"`, `"minutes"` — update frequency       |
| `cellAspect`  | number        | `0`            | analog dial cell height/width; `0` = auto-detect (fallback 2.0) |
| `format`      | object        | (see below)    | display options                                   |

**`format`**

| Field         | Type   | Default            | Notes                                                  |
| ------------- | ------ | ------------------ | ------------------------------------------------------ |
| `hour24`      | bool   | `true`             | 24-hour vs 12-hour                                     |
| `showSeconds` | bool   | `true`             | effective only at `granularity: "seconds"`             |
| `showDate`    | bool   | `true`             | date line below the clock                              |
| `dateFormat`  | string | `"Mon 2006-01-02"` | Go time layout                                         |
| `blinkColon`  | bool   | `false`            | seconds granularity only                               |
| `showAMPM`    | bool   | `true`             | only meaningful when `hour24` is `false`               |
| `font`        | string | `"block"`          | `"block"`, `"ascii"` (manual fallback), `"seg7"` (planned) |
| `showNumbers` | bool   | `true`             | analog: draw the hour numbers 1–12 on the face          |

**`customTheme`** (all fields optional; only non-empty fields override the
preset, so partial overrides are allowed). Colors are `#rgb` or `#rrggbb`.

| Field        | Role                                            |
| ------------ | ----------------------------------------------- |
| `primary`    | digits / hands                                  |
| `accent`     | colon / seconds / AM-PM                         |
| `secondary`  | date                                            |
| `muted`      | clock face / ticks (reserved; analog v1 is single-color) |
| `background` | optional surface — **not painted by default** (the terminal background is respected) |

### Behavior notes

- **Granularity vs seconds:** `granularity` is the *update frequency*;
  `showSeconds` is *display*. At `minutes`, seconds are auto-hidden and colon
  blink is disabled (so nothing looks frozen).
- **12-hour:** always zero-padded (`03:04:05`) so digit width never jitters.
  `AM`/`PM` is a small text label, never a giant glyph.
- **Degradation:** `NO_COLOR`, non-TTY, and `TERM=dumb` drop color; non-UTF-8
  locale / `TERM=dumb` (or `font: "ascii"`) switch to the ASCII fallback font.
- **Time source:** the system clock (`time.Now()`), rendered in the local
  timezone — no timezone code is needed. `TZ=…` is honored by Go.
- **Window scaling (digital):** the giant digits scale with the window in
  **integer steps**, anchored so a **54×10** window renders the original size
  (×1): `scale = max(1, floor(min(width/54, height/10)))` → 108×20 = ×2,
  162×30 = ×3. The font doubles only when *both* dimensions reach 2×; the date
  and AM/PM labels stay normal-size text.
- **Analog mode:** a braille dial (face, 12 ticks, hour numbers 1–12,
  hour/minute/second hands) drawn entirely in braille dots (2×4 per cell — the
  finest character resolution), so the numbers are pixelated to match the dial.
  Single-color, sized to fill the window. It **auto-falls back to digital** when
  the area is below its minimum or the terminal isn't UTF-8 — `mode` is
  preserved, so the dial returns when the window grows. Toggle with `m`; hide
  the numbers with `format.showNumbers: false`. The second hand follows
  `showSeconds`, so `s` hides/shows it (and it is absent at minute granularity).
- **Round dial:** the cell height/width ratio is auto-detected from the
  terminal's pixel size (`TIOCGWINSZ`) so the dial looks circular regardless of
  font; set `cellAspect` to override (e.g. `2.4`) if your terminal doesn't
  report pixels and the dial looks oval.
- **Live reload:** `r` re-reads the config file; the file is the source of
  truth, so runtime `s`/`t`/`m` toggles are discarded. A bad file is reported
  in the footer and never stops the clock.

### Example `config.json`

```json
{
  "mode": "digital",
  "theme": "tokyo-night",
  "customTheme": { "accent": "#ff79c6" },
  "granularity": "seconds",
  "format": {
    "hour24": false,
    "showSeconds": true,
    "showDate": true,
    "dateFormat": "Mon 02 Jan 2006",
    "blinkColon": true,
    "showAMPM": true,
    "font": "block"
  }
}
```

### Themes

Default is **tokyo-night**. `t` cycles through them in this order.

| Theme              | primary   | accent    | secondary | muted     | background |
| ------------------ | --------- | --------- | --------- | --------- | ---------- |
| `tokyo-night` ★    | `#c0caf5` | `#7aa2f7` | `#bb9af7` | `#565f89` | `#1a1b26`  |
| `dracula`          | `#f8f8f2` | `#ff79c6` | `#bd93f9` | `#6272a4` | `#282a36`  |
| `nord`             | `#d8dee9` | `#88c0d0` | `#81a1c1` | `#4c566a` | `#2e3440`  |
| `gruvbox`          | `#ebdbb2` | `#fe8019` | `#fabd2f` | `#928374` | `#282828`  |
| `catppuccin-mocha` | `#cdd6f4` | `#cba6f7` | `#fab387` | `#6c7086` | `#1e1e2e`  |
| `solarized-dark`   | `#839496` | `#b58900` | `#2aa198` | `#586e75` | `#002b36`  |
| `monochrome`       | `#ffffff` | `#ffffff` | `#bbbbbb` | `#666666` | `#000000`  |

## Architecture (project spec)

```
main.go                 flags (--config, --version) · config resolve · caps detect · run (no logic)
.goreleaser.yaml        cross-compile matrix (6 targets) · archives · GitHub Release
.github/workflows/release.yml  tag push → GoReleaser → stage + npm publish
scripts/stage-npm.mjs   dist/ binaries → npm packages (main + per-platform)
npm/tty-clock/          main npm package: bin launcher + optionalDependencies
internal/
  clock/   Mode/Granularity enums + TimeSnapshot  (framework-independent domain)
  config/  schema · defaults · Load (XDG + --config) · Validate · friendly errors
  theme/   Palette/Theme · 7 presets · Resolve (merge) · Next (cycle)
  caps/    explicit *lipgloss.Renderer + UTF-8/NO_COLOR heuristics
  render/  Renderer interface + RenderContext  (the tea-free rendering seam)
    digital/  giant block digits · ASCII fallback · font/glyphs
    analog/   braille dial · midpoint-circle face · Bresenham hands
  ui/      Model + Init/Update/View · keys · tick · layout · activeRenderer/reload  (only package importing bubbletea)
```

Key design points:

- **`render` never imports bubbletea** — it is a pure `RenderContext → string`,
  so output is golden-tested without starting a program. `ui` is the only
  package that imports bubbletea.
- **Tick:** `tea.Every` self-aligns to the wall-clock boundary (no drift,
  DST/TZ self-correct). `tickMsg` carries a **generation**; stale generations
  are dropped (bubbletea cannot cancel an in-flight timer). `now` is seeded in
  `New()` so the very first frame is correct.
- **Color:** styling goes through `caps.Renderer` (an explicit renderer), never
  the global lipgloss default — this keeps golden tests deterministic.

## Release & distribution

Releases are tag-driven and fully automated by
`.github/workflows/release.yml`:

1. **Tag** — push `vX.Y.Z`:
   ```bash
   git tag v0.1.0 && git push origin v0.1.0
   ```
2. **Build** — [GoReleaser](https://goreleaser.com) cross-compiles all six
   targets on a single Linux runner (pure Go, `CGO_ENABLED=0`), stamps
   `main.version` via `-ldflags -X`, writes archives + `checksums.txt`, and
   creates the **GitHub Release** with `dist/{artifacts,metadata}.json`.
3. **Stage** — `scripts/stage-npm.mjs` reads those JSON manifests and lays out
   `dist/npm/`: a per-platform package per binary (`tty-clock-<platform>-<arch>`
   with `os`/`cpu` pins, binary `chmod 755`) and the main `tty-clock` package
   (version + `optionalDependencies` injected).
4. **Publish** — platform packages publish first, then `tty-clock`, all with
   `npm publish --access public`, so `npx tty-clock` resolves a single matching
   binary. (Provenance is omitted because npm provenance requires a public
   repo; add `--provenance` + `id-token: write` back if this repo goes public.)

| Target (GOOS/GOARCH) | npm platform package      |
| -------------------- | ------------------------- |
| darwin/arm64         | `tty-clock-darwin-arm64`  |
| darwin/amd64         | `tty-clock-darwin-x64`    |
| linux/arm64          | `tty-clock-linux-arm64`   |
| linux/amd64          | `tty-clock-linux-x64`     |
| windows/arm64        | `tty-clock-win32-arm64`   |
| windows/amd64        | `tty-clock-win32-x64`     |

**Prerequisite:** an npm automation token stored as the `NPM_TOKEN` repository
secret (GitHub → Settings → Secrets and variables → Actions). `GITHUB_TOKEN` is
provided automatically. The GOOS/GOARCH → Node `process.platform`/`process.arch`
naming (`windows`→`win32`, `amd64`→`x64`) is what lets the launcher resolve
`tty-clock-${process.platform}-${process.arch}` at runtime.

## Verification

```bash
go build ./... && go vet ./... && go test ./...
go test ./internal/render/... -update       # regenerate digital + analog golden files
```

Interactive behavior needs a real TTY, so run `go run .` in your
terminal to confirm live updates and key handling.

## Spec changelog

- **Phase 1** — digital clock: JSON config (mode/theme/customTheme/granularity/
  format), 7 themes + partial override, block & ASCII fonts, seconds/blink/date,
  seconds vs minutes granularity, keys `s`/`t`/`?`/`q`. Help line toggles with
  `?` (hidden = clock fills the screen).
- **Phase 1.1** — window-proportional font scaling (integer steps, anchored at
  54×10 = ×1; `floor(min(w/54, h/10))`). Digits scale; date/AM-PM stay text.
- **Phase 2** — analog braille dial (`internal/render/analog`: integer
  midpoint-circle face, 12 ticks, Bresenham hands; single-color; sizes to the
  window). `m` toggles digital ⇄ analog, with auto-fallback to digital when the
  area is too small or the terminal isn't UTF-8 (mode preserved). `r` reloads
  the config file (file wins over runtime toggles; tick re-armed only when
  granularity changes; failures shown in the footer). `m`/`r` added to the help.
- **Phase 2.1** — analog dial made visually round (cell-aspect auto-detected via
  `TIOCGWINSZ`, parametric ellipse; `cellAspect` config override) and hour
  numbers 1–12 drawn on the face (`format.showNumbers`, default on).
- **Phase 2.2** — analog hour numbers drawn as braille dots (3×5 dot-matrix
  font) so they match the dial's pixel style; the analog second hand now honors
  `showSeconds` (hidden by `s` and at minute granularity).
- **Release pipeline** — tag-driven GoReleaser cross-compile (6 targets, pure
  Go) → GitHub Release; npm distribution via a main `tty-clock` launcher +
  per-platform `optionalDependencies` packages staged from `dist/` by
  `scripts/stage-npm.mjs`, so `npx tty-clock` runs a prebuilt binary with no Go
  toolchain. Added `main.version` + `--version` flag (linker-stamped).
