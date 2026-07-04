package main

import (
	"fmt"
	"net/http"
	m "textarea/models"
	t "textarea/views"

	v "github.com/srutherhub/web-app/views"
)

func BaseViewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexProps := v.IndexProps{Title: "textarea"}
		v.Index(t.HomeView(), indexProps).Render(r.Context(), w)
	}
}

func CreateArticleViewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexProps := v.IndexProps{Title: "Create article"}
		v.Index(t.CreateArticleView(), indexProps).Render(r.Context(), w)
	}
}

func CreateArticleHandler(app *AppService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("article_name")

		id, accessCode, err := app.CreateArticle(m.ArticleName(name))

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		props := t.CreateArticleSuccessProps{Id: id, Name: m.ArticleName(name), AccessCode: accessCode}

		t.CreateArticleSuccess(props).Render(r.Context(), w)
	}
}

func EditArticleViewHandler(app *AppService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if id == "" {
			fmt.Println("/article/id/edit missing id")
		}

		article, err := app.GetArticleById(id)

		if err != nil {
			fmt.Println(err.Error())
		}

		props := t.ArticleEditProps{Article: article}

		indexProps := v.IndexProps{Title: "Editing " + string(article.Name)}
		v.Index(t.ArticleEditView(props), indexProps).Render(r.Context(), w)
	}
}

func SaveArticleHandler(app *AppService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		mdContent := r.FormValue("md")

		if id == "" {
			http.Error(w, "missing article id", http.StatusBadRequest)
			return
		}

		err := app.SaveArticle(id, mdContent)

		if err != nil {
			http.Error(w, "failed to save article: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
