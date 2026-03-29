# Contributing

Contributions are welcome. Please follow the guidelines below.

## Getting Started

1. Fork the repository.
2. Create a feature branch from `main` (`feat/` or `fix/` prefix).
3. Install pre-commit hooks: `pre-commit install && pre-commit install --hook-type pre-push`
4. Make your changes in the relevant skill directory.
5. Ensure checks pass from the skill directory:
   - `go vet ./...`
   - `go test ./...`
   - `golangci-lint run`
6. Submit a pull request.

## Commit Convention

All commit messages must use a semantic prefix:

- `fix:` for bug fixes
- `feat:` for new features
- `breaking:` for breaking changes

Subject line must be under 50 characters.

## Pull Requests

- All PRs must target `main` and require at least one approval before merge.
- PRs are squash-merged.
- Keep PRs focused. One logical change per PR.

## Adding a New Skill

Each skill is a self-contained directory with its own Go module. The expected structure:

```
<skill-name>/
  SKILL.md              # clawhub manifest (YAML frontmatter + agent instructions)
  LICENSE               # MIT-0 (clawhub requirement)
  go.mod / go.sum       # Go module
  main.go               # entrypoint
  cmd/                  # Cobra CLI commands
  internal/             # internal packages
  .goreleaser.yml       # cross-compilation config
  README.md             # skill-specific documentation
```

New skills must include:

- `SKILL.md` with clawhub manifest frontmatter (see `google-workspace/SKILL.md` for the schema)
- `LICENSE` (MIT-0, required by clawhub)
- `.goreleaser.yml` for cross-compilation
- Tests for all internal packages

## Security

- Never commit credentials, tokens, or secrets.
- OAuth tokens must be encrypted at rest. No unencrypted storage paths.
- Read-only access is the default. Write operations must be opt-in.
- Report vulnerabilities via [GitHub Security Advisories](https://github.com/panthrocorp/openclaw-skills/security), not public issues.
