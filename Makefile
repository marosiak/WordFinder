ifeq (,$(wildcard ./mockery))
	go build -o="mockery" github.com/vektra/mockery/v2
endif

mock:
	./mockery --name=.*Provider --recursive

build:
	make mock
	go mod download && go build -o genius-cli cmd/cli/main.go && ./genius-cli