#!/usr/bin/env bash
set -euo pipefail

ROOT="${ROOT:-}"
MOCKGEN_VERSION="${MOCKGEN_VERSION:-v0.6.0}"
USE_INSTALLED_MOCKGEN="${USE_INSTALLED_MOCKGEN:-false}"

find_module_root() {
	local current="$1"

	while [[ "$current" != "/" ]]; do
		if [[ -f "$current/go.mod" ]]; then
			printf '%s\n' "$current"
			return 0
		fi

		current="$(dirname "$current")"
	done

	if [[ -f "/go.mod" ]]; then
		printf '/\n'
		return 0
	fi

	printf 'go.mod not found from %s\n' "$1" >&2
	return 1
}

to_go_path() {
	local path="$1"
	local relative="${path#"$ROOT"/}"
	printf '%s\n' "$relative"
}

go_package() {
	local path="$1"
	sed -nE 's/^[[:space:]]*package[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)[[:space:]]*$/\1/p' "$path" | head -n 1
}

has_interface() {
	local path="$1"
	grep -Eq '^[[:space:]]*type[[:space:]]+[A-Za-z_][A-Za-z0-9_]*[[:space:]]+interface[[:space:]]*\{' "$path"
}

run_mockgen() {
	if [[ "$USE_INSTALLED_MOCKGEN" == "true" ]] && command -v mockgen >/dev/null 2>&1; then
		mockgen "$@"
		return
	fi

	go run "go.uber.org/mock/mockgen@$MOCKGEN_VERSION" "$@"
}

new_mock() {
	local source_path="$1"
	local destination_path="${source_path%.go}_mock.go"
	local package_name
	package_name="$(go_package "$source_path")"

	if [[ -z "$package_name" ]]; then
		printf 'package declaration not found in %s\n' "$source_path" >&2
		return 1
	fi

	local source
	local destination
	source="$(to_go_path "$source_path")"
	destination="$(to_go_path "$destination_path")"

	printf 'mockgen %s -> %s\n' "$source" "$destination"

	run_mockgen \
		"-source=$source" \
		"-destination=$destination" \
		"-package=$package_name"
}

if [[ -z "$ROOT" ]]; then
	SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
	ROOT="$(find_module_root "$SCRIPT_DIR")"
fi

ROOT="$(cd "$ROOT" && pwd)"
cd "$ROOT"

sources=()

if [[ -f "$ROOT/pkg/base/usecase.go" ]]; then
	sources+=("$ROOT/pkg/base/usecase.go")
fi

if [[ -d "$ROOT/apps" ]]; then
	while IFS= read -r -d '' file; do
		if has_interface "$file"; then
			sources+=("$file")
		fi
	done < <(find "$ROOT/apps" -path '*/internal/application/port/*.go' -type f ! -name '*_mock.go' ! -name '*_test.go' -print0 | sort -z)
fi

for source in "${sources[@]}"; do
	new_mock "$source"
done

printf 'Mocks generated: %s\n' "${#sources[@]}"
