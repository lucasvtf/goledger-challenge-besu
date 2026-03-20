BESU_DIR     := besu
CONTRACT_DIR := SimpleStorage
ENV_FILE     := $(CONTRACT_DIR)/.env
ENV_EXAMPLE  := $(CONTRACT_DIR)/.env.example

.PHONY: devnet stop-devnet deploy devnet-deploy

devnet:
	cd $(BESU_DIR) && ./startBesu.sh

stop-devnet:
	cd $(BESU_DIR) && ./stopBesu.sh

deploy:
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "No .env found — copying from .env.example. Edit $(ENV_FILE) to set a custom PRIVATE_KEY."; \
		cp $(ENV_EXAMPLE) $(ENV_FILE); \
	fi
	@set -a && . ./$(ENV_FILE) && set +a && \
		cd $(CONTRACT_DIR) && forge script script/SimpleStorage.s.sol:SimpleStorageScript \
		--rpc-url besu \
		--broadcast

devnet-deploy: devnet deploy
