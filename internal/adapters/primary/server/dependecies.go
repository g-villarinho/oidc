package server

import (
	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres"
	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres/repositories"
	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/redis"
	"github.com/g-villarinho/oidc-server/pkg/injector"
	"go.uber.org/dig"
)

func InitializeContainer() *dig.Container {
	container := dig.New()

	provideInfraDependencies(container)
	provideRepositories(container)
	provideCache(container)

	return container
}

func provideInfraDependencies(container *dig.Container) {
	injector.Provide(container, redis.NewRedisClient)
	injector.Provide(container, postgres.NewPoolConnection)
}

func provideRepositories(container *dig.Container) {
	injector.Provide(container, repositories.NewClientRepository)
}

func provideCache(container *dig.Container) {
	injector.Provide(container, redis.NewCache)
}
