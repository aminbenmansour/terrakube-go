# Security policy

## Reporting a vulnerability

Report privately through GitHub's
[security advisories](https://github.com/aminbenmansour/terrakube-go/security/advisories/new)
or by email to `security@aminbenmansour.dev`. Do not open a public issue for a
suspected vulnerability.

Include the version or commit, a reproduction, and the impact you expect. We
acknowledge reports within **3 working days** and aim to publish a fix or a
mitigation within **30 days** for high-severity issues; lower-severity reports
follow in the next regular release. We credit reporters in the advisory unless
asked not to.

## Supported versions

Only the latest minor release and `main` receive security fixes. Older
releases are upgraded, not patched.

## Scope

In scope: this repository's API (`/api/v1`, `/graphql/api/v1`), TFE cloud API
(`/remote/tfe/v2`), state backend (`/tfstate/v1`), executor, registry, and the
deployment assets under `deploy/`.

Out of scope: vulnerabilities in upstream Terrakube's Java implementation, in
Terraform/OpenTofu itself, in the React UI image (report those to their
projects), and findings that need a compromised host or database.

## Our own checks

- `govulncheck` on every push and pull request (CI, `make vuln`).
- GitGuardian secret scanning on every pull request.
- A CycloneDX SBOM per build, published as a CI artifact.
- `go test -race` plus the executor, integration and smoke suites.

The current threat model, including accepted risks with an owner and an expiry,
is [docs/threat-model.md](docs/threat-model.md).
