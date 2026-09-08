.PHONY: test vet check

test:
	go -C skills/best-harness/scripts test ./...

vet:
	go -C skills/best-harness/scripts vet ./...

check: test vet
	go -C skills/best-harness/scripts run . check --staged
