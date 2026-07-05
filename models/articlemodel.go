package models

import (
	"time"
)

type ArticleModel struct {
	Id           int64
	Name         string
	Slug         *string
	MdContent    string
	HtmlContent  string
	AccessCode   string
	IsPublished  int64
	CreatedAt    time.Time
	LastEditedAt time.Time
	PublishedAt  *time.Time
}
