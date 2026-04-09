#!/usr/bin/env bash
set -euo pipefail

BINARY_NAME="sinister"
DIST_DIR="${DIST_DIR:-dist}"
CGO_VALUE="${CGO_ENABLED:-1}"

TARGETS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
)

usage() {
  cat <<'EOF'
Usage:
  ./build-release.sh [version]

Examples:
  ./build-release.sh v0.2.0
  VERSION=v0.2.0 ./build-release.sh

Notes:
  - Output tarballs are written to ./dist by default.
  - Set DIST_DIR to customize the output directory.
EOF
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command not found: $1" >&2
    exit 1
  fi
}

checksum_cmd() {
  if command -v sha256sum >/dev/null 2>&1; then
    echo "sha256sum"
    return
  fi
  if command -v shasum >/dev/null 2>&1; then
    echo "shasum -a 256"
    return
  fi
  echo ""
}

resolve_version() {
  local requested="${1:-${VERSION:-}}"
  if [[ -n "${requested}" ]]; then
    printf '%s\n' "${requested}"
    return
  fi

  if command -v git >/dev/null 2>&1; then
    local derived=""
    derived="$(git describe --tags --always --dirty 2>/dev/null || true)"
    if [[ -n "${derived}" ]]; then
      printf '%s\n' "${derived}"
      return
    fi
  fi

  echo "unversioned"
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

require_cmd go
require_cmd tar

VERSION_VALUE="$(resolve_version "${1:-}")"
mkdir -p "${DIST_DIR}"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

echo "building ${BINARY_NAME} release artifacts (version: ${VERSION_VALUE})"

for target in "${TARGETS[@]}"; do
  GOOS="${target%/*}"
  GOARCH="${target#*/}"
  ARTIFACT="${BINARY_NAME}_${VERSION_VALUE}_${GOOS}_${GOARCH}.tar.gz"

  STAGE_DIR="${TMP_DIR}/${BINARY_NAME}_${GOOS}_${GOARCH}"
  mkdir -p "${STAGE_DIR}"

  echo "- building ${GOOS}/${GOARCH} (CGO_ENABLED=${CGO_VALUE})"
  CGO_ENABLED="${CGO_VALUE}" GOOS="${GOOS}" GOARCH="${GOARCH}" \
    go build -trimpath -ldflags "-s -w -X main.VERSION=${VERSION_VALUE}" -o "${STAGE_DIR}/${BINARY_NAME}" .

  cp LICENSE "${STAGE_DIR}/LICENSE"
  cp README.md "${STAGE_DIR}/README.md"

  tar -C "${STAGE_DIR}" -czf "${DIST_DIR}/${ARTIFACT}" .
done

SUM_CMD="$(checksum_cmd)"
if [[ -n "${SUM_CMD}" ]]; then
  (
    cd "${DIST_DIR}"
    # shellcheck disable=SC2086
    ${SUM_CMD} *.tar.gz > checksums.txt
  )
  echo "wrote checksums: ${DIST_DIR}/checksums.txt"
else
  echo "warning: checksum tool not found (sha256sum/shasum), skipping checksums"
fi

echo "release artifacts written to ${DIST_DIR}"
