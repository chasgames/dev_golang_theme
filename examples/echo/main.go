package main

import (
    "html/template"
    "net/http"
    shoelacetheme "github.com/chasgames/dev_golang_theme" // Adjust to your actual import path
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func renderContent(route, title, content string) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Check if this is an AJAX request (client-side fetch)
        if c.Request().Header.Get("X-Requested-With") == "XMLHttpRequest" {
            // Return only the content fragment
            return c.HTML(http.StatusOK, content)
        }

        // Otherwise, render the full theme
        data := shoelacetheme.ThemeData{
            Title:        title,
            CurrentRoute: route,
            LogoText:     "MyApp",
            NavItems: []shoelacetheme.NavItem{
                {Route: "/", Icon: "house"},
                {Route: "/content/games", Icon: "hdd-stack"},
                {Route: "/content/profile", Icon: "person"},
                {Route: "/content/caaaahat", Icon: "chat"}, // Note: Kept your typo for consistency
            },
            InitialContent: template.HTML(content),
            HeadingFont:    shoelacetheme.Font{Family: "Inter", Weight: 900, Size: "1rem"},
            BodyFont:       shoelacetheme.Font{Family: "Inter", Weight: 400, Size: "1rem"},
        }

        // Echo requires manual rendering since RenderTheme writes to http.ResponseWriter
        w := c.Response().Writer
        if err := shoelacetheme.RenderTheme(w, data); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
        return nil
    }
}

func main() {
    // Initialize Echo instance
    e := echo.New()

    // Optional: Add middleware (e.g., logging, recovery)
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Root route
    e.GET("/", func(c echo.Context) error {
        if c.Request().URL.Path != "/" {
            return c.NoContent(http.StatusNotFound)
        }
        return renderContent("/", "My Dynamic App - Home", "<h1>Welcome</h1><p>This is the home page.</p>")(c)
    })

    // Content routes
    e.GET("/content/games", renderContent("/content/games", "My Dynamic App - Games", "<h1>Games Page</h1><p>Games content here.</p>"))
    e.GET("/content/profile", renderContent("/content/profile", "My Dynamic App - Profile", "<h1>Profile Page</h1><p>Profile content here.</p>"))
    e.GET("/content/caaaahat", renderContent("/content/caaaahat", "My Dynamic App - Chat", "<h1>Chat Page</h1><p>Chat content here.</p>")) // Adjusted route to match NavItem

    // Start server
    e.Logger.Fatal(e.Start(":8080"))
}
