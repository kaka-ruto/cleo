SHELL := /usr/bin/env bash

.PHONY: fmt lint shellcheck test smoke clean quality ci-status install-git-hooks

fmt:
	@./scripts/fmt.sh

lint:
	@./scripts/lint.sh

shellcheck:
	@./scripts/shellcheck.sh

test:
	@./scripts/test.sh

smoke:
	@./scripts/smoke.sh

clean:
	@./scripts/clean.sh

quality: fmt lint shellcheck test smoke

ci-status:
	@./scripts/ci-status.sh --latest

install-git-hooks:
	@./scripts/install-git-hooks.sh
