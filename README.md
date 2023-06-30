# GRPC JWT Auth Example

Quick example of JWT Auth in GRPC. It requires some touch-up. For instance, the server doesn't implement the refresh endpoint. I should also talk more here about what the generator scripts do and add a flag to generate an expired token so that I can demo the code path that refreshes the token.

## Prerequisites

```bash
sudo snap install go --classic
sudo snap install protobuf --classic
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Add to your .bashrc

```
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Run

```bash
go run .\certgen\main.go
go run .\tokengen\main.go
go run .\greeter_server\main.go
go run .\greeter_client\main.go
```

## Regenerate

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    helloworld/helloworld.proto
```

## Troubleshooting

If using docker desktop and WSL2, you might have docker mounting a Windows directory which causes and error mentioning too many tail fields.

- Stop Docker Desktop
- Restart wsl (wsl --shutdown) from PowerShell
- Run snap commands
- Restart Docker Desktop

## License

Most of this is new, but I started with some work in the GRPC examples after modifying https://grpc.io/docs/languages/go/quickstart/ so that it can handle JWT auth.