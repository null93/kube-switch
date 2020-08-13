package: ## Package into binary
	pyinstaller --clean --noconfirm --onefile src/kube-switch.py

clean: ## Remove generated files/folders
	rm -rf ./{dist,include,build} ./src/__pycache__

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
