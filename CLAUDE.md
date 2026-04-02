# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

This repository contains custom OpenClaw skills published to the Claw Skills registry (`clawskills.sh`) under the `panthrocorp` author. Each skill lives in its own subdirectory with an independent Go module, build pipeline, and SKILL.md manifest.

## Repository structure

Each skill is a self-contained directory:

```
<skill-name>/
  SKILL.md              # clawhub manifest (YAML frontmatter + agent instructions)
  LICENSE               # MIT-0 (clawhub requirement)
  go.mod / go.sum       # Go module
  main.go               # entrypoint
  cmd/                  # Cobra CLI commands
  internal/             # internal packages
  .goreleaser.yml       # cross-compilation config
```

## Build and test

From any skill directory:

```bash
go build -o <skill-name> .
go test ./...
```

Cross-compile for the OpenClaw EC2 target (Graviton ARM):

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o <skill-name> .
```

## Architecture

Each skill follows a layered pattern:

- `cmd/` uses Cobra for CLI subcommands. Each service (gmail, calendar, contacts, drive, docs, sheets) has its own file registering commands onto `rootCmd`. Helper functions like `gmailClient()` handle config loading, token decryption, and API client creation inline.
- `internal/config/` holds a JSON config model (`config.json`) that gates which services are active and at what access level. A unified `ServiceMode` type (`off`/`readonly`/`readwrite`) is used by Calendar, Drive, Docs, and Sheets. Gmail and Contacts are simple on/off booleans. Config drives OAuth scope selection via `Config.OAuthScopes()`.
- `internal/oauth/` handles the OAuth2 Desktop flow (localhost redirect) and encrypted token storage. Tokens are stored as `token.enc` (AES-256-GCM encrypted JSON). `SaveToken` refuses to write without an encryption key; there is no unencrypted fallback.
- `internal/crypto/` implements AES-256-GCM with HKDF-SHA256 key derivation. Wire format: `salt (16B) || nonce (12B) || ciphertext+tag`.
- `internal/google/` contains thin typed wrappers around Google API services. Each wrapper exposes only the operations the skill permits (e.g. `GmailClient` has no send/modify/delete methods).

Config and credential directory default: `~/.openclaw/credentials/google-workspace/` (overridable via `GOOGLE_WORKSPACE_CONFIG_DIR` or `--config-dir`).

## Required environment variables

| Variable | Purpose |
|----------|---------|
| `GOOGLE_WORKSPACE_TOKEN_KEY` | Passphrase for AES-256-GCM token encryption |
| `GOOGLE_CLIENT_ID` | Google OAuth2 client ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth2 client secret |

## Security principles

All skills in this repo follow these rules:

- Read-only operations are the default. Write operations are always opt-in via config.
- OAuth tokens are encrypted at rest (AES-256-GCM). No unencrypted fallback exists.
- OAuth scopes are hardcoded per service. Code paths that could escalate scope do not exist.
- Google Cloud projects backing each skill should only enable the APIs the skill uses.
- Every skill must pass clawhub moderation and VirusTotal scanning before publish.

## Publishing

Releases are automated via [release-please](https://github.com/googleapis/release-please). On each merge to `main`, release-please opens or updates a release PR per skill with a generated changelog. Merging that PR creates a GitHub release, tag, and triggers GoReleaser + `clawhub publish`.

Tags are scoped per skill: `google-workspace/v0.1.0`, not `v0.1.0`. GoReleaser config lives inside each skill directory.

Configuration lives in `.release-please-config.json` (package definitions) and `.release-please-manifest.json` (current versions).

### SKILL.md version management

Release-please updates the `version` field in each skill's `SKILL.md` via a `type: yaml` updater with `jsonpath: $.version` configured in `extra-files`. The `clawhub publish` step passes `--version` explicitly from the release tag rather than relying on clawhub's YAML parser (which cannot handle inline comments or quoted values in frontmatter).

### Release-please PR merging

Release-please bot PRs must be merged manually. GitHub does not trigger `pull_request` workflows for events created by other GitHub Actions, so the auto-merge workflow (`release-please-auto-merge.yml`) does not fire on bot-opened PRs. Dependabot auto-merge works because Dependabot is a first-party GitHub feature.

## Adding a new skill

1. Create a new directory at the repo root with its own `go.mod`, `main.go`, `cmd/`, and `internal/` packages.
2. Add a `SKILL.md` with YAML frontmatter (see `google-workspace/SKILL.md` for the schema).
3. Add a `.goreleaser.yml` for cross-compilation.
4. Include an MIT-0 `LICENSE` file (clawhub requirement).
5. Add a `gomod` entry for the new skill directory in `.github/dependabot.yml`.
6. Add the new skill path to `.release-please-config.json` and `.release-please-manifest.json`.

## Git conventions

- Never commit directly to `main`. Use feature branches and PRs.
- Commit messages: `fix:`, `feat:`, or `breaking:` prefix.
- Subject line under 50 characters, body uses bullet points with emojis.
- Branch protection requires the `gate` CI check to pass and signed commits.
- Dependabot PRs are auto-approved and squash-merged.
