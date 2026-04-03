.PHONY: all
all: help

default: help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Generate
.PHONY: gen-client
gen-client: ## gen client code of {svc}. example: make gen-client svc=ai
	@cd rpc_gen && cwgo client --type RPC --service ${svc} --module github.com/MoScenix/industrial-fault-tree-ai/rpc_gen -I ../idl --idl ../idl/${svc}.proto

.PHONY: gen-server
gen-server: ## gen service code of {svc}. example: make gen-server svc=ai
	@cd app/${svc} && cwgo server --type RPC --service ${svc} --module github.com/MoScenix/industrial-fault-tree-ai/app/${svc} --pass "-use github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/${svc}.proto
