#!/usr/bin/env bash
set -euo pipefail

payload="$(cat)"
cwd="$(jq -r '.cwd // ""' <<<"$payload")"
cmd="$(jq -r '.tool_input.command // ""' <<<"$payload")"

project_root="$(
  git -C "$cwd" rev-parse --show-toplevel 2>/dev/null || realpath -m "$cwd"
)"

allowed_root="$(realpath -m "$project_root/agent-work")"

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
    echo "Blocked: temporary outputs must be under project agent-work directory." >&2
    echo "Project root: $project_root" >&2
    echo "Allowed root: $allowed_root" >&2
    echo "Rejected path: $path" >&2
    exit 2
  fi
}

# For Go test caching
is_go_command() {
  [[ "$cmd" =~ (^|[[:space:];&|])go[[:space:]]+(test|tool|build|run|vet|generate|list)([[:space:]]|$) ]]
}

guard_go_temp_envs() {
  is_go_command || return 0

  for key in TMPDIR GOTMPDIR GOCACHE; do
    value="$(
      grep -oE "(^|[[:space:]])${key}=[^[:space:]]+" <<<"$cmd" \
        | head -1 \
        | sed -E "s/^[[:space:]]*${key}=//" \
        || true
    )"

    [[ -z "$value" ]] && continue
    block_if_outside_agents "$value"
  done
}

extract_go_test_compile_output() {
  awk '
    {
      for (i = 1; i <= NF; i++) {
        if ($i == "-o") {
          if (i + 1 <= NF) {
            print $(i + 1)
          }
          exit
        }
        if ($i ~ /^-o=.+/) {
          sub(/^-o=/, "", $i)
          print $i
          exit
        }
      }
    }
  ' <<<"$cmd"
}

guard_go_test_compile_output() {
  [[ "$cmd" =~ (^|[[:space:];&|])go[[:space:]]+test([[:space:]][^;&|]*)?[[:space:]]-c([[:space:]]|$) ]] || return 0

  local output_path
  output_path="$(extract_go_test_compile_output)"

  if [[ -z "$output_path" ]]; then
    echo "Blocked: go test -c must specify -o under project agent-work directory." >&2
    echo "Project root: $project_root" >&2
    echo "Allowed root: $allowed_root" >&2
    exit 2
  fi

  block_if_outside_agents "$output_path"
}

guard_go_temp_envs
guard_go_test_compile_output

# For command block
block_tmp_cache_args_for_command() {
  local command_name="$1"

  while read -r path; do
    [[ -z "$path" ]] && continue

    local lower_path
    lower_path="$(tr '[:upper:]' '[:lower:]' <<<"$path")"

    if [[ "$lower_path" == *tmp* || "$lower_path" == *cache* ]]; then
      echo "Blocked: ${command_name} arguments must not contain tmp or cache." >&2
      echo "Rejected argument: $path" >&2
      exit 2
    fi
  done < <(
    grep -oE "(^|[;&|[:space:]])${command_name}[[:space:]][^;&|]+" <<<"$cmd" \
      | sed -E "s/(^|[;&|[:space:]])${command_name}[[:space:]]+//" \
      | tr ' ' '\n' \
      | grep -v '^-'
  )
}

block_tmp_cache_args_for_command "mkdir"
block_tmp_cache_args_for_command "touch"
