#!/usr/bin/env bash
set -euo pipefail

NAME="${1:?image name is required}"
PLATFORM="${2:-linux/amd64}"
IMAGE="benzhi/${NAME}:latest"

docker buildx build \
  --platform "${PLATFORM}" \
  --file benzhi.Dockerfile \
  --tag "${IMAGE}" \
  --load \
  .
