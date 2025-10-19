package main

import (
	"fmt"
	"jet_solver/frontend"
	"net/http"
)

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	// Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, frontend.SolverList())
	})
	http.HandleFunc("/solver", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		fmt.Fprint(w, frontend.SolverInput(name))
	})

	http.HandleFunc("/output", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprint(w, frontend.SolverOutput(r.FormValue("name"), r.FormValue("json")))
	})

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
