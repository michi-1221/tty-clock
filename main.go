// Command tty-clock is a terminal clock built on Bubble Tea. main only parses
// flags, resolves/loads config, detects capabilities, and runs the program —
// all logic lives in internal/.
package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/michi-1221/tty-clock/internal/caps"
	"github.com/michi-1221/tty-clock/internal/config"
	"github.com/michi-1221/tty-clock/internal/ui"
)

// version is the build version, overwritten by the linker at release time
// (goreleaser: -X main.version={{.Version}}). Defaults to "dev" for local builds.
var version = "dev"

func main() {
	os.Exit(run())
}

func run() int {
	configFlag := flag.String("config", "", "path to config.json (default: XDG search)")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println("tty-clock", version)
		return 0
	}

	path, err := config.Resolve(*configFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tty-clock:", err)
		return 1
	}
	cfg, err := config.Load(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tty-clock:", err)
		return 1
	}

	model := ui.New(cfg, path, caps.Detect(os.Stdout))

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "tty-clock:", err)
		return 1
	}
	return 0
}
