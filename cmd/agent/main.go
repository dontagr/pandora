// Package main представляет собой основной пакет исполняемого агентского приложения Pandora.
// Этот пакет инициализирует приложение через Uber FX, управляет запуском и зависимостями с помощью DI.
package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/agent/bootstrap"
)

// main — это точка входа приложения, где инициализируется и запускается агентское приложение.
func main() {
	app := fx.New(CreateApp())

	startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}
}

// CreateApp возвращает fx.Option, которая конфигурирует приложение агента,
// инкапсулируя настройки конфигурации, зависимости и компоненты.
func CreateApp() fx.Option {
	return fx.Options(
		bootstrap.Config,    // Инициализация конфигурации приложения.
		bootstrap.Logger,    // Установка системы логирования.
		bootstrap.SqlLite,   // Настройка доступа к базе данных SQLite.
		bootstrap.Transport, // Настройка транспортного уровня.
		bootstrap.Service,   // Настройка службы приложения.
		bootstrap.CLI,       // Конфигурация интерфейса командной строки (CLI).
	)
}
