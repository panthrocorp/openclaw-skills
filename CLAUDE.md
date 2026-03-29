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

- `cmd/` uses Cobra for CLI subcommands. Each service (gmail, calendar, contacts) has its own file registering commands onto `rootCmd`. Helper functions like `gmailClient()` handle config loading, token decryption, and API client creation inline.
- `internal/config/` holds a JSON config model (`config.json`) that gates which services are active and at what access level (e.g. `CalendarMode` with off/readonly/readwrite). Config drives OAuth scope selection via `Config.OAuthScopes()`.
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

Skills are published to clawhub via GoReleaser and GitHub Actions on tag push:

```bash
git tag -a google-workspace/v0.1.0 -m "Initial release"
git push origin google-workspace/v0.1.0
```

The release workflow builds binaries for linux/arm64, linux/amd64, and darwin/arm64, then runs `clawhub publish`.

Tags are scoped per skill: `google-workspace/v0.1.0`, not `v0.1.0`. GoReleaser config lives inside each skill directory.

## Adding a new skill

1. Create a new directory at the repo root with its own `go.mod`, `main.go`, `cmd/`, and `internal/` packages.
2. Add a `SKILL.md` with YAML frontmatter (see `google-workspace/SKILL.md` for the schema).
3. Add a `.goreleaser.yml` for cross-compilation.
4. Include an MIT-0 `LICENSE` file (clawhub requirement).
5. Add a `gomod` entry for the new skill directory in `.github/dependabot.yml`.

## Git conventions

- Never commit directly to `main`. Use feature branches and PRs.
- Commit messages: `fix:`, `feat:`, or `breaking:` prefix.
- Subject line under 50 characters, body uses bullet points with emojis.
- Branch protection requires the `gate` CI check to pass and signed commits.
- Dependabot PRs are auto-approved and squash-merged.
