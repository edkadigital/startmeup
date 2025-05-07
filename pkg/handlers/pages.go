package handlers

import (
	"fmt"

	"github.com/edkadigital/startmeup/pkg/pager"
	"github.com/edkadigital/startmeup/pkg/routenames"
	"github.com/edkadigital/startmeup/pkg/services"
	"github.com/edkadigital/startmeup/pkg/ui/models"
	"github.com/edkadigital/startmeup/pkg/ui/pages"
	"github.com/labstack/echo/v4"
)

type Pages struct{}

func init() {
	Register(new(Pages))
}

func (h *Pages) Init(c *services.Container) error {
	return nil
}

func (h *Pages) Routes(g *echo.Group) {
	g.GET("/", h.Home).Name = routenames.Home
	g.GET("/about", h.About).Name = routenames.About
}

func (h *Pages) Home(ctx echo.Context) error {
	pgr := pager.NewPager(ctx, 4)

	return pages.Home(ctx, &models.Posts{
		Posts: h.fetchPosts(&pgr),
		Pager: pgr,
	})
}

// fetchPosts is a mock example of fetching posts to illustrate how paging works.
func (h *Pages) fetchPosts(pager *pager.Pager) []models.Post {
	pager.SetItems(20)
	posts := make([]models.Post, 20)

	for k := range posts {
		posts[k] = models.Post{
			Title: fmt.Sprintf("Post example #%d", k+1),
			Body:  fmt.Sprintf("Lorem ipsum example #%d ddolor sit amet, consectetur adipiscing elit. Nam elementum vulputate tristique.", k+1),
		}
	}
	return posts[pager.GetOffset() : pager.GetOffset()+pager.ItemsPerPage]
}

func (h *Pages) About(ctx echo.Context) error {
	return pages.About(ctx)
}
