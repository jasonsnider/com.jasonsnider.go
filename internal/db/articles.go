package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
	"github.com/jasonsnider/com.jasonsnider.go/pkg/inflection"
)

func (db *DB) CreateArticle(article types.Article) (string, error) {
	articleID := uuid.New().String()
	slug := inflection.Slugify(article.Title)

	sql := "INSERT INTO articles (id, title, slug) VALUES ($1, $2, $3)"
	_, err := db.DB.Exec(context.Background(), sql, articleID, article.Title, slug)
	if err != nil {
		return "", fmt.Errorf("query failed: %v", err)
	}

	return articleID, nil
}

func (db *DB) FetchArticles() ([]types.Article, error) {
	sql := "SELECT id, slug, title, description, keywords, body, type, format FROM articles"
	rows, err := db.DB.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var articles []types.Article
	for rows.Next() {
		var article types.Article
		err := rows.Scan(&article.ID, &article.Slug, &article.Title, &article.Description, &article.Keywords, &article.Body, &article.Type, &article.Format)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}
		articles = append(articles, article)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration failed: %v", rows.Err())
	}

	return articles, nil
}

func (db *DB) FetchArticlesByType(articleType string) ([]types.Article, error) {
	sql := "SELECT id, slug, title, description, keywords, body FROM articles WHERE type=$1"
	rows, err := db.DB.Query(context.Background(), sql, articleType)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var articles []types.Article
	for rows.Next() {
		var article types.Article
		err := rows.Scan(&article.ID, &article.Slug, &article.Title, &article.Description, &article.Keywords, &article.Body)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}
		articles = append(articles, article)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration failed: %v", rows.Err())
	}

	return articles, nil
}

func (db *DB) FetchArticleByID(id string) (types.Article, error) {
	var article types.Article
	sql := "SELECT id, title, slug, body, keywords, description, type, format, published FROM articles WHERE id=$1"
	err := db.DB.QueryRow(context.Background(), sql, id).Scan(&article.ID, &article.Title, &article.Slug, &article.Body, &article.Keywords, &article.Description, &article.Type, &article.Format, &article.Published)
	if err != nil {
		return article, fmt.Errorf("query failed: %v", err)
	}

	return article, nil
}

func (db *DB) FetchArticleBySlug(slug string) (types.Article, error) {
	var article types.Article
	sql := "SELECT id, title, slug, body, keywords, description FROM articles WHERE slug=$1"
	err := db.DB.QueryRow(context.Background(), sql, slug).Scan(&article.ID, &article.Title, &article.Slug, &article.Body, &article.Keywords, &article.Description)
	if err != nil {
		return article, fmt.Errorf("query failed: %v", err)
	}
	return article, nil
}

func (db *DB) FetchMetaDataBySlug(slug string) (types.Article, error) {
	var article types.Article
	sql := "SELECT id, slug, title, description, keywords FROM articles WHERE slug=$1"
	err := db.DB.QueryRow(context.Background(), sql, slug).Scan(&article.ID, &article.Slug, &article.Title, &article.Description, &article.Keywords)
	if err != nil {
		return article, fmt.Errorf("query failed: %v", err)
	}

	return article, nil
}

func (db *DB) DeleteArticle(id string) error {
	sql := "DELETE FROM articles WHERE id=$1"
	_, err := db.DB.Exec(context.Background(), sql, id)
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}

	return nil
}
