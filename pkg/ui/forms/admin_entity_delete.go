package forms

import (
	"net/http"

	"github.com/edkadigital/startmeup/pkg/routenames"
	"github.com/edkadigital/startmeup/pkg/ui"
	. "github.com/edkadigital/startmeup/pkg/ui/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AdminEntityDelete(r *ui.Request, entityTypeName string) Node {
	return Form(
		Method(http.MethodPost),
		P(
			Class("subtitle"),
			Textf("Are you sure you want to delete this %s?", entityTypeName),
		),
		ControlGroup(
			FormButton("is-link", "Delete"),
			ButtonLink(
				r.Path(routenames.AdminEntityList(entityTypeName)),
				"is-secondary",
				"Cancel",
			),
		),
		CSRF(r),
	)
}
