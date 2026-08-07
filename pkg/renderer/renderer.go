// Package renderer provides MJML template rendering functionality.
// It supports both embedded templates (for production) and live file system
// templates (for development). Templates can be organized in subdirectories.
package renderer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"mailtplsrv"

	"github.com/preslavrachev/gomjml/mjml"
	"github.com/vmkteam/embedlog"
)

const (
	TemplateDir    = "tmpl"       // Directory containing template files. Used in embed.FS
	StaticDir      = "static"     // Directory for static assets. Used in embed.FS
	templateSuffix = ".mjml"      // File extension for templates
	partialPrefix  = "partial-"   // Prefix for partial template files. Used only for skip in List()
	sharedDir      = "_shared"    // Directory for shared partials/styles
	fixtureDir     = "fixture"    // Directory for test fixtures
	debugFilename  = "debug.html" // File for debug output
)

// skipDirs contains directory names that are excluded from template listing.
var skipDirs = map[string]bool{
	sharedDir:  true,
	fixtureDir: true,
}

// skipSuffixes contains file name suffixes (without templateSuffix) that mark
// a file as a fragment of another template, not a standalone email.
var skipSuffixes = []string{"-body", "-styles"}

// Config holds configuration for the template renderer.
type Config struct {
	PublicWebPath string // Public web path for static assets
	WebPath       string // Local web path for static assets for local server
	IsEmbed       bool   // Used for Development mode
}

// Renderer handles MJML template parsing and rendering.
// It can operate in embedded mode (production) or live mode (development).
type Renderer struct {
	embedlog.Logger
	cfg Config

	parsedTemplates *template.Template // only for embed
}

// New creates a new Renderer instance with the given configuration and logger.
// It initializes templates in embed mode if configured.
func New(cfg Config, lg embedlog.Logger) *Renderer {
	r := &Renderer{
		cfg:    cfg,
		Logger: lg,
	}

	// init templates for embed mode
	t, err := r.parseTemplates(mailtplsrv.Templates)
	if err != nil {
		lg.Error(context.Background(), "failed to parse template", "err", err)
	} else {
		r.parsedTemplates = t
	}

	return r
}

// fsys returns the file system with templates: embedded one in production,
// the working directory in development.
func (r *Renderer) fsys() fs.FS {
	if r.cfg.IsEmbed {
		return mailtplsrv.Templates
	}

	return os.DirFS(".")
}

// parseTemplates walks TemplateDir in fsys and parses every .mjml file, including subdirectories.
// Template name is the file name only, so tmpl/welcome/welcome-body.mjml is referenced
// as {{ template "welcome-body.mjml" . }} regardless of its directory.
func (r *Renderer) parseTemplates(fsys fs.FS) (*template.Template, error) {
	var (
		t    = r.newTemplateWithFunc()
		seen = map[string]string{}
	)

	err := fs.WalkDir(fsys, TemplateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if !strings.HasSuffix(path, templateSuffix) {
			return nil
		}

		data, readErr := fs.ReadFile(fsys, path)
		if readErr != nil {
			return fmt.Errorf("read %s: %w", path, readErr)
		}

		// file names must be unique across directories, otherwise templates silently override each other
		name := filepath.Base(path)
		if prev, ok := seen[name]; ok {
			return fmt.Errorf("duplicate template name %q in %s and %s", name, prev, path)
		}
		seen[name] = path

		if _, parseErr := t.New(name).Parse(string(data)); parseErr != nil {
			return fmt.Errorf("parse %s: %w", path, parseErr)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}

// List returns the names of all available templates (without suffix), including subdirectories.
// Excludes partials, shared directory, fixture directory and `-body`/`-styles` fragments.
func (r *Renderer) List() ([]string, error) {
	var result []string

	err := fs.WalkDir(r.fsys(), TemplateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skip special directories
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return fs.SkipDir
			}

			return nil
		}

		name := d.Name()
		if !strings.HasSuffix(name, templateSuffix) || strings.HasPrefix(name, partialPrefix) {
			return nil
		}

		s, _ := strings.CutSuffix(name, templateSuffix)
		for _, suffix := range skipSuffixes {
			if strings.HasSuffix(s, suffix) {
				return nil
			}
		}

		result = append(result, s)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk templates failed: %w", err)
	}

	return result, nil
}

// Render renders the specified template with the given data.
// It automatically chooses between embed and live rendering based on configuration.
func (r *Renderer) Render(ctx context.Context, templateName string, data any) (string, error) {
	if r.cfg.IsEmbed {
		return r.RenderEmbed(ctx, templateName, data)
	}

	return r.RenderLive(ctx, templateName, data)
}

func (r *Renderer) newTemplateWithFunc() *template.Template {
	return template.New("").Funcs(template.FuncMap{
		"static": func() string {
			return r.cfg.PublicWebPath
		},
		"add":       func(a, b int) int { return a + b },
		"sub":       func(a, b int) int { return a - b },
		"pluralize": pluralize,
	})
}

// pluralize returns one of three forms based on Russian plural rules.
// Accepts any numeric type from JSON (float64) or Go (int).
// Example: pluralize(43, "товар", "товара", "товаров") = "товаров"
func pluralize(n any, one, two, five string) string {
	var v int
	switch x := n.(type) {
	case int:
		v = x
	case int64:
		v = int(x)
	case float64:
		v = int(x)
	case float32:
		v = int(x)
	}
	if v < 0 {
		v = -v
	}
	mod100 := v % 100
	if mod100 >= 5 && mod100 <= 20 {
		return five
	}
	switch mod100 % 10 {
	case 1:
		return one
	case 2, 3, 4:
		return two
	}
	return five
}

func (r *Renderer) render(ctx context.Context, t *template.Template, templateName string, data any) (string, error) {
	var b bytes.Buffer
	err := t.ExecuteTemplate(&b, templateName+templateSuffix, data)
	if err != nil {
		r.writeDebugFile(ctx, b.Bytes())
		return "", fmt.Errorf("execute template failed: %w", err)
	}

	html, err := mjml.Render(b.String())
	if err != nil {
		r.writeDebugFile(ctx, b.Bytes())
		return "", fmt.Errorf("render mjml file failed: %w", err)
	}

	return html, err
}

// RenderEmbed renders a template using embedded file system templates.
func (r *Renderer) RenderEmbed(ctx context.Context, templateName string, data any) (string, error) {
	if r.parsedTemplates == nil {
		return "", errors.New("failed to parse embed templates")
	}

	return r.render(ctx, r.parsedTemplates, templateName, data)
}

// RenderLive renders a template using live file system templates.
// Templates are re-read on every call, so edits are visible without restart.
func (r *Renderer) RenderLive(ctx context.Context, templateName string, data any) (string, error) {
	t, err := r.parseTemplates(os.DirFS("."))
	if err != nil {
		return "", err
	}

	return r.render(ctx, t, templateName, data)
}

// writeDebugFile writes debug output to a file when in development mode.
func (r *Renderer) writeDebugFile(ctx context.Context, result []byte) {
	if r.cfg.IsEmbed {
		return
	}

	err := os.WriteFile(debugFilename, result, os.ModePerm)
	r.PrintOrErr(ctx, "save debug output", err, "filename", debugFilename)
	if err != nil {
		r.Error(ctx, "debug output failed", "err", err)
	}
}
