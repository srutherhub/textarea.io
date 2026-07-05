package main

import (
	"fmt"
	"net/http"
	"strconv"
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

		id, accessCode, err := app.CreateArticle(name)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		props := t.CreateArticleSuccessProps{Id: id, Name: name, AccessCode: accessCode}

		t.CreateArticleSuccess(props).Render(r.Context(), w)
	}
}

func EditArticleViewHandler(app *AppService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		parsedId, ok := parseIDOrError(w, id)

		if !ok {
			return
		}

		article, err := app.GetArticleById(parsedId)

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

		parsedId, ok := parseIDOrError(w, id)

		if !ok {
			return
		}

		err := app.SaveArticle(parsedId, mdContent)

		if err != nil {
			http.Error(w, "failed to save article: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func PublishArticleHandler(app *AppService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		parsedId, ok := parseIDOrError(w, id)

		if !ok {
			return
		}

		err := app.PublishArticle(parsedId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
	}
}

func ArticleViewHandler(app *AppService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		if slug == "" {
			fmt.Println("no article")
		}

		article, err := app.GetArticleBySlug(slug)

		if err != nil {
			fmt.Println("no article with slug: " + slug)
			return
		}

		if article.IsPublished == 0 {
			fmt.Println("article is not published: " + slug)
			return
		}

		props := t.ArticleViewProps{Article: article}
		indexProps := v.IndexProps{Title: article.Name}
		v.Index(t.ArticleView(props), indexProps).Render(r.Context(), w)
	}
}

func parseIDOrError(w http.ResponseWriter, id string) (int64, bool) {
	if id == "" {
		http.Error(w, "missing article id", http.StatusBadRequest)
		return 0, false
	}

	res, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id %q: %v", id, err), http.StatusBadRequest)
		return 0, false
	}
	return res, true
}
