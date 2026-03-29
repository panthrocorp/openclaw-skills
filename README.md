# openclaw-skills

![Language](https://img.shields.io/badge/language-Go-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT--0-green)
![Platform](https://img.shields.io/badge/platform-OpenClaw-blueviolet)

Custom skills for [OpenClaw](https://openclaw.ai) agents, published to the [Claw Skills](https://clawskills.sh) registry.

## Skills

| Skill | Description | Status |
|-------|-------------|--------|
| [google-workspace](./google-workspace/) | Read-only Gmail and Contacts, configurable Calendar access | Published (`panthrocorp-google-workspace`) |

## Installation

Install any skill via clawhub:

```bash
clawhub install panthrocorp-google-workspace
```

## Security

Every skill in this repository is designed with a hostile-instance threat model in mind:

- OAuth tokens are encrypted at rest using AES-256-GCM with HKDF-SHA256 key derivation. No unencrypted fallback exists.
- Gmail and Contacts access is strictly read-only with no write code paths.
- Calendar write operations are opt-in and gated at both config and runtime level.
- Google Cloud projects are scoped to only the required APIs.

See each skill's README and [SECURITY.md](./SECURITY.md) for the full security policy.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on branching, commits, and adding new skills.

## Licence

This repository is released under the [MIT No Attribution](./LICENSE) licence.
