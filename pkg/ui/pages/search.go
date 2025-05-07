package pages

import (
	"github.com/edkadigital/startmeup/pkg/ui"
	"github.com/edkadigital/startmeup/pkg/ui/layouts"
	"github.com/edkadigital/startmeup/pkg/ui/models"
	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
)

func SearchResults(ctx echo.Context, results []*models.SearchResult) error {
	r := ui.NewRequest(ctx)

	g := make(Group, len(results))
	for i, result := range results {
		g[i] = result.Render()
	}

	return r.Render(layouts.Primary, g)
}
