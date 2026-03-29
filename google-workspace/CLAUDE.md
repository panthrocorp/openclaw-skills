# google-workspace

## Overview

Read-only Gmail and Contacts skill with configurable Calendar access. Security-first design: write code paths for Gmail and Contacts do not exist in the codebase.

## Architecture

- `cmd/` — Cobra CLI. Service commands (`gmail`, `calendar`, `contacts`, `auth`, `config`) register on `rootCmd`. Helper functions `gmailClient()`, `calendarClient()`, `contactsClient()` handle config loading, scope validation, token decryption, and API client construction inline.
- `internal/config/config.go` — JSON config model with three fields: `Gmail bool`, `Contacts bool`, and `Calendar CalendarMode` (`off`, `readonly`, `readwrite`). `CalendarMode` is the only multi-state gate; Gmail and Contacts are simple on/off booleans. `Config.OAuthScopes()` derives the OAuth scope list from the current config; scopes expand as services are enabled. Changing any value requires re-authentication to issue a new token with the updated scopes.
- `internal/oauth/` — OAuth2 Desktop flow with a localhost redirect. `InteractiveLogin` prompts for the **code value only** (the part after `code=` and before `&scope=` in the redirect URL). The browser fails to load `http://localhost`, and the operator copies the code parameter value from the URL bar. This differs from zoho-mail, which takes the full redirect URL.
- `internal/google/` — Thin typed wrappers around Google API services. Each wrapper exposes only the operations the skill permits: `GmailClient` has no send/modify/delete; `ContactsClient` has no write operations.
- `internal/crypto/` — AES-256-GCM with HKDF-SHA256 key derivation. Wire format: `salt (16B) || nonce (12B) || ciphertext+tag`.

## Scope enforcement (three layers)

1. Code: write methods do not exist in `internal/google/`.
2. Config: Calendar write operations check `config.CalendarMode == "readwrite"` at runtime and return an error otherwise.
3. Google Cloud project: only Gmail API, Calendar API, and People API should be enabled, providing server-side enforcement.

## Required environment variables

| Variable | Purpose |
|----------|---------|
| `GOOGLE_WORKSPACE_TOKEN_KEY` | Passphrase for AES-256-GCM token encryption |
| `GOOGLE_CLIENT_ID` | Google OAuth2 client ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth2 client secret |

Config directory defaults to `~/.openclaw/credentials/google-workspace/`. Override with `GOOGLE_WORKSPACE_CONFIG_DIR` or `--config-dir`.

## OAuth flow note

`InteractiveLogin` expects only the `code` value from the redirect URL, not the full URL. After the browser fails to load `http://localhost`, the operator copies the value of the `code=` parameter only and pastes it into the terminal. See `internal/oauth/oauth.go` for the implementation.

## Calendar mode and scope changes

If `CalendarMode` changes (e.g. `readonly` to `readwrite`), the stored token was issued with the old scopes. The operator must run `google-workspace auth login` again to issue a new token with the expanded scopes. The binary does not detect scope mismatches at runtime.
