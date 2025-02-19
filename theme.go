package chastheme

import (
    "embed"
    "html/template"
    "net/http"
)

//go:embed templates/base.html
var templates embed.FS

// Font represents font settings for the theme.
type Font struct {
    Family string // e.g., "Inter", "Open Sans"
    Weight int    // e.g., 400, 900
    Size   string // e.g., "1rem", "16px"
}

// NavItem represents a single navigation item in the navbar.
type NavItem struct {
    Route string // e.g., "/", "/games"
    Icon  string // e.g., "house", "hdd-stack"
}

// ThemeData holds all dynamic data for rendering the template.
type ThemeData struct {
    Title          string        // Page title
    CurrentRoute   string        // Current active route for highlighting
    LogoURL        string        // URL for logo image (if used)
    LogoText       string        // Text for logo (if used)
    NavItems       []NavItem     // List of navigation items
    InitialContent template.HTML // Initial content for the content area
    HeadingFont    Font          // Font settings for headings
    BodyFont       Font          // Font settings for body text
}

// RenderTheme renders the Shoelace theme template with the provided data.
func RenderTheme(w http.ResponseWriter, data ThemeData) error {
    tmpl, err := template.ParseFS(templates, "templates/base.html")
    if err != nil {
        return err
    }
    return tmpl.Execute(w, data)
}