build::
	go build -o alias-app ./cmd/cmd-app

fmt::
	go fmt ./...

run::
	go run ./cmd/cmd-app

vet::
	go vet ./...

clean::
	go clean
	rm -f myapp

mod::
	go mod tidy
