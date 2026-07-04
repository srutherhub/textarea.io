package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	m "textarea/models"

	"golang.org/x/crypto/bcrypt"
)

type AppService struct {
	query *DBService
}

func InitAppService(db *DBService) *AppService {
	return &AppService{query: db}
}

func (a *AppService) CreateArticle(title m.ArticleName) (m.ArticleId, m.ArticleAccessCode, error) {
	if title == "" {
		title = "My article"
	}

	code := a.GenerateAccessCode()

	codeHash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)

	if err != nil {
		return 0, "", fmt.Errorf("failed to hash access code: %w", err)
	}

	articleId, err := a.query.CreateNewArticle(title, m.ArticleAccessCode(codeHash))

	if err != nil {
		return 0, "", fmt.Errorf("failed to create article in db: %w", err)
	}

	return articleId, code, nil
}

func (a *AppService) GenerateAccessCode() m.ArticleAccessCode {
	b := make([]byte, 6)
	rand.Read(b)
	return m.ArticleAccessCode(strings.ToLower(base64.URLEncoding.EncodeToString(b)[:8]))
}

func (a *AppService) GetArticleById(id string) (m.ArticleModel, error) {
	parsedId, err := m.ArticleIdStrToInt(id)

	if err != nil {
		return m.ArticleModel{}, err
	}

	if parsedId == 0 {
		return m.ArticleModel{}, errors.New("GetArticleById: id is required")
	}

	article, err := a.query.GetArticleById(m.ArticleId(parsedId))

	if err != nil {
		return m.ArticleModel{}, errors.New("GetArticleById: " + err.Error())
	}

	return article, nil
}

func (a *AppService) SaveArticle(id, mdContent string) error {
	if id == "" {
		return errors.New("SaveArticle: id is required")
	}

	parsedId, err := m.ArticleIdStrToInt(id)

	if err != nil {
		return err
	}

	err = a.query.SaveArticle(parsedId, mdContent)

	if err != nil {
		return err
	}

	return nil
}
