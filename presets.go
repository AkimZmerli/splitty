package splitty

// Preset names for built-in layout configurations.
const (
	PresetSingle = "single" // One pane (default)
	PresetDev    = "dev"    // 60/40 vertical split
	PresetTriple = "triple" // Three columns
	PresetQuad   = "quad"   // 2x2 grid
)

// presetBuilder creates a layout tree for a preset.
type presetBuilder func(shell string, env []string, width, height, scrollbackSize int) node

// presetRegistry holds all registered presets.
var presetRegistry = map[string]presetBuilder{
	PresetSingle: buildPresetSingle,
	PresetDev:    buildPresetDev,
	PresetTriple: buildPresetTriple,
	PresetQuad:   buildPresetQuad,
}

// RegisterPreset adds a custom layout preset.
// The builder function receives the shell command and environment and should
// return a fully constructed node tree.
func RegisterPreset(name string, builder func(shell string, env []string, width, height, scrollbackSize int) node) {
	presetRegistry[name] = builder
}

func buildPresetSingle(shell string, env []string, width, height, scrollbackSize int) node {
	return &leafNode{pane: newPane(shell, env, width, height, scrollbackSize)}
}

func buildPresetDev(shell string, env []string, width, height, scrollbackSize int) node {
	return &splitNode{
		dir:   Vertical,
		ratio: 0.6,
		first: &leafNode{pane: newPane(shell, env, width, height, scrollbackSize)},
		second: &splitNode{
			dir:    Horizontal,
			ratio:  0.6,
			first:  &leafNode{pane: newPane(shell, env, width, height, scrollbackSize)},
			second: &leafNode{pane: newPane(shell, env, width, height, scrollbackSize)},
		},
	}
}

func buildPresetTriple(shell string, env []string, width, height, scrollbackSize int) node {
	return &splitNode{
		dir:   Vertical,
		ratio: 0.33,
		first: &leafNode{pane: newPane(shell, env, width, height, scrollbackSize)},
		second: &splitNode{
			dir:    Vertical,
			ratio:  0.5,
			first:  &leafNode{pane: newPane(shell, env, width, height, scrollbackSize)},
			second: &leafNode{pane: newPane(shell, env, width, height, scrollbackSize)},
		},
	}
}

func buildPresetQuad(shell string, env []string, width, height, scrollbackSize int) node {
	return &splitNode{
		dir:   Vertical,
		ratio: 0.5,
		first: &splitNode{
			dir:    Horizontal,
			ratio:  0.5,
			first:  &leafNode{pane: newPane(shell, env, width, height, scrollbackSize)},
			second: &leafNode{pane: newPane(shell, env, width, height, scrollbackSize)},
		},
		second: &splitNode{
			dir:    Horizontal,
			ratio:  0.5,
			first:  &leafNode{pane: newPane(shell, env, width, height, scrollbackSize)},
			second: &leafNode{pane: newPane(shell, env, width, height, scrollbackSize)},
		},
	}
}
