package main

import (
    "fmt"
    "net/http"
    "html/template"
    shoelacetheme "github.com/chasgames/dev_golang_theme" // Replace with your actual theme package
)

func main() {
    // Handler for initial page load
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := shoelacetheme.ThemeData{
            Title:        "My Dynamic App",
            CurrentRoute: "/",
            LogoText:     "MyApp",
            NavItems: []shoelacetheme.NavItem{
                {Route: "/content/home", Icon: "house"},
                {Route: "/content/games", Icon: "hdd-stack"},
                {Route: "/content/ok", Icon: "person"},
                {Route: "nonexistant", Icon: "chat"},
            },
            InitialContent: template.HTML("<h1>Welcome</h1><p>This is the home page.</p>"),
            HeadingFont:    shoelacetheme.Font{Family: "Inter", Weight: 900, Size: "1rem"},
            BodyFont:       shoelacetheme.Font{Family: "Inter", Weight: 400, Size: "1rem"},
        }
        if err := shoelacetheme.RenderTheme(w, data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })

    // Content fragment handlers
    http.HandleFunc("/content/home", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "<h1>Home Page</h1><p>This is the home content.</p>")
    })
    http.HandleFunc("/content/games", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "<h1>Games Page</h1><p>Games content here.</p>")
    })
    http.HandleFunc("/content/profile", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "<h1>Profile Page</h1><p>Profile content here.</p>")
    })
    http.HandleFunc("/content/chat", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "<h1>Chat Page</h1><p>Chat content here.</p>")
    })

    // Start server
    http.ListenAndServe(":8080", nil)
}