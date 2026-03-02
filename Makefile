SHELL := /usr/bin/env bash

.PHONY: fmt lint test smoke quality ci-status install-git-hooks

fmt:
	@./scripts/fmt.sh

lint:
	@./scripts/lint.sh

test:
	@./scripts/test.sh

smoke:
	@./scripts/smoke.sh

quality: fmt lint test smoke

ci-status:
	@./scripts/ci-status.sh --latest

install-git-hooks:
	@./scripts/install-git-hooks.sh
