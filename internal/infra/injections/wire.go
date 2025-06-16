package injections

import (
	"github.com/brienze1/notes-api/internal/domain/usecases"
	"github.com/brienze1/notes-api/internal/infra/connections/aws"
	"github.com/brienze1/notes-api/internal/infra/connections/database"
	injections2 "github.com/brienze1/notes-api/internal/infra/server"
	"github.com/brienze1/notes-api/internal/integration/adapters"
	controllers2 "github.com/brienze1/notes-api/internal/integration/entrypoint/controllers"
	"github.com/brienze1/notes-api/internal/integration/persistence/caches"
	repositories2 "github.com/brienze1/notes-api/internal/integration/persistence/repositories"
	"github.com/brienze1/notes-api/internal/integration/persistence/secrets"
	"github.com/brienze1/notes-api/internal/integration/queues"
	"gorm.io/gorm"
	"sync"
)

type wire struct {
	Db  *gorm.DB
	Sqs queues.SqsClient
}

var wireInit sync.Once
var wireInstance *wire

func Wire() *wire {
	if wireInstance == nil {
		wireInit.Do(
			func() {
				wireInstance = &wire{}
			},
		)
	}

	return wireInstance
}

func (w *wire) InitializeServer() (adapters.Server, error) {
	config := aws.NewAws()
	if w.Db == nil {
		cacheClientSet := database.NewCacheSet()
		cacheClientGet := database.NewCacheGet()
		cache := caches.NewCache(cacheClientGet, cacheClientSet)
		secretClient := aws.NewAwsSecretsManager(config)
		secret := secrets.NewSecret(secretClient)
		w.Db = database.NewDb(cache, secret)
	}
	notesRepository := repositories2.NewNote(w.Db)
	usersRepository := repositories2.NewUser(w.Db)
	if w.Sqs == nil {
		w.Sqs = aws.NewAwsSqs(config)
	}
	notesQueue := queues.NewNotesQueue(w.Sqs)
	errorHandler := controllers2.NewErrorHandler()
	usersController := controllers2.NewUsersController(
		usecases.CreateUserUseCase(usersRepository),
		usecases.DeleteUserUseCase(usersRepository),
		usecases.GetUserUseCase(usersRepository),
		usecases.ListUsersUseCase(usersRepository),
		errorHandler,
	)
	notesController := controllers2.NewNotesController(
		*usecases.NewCreateNoteUseCase(notesRepository, notesQueue),
		usecases.DeleteNoteUseCase(notesRepository),
		usecases.GetNoteUseCase(notesRepository),
		usecases.ListNotesUseCase(notesRepository),
		errorHandler,
	)
	server := injections2.NewServer(usersController, notesController)
	return server, nil
}
