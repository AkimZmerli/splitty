package splitty

import "testing"

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.shell == "" {
		t.Error("shell should be set")
	}
	if m.minWidth != 10 {
		t.Errorf("default minWidth: expected 10, got %d", m.minWidth)
	}
	if m.minHeight != 3 {
		t.Errorf("default minHeight: expected 3, got %d", m.minHeight)
	}
	if !m.statusBar {
		t.Error("statusBar should default to true")
	}
	if !m.mouse {
		t.Error("mouse should default to true")
	}
}

func TestWithShell(t *testing.T) {
	m := New(WithShell("/bin/zsh"))
	if m.shell != "/bin/zsh" {
		t.Errorf("expected /bin/zsh, got %s", m.shell)
	}
}

func TestWithTheme(t *testing.T) {
	m := New(WithTheme(TokyoNight))
	if m.theme.ZoomIndicator != TokyoNight.ZoomIndicator {
		t.Error("theme was not set correctly")
	}
}

func TestWithMinSize(t *testing.T) {
	m := New(WithMinSize(20, 10))
	if m.minWidth != 20 {
		t.Errorf("expected minWidth=20, got %d", m.minWidth)
	}
	if m.minHeight != 10 {
		t.Errorf("expected minHeight=10, got %d", m.minHeight)
	}
}

func TestWithStatusBar(t *testing.T) {
	m := New(WithStatusBar(false))
	if m.statusBar {
		t.Error("statusBar should be false")
	}
}

func TestWithMouse(t *testing.T) {
	m := New(WithMouse(false))
	if m.mouse {
		t.Error("mouse should be false")
	}
}

func TestWithPreset(t *testing.T) {
	m := New(WithPreset(PresetDev))
	if m.presetName != PresetDev {
		t.Errorf("expected preset %s, got %s", PresetDev, m.presetName)
	}
}

func TestWithEnv(t *testing.T) {
	env := []string{"FOO=bar", "BAZ=qux"}
	m := New(WithEnv(env))
	if len(m.env) != 2 {
		t.Fatalf("expected 2 env vars, got %d", len(m.env))
	}
	if m.env[0] != "FOO=bar" {
		t.Errorf("expected FOO=bar, got %s", m.env[0])
	}
}

func TestMultipleOptions(t *testing.T) {
	m := New(
		WithShell("/bin/bash"),
		WithStatusBar(false),
		WithMouse(false),
		WithMinSize(15, 5),
	)
	if m.shell != "/bin/bash" {
		t.Errorf("shell: expected /bin/bash, got %s", m.shell)
	}
	if m.statusBar {
		t.Error("statusBar should be false")
	}
	if m.mouse {
		t.Error("mouse should be false")
	}
	if m.minWidth != 15 || m.minHeight != 5 {
		t.Error("minSize not set correctly")
	}
}
