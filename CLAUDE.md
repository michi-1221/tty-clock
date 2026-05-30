# tty-clock — guidance for Claude

A terminal clock in **Go + Bubble Tea** (digital now; analog/braille planned).
Module: `github.com/michi-1221/tty-clock`. Go 1.26, darwin/arm64.

## Documentation convention (IMPORTANT)

**`README.md` is the spec, and it must be kept in sync in the same change as the
code.** Whenever you implement or change a feature, behavior, or the config:

- Update the **JSON configuration schema** in README (field tables: type,
  default, allowed values, notes) — add rows/sections for anything new.
- Update the **project spec / architecture** section to match reality.
- Append a line to the **Spec changelog** section describing what landed.

Do not let README drift from the implementation. If a change makes a documented
behavior wrong, fix the docs in the same edit. Treat README as authoritative for
the config format and user-visible behavior.

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
- **Block font draws digits and `:` only.** AM/PM and the date are small text
  labels, never giant glyphs. Colon-blink swaps `:` for a same-width blank.
- **Pure View:** never call `time.Now()` in `View`/render; read `m.now` /
  `TimeSnapshot`.
- **Enums live in `internal/clock`** to avoid `config`↔`render`↔`ui` import
  cycles. `config` stores strings and parses via `clock`.

## Status

Phase 1 complete (digital). Phase 2 complete (analog braille + `m`, live reload
`r`). Config lives at `~/.tty-clock/config.json` (scaffolded on first run).
Runtime toggles (`s`/`t`/`m`/`?`) are **persisted** back to that file via
`config.Save` (atomic). Later: `seg7` font, multi-color hands, fsnotify.
