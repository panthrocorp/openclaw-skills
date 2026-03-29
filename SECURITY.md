# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in this project, please report it through GitHub's private vulnerability reporting.

**Do not open a public issue.**

1. Navigate to the [Security tab](https://github.com/panthrocorp/openclaw-skills/security) of this repository.
2. Click **Report a vulnerability**.
3. Provide a clear description of the issue, steps to reproduce, and any relevant context.

We aim to acknowledge receipt within 48 hours and provide a fix or mitigation within 7 days for critical issues. For non-critical vulnerabilities, we aim to address them within 30 days.

## Scope

This policy covers all skills published from this repository. Each skill has its own security posture documented in its subdirectory.

## Security Design Principles

All skills in this repository follow these principles:

- Read-only operations are the default. Write operations are opt-in via configuration.
- OAuth tokens are encrypted at rest using AES-256-GCM with HKDF-SHA256 key derivation. No unencrypted fallback exists.
- OAuth scopes are hardcoded per service. No code path exists that could escalate scope at runtime.
- Google Cloud projects backing each skill enable only the APIs the skill requires.
- All published binaries pass clawhub moderation and VirusTotal scanning.

## Supported Versions

Only the latest published version of each skill receives security updates. Prior versions are not patched.
