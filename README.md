# tty-clock

A configurable **digital & analog clock** for your terminal — themeable
(including gradient themes), keyboard-controlled, and easy on the eyes. Runs on
macOS, Linux, and Windows.

## Install

Run it instantly with npm — no toolchain required:

```bash
npx tty-clock@latest              # run it right now
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
| `e` | open the config file in your editor |
| `q` · `ctrl+c` · `esc` | quit |

What you toggle is **saved to your config file** right away, so the clock starts
the same way next time.

## Configuration

tty-clock keeps its settings in **`~/.tty-clock/config.json`**. The file is
created for you (filled with the defaults) the first time you run the clock, so
there's always something to edit. Press **`e`** to open it in your editor
(`$VISUAL` / `$EDITOR`, falling back to your OS default app) — when you save and
close, the clock reloads automatically. You can also edit it by hand and press
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

---

# tty-clock（日本語）

ターミナル向けの設定可能な **デジタル & アナログ時計** — テーマ対応
（グラデーションテーマを含む）、キーボード操作、目にやさしい表示。macOS、Linux、
Windows で動作します。

## インストール

ツールチェーン不要、npm ですぐに実行できます:

```bash
npx tty-clock              # 今すぐ実行する
npm install -g tty-clock   # または `tty-clock` コマンドをインストールする
```

macOS、Linux、Windows（x64 / arm64）向けにビルド済みです。

Go ツールチェーンを使いたい場合は、ソースからインストールできます:

```bash
go install github.com/michi-1221/tty-clock@latest   # `tty-clock` をインストール
go run github.com/michi-1221/tty-clock@latest         # または一度だけ実行（インストールなし）
```

`tty-clock --version` でバージョンを表示します。

## キー操作

| キー | 動作 |
| --- | ------ |
| `s` | 秒の表示 / 非表示 |
| `t` | テーマを切り替える |
| `m` | デジタル ⇄ アナログを切り替える |
| `?` | ヘルプ行の表示 / 非表示 |
| `r` | 設定ファイルを再読み込みする |
| `e` | 設定ファイルをエディタで開く |
| `q` · `ctrl+c` · `esc` | 終了 |

切り替えた内容はすぐに **設定ファイルへ保存** されるので、次回起動時も同じ状態で
始まります。

## 設定

tty-clock は設定を **`~/.tty-clock/config.json`** に保存します。このファイルは
初回起動時に（デフォルト値が入った状態で）自動的に作成されるため、常に編集できる
ものが用意されています。**`e`** を押すとエディタ（`$VISUAL` / `$EDITOR`、なければ
OS のデフォルトアプリ）で開きます — 保存して閉じると、時計が自動的に再読み込み
します。手動で編集して **`r`** を押せば、その場で再読み込みすることもできます。

> 別のファイルを使いたい場合は `tty-clock --config /path/to/config.json` を実行します。

### オプション

**トップレベル**

| オプション      | デフォルト       | 内容                                                                         |
| ------------- | --------------- | --------------------------------------------------------------------------- |
| `mode`        | `"digital"`     | `"digital"` または `"analog"` の時計表示（`m` でライブ切り替えも可能）          |
| `theme`       | `"tokyo-night"` | カラーテーマ — [テーマ](#テーマ) を参照                                          |
| `customTheme` | `null`          | 個々のテーマ色を上書き（下記参照）                                              |
| `granularity` | `"seconds"`     | 更新頻度: `"seconds"` または `"minutes"`                                       |
| `showHelp`    | `true`          | 下部のヘルプ行を表示する（`?` で切り替え）                                       |
| `cellAspect`  | `0`             | アナログ専用: 形状補正。`0` = 自動。ダイヤルが楕円に見える場合は例えば `2.4` に設定 |
| `format`      | —               | 表示オプション（下記）                                                          |

**`format`**

| オプション      | デフォルト          | 内容                                                                |
| ------------- | ------------------ | ------------------------------------------------------------------ |
| `hour24`      | `true`             | 24 時間表示。`false` で 12 時間表示                                  |
| `showSeconds` | `true`             | 秒を表示する（`granularity` が `"seconds"` のときのみ）              |
| `showDate`    | `true`             | 時計の下に日付行を表示する                                          |
| `dateFormat`  | `"Mon 2006-01-02"` | 日付のレイアウト（[Go の時刻フォーマット](https://pkg.go.dev/time#pkg-constants)） |
| `blinkColon`  | `false`            | `:` を 1 秒ごとに点滅させる                                          |
| `font`        | `"block"`          | 数字のフォント: `"block"`（大きい）または `"ascii"`（シンプルな代替） |
| `showNumbers` | `true`             | アナログ専用: 文字盤に 1〜12 の時刻数字を描画する                     |

**`customTheme`** — 選択したテーマの色を上書きします。変更したいものだけを設定
すれば、省略したものはテーマの値を保ちます。色は `#rgb` または `#rrggbb` です。

| 色           | 用途                                                                  |
| ------------ | -------------------------------------------------------------------- |
| `primary`    | 数字 / 時計の針                                                       |
| `accent`     | コロン、秒、ハイライト                                                 |
| `secondary`  | 日付                                                                  |
| `muted`      | 文字盤 / 目盛り                                                       |
| `background` | 任意の背景塗りつぶし（デフォルトはオフ — ターミナルの背景が透けて見える） |
| `gradient`   | 数字 / ダイヤルをカラーランプで塗る（下記参照）                          |

**グラデーション** は、フラットな `primary` 色の代わりに、巨大な数字とアナログ
ダイヤルを複数ストップのカラーランプで塗ります（日付は `secondary` のまま）:

```json
"customTheme": {
  "gradient": {
    "stops": ["#ffe27a", "#ff9e5e", "#ff5e62", "#a23bc0"],
    "direction": "vertical"
  }
}
```

`stops` には少なくとも 2 色が必要です。`direction` は `"vertical"`（上→下、
デフォルト）または `"horizontal"` です。グラデーションはトゥルーカラー対応の
ターミナルで最もきれいに表示され、256 色ターミナルでは自動的にダウンサンプル
されます。

### 設定例 `~/.tty-clock/config.json`

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

### テーマ

`tokyo-night` がデフォルトです。`t` を押すと全 10 種類（グラデーションテーマを
含む）を切り替えられます。

| テーマ             | primary   | accent    | secondary | muted     | background |
| ------------------ | --------- | --------- | --------- | --------- | ---------- |
| `tokyo-night` ★    | `#c0caf5` | `#7aa2f7` | `#bb9af7` | `#565f89` | `#1a1b26`  |
| `dracula`          | `#f8f8f2` | `#ff79c6` | `#bd93f9` | `#6272a4` | `#282a36`  |
| `nord`             | `#d8dee9` | `#88c0d0` | `#81a1c1` | `#4c566a` | `#2e3440`  |
| `gruvbox`          | `#ebdbb2` | `#fe8019` | `#fabd2f` | `#928374` | `#282828`  |
| `catppuccin-mocha` | `#cdd6f4` | `#cba6f7` | `#fab387` | `#6c7086` | `#1e1e2e`  |
| `solarized-dark`   | `#839496` | `#b58900` | `#2aa198` | `#586e75` | `#002b36`  |
| `monochrome`       | `#ffffff` | `#ffffff` | `#bbbbbb` | `#666666` | `#000000`  |

**グラデーションテーマ** — 数字 / ダイヤルがカラーランプで塗られます:

| テーマ    | 方向       | 色                                            |
| --------- | ---------- | --------------------------------------------- |
| `sunset`  | vertical   | 金 → オレンジ → 赤 → 紫                         |
| `aurora`  | vertical   | 緑 → シアン → すみれ色                          |
| `rainbow` | horizontal | 赤 → オレンジ → 黄 → 緑 → 青 → 紫               |

### 知っておくと便利なこと

- **色** — tty-clock は `NO_COLOR` を尊重し、色をサポートしないターミナルでは色を
  無効にします。UTF-8 でないターミナル（または `font: "ascii"`）では、シンプルな
  ASCII の数字フォントを使います。
- **アナログモード**（`m`）— 丸いブライユ点字のダイヤルです。ウィンドウが小さ
  すぎる場合や、ターミナルが UTF-8 でない場合は、自動的にデジタル時計を表示し、
  スペースができると元に戻ります。
- **12 時間モード** — ゼロ埋めで表示され（例: `03:04:05`）、AM/PM のラベルは
  付きません。
- **タイムゾーン** — システムの時計とローカルのタイムゾーンを使います。`TZ=…` で
  上書きできます。

---

[Go](https://go.dev) と [Bubble Tea](https://github.com/charmbracelet/bubbletea) で作られています。
