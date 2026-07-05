package main

import (
	"database/sql"
	"fmt"
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
		 last_edited DATETIME DEFAULT CURRENT_TIMESTAMP,
		 published_at DATETIME DEFAULT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE NOT NULL DEFAULT '',
			parent_id INTEGER DEFAULT NULL,
		  FOREIGN KEY (parent_id) REFERENCES tags(id)
			)`,
		`CREATE TABLE IF NOT EXISTS article_tags (
			article_id INTEGER NOT NULL REFERENCES articles(id),
			tag_id INTEGER NOT NULL REFERENCES tags(id),
			PRIMARY KEY (article_id,tag_id)
		)`,
		`INSERT OR IGNORE INTO tags (name) VALUES ("news")`,
		`INSERT OR IGNORE INTO tags (name) VALUES ("technology")`,
		`INSERT OR IGNORE INTO tags (name) VALUES ("science")`,
		`INSERT OR IGNORE INTO tags (name) VALUES ("business")`,
		`INSERT OR IGNORE INTO tags (name) VALUES ("health")`,
		`INSERT OR IGNORE INTO tags (name) VALUES ("sports")`,
		`INSERT OR IGNORE INTO tags (name) VALUES ("finance")`,
		`INSERT OR IGNORE INTO tags (name) VALUES ("lifestyle")`,
		`INSERT OR IGNORE INTO tags (name) VALUES ("culture")`,
		`INSERT OR IGNORE INTO tags (name) VALUES ("education")`,
		`INSERT OR IGNORE INTO tags (name) VALUES ("other")`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('software', (SELECT id FROM tags WHERE name = 'technology'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('hardware', (SELECT id FROM tags WHERE name = 'technology'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('cars', (SELECT id FROM tags WHERE name = 'technology'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('startups', (SELECT id FROM tags WHERE name = 'business'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('industry', (SELECT id FROM tags WHERE name = 'business'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('travel', (SELECT id FROM tags WHERE name = 'lifestyle'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('food', (SELECT id FROM tags WHERE name = 'lifestyle'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('home', (SELECT id FROM tags WHERE name = 'lifestyle'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('research', (SELECT id FROM tags WHERE name = 'science'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('environment', (SELECT id FROM tags WHERE name = 'science'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('film', (SELECT id FROM tags WHERE name = 'culture'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('television', (SELECT id FROM tags WHERE name = 'culture'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('music', (SELECT id FROM tags WHERE name = 'culture'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('internet', (SELECT id FROM tags WHERE name = 'culture'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('art', (SELECT id FROM tags WHERE name = 'culture'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('gaming', (SELECT id FROM tags WHERE name = 'culture'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('fashion', (SELECT id FROM tags WHERE name = 'culture'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('design', (SELECT id FROM tags WHERE name = 'culture'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('us', (SELECT id FROM tags WHERE name = 'news'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('global', (SELECT id FROM tags WHERE name = 'news'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('fitness', (SELECT id FROM tags WHERE name = 'health'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('nutrition', (SELECT id FROM tags WHERE name = 'health'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('mental', (SELECT id FROM tags WHERE name = 'health'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('football', (SELECT id FROM tags WHERE name = 'sports'))`,
		`INSERT OR IGNORE INTO tags (name, parent_id) VALUES ('basketball', (SELECT id FROM tags WHERE name = 'sports'))`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			panic(err)
		}
	}
}

func (db *DBService) CreateNewArticle(title, accessCode string) (int64, error) {
	result, err := db.db.Exec(`INSERT INTO articles (name, access_code) VALUES (?, ?)`, title, accessCode)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db *DBService) GetArticleById(id int64) (m.ArticleModel, error) {
	var article m.ArticleModel
	err := db.db.QueryRow(`SELECT id, name, slug , md_content, html_content, access_code, published, created_at, last_edited, published_at FROM articles WHERE id = ?`, id).Scan(&article.Id, &article.Name, &article.Slug, &article.MdContent, &article.HtmlContent, &article.AccessCode, &article.IsPublished, &article.CreatedAt, &article.LastEditedAt, &article.PublishedAt)

	if err != nil {
		return m.ArticleModel{}, err
	}

	return article, nil
}

func (db *DBService) GetArticleBySlug(slug string) (m.ArticleModel, error) {
	var article m.ArticleModel
	err := db.db.QueryRow(`SELECT id, name, slug , md_content, html_content, access_code, published, created_at, last_edited, published_at FROM articles WHERE slug = ?`, slug).Scan(&article.Id, &article.Name, &article.Slug, &article.MdContent, &article.HtmlContent, &article.AccessCode, &article.IsPublished, &article.CreatedAt, &article.LastEditedAt, &article.PublishedAt)

	if err != nil {
		return m.ArticleModel{}, err
	}

	return article, nil
}

func (db *DBService) SaveArticle(id int64, mdContent string) error {
	_, err := db.db.Exec(`UPDATE articles SET md_content = ?, last_edited = CURRENT_TIMESTAMP WHERE id = ?`, mdContent, id)

	if err != nil {
		return err
	}

	return nil
}

func (db *DBService) PublishArticle(id int64) error {
	_, err := db.db.Exec(`UPDATE articles SET published = NOT published, published_at = CASE WHEN published_at IS NULL AND NOT published THEN CURRENT_TIMESTAMP ELSE published_at END WHERE id = ?`, id)

	if err != nil {
		return err
	}

	return nil
}

func (db *DBService) GenerateSlug(id int64, slug string) error {
	//Get articles with similar slugs and append count+1 to the new slug
	var count int
	err := db.db.QueryRow(`SELECT COUNT(*) FROM articles WHERE (slug = ? or slug GLOB ?) and id != ?`, slug, slug+"-[0-9]*", id).Scan(&count)

	if err != nil {
		return err
	}

	if count != 0 {
		slug = fmt.Sprintf("%s-%d", slug, count+1)
	}

	_, err = db.db.Exec(`UPDATE articles SET slug = ? WHERE id = ?`, slug, id)

	if err != nil {
		return err
	}

	return nil
}

func (db *DBService) SaveMdToHTML(id int64, htmlString string) error {
	_, err := db.db.Exec(`UPDATE articles SET html_content = ? WHERE id = ?`, htmlString, id)

	if err != nil {
		return err
	}

	return nil
}
