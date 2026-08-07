package renderer

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/vmkteam/embedlog"
)

// chdirRoot switches working directory to the repo root, so live mode can read tmpl/ and static/.
func chdirRoot(t *testing.T) {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to detect test file path")
	}

	t.Chdir(filepath.Join(filepath.Dir(file), "..", ".."))
}

func newTestRenderer(t *testing.T, isEmbed bool) *Renderer {
	t.Helper()
	chdirRoot(t)

	return New(Config{
		PublicWebPath: "http://localhost:8075/mail/",
		WebPath:       "/mail/",
		IsEmbed:       isEmbed,
	}, embedlog.NewDevLogger())
}

func TestRendererList(t *testing.T) {
	for _, isEmbed := range []bool{false, true} {
		r := newTestRenderer(t, isEmbed)

		list, err := r.List()
		if err != nil {
			t.Fatalf("isEmbed=%v: List failed: %v", isEmbed, err)
		}

		if !slices.Contains(list, "welcome") {
			t.Errorf("isEmbed=%v: welcome template not found in %v", isEmbed, list)
		}

		// _shared, fixture and partials must not be listed
		for _, name := range list {
			if strings.HasPrefix(name, partialPrefix) {
				t.Errorf("isEmbed=%v: partial %q must not be listed", isEmbed, name)
			}
		}
	}
}

func TestRendererRender(t *testing.T) {
	for _, isEmbed := range []bool{false, true} {
		r := newTestRenderer(t, isEmbed)

		html, err := r.Render(context.Background(), "welcome", map[string]any{
			"test":      "TEST HEADER",
			"discount":  50,
			"promoCode": "PROMO",
			"shopUrl":   "https://example.com/shop",
			"daysLeft":  3,
		})
		if err != nil {
			t.Fatalf("isEmbed=%v: Render failed: %v", isEmbed, err)
		}

		for _, want := range []string{"<!doctype html>", "TEST HEADER", "PROMO", "http://localhost:8075/mail/test.jpg"} {
			if !strings.Contains(strings.ToLower(html), strings.ToLower(want)) {
				t.Errorf("isEmbed=%v: rendered html does not contain %q", isEmbed, want)
			}
		}
	}

	// debug.html is written by live mode on failure only, keep the tree clean
	t.Cleanup(func() { _ = os.Remove(debugFilename) })
}

func TestRendererRenderUnknownTemplate(t *testing.T) {
	r := newTestRenderer(t, true)

	if _, err := r.Render(context.Background(), "no-such-template", nil); err == nil {
		t.Error("expected error for unknown template")
	}
}

func TestPluralize(t *testing.T) {
	const (
		one  = "товар"
		two  = "товара"
		five = "товаров"
	)

	tests := []struct {
		n    any
		want string
	}{
		{0, five},
		{1, one},
		{2, two},
		{4, two},
		{5, five},
		{11, five},
		{14, five},
		{21, one},
		{22, two},
		{43.0, two},
		{float64(101), one},
		{-2, two},
	}

	for _, tt := range tests {
		if got := pluralize(tt.n, one, two, five); got != tt.want {
			t.Errorf("pluralize(%v) = %q, want %q", tt.n, got, tt.want)
		}
	}
}
