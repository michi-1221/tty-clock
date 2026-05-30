# tty-clock — guidance for Claude

A terminal clock in **Go + Bubble Tea** (digital now; analog/braille planned).
Module: `github.com/michi-1221/tty-clock`. Go 1.26, darwin/arm64.

## Documentation convention (IMPORTANT)

**`README.md` is the user-facing manual** — how to install, run, and configure
the clock. Keep it to what an *end user* needs: install, keys, the config
options / example / themes, and user-visible behavior. **Whenever a change adds
or alters user-visible behavior or config, update the relevant README section in
the same change — but add only user-necessary information.**

Keep developer/internal detail **out** of the README: architecture, rendering
invariants, the release/npm pipeline, and per-change "spec changelog" notes
belong here in CLAUDE.md and in the code, not in the README. If a change makes a
documented user behavior wrong, fix the README in the same edit.

## Commands

```bash
go build ./... && go vet ./... && go test ./...
go test ./internal/render/digital -update   # regenerate View golden files
gofmt -w <files>
go run .                                     # needs a real TTY (won't run headless)
```

## Invariants to preserve (don't regress)

- **`internal/render` must not import bubbletea.** It is a pure
  `RenderContext → string` seam so it can be golden-tested. **`internal/ui` is
  the only package allowed to import bubbletea.**
- **Tick generation:** `tickMsg` carries a `gen`; `Update` re-arms only when
  `gen == m.tickGen` and drops stale ticks. On a granularity change, bump
  `m.tickGen`. (bubbletea cannot cancel an in-flight `tea.Every`.)
- **Seed `now` in `New()`** (`time.Now()`) — the first `View()` is drawn before
  the first tick; without this, minute granularity shows a wrong time for up to
  60s.
- **Styling goes through `caps.Renderer`** (an explicit `*lipgloss.Renderer`),
  never the global lipgloss default — keeps golden tests deterministic. Tests
  force the profile via `r.SetColorProfile(...)`.
- **Block font draws digits and `:` only.** The date is a small text label,
  never giant glyphs. Colon-blink swaps `:` for a same-width blank.
- **Pure View:** never call `time.Now()` in `View`/render; read `m.now` /
  `TimeSnapshot`.
- **Enums live in `internal/clock`** to avoid `config`↔`render`↔`ui` import
  cycles. `config` stores strings and parses via `clock`.

## Status

Phase 1 complete (digital). Phase 2 complete (analog braille + `m`, live reload
`r`). Config lives at `~/.tty-clock/config.json` (scaffolded on first run).
Runtime toggles (`s`/`t`/`m`/`?`) are **persisted** back to that file via
`config.Save` (atomic). Later: `seg7` font, multi-color hands, fsnotify.
