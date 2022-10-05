package main

//go:generate rm -rf mock/domain
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=domain/user.go -destination=mock/domain/user.go
//go:generate go run github.com/swaggo/swag/cmd/swag@latest init
