package chastheme

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"github.com/labstack/echo/v4"
)

//go:embed templates/base.html
var templates embed.FS

// Font represents font settings for the theme.
type Font struct {
	Family string // e.g., "Inter", "Open Sans"
	Weight int    // e.g., 400, 900
	Size   string // e.g., "1rem", "16px"
}

// ContentProvider defines how content is provided.
type ContentProvider interface {
	Render() (template.HTML, error)
}

// FileContent serves content from a file on the filesystem.
type FileContent string
func (f FileContent) Render() (template.HTML, error) {
	tmpl, err := template.ParseFiles(string(f))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// FuncContent serves content from a function.
type FuncContent func() template.HTML
func (f FuncContent) Render() (template.HTML, error) {
	return f(), nil
}

// NavItem represents a single navigation item or route.
type NavItem struct {
	Route   string         // e.g., "/", "/chat"
	Icon    string         // e.g., "house", "chat"
	Title   string         // e.g., "MyApp - Chat"
	Content ContentProvider // How to serve the content
}

// ThemeData holds all dynamic data for rendering the template.
type ThemeData struct {
	Title        string        // Default page title
	CurrentRoute string        // Current active route for highlighting
	LogoURL      string        // URL for logo image (if used, overrides LogoText)
	LogoText     string        // Text for logo (if used)
	NavItems     []NavItem     // List of navigation items and routes
	HeadingFont  Font          // Font settings for headings
	BodyFont     Font          // Font settings for body text
	Error404     NavItem       // Custom 404 handler
}

// RenderTheme renders the Chastheme template with the provided data.
func RenderTheme(w http.ResponseWriter, data ThemeData, content template.HTML) error {
	tmpl, err := template.ParseFS(templates, "templates/base.html")
	if err != nil {
		return err
	}
	renderData := struct {
		ThemeData
		InitialContent template.HTML
	}{
		ThemeData:      data,
		InitialContent: content,
	}
	return tmpl.Execute(w, renderData)
}

// RegisterRoutes registers all routes defined in NavItems with Echo.
func RegisterRoutes(e *echo.Echo, data ThemeData) {
	for _, item := range data.NavItems {
		e.GET(item.Route, func(c echo.Context) error {
			content, err := item.Content.Render()
			if err != nil {
				return c.HTML(http.StatusInternalServerError, "Error rendering content")
			}
			if c.Request().Header.Get("X-Requested-With") == "XMLHttpRequest" {
				return c.HTML(http.StatusOK, string(content))
			}
			data.CurrentRoute = item.Route
			data.Title = item.Title
			return RenderTheme(c.Response().Writer, data, content)
		})
	}

	// Custom 404 handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusNotFound {
			content, err := data.Error404.Content.Render()
			if err != nil {
				c.HTML(http.StatusInternalServerError, "Error rendering 404")
				return
			}
			data.CurrentRoute = c.Request().URL.Path
			data.Title = data.Error404.Title
			if c.Request().Header.Get("X-Requested-With") == "XMLHttpRequest" {
				c.HTML(http.StatusNotFound, string(content))
			} else {
				RenderTheme(c.Response().Writer, data, content)
			}
		}
	}
}