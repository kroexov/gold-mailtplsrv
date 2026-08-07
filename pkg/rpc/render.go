package rpc

import (
	"context"

	"mailtplsrv/pkg/renderer"

	"github.com/vmkteam/zenrpc/v2"
)

type RenderService struct {
	r *renderer.Renderer
	zenrpc.Service
}

func NewRenderService(r *renderer.Renderer) RenderService {
	return RenderService{r: r}
}

// RenderTemplate renders template from tmpl folder.
func (r RenderService) RenderTemplate(ctx context.Context, templateName string, data any) (string, error) {
	html, err := r.r.Render(ctx, templateName, data)
	if err != nil {
		return "", newInternalError(err)
	}

	return html, nil
}

// List returns all available templates.
func (r RenderService) List() ([]string, error) {
	list, err := r.r.List()
	if err != nil {
		return nil, newInternalError(err)
	}

	return list, nil
}
