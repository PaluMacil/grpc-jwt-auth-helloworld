# GRPC JWT Auth Example

Quick example of JWT Auth in GRPC. It requires some touch-up.

## Run

```
go run .\certgen\main.go
go run .\tokengen\main.go
go run .\greeter_server\main.go
go run .\greeter_client\main.go
```