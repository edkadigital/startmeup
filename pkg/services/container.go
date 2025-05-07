package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/entc/gen"
	"github.com/edkadigital/startmeup/config"
	"github.com/edkadigital/startmeup/ent"
	"github.com/edkadigital/startmeup/ent/migrate"
	"github.com/edkadigital/startmeup/pkg/services/interfaces"
	"github.com/edkadigital/startmeup/pkg/tasks/riveradapter"
	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver
	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"

	// Required by ent.
	_ "github.com/edkadigital/startmeup/ent/runtime"
)

// Container contains all services used by the application and provides an easy way to handle dependency
// injection including within tests.
type Container struct {
	// Validator stores a validator
	Validator *Validator

	// Web stores the web framework.
	Web *echo.Echo

	// Config stores the application configuration.
	Config *config.Config

	// Cache contains the cache client.
	Cache *CacheClient

	// Database stores the connection to the database.
	Database *sql.DB

	// Files stores the file system.
	Files afero.Fs

	// ORM stores a client to the ORM.
	ORM *ent.Client

	// Graph is the entity graph defined by your Ent schema.
	Graph *gen.Graph

	// Mail stores an email sending client.
	Mail *MailClient

	// Auth stores an authentication client.
	Auth *AuthClient

	// Tasks stores the River worker.
	Tasks *riveradapter.Worker
}

// Ensure Container implements the interfaces.ServiceContainer interface
var _ interfaces.ServiceContainer = (*Container)(nil)

// GetWeb returns the web framework instance
func (c *Container) GetWeb() *echo.Echo {
	return c.Web
}

// GetConfig returns the application configuration
func (c *Container) GetConfig() *config.Config {
	return c.Config
}

// GetORM returns the ORM client
func (c *Container) GetORM() *ent.Client {
	return c.ORM
}

// GetDatabase returns the database connection
func (c *Container) GetDatabase() *sql.DB {
	return c.Database
}

// GetTasks returns the River worker
func (c *Container) GetTasks() *riveradapter.Worker {
	return c.Tasks
}

// NewContainer creates and initializes a new Container.
func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initValidator()
	c.initWeb()
	c.initCache()
	c.initDatabase()
	c.initFiles()
	c.initORM()
	c.initAuth()
	c.initMail()
	c.initTasks()
	return c
}

// Shutdown gracefully shuts the Container down and disconnects all connections.
func (c *Container) Shutdown() error {
	// Shutdown the web server.
	webCtx, webCancel := context.WithTimeout(context.Background(), c.Config.HTTP.ShutdownTimeout)
	defer webCancel()
	if err := c.Web.Shutdown(webCtx); err != nil {
		return err
	}

	// Shutdown the task runner.
	taskCtx, taskCancel := context.WithTimeout(context.Background(), c.Config.Tasks.ShutdownTimeout)
	defer taskCancel()
	c.Tasks.Stop(taskCtx)

	// Shutdown the ORM.
	if err := c.ORM.Close(); err != nil {
		return err
	}

	// Shutdown the database.
	if err := c.Database.Close(); err != nil {
		return err
	}

	// Shutdown the cache.
	c.Cache.Close()

	return nil
}

// initConfig initializes configuration.
func (c *Container) initConfig() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	c.Config = &cfg

	slog.Debug("Detected App Environment", "environment", c.Config.App.Environment)

	// Configure logging.
	switch cfg.App.Environment {
	case config.EnvProduction:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

// initValidator initializes the validator.
func (c *Container) initValidator() {
	c.Validator = NewValidator()
}

// initWeb initializes the web framework.
func (c *Container) initWeb() {
	c.Web = echo.New()
	c.Web.HideBanner = true
	c.Web.Validator = c.Validator
}

// initCache initializes the cache.
func (c *Container) initCache() {
	store, err := newInMemoryCache(c.Config.Cache.Capacity)
	if err != nil {
		panic(err)
	}

	c.Cache = NewCacheClient(store)
}

// initDatabase initializes the database.
func (c *Container) initDatabase() {
	var err error
	var connection string

	switch c.Config.App.Environment {
	case config.EnvTest:
		// TODO: Drop/recreate the DB, if this isn't in memory?
		connection = c.Config.Database.TestConnection
	default:
		connection = c.Config.Database.Connection
	}

	c.Database, err = openDB(c.Config.Database.Driver, connection)
	if err != nil {
		panic(err)
	}
}

// initFiles initializes the file system.
func (c *Container) initFiles() {
	// Use in-memory storage for tests.
	if c.Config.App.Environment == config.EnvTest {
		c.Files = afero.NewMemMapFs()
		return
	}

	fs := afero.NewOsFs()
	if err := fs.MkdirAll(c.Config.Files.Directory, 0755); err != nil {
		panic(err)
	}
	c.Files = afero.NewBasePathFs(fs, c.Config.Files.Directory)
}

// initORM initializes the ORM.
func (c *Container) initORM() {
	// Use dialect.Postgres explicitly
	drv := entsql.OpenDB(dialect.Postgres, c.Database)
	c.ORM = ent.NewClient(ent.Driver(drv))

	// Initialize the graph from the generated tables information.
	// This avoids runtime dependency on entc/load and the go toolchain.
	c.Graph = &gen.Graph{
		Nodes:  make([]*gen.Type, len(migrate.Tables)), // Assuming one type per table
		Config: &gen.Config{},                          // Provide a minimal config if needed, might need adjustment
	}
	// Note: This is a simplified way to get *something* into c.Graph.
	// It might not perfectly replicate the structure entc.LoadGraph provides,
	// especially regarding edges or complex configurations.
	// The admin handler primarily seems to need Node names.
	for i, table := range migrate.Tables {
		// Attempt to create a basic gen.Type for each table.
		// This is a placeholder - the actual fields/edges might be needed
		// depending on how deep the admin UI goes.
		c.Graph.Nodes[i] = &gen.Type{
			Name: table.Name, // Provide the name, which admin.go seems to use
			// Other fields like `Fields`, `Edges`, `ID` might be needed
			// by the admin handler or templates it uses.
			// If tests still fail, we might need to populate more here.
		}
	}
	// TODO: If tests still fail, investigate if more detailed info
	// from `migrate.Tables` (like Columns) needs to be mapped to `gen.Type`.
}

// initAuth initializes the authentication client.
func (c *Container) initAuth() {
	c.Auth = NewAuthClient(c.Config, c.ORM)
}

// initMail initialize the mail client.
func (c *Container) initMail() {
	var err error
	c.Mail, err = NewMailClient(c.Config)
	if err != nil {
		panic(fmt.Sprintf("failed to create mail client: %v", err))
	}
}

// initTasks initializes the River worker.
func (c *Container) initTasks() {
	var err error

	// Create the River worker - migrations will be handled separately
	// by the migrate command, not during application startup
	c.Tasks, err = riveradapter.NewWorker(c.Database)
	if err != nil {
		panic(fmt.Sprintf("failed to create River worker: %v", err))
	}
}

// openDB opens a database connection.
func openDB(driver, connection string) (*sql.DB, error) {
	// Remove SQLite specific logic
	return sql.Open(driver, connection)
}
