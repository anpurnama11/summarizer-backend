package repository

import (
	"context"
)

type HistoryRepository interface {
	Create(ctx context.Context, history *History) error
	GetByID(ctx context.Context, id int) (*History, error)
	GetWithStyle(ctx context.Context, id int) (*History, error)
	List(ctx context.Context, limit, offset int) ([]History, error)
	ListWithStyles(ctx context.Context, limit, offset int) ([]History, error)
	Search(ctx context.Context, query string, limit, offset int) ([]History, error)
	Count(ctx context.Context) (int, error)
}

type StyleRepository interface {
	Create(ctx context.Context, style *Style) error
	GetByID(ctx context.Context, id int) (*Style, error)
	GetByName(ctx context.Context, name string) (*Style, error)
	List(ctx context.Context) ([]Style, error)
}
