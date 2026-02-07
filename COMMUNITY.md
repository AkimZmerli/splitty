# Splitty Community Growth Strategy

This document outlines the philosophy behind Splitty, how to pitch it, how to document it, and the steps toward building a developer community around it.

---

## Philosophy

### The Problem

Agentic coding is becoming the default way developers work. Running multiple AI agents in parallel -- a coding agent, a test runner, a reviewer -- is increasingly normal. But the tools developers use to manage these sessions were designed decades ago for a different workflow.

Tmux was built for managing shell sessions. iTerm2 was built for tabbed terminal windows. Neither was designed for **observing multiple concurrent AI agents** where fast spatial navigation is the primary interaction.

### The Insight

Developers are gamers. The most natural spatial navigation for a developer's left hand is WASD -- the same movement pattern used in every first-person game for the last 30 years. This muscle memory is deeply embedded and transfers instantly.

When you're watching 4 AI agents work in parallel, switching between them should feel like moving through a game world: reflexive, spatial, instant. Not a prefix key followed by an arrow. Not clicking a tab. Just `Ctrl+W/A/S/D`.

### The Vision

Splitty is the terminal multiplexer for the agentic era. It prioritizes:

1. **Observability** -- See what all your agents are doing at a glance
2. **Intuitive navigation** -- WASD means your hands already know how to use it
3. **Low cognitive load** -- Monitoring multiple sessions should feel effortless
4. **Developer-native UX** -- Built by developers, for developers, with gaming DNA

### Core Principles

- **Left hand navigates, right hand is free.** Every navigation action is reachable with the left hand in its natural resting position.
- **No prefix keys.** Tmux's `Ctrl+B` then action pattern adds latency and cognitive overhead. Splitty uses direct keybindings.
- **Visual clarity.** Pane borders, themes, and status indicators tell you what's happening without reading. Color and layout communicate state.
- **Simple by default, powerful when needed.** Works out of the box with sane defaults. Customization is available but never required.
- **Embeddable.** Splitty is a library first. It can be integrated into any Bubble Tea application, not just used as a standalone tool.

---

## Pitching Splitty

### Elevator Pitch (10 seconds)

> "Splitty is a terminal multiplexer with WASD navigation, built for developers who run multiple AI agents in parallel. Navigate panes like you navigate games."

### Short Pitch (30 seconds)

> "If you're using AI coding agents, you're probably running multiple sessions at once -- a code agent, a test runner, maybe a reviewer. Switching between them in tmux or iTerm2 is clunky. Splitty uses WASD navigation so your left hand moves between panes the same way you move in games. It's instant, spatial, and built for observing concurrent agents."

### Conference/Meetup Pitch (2 minutes)

> "Agentic coding is changing how we work. We're not typing into one terminal anymore -- we're orchestrating. Three agents, four panes, all working at the same time. And we need to watch all of them.
>
> The problem is, our multiplexers weren't designed for this. Tmux is great for managing shell sessions, but switching between panes requires a prefix key, then a direction. That's two keystrokes and a mental context switch every time you glance at another agent.
>
> Splitty fixes this with one design decision: WASD navigation. Ctrl+W goes up, Ctrl+A goes left, Ctrl+S goes down, Ctrl+D goes right. Your left hand never leaves its natural position. It feels like playing a game, and that's the point -- developers are gamers, and this muscle memory transfers immediately.
>
> It's built on Bubble Tea, it's a Go library you can embed in your own TUI apps, and it's open source. If you run multiple agents, try it. Your hands will thank you."

### Written Pitch (README/Social Media)

> **Splitty**: Terminal multiplexer with WASD navigation for agentic coding workflows.
>
> - `Ctrl+W/A/S/D` to navigate panes (left hand never moves)
> - Built for watching multiple AI agents work in parallel
> - 5 themes, layout presets, scrollback, broadcast mode
> - Embeddable Bubble Tea component
> - No prefix keys. No clicking tabs. Just WASD.

### Key Messaging Points

1. **Lead with the use case**: "Built for agentic coding" -- not "tmux alternative"
2. **Lead with the UX**: WASD navigation is the hook that makes people try it
3. **Gaming angle**: "Navigate panes like you navigate games" -- instantly communicates the feel
4. **No prefix keys**: This is a concrete tmux pain point that developers immediately relate to
5. **Show, don't tell**: A 15-second GIF of navigating 4 agent panes with WASD communicates more than any pitch

---

## Documentation Strategy

### Documentation Tiers

**Tier 1: README (First Impression)**
- Must communicate the value proposition in 5 seconds
- WASD + agentic coding in the first paragraph
- Quick start that works in under a minute
- Keybindings table front and center
- Link to deeper docs for everything else

**Tier 2: In-Repo Docs (Active Users)**
- `ARCHITECTURE.md` -- How the code is organized (exists)
- `DEPLOYMENT.md` -- How to install and run (exists)
- `ROADMAP.md` -- Where the project is going (exists)
- `COMMUNITY.md` -- This document
- `CONTRIBUTING.md` -- How to contribute (future)
- `CHANGELOG.md` -- Version history (future)

**Tier 3: External Docs (Growth)**
- Dedicated documentation site (GitHub Pages)
- Tutorial: "Setting Up Splitty for Multi-Agent Workflows"
- Tutorial: "Embedding Splitty in Your Bubble Tea App"
- Video: "Why WASD Navigation Changes Everything"
- Blog post: "Agentic Coding Needs Better Terminal Tools"

### Documentation Principles

1. **Show the workflow, not the API.** Developers care about "how do I watch 3 agents at once?" not "what does FocusPane() do?"
2. **Lead with keybindings.** The keybinding table is the most important reference in the entire project.
3. **Use real agent examples.** Don't show `ls -la` in demo screenshots. Show Claude, GPT, or AutoGen agents working.
4. **Keep it concise.** Developers skim. Short paragraphs. Tables over prose. Code over explanation.
5. **Update docs with code.** Every feature PR should include doc updates. Stale docs are worse than no docs.

### Content Calendar (Suggested)

**Month 1: Foundation**
- Finalize README with agentic positioning
- Record a demo GIF/video of WASD navigation with 4 panes
- Write "Getting Started" tutorial
- Publish initial blog post introducing Splitty

**Month 2: Depth**
- Write "Embedding Splitty" tutorial
- Create video walkthrough
- Document all configuration options
- Write comparison piece: "Splitty vs tmux for Agent Workflows"

**Month 3: Community**
- Publish "Contributing" guide
- Create issue templates (bug, feature, roadmap idea)
- Write "Building Custom Themes" tutorial
- Start collecting user stories

---

## Community Growth Plan

### Phase 1: Seed (Months 1-2)

**Goal**: Get the first 50 stars and 5 contributors.

**Actions:**
- Post on Go subreddit (r/golang) with demo GIF
- Post on Hacker News (Show HN) with the agentic angle
- Share in Bubble Tea Discord community
- Share in AI/agent-focused communities (r/LocalLLaMA, r/ChatGPT, AI Discord servers)
- Tweet/post with the gaming angle ("Navigate your terminal like a game")
- Reach out to Go/TUI newsletter curators

**Channels (prioritized):**
1. Charm community (Bubble Tea Discord, GitHub) -- most aligned audience
2. Go community (r/golang, Go weekly newsletter, GopherCon Slack)
3. AI agent communities (where multi-agent users are)
4. Hacker News (broader developer audience)
5. Twitter/X, Mastodon, Bluesky (developer accounts)

**Key Assets Needed:**
- 15-second demo GIF showing WASD navigation across 4 agent panes
- Comparison screenshot: tmux keybinding vs Splitty keybinding
- Blog post explaining the "why" behind WASD

### Phase 2: Grow (Months 3-6)

**Goal**: 200+ stars, 10+ contributors, first external blog posts about Splitty.

**Actions:**
- Submit talk proposals to Go meetups/conferences
- Write guest posts for Go/developer blogs
- Create YouTube tutorial series
- Engage with issues and PRs promptly (< 24 hour response)
- Feature community contributions in release notes
- Start a "Splitty Showcase" for creative uses

**Community Health:**
- Clear issue labels (good-first-issue, help-wanted, roadmap)
- PR review turnaround under 48 hours
- Monthly release cadence
- Transparent roadmap (ROADMAP.md)

### Phase 3: Sustain (Months 6+)

**Goal**: Self-sustaining community with regular external contributions.

**Actions:**
- Establish contributor recognition (CONTRIBUTORS.md, release credits)
- Consider sponsoring/funding model if adoption warrants it
- Explore integration partnerships with agent frameworks
- Present at GopherCon or similar conferences
- Build maintainer team (2-3 trusted contributors with merge access)

---

## Competitive Positioning

### What Splitty is NOT

- Not a tmux replacement for all use cases (tmux is battle-tested for server session management)
- Not a terminal emulator (Splitty runs inside a terminal)
- Not a shell (it manages shells, doesn't replace them)
- Not trying to do everything (focused on the WASD + observability thesis)

### What Splitty IS

- The best multiplexer for watching multiple concurrent processes
- The most intuitive pane navigation available (WASD)
- A modern, embeddable Go library (not a legacy C codebase)
- Purpose-built for the agentic coding era

### Differentiation Matrix

| Feature | Splitty | tmux | Zellij | screen |
|---------|---------|------|--------|--------|
| WASD navigation | Yes | No | No | No |
| No prefix keys | Yes | No | No (has modes) | No |
| Embeddable library | Yes | No | No | No |
| Agentic workflow focus | Yes | No | No | No |
| Gaming-familiar UX | Yes | No | Partial | No |
| Bubble Tea ecosystem | Yes | No | No | No |
| Themes built-in | 5 | 0 (manual) | 1 | 0 |
| Modern language (Go) | Yes | C | Rust | C |

---

## Metrics to Track

### Adoption
- GitHub stars (vanity but useful for social proof)
- Go module downloads (pkg.go.dev)
- Clones and unique visitors (GitHub Insights)

### Engagement
- Issues opened (shows people are using it)
- PRs from external contributors
- Discord/community messages
- Blog posts and mentions by others

### Quality
- Issue close time
- PR review time
- Test coverage
- Release frequency

---

## Next Steps (Immediate)

1. Record a demo GIF showing WASD navigation with 4 panes running agents
2. Prepare a Show HN post with the agentic coding angle
3. Share in Charm/Bubble Tea Discord
4. Write an introductory blog post
5. Set up issue templates and labels for community contributions
