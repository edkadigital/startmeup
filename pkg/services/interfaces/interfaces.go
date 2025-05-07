package interfaces

import (
	"database/sql"

	"github.com/edkadigital/startmeup/config"
	"github.com/edkadigital/startmeup/ent"
	"github.com/edkadigital/startmeup/pkg/tasks/riveradapter"
	"github.com/labstack/echo/v4"
)

// ServiceContainer defines the minimum interface that tasks need to access services
type ServiceContainer interface {
	// GetWeb returns the web framework instance
	GetWeb() *echo.Echo

	// GetConfig returns the application configuration
	GetConfig() *config.Config

	// GetORM returns the ORM client
	GetORM() *ent.Client

	// GetDatabase returns the database connection
	GetDatabase() *sql.DB

	// GetTasks returns the River worker
	GetTasks() *riveradapter.Worker
}
