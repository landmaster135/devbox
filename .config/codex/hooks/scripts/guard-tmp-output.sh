#!/usr/bin/env bash
set -euo pipefail

payload="$(cat)"
cwd="$(jq -r '.cwd // ""' <<<"$payload")"
cmd="$(jq -r '.tool_input.command // ""' <<<"$payload")"

project_root="$(
  git -C "$cwd" rev-parse --show-toplevel 2>/dev/null || realpath -m "$cwd"
)"

allowed_root="$(realpath -m "$project_root/.agents")"

is_under_allowed_root() {
  local path="$1"
  local resolved

  if [[ "$path" = /* ]]; then
    resolved="$(realpath -m "$path")"
  else
    resolved="$(realpath -m "$cwd/$path")"
  fi

  [[ "$resolved" == "$allowed_root" || "$resolved" == "$allowed_root/"* ]]
}

block_if_outside_agents() {
  local path="$1"

  if ! is_under_allowed_root "$path"; then
    echo "Blocked: temporary outputs must be under project .agents directory." >&2
    echo "Project root: $project_root" >&2
    echo "Allowed root: $allowed_root" >&2
    echo "Rejected path: $path" >&2
    exit 2
  fi
}

for key in TMPDIR GOTMPDIR GOCACHE; do
	value="$(
    grep -oE "(^|[[:space:]])${key}=[^[:space:]]+" <<<"$cmd" \
      | head -1 \
      | sed -E "s/^[[:space:]]*${key}=//" \
      || true
  )"
  if [[ -n "$value" ]]; then
    block_if_outside_agents "$value"
  fi
done

while read -r path; do
  [[ -z "$path" ]] && continue
  block_if_outside_agents "$path"
done < <(
  grep -oE 'mkdir[[:space:]]+(-p[[:space:]]+)?[^;&|]+' <<<"$cmd" \
    | sed -E 's/^mkdir[[:space:]]+(-p[[:space:]]+)?//' \
    | tr ' ' '\n' \
    | grep -v '^-'
)
