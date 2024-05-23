package routes

import (
	"net/http"

	"axiata_test/handler"
)

func Register() {
	handler := handler.NewHandler()
	// router for posts
	http.HandleFunc("/posts", handler.PostsHandler)
	http.HandleFunc("/posts/", handler.PostsByIDHandler)
	http.HandleFunc("/posts/tags", handler.PostByTagHandler)

	// router for account
	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)

	// router for tags
	http.HandleFunc("/tags", handler.TagsHandler)
	http.HandleFunc("/tags/", handler.TagsByIDHandler)
}