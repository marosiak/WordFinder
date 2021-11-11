MOCKERY_VERSION=mockery_2.9.4_Linux_x86_64

install_mockery:
ifeq (,$(wildcard /usr/local/bin/mockery))
ifeq ($(GOOS),darwin)
	brew install mockery
else
	mkdir tmp
	wget -O tmp/mockery.tar.gz https://github.com/vektra/mockery/releases/download/v2.9.4/$(MOCKERY_VERSION).tar.gz

	sudo tar -C tmp -xzf tmp/mockery.tar.gz
	sudo cp tmp/mockery /usr/local/bin/mockery
	rm -rf tmp
endif
endif

requirements: install_mockery
	go mod download

mock: requirements
	mockery --recursive --name="^.*?Database$$|^.*?Provider|^.*?Service$$"

build: mock
	go mod download && go build -o genius-cli cmd/cli/main.go && ./genius-cli

test:
	go test -v ./... -short
