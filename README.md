# openclaw-skills

![Language](https://img.shields.io/badge/language-Go-00ADD8?logo=go&logoColor=white)
![Language](https://img.shields.io/badge/language-Node.js-5FA04E?logo=nodedotjs&logoColor=white)
![License](https://img.shields.io/badge/license-MIT--0-green)
![Platform](https://img.shields.io/badge/platform-OpenClaw-blueviolet)

Custom skills for [OpenClaw](https://openclaw.ai) agents, published to the [Claw Skills](https://clawskills.sh) registry.

## Skills

| Skill | Description | Status |
|-------|-------------|--------|
| [google-workspace](./google-workspace/) | Gmail, Calendar, Contacts, Drive (with comments), Docs, and Sheets | Published (`panthrocorp-google-workspace`) |
| [zoho-mail](./zoho-mail/) | Full read/write Zoho Mail access (EU data centre) | Published (`panthrocorp-zoho-mail`) |
| [aws-s3](./aws-s3/) | Self-contained AWS S3 SDK bundle for gateway containers | Pre-release (`panthrocorp-aws-s3`) |

## Installation

Install any skill via clawhub:

```bash
clawhub install panthrocorp-google-workspace
clawhub install panthrocorp-zoho-mail
clawhub install panthrocorp-aws-s3
```

## Security

Every skill in this repository is designed with a hostile-instance threat model in mind:

- OAuth tokens are encrypted at rest using AES-256-GCM with HKDF-SHA256 key derivation. No unencrypted fallback exists.
- Gmail and Contacts access is strictly read-only with no write code paths.
- Calendar, Drive (comments), Docs, and Sheets write operations are opt-in and gated at both config and runtime level.
- Zoho Mail provides full read/write access, scoped to a dedicated agent mailbox.
- AWS S3 skill uses the default credential provider chain (IMDS). No static credentials are stored.
- Google Cloud and Zoho API projects are scoped to only the required APIs.

See each skill's README and [SECURITY.md](./SECURITY.md) for the full security policy.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on branching, commits, and adding new skills.

## Licence

This repository is released under the [MIT No Attribution](./LICENSE) licence.
