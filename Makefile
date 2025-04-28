.PHONY: all build clean install

BINARY_NAME=browsir
INSTALL_PATH=/usr/local/bin
CONFIG_PATH=/etc/browsir

all: build

dev:
	go run cmd/browsir/main.go

build:
	go build -o dist/$(BINARY_NAME) ./cmd/browsir

clean:
	go clean
	rm -f $(BINARY_NAME)

install: build
	sudo mkdir -p $(CONFIG_PATH)
	sudo ln -sf "$(PWD)/config.yml" "$(CONFIG_PATH)/config.yml" || true
	sudo ln -sf "$(PWD)/shortcuts" "$(CONFIG_PATH)/shortcuts" || true
	sudo ln -sf "$(PWD)/links" "$(CONFIG_PATH)/links" || true
	sudo ln -sf "$(PWD)/dist/$(BINARY_NAME)" "$(INSTALL_PATH)/$(BINARY_NAME)"
	@echo "Created symlink to $(BINARY_NAME) in $(INSTALL_PATH)"
	@echo "Created symlinks to config files in $(CONFIG_PATH)"
	@echo "You can now run 'browsir' from anywhere"
