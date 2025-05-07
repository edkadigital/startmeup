package pages

import (
	"github.com/edkadigital/startmeup/pkg/ui"
	"github.com/edkadigital/startmeup/pkg/ui/forms"
	"github.com/edkadigital/startmeup/pkg/ui/layouts"
	"github.com/labstack/echo/v4"
)

func UpdateCache(ctx echo.Context, form *forms.Cache) error {
	r := ui.NewRequest(ctx)
	r.Title = "Set a cache entry"

	return r.Render(layouts.Primary, form.Render(r))
}
