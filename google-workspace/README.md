# google-workspace

![Language](https://img.shields.io/badge/language-Go-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT--0-green)
![Platform](https://img.shields.io/badge/platform-linux%2Farm64-lightgrey)

A custom OpenClaw skill providing read-only Gmail, Contacts, and Drive access, plus configurable Calendar, via the Google APIs. Built for environments where the agent instance is treated as potentially hostile.

## Why not use an existing skill?

Every existing Google-related skill on clawskills.sh requests broad read-write OAuth scopes and most are flagged suspicious by VirusTotal or OpenClaw moderation. This skill enforces strict scope boundaries at three levels: code (no write functions for Gmail/Contacts/Drive), config (Calendar write is opt-in), and Google Cloud project (only required APIs enabled).

## Services

| Service | OAuth Scope | Mode | Write code paths |
|---------|------------|------|-----------------|
| Gmail | `gmail.readonly` | Read-only | None |
| Calendar | `calendar.readonly` or `calendar.events` | Configurable: `off`, `readonly`, `readwrite` | Gated by config check |
| Contacts | `contacts.readonly` | Read-only | None |
| Drive | `drive.readonly` | Read-only | None |

## Installation

```bash
clawhub install panthrocorp-google-workspace
```

Or build from source:

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o google-workspace .
```

## Deployment guide

Follow these steps in order to deploy the skill to an OpenClaw instance.

### 1. Create a Google Cloud project

1. Create a new project in the [Google Cloud Console](https://console.cloud.google.com/)
2. Enable the following APIs only:
   - Gmail API
   - Google Calendar API
   - People API (Contacts)
   - Google Drive API
3. Configure the OAuth consent screen (External, but only used by the operator's own account)
4. Create an OAuth 2.0 client ID with application type **Desktop**
5. Note the **Client ID** and **Client Secret**

### 2. Create secrets in Bitwarden Secrets Manager

Create three secrets in Bitwarden SM (EU region, PanthroCorp project):

| Secret key | Value |
|------------|-------|
| `GOOGLE_WORKSPACE_TOKEN_KEY_OPENCLAW_AWS_SANDBOX_<INSTANCE>` | Random 64-character hex string (`openssl rand -hex 32`) |
| `GOOGLE_CLIENT_ID_OPENCLAW_AWS_SANDBOX_<INSTANCE>` | Client ID from step 1 |
| `GOOGLE_CLIENT_SECRET_OPENCLAW_AWS_SANDBOX_<INSTANCE>` | Client secret from step 1 |

### 3. Add Terraform configuration

In the `openclaw` repo, add the Bitwarden secret key references to the instance's tfvars file:

```hcl
bitwarden_secret_key_google_workspace_token_key = "GOOGLE_WORKSPACE_TOKEN_KEY_OPENCLAW_AWS_SANDBOX_<INSTANCE>"  # pragma: allowlist secret
bitwarden_secret_key_google_client_id           = "GOOGLE_CLIENT_ID_OPENCLAW_AWS_SANDBOX_<INSTANCE>"  # pragma: allowlist secret
bitwarden_secret_key_google_client_secret       = "GOOGLE_CLIENT_SECRET_OPENCLAW_AWS_SANDBOX_<INSTANCE>"  # pragma: allowlist secret
```

### 4. Deploy

Merge the Terraform changes to `main`. The deploy workflow will:
- Create SSM parameters for the three secrets
- Update the `.env` file on the instance with `GOOGLE_WORKSPACE_TOKEN_KEY`, `GOOGLE_CLIENT_ID`, and `GOOGLE_CLIENT_SECRET`
- Create the `config/credentials/google-workspace/` directory

### 5. Install the skill on the instance

SSH in via Twingate and install:

```bash
ssh -i ~/.ssh/openclaw-operator ubuntu@<domain>
sudo -u openclaw docker exec -it openclaw-gateway clawhub install panthrocorp-google-workspace
```

### 6. Configure scopes

```bash
sudo -u openclaw docker exec -it openclaw-gateway \
  google-workspace config set --gmail=true --calendar=readonly --contacts=true --drive=true
```

### 7. Authenticate with Google

```bash
sudo -u openclaw docker exec -it openclaw-gateway \
  google-workspace auth login
```

1. Copy the URL printed to the terminal
2. Open it in your local browser and authenticate with your Google account
3. After authorisation, the browser will redirect to `http://localhost` (which will fail to load). Copy only the `code` parameter value from the URL bar (the part after `code=` and before `&scope=`)
4. Paste only the code back into the terminal (not the full URL)

**Google Advanced Protection users:** If you see `Error 400: policy_enforced`, Advanced Protection blocks unverified third-party OAuth apps. Temporarily unenroll from Advanced Protection at `https://myaccount.google.com/advanced-protection`, complete the OAuth flow, then re-enroll. The refresh token persists on the EBS volume so you only need to do this once.

### 8. Verify

```bash
sudo -u openclaw docker exec -it openclaw-gateway google-workspace auth status
sudo -u openclaw docker exec -it openclaw-gateway google-workspace gmail labels
```

No container restart is needed. The token is persisted on the EBS volume and the binary reads it fresh on each invocation.

## Upgrading from a previous version

If you are upgrading from a version that did not include Drive support, you need to:

1. Enable the **Google Drive API** in your Google Cloud project at [console.cloud.google.com/apis/library/drive.googleapis.com](https://console.cloud.google.com/apis/library/drive.googleapis.com)
2. Re-authenticate to obtain a token with the new `drive.readonly` scope:
   ```bash
   google-workspace auth login
   ```

The new `drive` config field defaults to `true`. If you do not want Drive access, disable it before re-authenticating:

```bash
google-workspace config set --drive=false
```

## Prerequisites

- A Google Cloud project with Gmail API, Calendar API, People API, and Google Drive API enabled (see deployment guide above)
- An OAuth 2.0 "Desktop" client configured in that project
- Three environment variables on the host:
  - `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` from the OAuth client
  - `GOOGLE_WORKSPACE_TOKEN_KEY` for encrypting the stored OAuth token (random 64-char hex string)

## Usage

### Configure scopes

```bash
google-workspace config set --gmail=true --calendar=readonly --contacts=true --drive=true
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

### Drive (read-only)

```bash
google-workspace drive list --max-results 20
google-workspace drive list --query "name contains 'report'" --max-results 10
google-workspace drive get --id FILE_ID
google-workspace drive download --id FILE_ID
```

Google Docs are exported as plain text, Sheets as CSV, and Slides as plain text. All other files download as raw bytes.

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
- Drive has no create, update, or delete code paths. Only file listing, metadata retrieval, and content download/export are supported.
- Calendar write operations check `config.CalendarMode == "readwrite"` at runtime and return an error if the mode is `readonly`.
- The Google Cloud project should only have Gmail API, Calendar API, People API, and Google Drive API enabled, providing server-side scope enforcement.
- Token encryption uses a random salt per encryption, preventing identical tokens from producing identical ciphertext.

## Development

```bash
go build -o google-workspace .
go test ./...
go vet ./...
```

## License

[MIT No Attribution](./LICENSE)
