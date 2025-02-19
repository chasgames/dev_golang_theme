package main

import (
    "fmt"
    "html/template"
    "net/http"
    shoelacetheme "github.com/chasgames/dev_golang_theme" // Replace with your actual import path
)

func renderContent(route, title, content string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check if this is an AJAX request (client-side fetch)
        if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
            // Return only the content fragment
            fmt.Fprint(w, content)
            return
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
                {Route: "/content/caaaahat", Icon: "chat"},
            },
            InitialContent: template.HTML(content),
            HeadingFont:    shoelacetheme.Font{Family: "Inter", Weight: 900, Size: "1rem"},
            BodyFont:       shoelacetheme.Font{Family: "Inter", Weight: 400, Size: "1rem"},
        }
        if err := shoelacetheme.RenderTheme(w, data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        renderContent("/content/home", "My Dynamic App - Home", "<h1>Welcome</h1><p>This is the home page.</p>")(w, r)
    })

    http.HandleFunc("/content/games", renderContent("/content/games", "My Dynamic App - Games", "<h1>Games Page</h1><p>Games content here.</p>"))
    http.HandleFunc("/content/profile", renderContent("/content/profile", "My Dynamic App - Profile", "<h1>Profile Page</h1><p>Profile content here.</p>"))
    http.HandleFunc("/content/chat", renderContent("/content/chat", "My Dynamic App - Chat", "<h1>Chat Page</h1><p>Chat content here.</p>"))

    http.ListenAndServe(":8080", nil)
}
