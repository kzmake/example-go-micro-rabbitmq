.DEFAULT_GOAL := help

.PHONY: proto
proto: ## protoファイルからgoファイルを生成します
	@for f in proto/*.proto; do \
		protoc \
		--proto_path=.:. \
		--proto_path=.:${GOPATH}/src \
		--go_out=paths=source_relative:. \
		--micro_out=paths=source_relative:. \
		$$f; \
		echo "generating $$f"; \
	done

.PHONY: __
__:
	@echo "\033[33m"
	@echo "kzmake/example-go-micro-rabbitmq"
	@echo "\033[0m"

.PHONY: help
help: __ ## ヘルプを表示します
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@cat $(MAKEFILE_LIST) \
	| grep -e "^[a-zA-Z_/\-]*: *.*## *" \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-24s\033[0m %s\n", $$1, $$2}'
