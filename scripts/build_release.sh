#!/usr/bin/env bash
set -euo pipefail

APP_NAME=${APP_NAME:-verifier}
PKG=${PKG:-./cmd/verifier}
OUT_DIR=${OUT_DIR:-dist}

# Space-separated GOOS/GOARCH pairs (override via TARGETS env var, e.g. TARGETS="linux/amd64 windows/arm64")
TARGETS=${TARGETS:-"darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64"}

mkdir -p "${OUT_DIR}"

echo "Building release binaries into ${OUT_DIR}/" >&2
for target in ${TARGETS}; do
  IFS=/ read -r goos goarch <<<"${target}"
  output_path="${OUT_DIR}/${APP_NAME}-${goos}-${goarch}"
  [[ "${goos}" == "windows" ]] && output_path+=".exe"

  echo "→ ${goos}/${goarch}" >&2
  env CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" \
    go build -trimpath -ldflags="-s -w" -o "${output_path}" "${PKG}"
done

echo "Done." >&2
