package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsnider/go.jasonsnider.com/internal/types"
)

type DB struct {
	DB *pgxpool.Pool
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

func (db *DB) FetchArticleBySlug(slug string) (types.Article, error) {
	var article types.Article
	sql := "SELECT id, slug, title, description, keywords, body FROM articles WHERE slug=$1"
	err := db.DB.QueryRow(context.Background(), sql, slug).Scan(&article.ID, &article.Slug, &article.Title, &article.Description, &article.Keywords, &article.Body)
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
