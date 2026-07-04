package main

import (
	"database/sql"
	"os"
	"path/filepath"
	m "textarea/models"

	_ "modernc.org/sqlite"
)

type DBService struct {
	db *sql.DB
}

func InitDB() *DBService {
	location := os.Getenv("DB_LOCATION")

	if err := os.MkdirAll(filepath.Dir(location), 0777); err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite", location)

	if err != nil {
		panic(err)
	}

	db.Exec("PRAGMA foreign_keys = ON")
	db.Exec("PRAGMA journal_mode = WAL")

	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	service := &DBService{db: db}

	service.InitTables()

	return service
}

func (d *DBService) InitTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS articles (
		 id INTEGER PRIMARY KEY,
		 name TEXT NOT NULL DEFAULT '',
		 slug TEXT UNIQUE DEFAULT NULL,
		 md_content TEXT NOT NULL DEFAULT '',
		 html_content TEXT NOT NULL DEFAULT '',
		 access_code TEXT UNIQUE NOT NULL DEFAULT '',
		 published INTEGER NOT NULL DEFAULT 0,
		 created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		 last_edited DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS article_tags (
			article_id INTEGER NOT NULL REFERENCES articles(id),
			tag_id INTEGER NOT NULL REFERENCES tags(id),
			PRIMARY KEY (article_id,tag_id)
		)`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			panic(err)
		}
	}
}

func (db *DBService) CreateNewArticle(title m.ArticleName, accessCode m.ArticleAccessCode) (m.ArticleId, error) {
	result, err := db.db.Exec(`INSERT INTO articles (name, access_code) VALUES (?, ?)`, title, accessCode)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return m.ArticleId(id), nil
}

func (db *DBService) GetArticleById(id m.ArticleId) (m.ArticleModel, error) {
	var article m.ArticleModel
	err := db.db.QueryRow(`SELECT * FROM articles WHERE id = ?`, id).Scan(&article.Id, &article.Name, &article.Slug, &article.MdContent, &article.HtmlContent, &article.AccessCode, &article.IsPublished, &article.CreatedAt, &article.LastEditedAt)

	if err != nil {
		return m.ArticleModel{}, err
	}

	return article, nil

}

func (db *DBService) SaveArticle(id m.ArticleId, mdContent string) error {
	_, err := db.db.Exec(`UPDATE articles SET md_content = ?, last_edited = CURRENT_TIMESTAMP WHERE id = ?`, mdContent, id)

	if err != nil {
		return err
	}

	return nil
}
