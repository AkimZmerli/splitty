# Splitty Deployment Strategy: Options Analysis

This document contemplates two primary deployment strategies for Splitty as a production-ready tool.

## Option 1: Terminal Multiplexer (Standalone Binary)

### Implementation

Splitty operates as a traditional terminal multiplexer, similar to `tmux` or `screen`:

```bash
# Build
go build -o splitty .

# Install to PATH
sudo mv splitty /usr/local/bin/

# Usage
splitty                    # Start a new session
splitty -s sessionname     # Create named session
splitty -a sessionname     # Attach to session
```

### Advantages

**User Control**
- Users explicitly invoke Splitty when desired
- Non-invasive - doesn't interfere with existing shell workflows
- Works seamlessly in any environment (Docker, SSH, CI/CD, non-interactive shells)
- No startup overhead or initialization complexity

**Technical Simplicity**
- Single binary, standard distribution (Homebrew, pkg managers)
- No shell integration code to maintain
- Works with any shell: bash, zsh, fish, ksh, tcsh
- Minimal surface area for bugs and edge cases

**Industry Precedent**
- Users already understand this model (tmux, screen, etc.)
- Well-defined conventions for sessions, naming, nesting
- Existing documentation templates and best practices

**Predictability**
- Behavior is consistent across environments
- No conflicts with existing shell initialization files
- Clear lifecycle: `splitty` to start, `exit` to close

### Disadvantages

**User Friction**
- Requires explicit invocation - users must remember to use it
- Additional mental overhead - one more tool to manage
- Context switching between environments where Splitty is/isn't active

**Session Persistence Gap**
- Unlike a native shell, losing the Splitty window loses the session
- Terminal close = forced session disconnect (unless backgrounded via tmux/nohup pattern)
- Less integrated with OS window management

### Implementation Considerations

**Distribution Channels**
```bash
# Homebrew (macOS/Linux)
brew install splitty

# Direct download
https://github.com/splitty/releases

# Docker
docker run splitty:latest

# Manual compilation
go build && cp splitty /usr/local/bin/
```

**Session Management**
- Store sessions in temp files or dedicated directory
- Handle graceful cleanup on terminal close
- Consider session listing/management CLI flags

**Shell Integration (Optional)**
- Provide `.zshrc`/`.bashrc` snippets for power users
- Optional alias: `alias st=splitty`
- Optional keybinding for quick access

---

## Option 2: Shell Integration (Auto-Launch)

### Implementation

Splitty auto-launches when terminal opens by adding initialization code:

**~/.zshrc**
```bash
# Auto-start Splitty for interactive shells
if [[ -o interactive && -z "$SPLITTY_SESSION" ]]; then
  export SPLITTY_SESSION=1
  exec splitty
fi
```

**~/.bashrc**
```bash
# Auto-start Splitty for interactive shells
if [[ $- == *i* && -z "$SPLITTY_SESSION" ]]; then
  export SPLITTY_SESSION=1
  exec splitty
fi
```

### Advantages

**Seamless User Experience**
- Zero friction - Splitty is *always* available without thinking
- Users open terminal and immediately get split-pane functionality
- Natural entry point for new users discovering the tool

**Deep Integration**
- Splitty becomes the default "terminal experience"
- Shell history, aliases, and settings available immediately
- Feels native because initialization is automatic

**Consistency**
- All interactive terminal sessions use Splitty
- Uniform UX across local and remote (SSH) connections
- No mental overhead about when to use it

### Disadvantages

**Startup Overhead**
- Every terminal open incurs Splitty initialization cost
- Noticeable lag on slow machines or SSH connections
- May break fast copy-paste workflows that expect instant shells

**Shell Assumption Violations**
- Non-interactive shells (scripts, cron jobs, CI/CD) break unexpectedly
  ```bash
  # This would hang waiting for Splitty
  echo "pwd" | ssh user@host
  ```
- Requires careful guards (`[[ -o interactive ]]`) to prevent breakage
- Edge cases: `tmux new-session -c dir -x cols -y lines`, `nohup`, subshells

**Configuration Complexity**
- Each user needs custom init code
- Different for bash, zsh, fish, etc.
- Documentation burden: teaching users how to install vs. just running binary
- Maintenance: tracking shell updates and compatibility

**Nested Multiplexer Hell**
- If user already uses tmux, now you have: SSH → tmux → Splitty
- Keybinding conflicts and confusion
- Escape sequence interpretation issues

**Lack of Explicit Control**
- Hard to opt-out for specific scenarios
- Users must debug why `splitty` env var exists
- Difficult to use Splitty alongside other multiplexers

**Scripts and Tooling Break**
```bash
# What happens here?
docker exec splitty-container bash -c "ls -la"

# Or here?
ssh -t user@host 'echo $HOME'

# Or in CI?
- run: make test
```
All of these would try to launch Splitty and potentially hang or fail.

---

## Comparative Analysis

| Aspect | Option 1: Multiplexer | Option 2: Shell Integration |
|--------|----------------------|---------------------------|
| **User Friction** | Medium (must invoke) | Low (automatic) |
| **Setup Complexity** | Simple (binary only) | Medium (init files) |
| **Breaking Changes** | None | High (scripts, SSH, CI) |
| **Distribution** | Standard (Homebrew, etc.) | Per-user (shell config) |
| **Environment Compatibility** | Excellent (all shells, Docker, CI) | Poor (breaks non-interactive) |
| **Startup Cost** | Per-invocation | Per-shell-open |
| **Nested Multiplexers** | Supported | Problematic |
| **Industry Precedent** | tmux, screen | Uncommon |
| **Discoverability** | Requires documentation | Very high (always available) |

---

## Recommendation: Option 1 as Primary, Option 2 as Optional

**Start with Option 1** (Terminal Multiplexer) because:

1. **Zero Risk** - No breaking changes to existing workflows
2. **Broader Appeal** - Works in all environments (local, SSH, Docker, CI/CD)
3. **Industry Alignment** - Users already understand this mental model
4. **Easier Distribution** - Standard binary installation via package managers
5. **Lower Support Burden** - Fewer edge cases and environment-specific issues

**Offer Option 2** as an opt-in shell integration for power users:
- Provide sample `~/.zshrc` and `~/.bashrc` snippets
- Clear documentation on when/why to use it
- Guard clauses to prevent non-interactive shell breakage
- Warning about SSH and CI/CD implications

### Phase 1: Stable Release
- Focus on Option 1 as the primary distribution
- Compile to binary, test thoroughly
- Release via Homebrew and GitHub releases
- Build user base with traditional multiplexer UX

### Phase 2: Optional Enhancement
- If user feedback requests auto-launch, provide documented shell integration
- Make it explicitly opt-in with clear warnings
- Provide detection/setup scripts to help users integrate safely

### Phase 3: Advanced
- Consider per-environment configuration (disable in Docker, SSH, etc.)
- Build session persistence layer if needed
- Explore true "native" integration only if Splitty shows strong product-market fit

---

## Implementation Roadmap

### Immediate (Option 1)
- [x] Core multiplexer functionality (split panes, keybindings, themes)
- [ ] Session persistence (basic file-based storage)
- [ ] Configuration file support (~/.splittyrc)
- [ ] Homebrew formula
- [ ] Release pipeline automation

### Future (Option 2 - If Justified)
- [ ] Shell integration snippets (with warnings)
- [ ] Automatic environment detection
- [ ] Safe guards for non-interactive shells
- [ ] Documentation and tutorials

### Never (Unless Massive Evidence)
- Replace shell entirely
- Force integration without opt-in
- Break existing workflows

---

## Conclusion

Splitty's value proposition is **split-pane terminal management**. Option 1 (Terminal Multiplexer) delivers that value clearly and safely. Option 2 (Shell Integration) adds convenience at significant complexity and compatibility cost.

**Recommendation: Pursue Option 1 as the production deployment strategy**, with Option 2 available as documented opt-in for power users willing to accept the tradeoffs.

This approach maximizes reach, minimizes support burden, and maintains the Unix philosophy of "do one thing well."
