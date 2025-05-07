package pages

import (
	"github.com/edkadigital/startmeup/pkg/ui"
	"github.com/edkadigital/startmeup/pkg/ui/components"
	"github.com/edkadigital/startmeup/pkg/ui/forms"
	"github.com/edkadigital/startmeup/pkg/ui/layouts"
	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AddTask(ctx echo.Context, form *forms.Task) error {
	r := ui.NewRequest(ctx)
	r.Title = "Create a task"
	r.Metatags.Description = "Test creating a task to see how it works."

	g := Group{
		Iff(r.Htmx.Target != "task", func() Node {
			return components.Message(
				"is-link",
				"",
				Group{
					P(Raw("Submitting this form will create an <i>ExampleTask</i> in the task queue. After the specified delay, the message will be logged by the queue processor.")),
					P(Text("See pkg/tasks and the README for more information.")),
				})
		}),
		form.Render(r),
	}

	return r.Render(layouts.Primary, g)
}
