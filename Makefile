install_mockery:
ifeq (,$(wildcard bin/mockery))
	mkdir bin
	mkdir tmp
ifeq ($(shell uname -s),Linux)
	wget -O "tmp/mockery.tar.gz" "https://github.com/vektra/mockery/releases/download/v2.9.4/mockery_2.9.4_Linux_x86_64.tar.gz"
endif
ifeq ($(shell uname -s),Darwin)
	wget -O "tmp/mockery.tar.gz" "https://github.com/vektra/mockery/releases/download/v2.9.4/mockery_2.9.4_Darwin_x86_64.tar.gz"
endif
	sudo tar -C tmp -xzf tmp/mockery.tar.gz
	sudo cp tmp/mockery bin/mockery
	rm -rf tmp
endif

requirements: install_mockery
	go mod download

mock: requirements
	bin/mockery --recursive --name="^.*?Database$$|^.*?Provider|^.*?Service$$"

build: mock
	go mod download && go build -o genius-cli cmd/cli/main.go && ./genius-cli

test:
	go test -v ./... -short
