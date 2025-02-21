package main

import (
	"fmt"
	"html/template"
	"net/http"
	// shoelacetheme "github.com/chasgames/dev_golang_theme"
	shoelacetheme "../../theme"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// baseThemeData creates a default ThemeData struct with common settings
func baseThemeData(title, route, content string) shoelacetheme.ThemeData {
	return shoelacetheme.ThemeData{
		Title:        title,
		CurrentRoute: route,
		LogoText:     "MyApp",
		NavItems: []shoelacetheme.NavItem{
			{Route: "/", Icon: "house"},
			{Route: "/content/games", Icon: "hdd-stack"},
			{Route: "/content/profile", Icon: "person"},
		},
		InitialContent: template.HTML(content),
		HeadingFont:    shoelacetheme.Font{Family: "Inter", Weight: 900, Size: "1rem"},
		BodyFont:       shoelacetheme.Font{Family: "Inter", Weight: 400, Size: "1rem"},
	}
}

// renderContent handles both AJAX and full page requests
func renderContent(route, title, content string) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("X-Requested-With") == "XMLHttpRequest" {
			return c.HTML(http.StatusOK, content)
		}
		data := baseThemeData(title, route, content)
		return shoelacetheme.RenderTheme(c.Response().Writer, data)
	}
}

// renderErrorPage displays an error within the theme
func renderErrorPage(c echo.Context, code int, msg string) {
	content := fmt.Sprintf(`<div class="error-page"><h1>Error %d</h1><p>%s</p><p><a href="/">Home</a></p></div>`, code, msg)
	data := baseThemeData(fmt.Sprintf("Error %d - My Dynamic App", code), c.Request().URL.Path, content)
	c.Response().Status = code
	if err := shoelacetheme.RenderTheme(c.Response().Writer, data); err != nil {
		c.Logger().Errorf("Render error: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())

	// Define routes
	e.GET("/", renderContent("/", "My Dynamic App - Home", "<h1>Welcome</h1><p>This is the home page.</p>"))
	e.GET("/content/games", renderContent("/content/games", "My Dynamic App - Games", "<h1>Games Page</h1><p>Games content here.</p>"))
	e.GET("/content/profile", renderContent("/content/profile", "My Dynamic App - Profile", "<h1>Profile Page</h1><p>Profile content here.</p>"))

	// Custom error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == http.StatusNotFound && c.Request().Header.Get("X-Requested-With") == "XMLHttpRequest" {
				c.HTML(http.StatusNotFound, "<h1>404 Not Found</h1><p>Page not found.</p>")
				return
			}
			renderErrorPage(c, he.Code, fmt.Sprintf("%v", he.Message))
		} else {
			renderErrorPage(c, http.StatusInternalServerError, "Unexpected error")
		}
	}

	e.Logger.Fatal(e.Start(":8080"))
}