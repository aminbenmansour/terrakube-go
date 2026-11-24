#!/usr/bin/env python3
"""Mints an upstream Terrakube PAT so JSON:API responses can be captured.

Upstream accepts two kinds of bearer token: a Dex token, or a PAT (iss
"Terrakube", HS256 signed with the configured PatSecret, jti present in the pat
table). Dex static passwords carry no groups claim, and every organization read
requires the instance owner group, so the PAT path is the only way in without
an LDAP server.

Run after start.sh:  ./pat.py > token.txt; curl -H "Authorization: Bearer $(cat token.txt)" ...

Set PAT_GROUPS to a comma-separated list to mint a token for other groups, e.g.
PAT_GROUPS=capture-team to capture team-derived permissions.
"""

import base64
import hashlib
import hmac
import json
import os
import pathlib
import subprocess
import time
import uuid

ISSUER = "Terrakube"


def pat_secret() -> str:
    """The secret start.sh gave the API; ``.secrets.env`` or PAT_SECRET."""
    if secret := os.environ.get("PAT_SECRET"):
        return secret
    env = pathlib.Path(__file__).with_name(".secrets.env")
    for line in env.read_text().splitlines():
        key, _, value = line.partition("=")
        if key == "PAT_SECRET":
            return value
    raise SystemExit("PAT_SECRET is not set: run start.sh first")


def b64url(raw: bytes) -> str:
    return base64.urlsafe_b64encode(raw).rstrip(b"=").decode()


def main() -> None:
    jti = str(uuid.uuid4())
    subprocess.run(
        [
            "podman",
            "exec",
            "tk-upstream-postgres",
            "psql",
            "-p",
            "5433",
            "-U",
            "terrakube",
            "-d",
            "terrakubedb",
            "-q",
            "-c",
            f"INSERT INTO pat (id, days, description, deleted, source) "
            f"VALUES ('{jti}', 1, 'differential capture', false, 'API')",
        ],
        check=True,
        stdout=subprocess.DEVNULL,
    )

    groups = [group.strip() for group in os.environ.get("PAT_GROUPS", "TERRAKUBE_ADMIN").split(",") if group.strip()]
    now = int(time.time())
    claims = {
        "iss": ISSUER,
        "sub": "admin@example.com (Token)",
        "aud": [ISSUER],
        "jti": jti,
        "email": "admin@example.com",
        "email_verified": True,
        "name": "admin@example.com (Token)",
        "groups": groups,
        "iat": now,
        "exp": now + 24 * 3600,
    }
    # Upstream decodes PatSecret as base64url, not base64.
    secret = pat_secret()
    key = base64.urlsafe_b64decode(secret + "=" * (-len(secret) % 4))
    signing_input = ".".join(
        [
            b64url(json.dumps({"alg": "HS256"}).encode()),
            b64url(json.dumps(claims).encode()),
        ]
    )
    signing_input = signing_input.encode()
    signature = hmac.new(key, signing_input, hashlib.sha256).digest()
    print(signing_input.decode() + "." + b64url(signature))


if __name__ == "__main__":
    main()
