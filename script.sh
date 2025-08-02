#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<EOF
Usage: $0 RAW_PASSWORD PEPPER [EXISTING_SALT]

  RAW_PASSWORD   – your admin’s plaintext password  
  PEPPER         – your secret pepper (don’t commit this in repo!)  
  EXISTING_SALT  – optional 32-char hex salt; if omitted, a new one is generated

This outputs:

  export PASSWORD_SALT=...
  export PASSWORD_PEPPER=...
  export ADMIN_PASSWORD_HASH=...

which you can copy into your .env or deployment config.
EOF
  exit 1
}

if [[ $# -lt 2 || $# -gt 3 ]]; then
  usage
fi

RAW_PASS="$1"
PEPPER="$2"
SALT="${3:-}"

# generate a new salt if none provided
if [[ -z "$SALT" ]]; then
  SALT=$(openssl rand -hex 16)
fi

# compute HMAC-SHA256 of (raw||salt) with key=pepper
HASH=$(printf '%s' "${RAW_PASS}${SALT}" \
       | openssl dgst -sha256 -hmac "${PEPPER}" -hex \
       | sed 's/^.* //')

cat <<EOF
# Add these to your environment or .env:
export PASSWORD_SALT=${SALT}
export PASSWORD_PEPPER=${PEPPER}
export ADMIN_PASSWORD_HASH=${HASH}
EOF
