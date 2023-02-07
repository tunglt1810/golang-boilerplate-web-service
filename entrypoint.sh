#!/usr/bin/env bash

set -e

main() {
	# Run service.
	./app "$@"
}

main "$@"
