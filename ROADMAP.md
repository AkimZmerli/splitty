# Splitty Roadmap

Future feature ideas and development direction for Splitty.

Splitty's core thesis: **developers running multiple AI agents need a terminal multiplexer built for observability, not just shell management.** Everything on this roadmap serves that thesis.

---

## Phase 1: Agent Observability

Features that make monitoring multiple concurrent agent sessions effortless.

### Pane Labels and Identification
- Assign custom names to panes (`agent-1`, `code-review`, `test-runner`)
- Display labels in pane borders or a minimal header bar
- Color-code panes by role or agent type
- API: `WithPaneLabel(name string)`, `SetLabel(paneID, name string)`

### Activity Indicators
- Visual indicator when a pane has new output (dot, color pulse, border highlight)
- "Idle" vs "active" state per pane - know which agent is still working
- Configurable idle timeout threshold
- Aggregate activity summary in status bar ("3 active, 1 idle")

### Follow Mode
- Lock view to a specific pane's output stream (auto-scroll pinned)
- Picture-in-picture: small preview of followed pane while navigating others
- Quick-follow: double-tap a direction key to lock onto that pane
- Unfollow on any navigation input

### Output Search
- Search across all pane output with a single keybinding
- Highlight matches across visible panes simultaneously
- Jump-to-match navigation
- Regex support for pattern matching agent output

---

## Phase 2: Agentic Workflow Features

Purpose-built features for AI agent orchestration.

### Agent Session Templates
- Pre-configured layouts for common agent workflows:
  - `agent-pair`: Two panes, one for coding agent, one for review agent
  - `agent-trio`: Three panes, code + test + review
  - `agent-monitor`: One large pane + multiple small observer panes
- Save custom agent layouts as reusable templates

### Output Capture and Logging
- Record all pane output to timestamped log files
- Per-pane log files for post-session analysis
- Structured log format (timestamp, pane ID, content) for machine parsing
- Export session transcript for debugging agent interactions

### Cross-Pane Awareness
- Detect when an agent in one pane is waiting for input
- Visual notification when any agent errors or stops
- Pattern-based alerts (match on "error", "failed", "waiting for input")
- Configurable alert rules per pane or globally

### Command Injection
- Send a command to a specific pane by name: `SendTo("agent-1", "continue\n")`
- Queue commands for panes (send when agent finishes current task)
- Macro support: send a sequence of commands across multiple panes

---

## Phase 3: UX Polish

General improvements that benefit all users.

### Improved Presets
- `PresetAgent`: Optimized for 2-4 agent sessions with status panel
- `PresetMonitor`: One large pane with 3 small panes along the bottom
- `PresetWide`: Side-by-side for ultrawide monitors
- User-defined presets saved to `~/.config/splitty/presets/`

### Enhanced Theming
- Theme hot-reload without restart
- Per-pane theme overrides (dark pane for code agent, light for logs)
- Import themes from popular terminal emulators (iTerm2, Alacritty)
- Community theme repository

### Status Bar Evolution
- Show pane count, active agent count, session duration
- Configurable status bar segments (left/center/right)
- Plugin system for custom status bar widgets
- Minimal mode: hide status bar, show only on keypress

### Resize Improvements
- Drag borders with mouse to resize
- Preset ratios (60/40, 70/30, equal) via keybinding
- "Focus mode": temporarily expand active pane, shrink others
- Snap-to-grid for clean layouts

---

## Phase 4: Integration and Ecosystem

Connect Splitty to the broader developer toolchain.

### CLI Interface
- `splitty new` - start new session
- `splitty list` - list active sessions
- `splitty attach <name>` - reattach to session
- `splitty send <pane> <command>` - send command to running pane
- `splitty layout save/load` - manage layouts from CLI

### Framework Integration
- Hooks for agent frameworks (LangChain, AutoGen, CrewAI)
- Automatic pane spawning when agents start
- Structured output parsing from known agent formats
- Event bus for external tools to subscribe to pane activity

### Remote Sessions
- Share a Splitty session over SSH (like tmux attach)
- Read-only spectator mode for pair debugging
- Session URL sharing for quick collaboration

### Plugin System
- Lua or Go plugin interface for custom behaviors
- Community plugins directory
- Example plugins: git status per pane, resource monitor, agent health checker

---

## Phase 5: Distribution

Get Splitty into developers' hands.

### Package Managers
- Homebrew formula (`brew install splitty`)
- APT/DNF packages for Linux
- AUR package for Arch
- Scoop/Chocolatey for Windows
- Nix flake

### Binary Releases
- Automated release pipeline (GoReleaser)
- Cross-compilation: macOS (arm64, amd64), Linux (arm64, amd64), Windows
- Signed binaries
- Checksums for verification

### Container Support
- Official Docker image
- Docker Compose example for agent workflows
- Devcontainer configuration

### Documentation Site
- Dedicated docs site (GitHub Pages or similar)
- Interactive tutorials
- Video walkthroughs of agent workflows
- API reference auto-generated from Go docs

---

## Backlog (Ideas Worth Exploring)

These are speculative ideas that need more thought before committing to a phase.

- **Pane scripting**: Write simple scripts that react to pane output (if agent prints X, do Y)
- **Session recording and playback**: Record a full multi-agent session, replay it later
- **Metrics dashboard**: Built-in resource usage per pane (CPU, memory of child process)
- **Voice control**: Navigate panes with voice commands (accessibility + hands-free monitoring)
- **AI-powered routing**: Automatically route focus to the pane that needs attention most
- **Splitty Cloud**: Persist sessions across machines, share layouts with team
- **Terminal within terminal**: Nest Splitty instances for complex hierarchical workflows

---

## Contributing to the Roadmap

Have an idea? Open an issue tagged `roadmap` or `feature-request`. The best features come from real workflows, so describe your use case alongside your suggestion.

Priority is determined by:
1. Does it serve the agentic workflow thesis?
2. How many users benefit?
3. Implementation complexity vs. value delivered
4. Does it maintain Splitty's simplicity?
