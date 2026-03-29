# openclaw-skills

![Language](https://img.shields.io/badge/language-Go-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT--0-green)
![Platform](https://img.shields.io/badge/platform-OpenClaw-blueviolet)

Custom skills for [OpenClaw](https://openclaw.ai) agents, published to the [Claw Skills](https://clawskills.sh) registry.

## Skills

| Skill | Description | Status |
|-------|-------------|--------|
| [google-workspace](./google-workspace/) | Read-only Gmail and Contacts, configurable Calendar access | In development |

## Installation

Install any skill via clawhub:

```bash
clawhub install panthrocorp/google-workspace
```

## Security

Every skill in this repository is designed with a hostile-instance threat model in mind:

- OAuth tokens are encrypted at rest using AES-256-GCM with keys stored in external secret managers
- Gmail and Contacts access is strictly read-only with no write code paths
- Calendar write operations are opt-in and gated at both config and runtime level
- Google Cloud projects are scoped to only the required APIs

See each skill's readme for its specific security posture.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Ensure `go test ./...` passes in the skill directory
4. Submit a pull request

## License

All skills are released under the [MIT No Attribution](./google-workspace/LICENSE) licence.
