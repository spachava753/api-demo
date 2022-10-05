package cmd

import (
	"api-demo/api"
	"api-demo/domain"
	"api-demo/repo"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func getPort(portInput string) string {
	if portInput[0] != ':' {
		return ":" + portInput
	}
	return portInput
}

// Run is the entrypoint of the function
func Run() error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var port string
	flag.StringVar(&port, "port", ":3000", "Port to use")
	flag.Parse()

	// create repo
	userRepo := repo.NewInMemUserRepo()

	// create service
	service, err := domain.NewUserServiceImpl(&userRepo)
	if err != nil {
		return fmt.Errorf("could not create service: %w", err)
	}

	// create routes
	userApi, err := api.NewUserApi(&service)
	if err != nil {
		return fmt.Errorf("could not create api: %w", err)
	}

	app := fiber.New()

	userApi.AddRoutes(app)

	return app.Listen(getPort(port))
}
