package splitty

import "github.com/charmbracelet/lipgloss"

// Theme defines the visual styling for split panes.
type Theme struct {
	BorderActive           lipgloss.Style
	BorderInactive         lipgloss.Style
	BorderScrollback       lipgloss.Style
	BorderScrollbackFocused lipgloss.Style
	Divider                lipgloss.Style
	DividerChar            string
	StatusBar              lipgloss.Style
	StatusText             lipgloss.Style
	ZoomIndicator          string
	BroadcastIndicator     string
}

// DefaultTheme works in all terminals with ANSI 256 colors.
var DefaultTheme = Theme{
	BorderActive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")),
	BorderInactive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")),
	BorderScrollback: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("178")),
	BorderScrollbackFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("214")),
	Divider: lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")),
	DividerChar: "│",
	StatusBar: lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("252")).
		Padding(0, 1),
	StatusText: lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")),
	ZoomIndicator:      "[ZOOM]",
	BroadcastIndicator: "[BROADCAST]",
}

// TokyoNight is a dark theme inspired by the Tokyo Night color scheme.
var TokyoNight = Theme{
	BorderActive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7AA2F7")),
	BorderInactive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#414868")),
	BorderScrollback: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#E0AF68")),
	BorderScrollbackFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF9E64")),
	Divider: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")),
	DividerChar: "│",
	StatusBar: lipgloss.NewStyle().
		Background(lipgloss.Color("#1A1B26")).
		Foreground(lipgloss.Color("#A9B1D6")).
		Padding(0, 1),
	StatusText: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A9B1D6")),
	ZoomIndicator:      "◉ ZOOM",
	BroadcastIndicator: "◉ BROADCAST",
}

// Dracula is a dark theme inspired by the Dracula color scheme.
var Dracula = Theme{
	BorderActive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#BD93F9")),
	BorderInactive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#44475A")),
	BorderScrollback: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#F1FA8C")),
	BorderScrollbackFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFB86C")),
	Divider: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")),
	DividerChar: "│",
	StatusBar: lipgloss.NewStyle().
		Background(lipgloss.Color("#282A36")).
		Foreground(lipgloss.Color("#F8F8F2")).
		Padding(0, 1),
	StatusText: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2")),
	ZoomIndicator:      "◉ ZOOM",
	BroadcastIndicator: "◉ BROADCAST",
}

// Nord is a dark theme inspired by the Nord color scheme.
var Nord = Theme{
	BorderActive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#88C0D0")),
	BorderInactive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4C566A")),
	BorderScrollback: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#EBCB8B")),
	BorderScrollbackFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#D08770")),
	Divider: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4C566A")),
	DividerChar: "│",
	StatusBar: lipgloss.NewStyle().
		Background(lipgloss.Color("#2E3440")).
		Foreground(lipgloss.Color("#ECEFF4")).
		Padding(0, 1),
	StatusText: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ECEFF4")),
	ZoomIndicator:      "◉ ZOOM",
	BroadcastIndicator: "◉ BROADCAST",
}

// Catppuccin is a dark theme inspired by the Catppuccin Mocha color scheme.
var Catppuccin = Theme{
	BorderActive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#CBA6F7")),
	BorderInactive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#585B70")),
	BorderScrollback: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#F9E2AF")),
	BorderScrollbackFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FAB387")),
	Divider: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6C7086")),
	DividerChar: "│",
	StatusBar: lipgloss.NewStyle().
		Background(lipgloss.Color("#1E1E2E")).
		Foreground(lipgloss.Color("#CDD6F4")).
		Padding(0, 1),
	StatusText: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CDD6F4")),
	ZoomIndicator:      "◉ ZOOM",
	BroadcastIndicator: "◉ BROADCAST",
}
