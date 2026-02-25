---
tags: ["chore", "ci"]
type: chore
---

# Setup GitHub Actions

## Implementation Plan

- **CI/CD Workflow (`.github/workflows/ci.yml`)**:
  - **Trigger**: Run on push to `main` and on pull requests.
  - **Environment**: Use `actions/setup-go` to provision the Go environment.
  - **Linting**: Execute `golangci-lint` to maintain code quality.
  - **Security**: Run `govulncheck` to scan for known vulnerabilities in dependencies.
  - **Testing**: Run unit tests and integration tests (`go test -tags=integration ./...`).
  - **Code Coverage**: Add a step to calculate test coverage (e.g., `go test -coverprofile=coverage.out ./...`), but temporarily comment it out until Codecov integration is set up.

- **Release Workflow (`.github/workflows/release.yml`)**:
  - **Trigger**: Run exclusively on new git tags (e.g., `v*`).
  - **GoReleaser Integration**: Use the official `goreleaser/goreleaser-action` with non-premium features.
  - **Cross-Platform Compilation**: Configure `.goreleaser.yaml` to compile binaries targeting Linux, Windows, and macOS across `amd64` and `arm64` architectures.
  - **Package Management Distribution**: Configure `goreleaser` to build and publish packages for:
    - Homebrew (`brew`)
    - Chocolatey (`choco`)
    - RPM, DEB, Snap, Flatpak, and AppImage.
  - **Announcements**: Enable release announcements via GitHub Releases (changelog generation based on commits).
  - **Docker Placeholder**: Add Docker build/push configuration to `.goreleaser.yaml`, but leave it commented out/disabled until the Docker Hub account is fully established.
