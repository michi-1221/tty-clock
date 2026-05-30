# tty-clock

A configurable, themeable **digital + analog** terminal clock (TUI), built with
[Go](https://go.dev) + [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Run it

```bash
npx tty-clock
```

This downloads a single prebuilt binary for your platform and runs it — no Go
toolchain required. Needs a real terminal (it's a full-screen TUI).

Install it globally instead:

```bash
npm install -g tty-clock
tty-clock
```

## How it works

`tty-clock` is a small Node launcher. The real binaries ship in per-platform
packages (`tty-clock-darwin-arm64`, `tty-clock-linux-x64`, …) declared as
`optionalDependencies` with `os`/`cpu` constraints, so npm fetches only the one
matching your machine. The launcher resolves that binary and execs it.

Supported: macOS, Linux, Windows on x64 and arm64.

## Configuration & keys

See the full docs and JSON config schema at
**https://github.com/michi-1221/tty-clock**.

## License

MIT
