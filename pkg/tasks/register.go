package tasks

import (
	"context"

	"github.com/edkadigital/startmeup/pkg/log"
	"github.com/edkadigital/startmeup/pkg/routenames"
	"github.com/edkadigital/startmeup/pkg/services"
	"github.com/edkadigital/startmeup/pkg/tasks/riveradapter"
)

// Register registers all task workers with the task client.
func Register(c *services.Container) {
	// Register the example task worker
	err := c.Tasks.RegisterExampleTask(func(ctx context.Context, task riveradapter.ExampleTask) error {
		log.Default().Info("Example task received",
			"message", task.Message,
		)
		log.Default().Info("This can access the container for dependencies",
			"echo", c.Web.Reverse(routenames.Home),
		)
		return nil
	})

	if err != nil {
		panic(err)
	}
}
