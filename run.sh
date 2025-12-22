#!/usr/bin/env bash

# script for Linux (WSL-compatible) dev workflow
# @Brandon Blanker Lim-it

set -euf -o pipefail

CLIENTPORT="${PORT:-3001}"
ADDRESS="${ADDRESS:-localhost}"

BROWSER="${BROWSER:-vivaldi}"
ISWSL=false
ISMAC=false
if [[ $(uname) == "Darwin" ]]; then
	ISMAC=true
elif grep -qi Microsoft /proc/version; then
	ISWSL=true
fi

client() {
	genall
	if "${ISWSL}"; then
		cmd.exe /c "start ${BROWSER} http://${ADDRESS}:7331/youth-convention-2025/"
	elif "${ISMAC}"; then
		open -a ${BROWSER} "http://${ADDRESS}:7331/youth-convention-2025/"
	fi
	templ generate --watch --proxy="http://localhost:${CLIENTPORT}" --open-browser=false &
	air -c ".air.client.toml" serve_client -p=":${CLIENTPORT}" -d=true -a="${ADDRESS}" -s=false
}

prod() {
	genall
	templ generate
	air -c ".air.client.toml" serve_client -p=":${CLIENTPORT}" -d=true -a="${ADDRESS}" -t=true
}

customrun() {
	gentempl
	go build -o ./tmp/main .
	./tmp/main
}

deps() {
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/air-verse/air@latest
	go install go.uber.org/nilaway/cmd/nilaway@latest
	go install github.com/kisielk/errcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

	# Tailwind https://tailwindcss.com/docs/installation/tailwind-cli
	if [[ ! -f tailwindcss ]]; then
		local BIN
		if "${ISWSL}"; then
			BIN="tailwindcss-linux-x64"
		elif "${ISMAC}"; then
			BIN="tailwindcss-macos-arm64"
		else
			BIN="tailwindcss-linux-x64"
		fi

		curl -LO "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/${BIN}"
		chmod +x "${BIN}"
		mv "${BIN}" tailwindcss
	else
		echo "You already have tailwindcss binary"
	fi
}

gentempl() {
	./tailwindcss -m -i client/static/css/style.css -o client/static/css/tailwind.css -m
	templ generate templ -v
}

genall() {
	go generate ./...
	gentempl
}

check() {
	go mod tidy
	templ fmt ./client/components

	set +f
	local gofiles=( internal/**/*.go cmd/*.go client/*.go client/**/*.go )
	for file in "${gofiles[@]}"; do
		if [[ ! $file == *_templ.go ]]; then
			goimports -w -local -v "$file"
		fi
	done
	set -f

	go vet ./...
	nilaway ./...
	errcheck ./...
	govulncheck ./...
}

if [ "$#" -eq 0 ]; then
	echo "First use: chmod +x ${0}"
	echo "Usage: ${0}"
	echo "Commands:"
	echo "    check"
	echo "    client"
	echo "    deps"
	echo "    genall"
	echo "    gentempl"
else
	echo "Running ${1}"
	"$1" "$@"
fi
