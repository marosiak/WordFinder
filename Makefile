requirements:
	go mod download
	wget https://github.com/vektra/mockery/releases/download/v2.9.4/mockery_2.9.4_Linux_arm64.tar.gz
	tar -xf mockery_2.9.4_Linux_x86_64.tar.gz

mock: requirements
	./mockery --name=.*Provider --recursive

build: mock
	go mod download && go build -o genius-cli cmd/cli/main.go && ./genius-cli

fmt:
	go fmt ./...
	goimports -w $(FILES)

lint:
	golint $(PACKAGES)

test:
	go test -v ./... -short