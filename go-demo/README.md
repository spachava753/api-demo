# golang api demo

This api demo implements the specified [API](https://github.com/spachava753/api-demo/issues/1)
using [golang](https://go.dev/) 1.19. It has the following features

1. Clean architecture
2. Mock tests for service
3. Mock tests for api
4. Swagger documentation
5. Utilizing `go generate` to automate generation mocks and swagger docs

## Getting Started

### Prerequisites

#### Golang

You must have Golang installed. You can follow instructions for downloading Golang here: https://go.dev/doc/install

Once you have golang installed, all you need to do is run these commands:

```shell
go generate && go run main.go
```

This will generate the mocks and swagger documentation and then start up the server.

### Tests

If you want to run the tests, you must have had run `go generate` beforehand at least once. To run the tests, simpy run

```shell
go test ./...
```

### Swagger

Swagger docs are generated in the `docs` folder, which only appear after running `go generate`. It utilizes [`github.com/swaggo/swag/cmd/swag`](github.com/swaggo/swag) command to generate the docs. 

### Mocks

mocks are generated in the `mocks` folder, which only appear after running `go generate`. It utilizes [`github.com/golang/mock/mockgen`](github.com/golang/mock) command to generate the docs.

### Go generate

If you want to take a look at how `go generate` is generating the swagger docs and mocks, take a look at [`gen.go`](gen.go). This file showcases how we utilize the power of `go generate` to automate code generation.  