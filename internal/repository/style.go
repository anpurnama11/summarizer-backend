package repository

import (
	"anpurnama/summarizer-backend/internal/database"
	"context"
	"database/sql"
	"sync"
)

type styleRepository struct {
	db    *database.DB
	cache struct {
		byID   map[int]*Style
		byName map[string]*Style
		mu     sync.RWMutex
	}
}

func NewStyleRepository(db *database.DB) StyleRepository {
	return &styleRepository{
		db: db,
		cache: struct {
			byID   map[int]*Style
			byName map[string]*Style
			mu     sync.RWMutex
		}{
			byID:   make(map[int]*Style),
			byName: make(map[string]*Style),
		},
	}
}

func (r *styleRepository) Create(ctx context.Context, style *Style) error {
	if err := style.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO summarization_styles (
			name, description, prompt_template
		) VALUES (?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		style.Name, style.Description, style.PromptTemplate,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	style.ID = int(id)

	// Update cache
	r.cache.mu.Lock()
	r.cache.byID[style.ID] = style
	r.cache.byName[style.Name] = style
	r.cache.mu.Unlock()

	return nil
}

func (r *styleRepository) GetByID(ctx context.Context, id int) (*Style, error) {
	// Check cache first
	r.cache.mu.RLock()
	if style, ok := r.cache.byID[id]; ok {
		r.cache.mu.RUnlock()
		return style, nil
	}
	r.cache.mu.RUnlock()

	query := `
		SELECT id, name, description, prompt_template, created_at
		FROM summarization_styles WHERE id = ?
	`
	style := &Style{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&style.ID, &style.Name, &style.Description,
		&style.PromptTemplate, &style.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Update cache
	r.cache.mu.Lock()
	r.cache.byID[style.ID] = style
	r.cache.byName[style.Name] = style
	r.cache.mu.Unlock()

	return style, nil
}

func (r *styleRepository) GetByName(ctx context.Context, name string) (*Style, error) {
	// Check cache first
	r.cache.mu.RLock()
	if style, ok := r.cache.byName[name]; ok {
		r.cache.mu.RUnlock()
		return style, nil
	}
	r.cache.mu.RUnlock()

	query := `
		SELECT id, name, description, prompt_template, created_at
		FROM summarization_styles WHERE name = ?
	`
	style := &Style{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&style.ID, &style.Name, &style.Description,
		&style.PromptTemplate, &style.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Update cache
	r.cache.mu.Lock()
	r.cache.byID[style.ID] = style
	r.cache.byName[style.Name] = style
	r.cache.mu.Unlock()

	return style, nil
}

func (r *styleRepository) List(ctx context.Context) ([]Style, error) {
	query := `
		SELECT id, name, description, prompt_template, created_at
		FROM summarization_styles
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var styles []Style
	for rows.Next() {
		var s Style
		err := rows.Scan(
			&s.ID, &s.Name, &s.Description,
			&s.PromptTemplate, &s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		styles = append(styles, s)

		// Update cache
		r.cache.mu.Lock()
		r.cache.byID[s.ID] = &s
		r.cache.byName[s.Name] = &s
		r.cache.mu.Unlock()
	}
	return styles, nil
}
