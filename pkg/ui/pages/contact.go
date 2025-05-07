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

func ContactUs(ctx echo.Context, form *forms.Contact) error {
	r := ui.NewRequest(ctx)
	r.Title = "Contact us"
	r.Metatags.Description = "Get in touch with us."

	g := Group{
		Iff(r.Htmx.Target != "contact", func() Node {
			return components.Message(
				"is-link",
				"",
				Group{
					P(Text("This is an example of a form with inline, server-side validation and HTMX-powered AJAX submissions without writing a single line of JavaScript.")),
					P(Text("Only the form below will update async upon submission.")),
				},
			)
		}),
		Iff(form.IsDone(), func() Node {
			return components.Message(
				"is-large is-success",
				"Thank you!",
				Text("No email was actually sent but this entire operation was handled server-side and degrades without JavaScript enabled."),
			)
		}),
		Iff(!form.IsDone(), func() Node {
			return form.Render(r)
		}),
	}

	return r.Render(layouts.Primary, g)
}
