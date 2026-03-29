---
name: Google Workspace
description: Read-only Gmail and Contacts access with configurable Calendar (readonly or readwrite) for OpenClaw agents
version: 0.2.0
author: panthrocorp
license: MIT-0
metadata:
  openclaw:
    requires:
      env:
        - GOOGLE_WORKSPACE_TOKEN_KEY
        - GOOGLE_CLIENT_ID
        - GOOGLE_CLIENT_SECRET
      bins:
        - google-workspace
    emoji: "📧"
    homepage: https://github.com/panthrocorp/openclaw-skills
    os: ["linux"]
---

# Google Workspace Skill

Access Gmail (read-only), Google Calendar (configurable), and Google Contacts (read-only).

## Important

- Gmail is strictly read-only. You cannot send, modify, or delete emails.
- Contacts is strictly read-only. You cannot create, modify, or delete contacts.
- Calendar access depends on the configured mode. Check with `google-workspace config show`.

## Check configuration

Before using any commands, verify what is enabled:

```
google-workspace config show
```

## Gmail commands

Search messages:
```
google-workspace gmail search --query "from:someone@example.com" --max-results 10
```

Read a message by ID:
```
google-workspace gmail read --id MESSAGE_ID
```

List labels:
```
google-workspace gmail labels
```

Search or read threads:
```
google-workspace gmail threads --query "subject:meeting"
google-workspace gmail threads --id THREAD_ID
```

## Calendar commands

List available calendars:
```
google-workspace calendar list
```

List upcoming events:
```
google-workspace calendar events --from 2026-03-29T00:00:00Z --to 2026-04-05T23:59:59Z
```

Get a specific event:
```
google-workspace calendar event --id EVENT_ID
```

Create an event (only if calendar mode is readwrite):
```
google-workspace calendar create --summary "Team sync" --start "2026-04-01T10:00:00Z" --end "2026-04-01T11:00:00Z"
```

Update an event (only if calendar mode is readwrite):
```
google-workspace calendar update --id EVENT_ID --summary "Updated title"
```

Delete an event (only if calendar mode is readwrite):
```
google-workspace calendar delete --id EVENT_ID
```

## Contacts commands

List contacts:
```
google-workspace contacts list --max-results 50
```

Search contacts:
```
google-workspace contacts search --query "John"
```

Get a specific contact:
```
google-workspace contacts get --id "people/c1234567890"
```

## Authentication status

Check if the token is valid:
```
google-workspace auth status
```

If the token has expired, ask the operator to re-authenticate by running `google-workspace auth login` on the host.

If authentication fails with `Error 400: policy_enforced`, the operator's Google account likely has Advanced Protection enabled. They will need to temporarily unenroll, complete the OAuth flow, then re-enroll. The refresh token persists across sessions.

## Output format

All commands output JSON by default. Use `--output text` for plain text where supported.
