.PHONY: generate install

install:
	go install -v github.com/nicksnyder/go-i18n/v2/goi18n@v2.2.1
	go install golang.org/x/tools/cmd/goimports@latest

generate:
	rm -rf internal/i18n/compiled
	mkdir -p internal/i18n/compiled
	goi18n merge -outdir ./internal/i18n/compiled ./internal/i18n/locales/translate.*.toml
	go generate internal/i18n/i18n.go
	goimports -w ./internal/i18n/msg_generated.go
