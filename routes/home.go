package routes

import (
	"os"
	"fmt"
	"strings"
	"net/http"
)

type Files struct {
	Title string
	Type string
	Name  []string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	folder, err := os.ReadDir("./data")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	files := Files{
		Title: "Files",
		Type:   "txt",
		Name: []string{},
	}
	
	for _, file := range folder {
		fmt.Println(file.Name())
		files.Name = append(files.Name, strings.TrimSuffix(file.Name(), ".txt"))
	}

	renderTemplate(w, "home", files)
}