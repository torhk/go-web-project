package routes

import (
	"net/http"
	"fmt"
	"os"
	"strings"
)

type Files struct {
	Title string
	Type string
	Name  []string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	dir := "./data"

	// ReadDir returns a slice of DirEntry
	folder, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("homeHandler")
	files := Files{
		Title: "Files",
		Type:   "txt",
		Name: []string{},
	}
	// Loop through entries
	for _, file := range folder {
		fmt.Println(file.Name())
		files.Name = append(files.Name, strings.TrimSuffix(file.Name(), ".txt"))
	}

	err = templates.ExecuteTemplate(w, "home.html", files)
	if err != nil {
		fmt.Println("homeHandler: nil")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
}