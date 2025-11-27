package server

import (
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/handlers"
	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/argon2"
	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres"
	postgresRepo "github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres/repositories"
	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/redis"
	redisRepo "github.com/g-villarinho/oidc-server/internal/adapters/secondary/redis/repositories"
	"github.com/g-villarinho/oidc-server/internal/core/services"
	"github.com/g-villarinho/oidc-server/pkg/injector"
	"go.uber.org/dig"
)

func InitializeContainer() *dig.Container {
	container := dig.New()

	provideInfraDependencies(container)
	provideRepositories(container)
	provideCache(container)
	provideServices(container)
	provideHandlers(container)
	porviderCrypto(container)

	return container
}

func provideInfraDependencies(container *dig.Container) {
	injector.Provide(container, redis.NewRedisClient)
	injector.Provide(container, postgres.NewPoolConnection)
}

func provideRepositories(container *dig.Container) {
	injector.Provide(container, postgresRepo.NewClientRepository)
	injector.Provide(container, postgresRepo.NewUserRepository)
	injector.Provide(container, redisRepo.NewSessionRepository)
}

func provideCache(container *dig.Container) {
	injector.Provide(container, redis.NewCache)
}

func provideServices(container *dig.Container) {
	injector.Provide(container, services.NewAuthService)
	injector.Provide(container, services.NewClientService)
}

func provideHandlers(container *dig.Container) {
	injector.Provide(container, handlers.NewClientHandler)
	injector.Provide(container, handlers.NewAuthHandler)
	injector.Provide(container, handlers.NewCookieHandler)
}

func porviderCrypto(container *dig.Container) {
	injector.Provide(container, argon2.NewHasher)
}
