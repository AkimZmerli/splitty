package splitty

import "charm.land/lipgloss/v2"

// themeEntry pairs a theme with its display name.
type themeEntry struct {
	Name  string
	Theme Theme
}

// themeList is the ordered list of built-in themes for cycling.
var themeList = []themeEntry{
	{"Tokyo Night", TokyoNight},
	{"Nightshade", Nightshade},
	{"Glacier", Glacier},
	{"Sorbet", Sorbet},
}

// Theme defines the visual styling for split panes.
type Theme struct {
	BorderActive           lipgloss.Style
	BorderInactive         lipgloss.Style
	BorderScrollback       lipgloss.Style
	BorderScrollbackFocused lipgloss.Style
	BorderResize           lipgloss.Style
	BorderCopyMode         lipgloss.Style
	BorderCopyModeFocused  lipgloss.Style
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
	BorderResize: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("44")),
	BorderCopyMode: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("76")),
	BorderCopyModeFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("118")),
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
	BorderResize: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#2AC3DE")),
	BorderCopyMode: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#73DACA")),
	BorderCopyModeFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9ECE6A")),
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

// Nightshade is a dark purple theme with vivid accents.
var Nightshade = Theme{
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
	BorderResize: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8BE9FD")),
	BorderCopyMode: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#50FA7B")),
	BorderCopyModeFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#69FF94")),
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

// Glacier is an arctic blue theme with cool, muted tones.
var Glacier = Theme{
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
	BorderResize: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#88C0D0")),
	BorderCopyMode: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#A3BE8C")),
	BorderCopyModeFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#B3CE9C")),
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

// Sorbet is a warm pastel theme with soft lavender and peach tones.
var Sorbet = Theme{
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
	BorderResize: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#94E2D5")),
	BorderCopyMode: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#A6E3A1")),
	BorderCopyModeFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#B6F3B1")),
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
