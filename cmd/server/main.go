// Package main представляет собой основную точку входа для серверной части приложения Pandora.
// Это приложение инициализируется и управляется через Uber FX, предоставляя полный набор сервисов и инфраструктуры.
package main

import (
	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/server/bootstrap"
)

// main — это точка входа приложения. Здесь инициализируется и запускается серверное приложение.
func main() {
	app := fx.New(CreateApp())

	app.Run()
}

// CreateApp создает и настраивает все компоненты серверного приложения используя fx.Options.
// Он возвращает fx.Option, которая включает все необходимые зависимости и компоненты.
func CreateApp() fx.Option {
	return fx.Options(
		bootstrap.Server,   // Основной серверный компонент.
		bootstrap.Config,   // Конфигурация приложения.
		bootstrap.Logger,   // Логирование для приложения.
		bootstrap.Postgres, // Настройка подключения к базе данных Postgres.
		bootstrap.Store,    // Слой доступа к данным для приложения.
		bootstrap.Route,    // Настройка маршрутизации HTTP-запросов для сервера.
		bootstrap.Service,  // Определение и регистрация бизнес-логаки приложения.
	)
}
