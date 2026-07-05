package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	m "textarea/models"

	"github.com/yuin/goldmark"
	"golang.org/x/crypto/bcrypt"
)

type AppService struct {
	query *DBService
}

func InitAppService(db *DBService) *AppService {
	return &AppService{query: db}
}

func (a *AppService) CreateArticle(title string) (int64, string, error) {
	if title == "" {
		title = "My article"
	}

	code := a.GenerateAccessCode()

	codeHash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)

	if err != nil {
		return 0, "", fmt.Errorf("failed to hash access code: %w", err)
	}

	articleId, err := a.query.CreateNewArticle(title, string(codeHash))

	if err != nil {
		return 0, "", fmt.Errorf("failed to create article in db: %w", err)
	}

	return articleId, code, nil
}

func (a *AppService) GenerateAccessCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return strings.ToLower(base64.URLEncoding.EncodeToString(b)[:8])
}

func (a *AppService) GetArticleById(id int64) (m.ArticleModel, error) {
	if id == 0 {
		return m.ArticleModel{}, errors.New("GetArticleById: id is required")
	}

	article, err := a.query.GetArticleById(id)

	if err != nil {
		return m.ArticleModel{}, errors.New("GetArticleById: " + err.Error())
	}

	return article, nil
}

func (a *AppService) GetArticleBySlug(slug string) (m.ArticleModel, error) {
	if slug == "" {
		return m.ArticleModel{}, errors.New("GetArticleBySlug: slug is required")
	}

	article, err := a.query.GetArticleBySlug(slug)

	if err != nil {
		return m.ArticleModel{}, err
	}

	return article, nil
}

func (a *AppService) SaveArticle(id int64, mdContent string) error {
	if id == 0 {
		return errors.New("SaveArticle: id is required")
	}

	err := a.query.SaveArticle(id, mdContent)

	if err != nil {
		return err
	}

	err = a.SaveMdToHTML(id, mdContent)

	if err != nil {
		return fmt.Errorf("PublishArticle: failed to render article to html %d"+err.Error(), id)
	}

	return nil
}

func (a *AppService) PublishArticle(id int64) error {
	article, err := a.GetArticleById(id)

	if article.Slug == nil {
		err = a.GenerateSlug(id, article.Name)
	}

	if err != nil {
		return fmt.Errorf("%s", "PublishArticle: failed to generate slug "+err.Error())
	}

	err = a.query.PublishArticle(id)

	if err != nil {
		return fmt.Errorf("PublishArticle: failed to publish %d"+err.Error(), id)
	}

	err = a.SaveMdToHTML(id, article.MdContent)

	if err != nil {
		return fmt.Errorf("PublishArticle: failed to render article to html %d"+err.Error(), id)
	}

	return nil
}

func (a *AppService) GenerateSlug(id int64, name string) error {
	slugInvalidChars := regexp.MustCompile(`[^a-z0-9]+`)
	slugTrimHyphens := regexp.MustCompile(`^-+|-+$`)
	s := strings.ToLower(name)
	s = slugInvalidChars.ReplaceAllString(s, "-")
	s = slugTrimHyphens.ReplaceAllString(s, "")

	err := a.query.GenerateSlug(id, s)

	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) SaveMdToHTML(id int64, content string) error {
	var buf bytes.Buffer
	contentBytes := []byte(content)

	if err := goldmark.Convert(contentBytes, &buf); err != nil {
		return err
	}

	a.query.SaveMdToHTML(id, buf.String())

	return nil
}
