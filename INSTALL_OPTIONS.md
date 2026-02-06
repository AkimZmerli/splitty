# Install Options Analysis for Splitty

## Current State

Splitty is a **Go library** for split-pane terminal multiplexing with Bubble Tea.

### Existing Installation
- **Module path**: `github.com/AkimZmerli/splitty`
- **Primary method**: `go get github.com/AkimZmerli/splitty`
- **Go version**: 1.24.2

### Current Makefile Targets
- `make test` - Run tests with race detection
- `make test-short` - Run short tests
- `make test-cover` - Generate coverage reports
- `make lint` - Run golangci-lint
- `make vet` - Run go vet
- `make build` - Build all packages
- `make examples` - Build example programs
- `make clean` - Remove coverage artifacts

### CI/CD
- Tests on Ubuntu and macOS (Go 1.24 & 1.25)
- Runs build, vet, test, and lint
- **No release automation** (no GoReleaser, no binary publishing)
- **No version tags** (no semver releases)

### Example Programs
Located in `examples/`:
- basic
- custom-theme
- custom-keys
- embedded
- presets
- ui-components-test
- ui-styles-showcase
- ui-forms-demo
- ui-components-comprehensive
- ui-components-advanced

## Missing Install Options

### High Priority

#### 1. **Version Tags** (Semver)
**Status**: ❌ Not implemented
**Benefit**: Users can depend on specific versions

```bash
# Example usage after tagging
go get github.com/AkimZmerli/splitty@v0.1.0
```

#### 2. **Developer Tooling Setup**
**Status**: ❌ Not implemented
**Benefit**: One-command setup for contributors

**Option A: Makefile target**
```makefile
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

dev-setup: install-tools
	go mod download
	go mod verify
```

**Option B: tools.go pattern**
```go
//go:build tools

package tools

import _ "github.com/golangci/golangci-lint/cmd/golangci-lint"
```

Then run: `go install $(go list -f '{{join .Imports " "}}' tools.go)`

#### 3. **CONTRIBUTING.md**
**Status**: ❌ Not implemented
**Benefit**: Clear onboarding for new contributors

Should document:
- How to set up dev environment (`make install-tools`)
- How to run tests (`make test`)
- How to build examples (`make examples`)
- Code style and linting requirements

### Medium Priority

#### 4. **GoReleaser Configuration**
**Status**: ❌ Not implemented
**Benefit**: Automated binary releases for example programs

Users could try examples without cloning:
```bash
# After GoReleaser setup
curl -sfL https://github.com/AkimZmerli/splitty/releases/download/v0.1.0/splitty-basic-linux-amd64 -o splitty-basic
chmod +x splitty-basic
./splitty-basic
```

Example `.goreleaser.yml`:
```yaml
builds:
  - id: basic
    main: ./examples/basic
    binary: splitty-basic
  - id: presets
    main: ./examples/presets
    binary: splitty-presets
```

#### 5. **GitHub Release Workflow**
**Status**: ❌ Not implemented
**Benefit**: Automated releases on git tags

`.github/workflows/release.yml` to trigger GoReleaser on version tags.

#### 6. **Example Installation Target**
**Status**: ❌ Not implemented
**Benefit**: Easy local installation of examples

```makefile
examples-install:
	go install ./examples/basic
	go install ./examples/presets
	go install ./examples/custom-theme
```

### Low Priority

#### 7. **Homebrew Tap**
**Status**: ❌ Not implemented
**Benefit**: `brew install akimzmerli/tap/splitty-examples`

Only worth it if examples become popular standalone tools.

#### 8. **Install Script**
**Status**: ❌ Not implemented
**Benefit**: One-liner installation

```bash
# Example
curl -sfL https://raw.githubusercontent.com/AkimZmerli/splitty/main/install.sh | sh
```

Not really needed for a Go library.

#### 9. **Docker Examples**
**Status**: ❌ Not implemented
**Benefit**: Isolated demo environment

Probably overkill for this project.

## Recommendations

### For Library Consumers
1. ✅ Current `go get` installation is correct
2. Add version tags (v0.1.0, v0.2.0, etc.) for better dependency management
3. Document version pinning in README

### For Contributors
1. Add `make install-tools` target
2. Consider `tools.go` for dependency tracking
3. Add `CONTRIBUTING.md` with setup instructions
4. Document required Go version and dependencies

### For Example Programs
1. Add `make examples-install` for local installation
2. Consider GoReleaser for publishing example binaries
3. Document how to run examples in README

## Implementation Priority

1. **Week 1**: Version tagging + `make install-tools` + CONTRIBUTING.md
2. **Week 2**: GoReleaser setup + release workflow
3. **Week 3**: Homebrew tap (if desired)

## Key Insight

Since Splitty is a **library** (not a CLI tool), the focus should be:
- ✅ Library consumer experience (version tags, clear imports)
- ⚠️ Contributor experience (automate dev setup)
- ⚠️ Example accessibility (consider publishing binaries)

The current `go get` approach is appropriate. Main gaps are developer tooling automation and version management.
