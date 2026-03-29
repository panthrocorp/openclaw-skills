# google-workspace

![Language](https://img.shields.io/badge/language-Go-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT--0-green)
![Platform](https://img.shields.io/badge/platform-linux%2Farm64-lightgrey)

A custom OpenClaw skill providing read-only Gmail, configurable Calendar, and read-only Contacts access via the Google APIs. Built for environments where the agent instance is treated as potentially hostile.

## Why not use an existing skill?

Every existing Google-related skill on clawskills.sh requests broad read-write OAuth scopes and most are flagged suspicious by VirusTotal or OpenClaw moderation. This skill enforces strict scope boundaries at three levels: code (no write functions for Gmail/Contacts), config (Calendar write is opt-in), and Google Cloud project (only required APIs enabled).

## Services

| Service | OAuth Scope | Mode | Write code paths |
|---------|------------|------|-----------------|
| Gmail | `gmail.readonly` | Read-only | None |
| Calendar | `calendar.readonly` or `calendar.events` | Configurable: `off`, `readonly`, `readwrite` | Gated by config check |
| Contacts | `contacts.readonly` | Read-only | None |

## Installation

```bash
clawhub install panthrocorp/google-workspace
```

Or build from source:

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o google-workspace .
```

## Prerequisites

- A Google Cloud project with Gmail API, Calendar API, and People API enabled
- An OAuth 2.0 "Desktop" client configured in that project
- Three environment variables on the host:
  - `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` from the OAuth client
  - `GOOGLE_WORKSPACE_TOKEN_KEY` for encrypting the stored OAuth token (random 64-char hex string)

## Usage

### Configure scopes

```bash
google-workspace config set --gmail=true --calendar=readonly --contacts=true
google-workspace config show
```

### Authenticate

```bash
google-workspace auth login
```

This prints a URL. Open it in your browser, authorise the requested scopes, then paste the authorisation code back into the terminal.

### Gmail (read-only)

```bash
google-workspace gmail search --query "from:boss@example.com" --max-results 5
google-workspace gmail read --id MESSAGE_ID
google-workspace gmail labels
google-workspace gmail threads --query "subject:meeting"
```

### Calendar

```bash
google-workspace calendar list
google-workspace calendar events --from 2026-04-01T00:00:00Z --to 2026-04-07T23:59:59Z
google-workspace calendar event --id EVENT_ID

# Only available in readwrite mode:
google-workspace calendar create --summary "Standup" --start "2026-04-01T09:00:00Z" --end "2026-04-01T09:15:00Z"
google-workspace calendar update --id EVENT_ID --summary "Updated title"
google-workspace calendar delete --id EVENT_ID
```

### Contacts (read-only)

```bash
google-workspace contacts list --max-results 50
google-workspace contacts search --query "John"
google-workspace contacts get --id "people/c1234567890"
```

### Check auth status

```bash
google-workspace auth status
```

## Token storage

OAuth tokens are encrypted at rest using AES-256-GCM with a key derived via HKDF-SHA256. The encryption key comes from the `GOOGLE_WORKSPACE_TOKEN_KEY` environment variable, which should be injected from an external secret manager (e.g. Bitwarden SM, AWS SSM Parameter Store). The binary refuses to store tokens if this variable is not set.

Default token location: `~/.openclaw/credentials/google-workspace/token.enc`

## Security

- Gmail has no send, modify, or delete code paths. The `internal/google/gmail.go` file only contains `messages.list`, `messages.get`, `labels.list`, and `threads` operations.
- Contacts has no create, update, or delete code paths.
- Calendar write operations check `config.CalendarMode == "readwrite"` at runtime and return an error if the mode is `readonly`.
- The Google Cloud project should only have Gmail API, Calendar API, and People API enabled, providing server-side scope enforcement.
- Token encryption uses a random salt per encryption, preventing identical tokens from producing identical ciphertext.

## Development

```bash
go build -o google-workspace .
go test ./...
go vet ./...
```

## License

[MIT No Attribution](./LICENSE)
