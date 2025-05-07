package forms

import (
	"net/http"

	"github.com/edkadigital/startmeup/pkg/form"
	"github.com/edkadigital/startmeup/pkg/routenames"
	"github.com/edkadigital/startmeup/pkg/ui"
	. "github.com/edkadigital/startmeup/pkg/ui/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type ForgotPassword struct {
	Email string `form:"email" validate:"required,email"`
	form.Submission
}

func (f *ForgotPassword) Render(r *ui.Request) Node {
	return Form(
		ID("forgot-password"),
		Method(http.MethodPost),
		HxBoost(),
		Action(r.Path(routenames.ForgotPasswordSubmit)),
		InputField(InputFieldParams{
			Form:      f,
			FormField: "Email",
			Name:      "email",
			InputType: "email",
			Label:     "Email address",
			Value:     f.Email,
		}),
		ControlGroup(
			FormButton("is-primary", "Reset password"),
			ButtonLink(r.Path(routenames.Home), "is-light", "Cancel"),
		),
		CSRF(r),
	)
}
