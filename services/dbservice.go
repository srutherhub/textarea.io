package services

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Space struct {
	ID           string
	Name         string
	Content      string
	Key          string
	CreatedAt    string
	LastAccessed string
	Files        []File
}

type File struct {
	ID           string
	Path         string
	CreatedAt    string
	LastAccessed string
}

type DBService struct {
	db     *sql.DB
	Spaces *Spaces
	Files  *Files
}

func InitDB() *DBService {
	location := os.Getenv("DB_LOCATION")

	if err := os.MkdirAll(filepath.Dir(location), 0777); err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite", location)
	db.Exec("PRAGMA foreign_keys = ON")
	db.Exec("PRAGMA journal_mode = WAL")

	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	sp := initSpaceQueries(db)
	fs := initFileQueries(db)

	return &DBService{db: db, Spaces: sp, Files: fs}
}

func (d *DBService) InitTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS spaces (
			id           TEXT PRIMARY KEY,
			name         TEXT NOT NULL DEFAULT '',
			content      TEXT NOT NULL DEFAULT '',
			key          TEXT NOT NULL,
			created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_accessed DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS files (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			path          TEXT NOT NULL DEFAULT '',
			display_name  TEXT NOT NULL DEFAULT '',
			space_id      INTEGER NOT NULL,
			created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_accessed DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_files_space_id ON files(space_id)`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			panic(err)
		}
	}
}

type Spaces struct {
	db *sql.DB
}

func initSpaceQueries(db *sql.DB) *Spaces {
	return &Spaces{db: db}
}

func (sp *Spaces) CreateSpace(id, name, key string) error {
	_, err := sp.db.Exec(`
        INSERT INTO spaces (id,name, key) VALUES (?,?, ?)
    `, id, name, key)
	return err
}

func (sp *Spaces) DeleteSpace(id string) error {
	_, err := sp.db.Exec(`DELETE FROM spaces WHERE id = ?`, id)
	return err
}

func (sp *Spaces) GetByID(id string) (Space, error) {
	var s Space
	err := sp.db.QueryRow(`
        SELECT id, name, content, key, created_at, last_accessed 
        FROM spaces WHERE id = ?
    `, id).Scan(&s.ID, &s.Name, &s.Content, &s.Key, &s.CreatedAt, &s.LastAccessed)
	if err != nil {
		return s, err
	}

	rows, err := sp.db.Query(`
        SELECT id, path, created_at, last_accessed
        FROM files WHERE space_id = ?
    `, id)
	if err != nil {
		return s, err
	}
	defer rows.Close()

	for rows.Next() {
		var f File
		if err := rows.Scan(&f.ID, &f.Path, &f.CreatedAt, &f.LastAccessed); err != nil {
			return s, err
		}
		s.Files = append(s.Files, f)
	}

	return s, nil
}

func (sp *Spaces) UpdateContent(id, content string) error {
	_, err := sp.db.Exec(`UPDATE spaces SET content = ? WHERE id = ?`, content, id)

	if err != nil {
		return err
	}

	return nil
}

func (sp *Spaces) UpdateLastAccessed(id string) error {
	_, err := sp.db.Exec(`UPDATE spaces SET last_accessed = CURRENT_TIMESTAMP WHERE id = ?`, id)

	if err != nil {
		return err
	}

	return nil
}

func (sp *Spaces) Authenticate(id, key string) error {
	_, err := sp.db.Exec(`SELECT id FROM spaces WHERE id = ? AND key = ?`, id, key)

	if err != nil {
		return err
	}

	return nil
}

type Files struct {
	db *sql.DB
}

func initFileQueries(db *sql.DB) *Files {
	return &Files{db: db}
}

func (fs *Files) Save(id, filepath string) {
}
