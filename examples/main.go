package main

import (
    "html/template"
    "net/http"
    shoelacetheme "github.com/chasgames/dev_golang_theme" // Replace with your actual import path
)

func renderContent(route, title, content string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        data := shoelacetheme.ThemeData{
            Title:        title,
            CurrentRoute: route,
            LogoText:     "MyApp",
            NavItems: []shoelacetheme.NavItem{
                {Route: "/content/home", Icon: "house"},
                {Route: "/content/games", Icon: "hdd-stack"},
                {Route: "/sffoffile", Icon: "person"},
                {Route: "/content/cfhat", Icon: "chat"},
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
    // Root handler (initial load)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r) // Handle non-root paths explicitly
            return
        }
        renderContent("/content/home", "My Dynamic App - Home", "<h1>Welcome</h1><p>This is the home page.</p>")(w, r)
    })

    // Content handlers
    http.HandleFunc("/content/home", renderContent("/content/home", "My Dynamic App - Home", "<h1>Home Page</h1><p>This is the home content.</p>"))
    http.HandleFunc("/content/games", renderContent("/content/games", "My Dynamic App - Games", "<h1>Games Page</h1><p>Games content here.</p>"))
    http.HandleFunc("/content/profile", renderContent("/content/profile", "My Dynamic App - Profile", "<h1>Profile Page</h1><p>Profile content here.</p>"))
    http.HandleFunc("/content/chat", renderContent("/content/chat", "My Dynamic App - Chat", "<h1>Chat Page</h1><p>Chat content here.</p>"))

    // Start server
    http.ListenAndServe(":8080", nil)
}