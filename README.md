# tty-clock

A terminal clock built with **Go + [Bubble Tea](https://github.com/charmbracelet/bubbletea)**.
Configuration-driven, theme-aware, keyboard-controlled. Aims for the "always nice to
look at" feel of k9s / lazygit.

> **This README is the spec.** The JSON configuration schema and the project
> spec below are kept in sync with the implementation — every feature that lands
> is documented here.

## Status

- **Phase 1 (done):** configurable, themeable **digital** clock — giant block
  digits, ASCII fallback, 7 themes, keyboard controls, JSON config.
- **Phase 2 (planned):** analog braille renderer (`m` toggle), live reload
  (`r`), `seg7` font, multi-color hands.

## Run

```bash
go run ./cmd/tty-clock                                   # XDG config or defaults
go run ./cmd/tty-clock --config path/to/config.json      # explicit config
go build -o tty-clock ./cmd/tty-clock                    # build a binary
```

## Keybindings

| Key                 | Action                                   |
| ------------------- | ---------------------------------------- |
| `s`                 | show / hide seconds                      |
| `t`                 | cycle theme (clean presets only)         |
| `?`                 | toggle the help line (hidden → clock fills the screen) |
| `q` / `ctrl+c` / `esc` | quit                                  |
| `m`                 | mode digital ⇄ analog *(phase 2; no-op now)* |
| `r`                 | reload config *(phase 2)*                |

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
| `mode`        | string        | `"digital"`    | `"digital"`, `"analog"` (analog is phase 2)       |
| `theme`       | string        | `"tokyo-night"`| a preset name (see **Themes**)                    |
| `customTheme` | object \| null| `null`         | partial palette override (see below)              |
| `granularity` | string        | `"seconds"`    | `"seconds"`, `"minutes"` — update frequency       |
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
| `font`        | string | `"block"`          | `"block"`, `"ascii"` (manual fallback), `"seg7"` (phase 2) |

**`customTheme`** (all fields optional; only non-empty fields override the
preset, so partial overrides are allowed). Colors are `#rgb` or `#rrggbb`.

| Field        | Role                                            |
| ------------ | ----------------------------------------------- |
| `primary`    | digits / hands                                  |
| `accent`     | colon / seconds / AM-PM                         |
| `secondary`  | date                                            |
| `muted`      | clock face / ticks (phase-2 analog)             |
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
- **Window scaling:** the giant digits scale with the window in **integer
  steps**, anchored so a **54×10** window renders the original size (×1):
  `scale = max(1, floor(min(width/54, height/10)))` → 108×20 = ×2, 162×30 = ×3.
  The font doubles only when *both* dimensions reach 2×; the date and AM/PM
  labels stay normal-size text.

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
cmd/tty-clock/main.go   flags · config resolve · caps detect · run (no logic)
internal/
  clock/   Mode/Granularity enums + TimeSnapshot  (framework-independent domain)
  config/  schema · defaults · Load (XDG + --config) · Validate · friendly errors
  theme/   Palette/Theme · 7 presets · Resolve (merge) · Next (cycle)
  caps/    explicit *lipgloss.Renderer + UTF-8/NO_COLOR heuristics
  render/  Renderer interface + RenderContext  (the tea-free rendering seam)
    digital/  giant block digits · ASCII fallback · font/glyphs
  ui/      Model + Init/Update/View · keys · tick · layout  (only package importing bubbletea)
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

## Verification

```bash
go build ./... && go vet ./... && go test ./...
go test ./internal/render/digital -update   # regenerate View golden files
```

Interactive behavior needs a real TTY, so run `go run ./cmd/tty-clock` in your
terminal to confirm live updates and key handling.

## Spec changelog

- **Phase 1** — digital clock: JSON config (mode/theme/customTheme/granularity/
  format), 7 themes + partial override, block & ASCII fonts, seconds/blink/date,
  seconds vs minutes granularity, keys `s`/`t`/`?`/`q`. Help line toggles with
  `?` (hidden = clock fills the screen).
- **Phase 1.1** — window-proportional font scaling (integer steps, anchored at
  54×10 = ×1; `floor(min(w/54, h/10))`). Digits scale; date/AM-PM stay text.
