# tty-clock

A configurable **digital & analog clock** for your terminal — themeable
(including gradient themes), keyboard-controlled, and easy on the eyes. Runs on
macOS, Linux, and Windows.

## Install

Run it instantly with npm — no toolchain required:

```bash
npx tty-clock              # run it right now
npm install -g tty-clock   # or install the `tty-clock` command
```

Prebuilt for macOS, Linux, and Windows (x64 / arm64).

Prefer the Go toolchain? Install from source instead:

```bash
go install github.com/michi-1221/tty-clock@latest   # installs `tty-clock`
go run github.com/michi-1221/tty-clock@latest         # or run once, no install
```

`tty-clock --version` prints the version.

## Keys

| Key | Action |
| --- | ------ |
| `s` | show / hide seconds |
| `t` | cycle through themes |
| `m` | switch digital ⇄ analog |
| `?` | show / hide the help line |
| `r` | reload the config file |
| `q` · `ctrl+c` · `esc` | quit |

What you toggle is **saved to your config file** right away, so the clock starts
the same way next time.

## Configuration

tty-clock keeps its settings in **`~/.tty-clock/config.json`**. The file is
created for you (filled with the defaults) the first time you run the clock, so
there's always something to edit. After editing, restart the clock or press
**`r`** to reload it live.

> Want a different file? Run `tty-clock --config /path/to/config.json`.

### Options

**Top level**

| Option        | Default         | What it does                                                                 |
| ------------- | --------------- | --------------------------------------------------------------------------- |
| `mode`        | `"digital"`     | `"digital"` or `"analog"` clock face (also switch live with `m`)            |
| `theme`       | `"tokyo-night"` | colour theme — see [Themes](#themes)                                         |
| `customTheme` | `null`          | override individual theme colours (see below)                               |
| `granularity` | `"seconds"`     | how often it updates: `"seconds"` or `"minutes"`                            |
| `showHelp`    | `true`          | show the help line at the bottom (toggle with `?`)                          |
| `cellAspect`  | `0`             | analog only: shape correction; `0` = auto. Set e.g. `2.4` if the dial looks oval |
| `format`      | —               | display options (below)                                                     |

**`format`**

| Option        | Default            | What it does                                                        |
| ------------- | ------------------ | ------------------------------------------------------------------ |
| `hour24`      | `true`             | 24-hour clock; `false` for 12-hour                                 |
| `showSeconds` | `true`             | show seconds (only when `granularity` is `"seconds"`)              |
| `showDate`    | `true`             | show the date line under the clock                                 |
| `dateFormat`  | `"Mon 2006-01-02"` | date layout ([Go time format](https://pkg.go.dev/time#pkg-constants)) |
| `blinkColon`  | `false`            | blink the `:` once a second                                        |
| `font`        | `"block"`          | digit font: `"block"` (big) or `"ascii"` (simple fallback)         |
| `showNumbers` | `true`             | analog only: draw the hour numbers 1–12 on the face                |

**`customTheme`** — override the chosen theme's colours. Set only what you want
to change; anything you leave out keeps the theme's value. Colours are `#rgb`
or `#rrggbb`.

| Colour       | Used for                                                              |
| ------------ | -------------------------------------------------------------------- |
| `primary`    | the digits / clock hands                                             |
| `accent`     | colon, seconds, highlights                                           |
| `secondary`  | the date                                                             |
| `muted`      | clock face / ticks                                                   |
| `background` | optional background fill (off by default — your terminal shows through) |
| `gradient`   | paint the digits/dial with a colour ramp (see below)                |

A **gradient** paints the giant digits and the analog dial with a multi-stop
colour ramp instead of the flat `primary` colour (the date stays `secondary`):

```json
"customTheme": {
  "gradient": {
    "stops": ["#ffe27a", "#ff9e5e", "#ff5e62", "#a23bc0"],
    "direction": "vertical"
  }
}
```

`stops` needs at least two colours; `direction` is `"vertical"` (top→bottom, the
default) or `"horizontal"`. Gradients look best in a true-colour terminal and
downsample automatically on 256-colour terminals.

### Example `~/.tty-clock/config.json`

```json
{
  "mode": "digital",
  "theme": "tokyo-night",
  "customTheme": { "accent": "#ff79c6" },
  "granularity": "seconds",
  "showHelp": true,
  "cellAspect": 0,
  "format": {
    "hour24": false,
    "showSeconds": true,
    "showDate": true,
    "dateFormat": "Mon 02 Jan 2006",
    "blinkColon": true,
    "font": "block",
    "showNumbers": true
  }
}
```

### Themes

`tokyo-night` is the default; press `t` to cycle through all 10 (gradient themes
included).

| Theme              | primary   | accent    | secondary | muted     | background |
| ------------------ | --------- | --------- | --------- | --------- | ---------- |
| `tokyo-night` ★    | `#c0caf5` | `#7aa2f7` | `#bb9af7` | `#565f89` | `#1a1b26`  |
| `dracula`          | `#f8f8f2` | `#ff79c6` | `#bd93f9` | `#6272a4` | `#282a36`  |
| `nord`             | `#d8dee9` | `#88c0d0` | `#81a1c1` | `#4c566a` | `#2e3440`  |
| `gruvbox`          | `#ebdbb2` | `#fe8019` | `#fabd2f` | `#928374` | `#282828`  |
| `catppuccin-mocha` | `#cdd6f4` | `#cba6f7` | `#fab387` | `#6c7086` | `#1e1e2e`  |
| `solarized-dark`   | `#839496` | `#b58900` | `#2aa198` | `#586e75` | `#002b36`  |
| `monochrome`       | `#ffffff` | `#ffffff` | `#bbbbbb` | `#666666` | `#000000`  |

**Gradient themes** — the digits/dial are painted with a colour ramp:

| Theme     | direction  | colours                                       |
| --------- | ---------- | --------------------------------------------- |
| `sunset`  | vertical   | gold → orange → red → purple                  |
| `aurora`  | vertical   | green → cyan → violet                         |
| `rainbow` | horizontal | red → orange → yellow → green → blue → purple |

### Good to know

- **Colours** — tty-clock respects `NO_COLOR` and drops colour on terminals that
  don't support it. On non-UTF-8 terminals (or with `font: "ascii"`) it uses a
  simple ASCII digit font.
- **Analog mode** (`m`) — a round braille dial. If the window is too small, or
  the terminal isn't UTF-8, it automatically shows the digital clock instead and
  switches back once there's room.
- **12-hour mode** — shown zero-padded (e.g. `03:04:05`), without an AM/PM label.
- **Timezone** — uses your system clock and local timezone; set `TZ=…` to
  override.

---

Built with [Go](https://go.dev) and [Bubble Tea](https://github.com/charmbracelet/bubbletea).
