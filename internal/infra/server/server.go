package injections

import (
	adapters2 "github.com/brienze1/notes-api/internal/integration/adapters"
	"os"

	"github.com/gofiber/fiber/v3"
)

type server struct {
	usersController adapters2.UsersController
	notesController adapters2.NotesController
}

func (s server) Serve() error {
	app := fiber.New()
	s.route(app)
	return app.Listen(":" + os.Getenv("SERVER_PORT"))
}

func NewServer(usersController adapters2.UsersController, notesController adapters2.NotesController) adapters2.Server {
	return &server{
		usersController: usersController,
		notesController: notesController,
	}
}
