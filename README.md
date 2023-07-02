# GRPC JWT Auth Example

Quick example of JWT Auth in GRPC. It oversimplifies a couple things. First, it stores the refresh token in the access 
token which would defeat the purpose of having a long lived and short lived token, and the hmac key, user, and a couple 
other things are hard coded. Also, when it refreshes the token, it doesn't write the new one to disk. However, this 
provided me with a quick experiment into using JWT in GRPC auth.

## Prerequisites

Installing with `go install` requires your GOBIN to be on your path for execution of the command later. For go and protobuf, snap can work on Ubuntu and Windows via WSL2, but it's worth looking into the official websites and considering the options.

```bash
sudo snap install go --classic
sudo snap install protobuf --classic
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go install github.com/go-task/task/v3/cmd/task@v3.27.1
```

Add to your .bashrc

```
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Build

```bash
task build
```

## Run

The token generator creates a new token or an expired one and writes it to standard out for viewing purposes as well as
to disk to simulate it being a token that the client has already stored in the past via some sort of login process. It
contains a refresh token that the server considers to be valid for creating a new token with a day till expiration.

```bash
task newtoken
```

or

```bash
task expiredtoken
```

Generate a cert:

```bash
task cert
```

Run server (blocking)

```bash
task serve
```

Run client and greet:

```bash
task greet
```

## Regenerate

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pb/helloworld.proto
```

## Troubleshooting

If using docker desktop and WSL2, you might have docker mounting a Windows directory which causes and error mentioning too many tail fields.

- Stop Docker Desktop
- Restart wsl (wsl --shutdown) from PowerShell
- Run snap commands
- Restart Docker Desktop

## License

Most of this is new, but I started with some work in the GRPC examples after modifying https://grpc.io/docs/languages/go/quickstart/ so that it can handle JWT auth.