package api

import (
	"api-demo/domain"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/google/uuid"
	"log"
	"net/http"

	_ "api-demo/docs"
)

type CreateUserDto struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type UpdateUserDto struct {
	Name *string `json:"name,omitempty"`
	Role *string `json:"role,omitempty"`
}

type UserDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func userToDto(u domain.User) UserDto {
	return UserDto{
		u.Id.String(),
		u.Name,
		u.Role,
	}
}

type UserApi struct {
	service domain.UserService
}

func NewUserApi(service domain.UserService) (UserApi, error) {
	if service == nil {
		return UserApi{}, fmt.Errorf("cannot create user api")
	}
	return UserApi{service: service}, nil
}

// @Summary      Get a user by their ID
// @Description  Get a user by their ID. ID must be a valid UUID
// @ID           get-user-by-id
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  UserDto
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /users/{id} [get]
func (u *UserApi) getUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		if writeErr := c.Status(http.StatusBadRequest).SendString(err.Error()); writeErr != nil {
			log.Println("could not write to response body", err.Error())
		}
		return nil
	}
	user, err := u.service.GetUserById(c.UserContext(), parsedId)
	if err != nil {
		if writeErr := c.Status(http.StatusBadRequest).SendString(err.Error()); writeErr != nil {
			log.Println("could not write to response body", err.Error())
		}
		return nil
	}
	return dtoResponse(c, user)
}

// @Summary      Get a users filtering on properties
// @Description  Get a users filtering on properties
// @ID           get-users-by-property
// @Tags         users
// @Produce      json
// @Param        name   query      string  false  "User's name"
// @Param        role   query      string  false  "User's role"
// @Success      200  {object}  []UserDto
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /users [get]
func (u *UserApi) getUserByProperty(c *fiber.Ctx) error {
	var prop domain.UserProperties
	if name := c.Query("name"); name != "" {
		prop.Name = &name
	}
	if role := c.Query("role"); role != "" {
		prop.Role = &role
	}
	users, err := u.service.GetByProperty(c.UserContext(), &prop)
	if err != nil {
		if writeErr := c.Status(http.StatusInternalServerError).SendString(err.Error()); writeErr != nil {
			log.Println("could not write to response body", err.Error())
		}
		return nil
	}

	resp := make([]UserDto, 0, len(users))
	for _, user := range users {
		resp = append(resp, userToDto(user))
	}
	b, err := json.Marshal(resp)
	if err != nil {
		if writeErr := c.Status(http.StatusInternalServerError).SendString(err.Error()); writeErr != nil {
			log.Println("could not write to response body", err.Error())
		}
		return nil
	}
	if _, writeErr := c.Write(b); writeErr != nil {
		log.Println("could not write to response body", err.Error())
	}
	return nil
}

// @Summary      Create a user
// @Description  Create a user by passing their name and role
// @ID           create-user
// @Tags         users
// @Produce      json
// @Param        body   body    CreateUserDto  true  "User's name and role"
// @Success      200  {object}  UserDto
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /users [post]
func (u *UserApi) createUser(c *fiber.Ctx) error {
	var createDto CreateUserDto
	if err := json.Unmarshal(c.Body(), &createDto); err != nil {
		if writeErr := c.Status(http.StatusBadRequest).SendString(err.Error()); writeErr != nil {
			log.Println("could not write to response body", err.Error())
		}
		return nil
	}
	user, err := u.service.CreateUser(c.UserContext(), domain.NewUser(createDto.Name, createDto.Role))
	if err != nil {
		if writeErr := c.Status(http.StatusInternalServerError).SendString(err.Error()); writeErr != nil {
			log.Println("could not write to response body", err.Error())
		}
		return nil
	}
	return dtoResponse(c, user)
}

// @Summary      Update a user
// @Description  Update a user by passing their name and role. Name or role can be omitted, but not both
// @ID           update-user
// @Tags         users
// @Produce      json
// @Param        id   path    string  true  "User's ID"
// @Param        body body    UpdateUserDto  true  "User's updated name and role"
// @Success      200  {object}  UserDto
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /users/{id} [put]
func (u *UserApi) updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, parseErr := uuid.Parse(id)
	if parseErr != nil {
		if writeErr := c.Status(http.StatusBadRequest).SendString(parseErr.Error()); writeErr != nil {
			log.Println("could not write to response body", parseErr.Error())
		}
		return nil
	}

	var updateDto UpdateUserDto
	if err := json.Unmarshal(c.Body(), &updateDto); err != nil {
		if writeErr := c.Status(http.StatusBadRequest).SendString(err.Error()); writeErr != nil {
			log.Println("could not write to response body", err.Error())
		}
		return nil
	}
	var updateUser domain.UpdateUser
	if updateDto.Name != nil {
		updateUser.Name = updateDto.Name
	}
	if updateDto.Role != nil {
		updateUser.Role = updateDto.Role
	}

	user, err := u.service.UpdateUser(c.UserContext(), parsedId, updateUser)
	if err != nil {
		if writeErr := c.Status(http.StatusInternalServerError).SendString(err.Error()); writeErr != nil {
			log.Println("could not write to response body", err.Error())
		}
		return nil
	}
	return dtoResponse(c, user)
}

// @Summary      Delete a user
// @Description  Delete a user by passing their ID. ID must be valid
// @ID           delete-user
// @Tags         users
// @Produce      json
// @Param        id   path    string  true  "User's ID"
// @Success      200
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /users/{id} [delete]
func (u *UserApi) deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, parseErr := uuid.Parse(id)
	if parseErr != nil {
		if writeErr := c.Status(http.StatusBadRequest).SendString(parseErr.Error()); writeErr != nil {
			log.Println("could not write to response body", parseErr.Error())
		}
		return nil
	}

	if err := u.service.Delete(c.UserContext(), parsedId); err != nil {
		if writeErr := c.Status(http.StatusBadRequest).SendString(parseErr.Error()); writeErr != nil {
			log.Println("could not write to response body", parseErr.Error())
		}
		return nil
	}

	return nil
}

func dtoResponse(c *fiber.Ctx, user domain.User) error {
	dto := userToDto(user)
	b, err := json.Marshal(dto)
	if err != nil {
		if writeErr := c.Status(http.StatusInternalServerError).SendString(err.Error()); writeErr != nil {
			log.Println("could not write to response body", err.Error())
		}
		return nil
	}
	if _, writeErr := c.Write(b); writeErr != nil {
		log.Println("could not write to response body", err.Error())
	}
	return nil
}

// AddRoutes add routes to fiber.App
// @title User Api Demo
// @version 1.0
// @description This is an endpoint to get a user by their ID
// @termsOfService http://swagger.io/terms/
// @contact.name Shashank Pachava
// @contact.url https://github.com/spachava753/api-demo
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func (u *UserApi) AddRoutes(app *fiber.App) {
	app.Get("/users", func(c *fiber.Ctx) error {
		return u.getUserByProperty(c)
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return u.getUserById(c)
	})

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		return u.updateUser(c)
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		return u.createUser(c)
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		return u.deleteUser(c)
	})

	// swagger
	app.Get("/docs/*", swagger.HandlerDefault)

	// middleware
	app.Use(compress.New())
	app.Use(recover.New())
}
