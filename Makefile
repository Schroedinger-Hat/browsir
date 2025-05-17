.PHONY: build clean install

BINARY_NAME=browsir
INSTALL_PATH=$(HOME)/bin
CONFIG_PATH=$(HOME)/.config/browsir

dev:
	go run cmd/browsir/main.go

build:
	go build -o dist/$(BINARY_NAME) ./cmd/browsir

clean:
	go clean
	rm -f $(BINARY_NAME)

install: build
	@echo "Installing Browsir..."
	@echo "This will create and install the necessary symlinks to run Browsir from anywhere."
	@echo "It will also create a config directory in $(CONFIG_PATH) if it doesn't exist."
	@echo "If this is not your first time using Browsir on this system, please run 'make update' instead, to avoid overwriting symlinks."
	@echo ""
	@read -p "Press ENTER to continue with installation or Ctrl+C to abort..." dummy
	mkdir -p $(CONFIG_PATH)
	cp "$(PWD)/config.example.yml" "$(CONFIG_PATH)/config.yml" || true
	cp "$(PWD)/shortcuts.example" "$(CONFIG_PATH)/shortcuts" || true
	cp "$(PWD)/links.example" "$(CONFIG_PATH)/links" || true
	sudo ln -sf "$(PWD)/dist/$(BINARY_NAME)" "$(INSTALL_PATH)/$(BINARY_NAME)"
	@echo "Created symlink to $(BINARY_NAME) in $(INSTALL_PATH)"
	@echo "Created symlinks to config files in $(CONFIG_PATH)"
	@echo "You can now run 'browsir' from anywhere"

update:
	@echo "Updating Browsir..."
	git pull origin main
	go build -o dist/$(BINARY_NAME) ./cmd/browsir
	sudo ln -sf "$(PWD)/dist/$(BINARY_NAME)" "$(INSTALL_PATH)/$(BINARY_NAME)"
	@echo "You can now run 'browsir' from anywhere"