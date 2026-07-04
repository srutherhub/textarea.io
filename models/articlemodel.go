package models

import (
	"strconv"
	"time"
)

type ArticleModel struct {
	Id           ArticleId
	Name         ArticleName
	Slug         ArticleSlug
	MdContent    ArticleMdContent
	HtmlContent  ArticleHtmlContent
	AccessCode   ArticleAccessCode
	IsPublished  ArticleisPublished
	CreatedAt    ArticleCreatedAt
	LastEditedAt ArticleLastEditedAt
}

type ArticleAccessCode string

type ArticleId int64

type ArticleName string

type ArticleSlug *string

type ArticleMdContent string

type ArticleHtmlContent string

type ArticleisPublished int64

type ArticleCreatedAt time.Time

type ArticleLastEditedAt time.Time

func ArticleIdStrToInt(id string) (ArticleId, error) {
	parsedId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		return 0, err
	}
	return ArticleId(parsedId), nil
}
