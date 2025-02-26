package repository

import (
	"context"
	"database/sql"

	"anpurnama/summarizer-backend/internal/database"
)

type historyRepository struct {
	db *database.DB
}

func NewHistoryRepository(db *database.DB) HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) GetWithStyle(ctx context.Context, id int) (*History, error) {
	query := `
        SELECT h.id, h.url, h.title, h.content, h.summary,
            h.style_id, h.language, h.created_at,
            s.id, s.name, s.description, s.prompt_template, s.created_at
        FROM history h
        LEFT JOIN summarization_styles s ON h.style_id = s.id
        WHERE h.id = ?
    `

	history := &History{}
	style := &Style{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&history.ID, &history.URL, &history.Title,
		&history.Content, &history.Summary, &history.StyleID,
		&history.Language, &history.CreatedAt,
		&style.ID, &style.Name, &style.Description,
		&style.PromptTemplate, &style.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if style.ID != 0 {
		history.Style = style
	}

	return history, nil
}

func (r *historyRepository) ListWithStyles(ctx context.Context, limit, offset int) ([]History, error) {
	query := `
        SELECT h.id, h.url, h.title, h.content, h.summary,
            h.style_id, h.language, h.created_at,
            s.id, s.name, s.description, s.prompt_template, s.created_at
        FROM history h
        LEFT JOIN summarization_styles s ON h.style_id = s.id
        ORDER BY h.created_at DESC
        LIMIT ? OFFSET ?
    `
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []History
	for rows.Next() {
		h := History{}
		s := Style{}

		err := rows.Scan(
			&h.ID, &h.URL, &h.Title, &h.Content,
			&h.Summary, &h.StyleID, &h.Language, &h.CreatedAt,
			&s.ID, &s.Name, &s.Description,
			&s.PromptTemplate, &s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if s.ID != 0 {
			h.Style = &s
		}

		histories = append(histories, h)
	}
	return histories, nil
}

// Update Create method to include validation
func (r *historyRepository) Create(ctx context.Context, history *History) error {
	if err := history.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO history (
			url, title, content, summary, style_id,
			language, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		history.URL, history.Title, history.Content,
		history.Summary, history.StyleID, history.Language,
		history.CreatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	history.ID = int(id)
	return nil
}

func (r *historyRepository) GetByID(ctx context.Context, id int) (*History, error) {
	query := `
        SELECT id, url, title, content, summary, style_id, language, created_at
        FROM history
        WHERE id = ?
    `

	history := &History{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&history.ID, &history.URL, &history.Title,
		&history.Content, &history.Summary, &history.StyleID,
		&history.Language, &history.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (r *historyRepository) List(ctx context.Context, limit, offset int) ([]History, error) {
	query := `
		SELECT id, url, title, content, summary,
			style_id, language, created_at
		FROM history
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []History
	for rows.Next() {
		var h History
		err := rows.Scan(
			&h.ID, &h.URL, &h.Title, &h.Content,
			&h.Summary, &h.StyleID, &h.Language,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, h)
	}
	return histories, nil
}

func (r *historyRepository) Search(ctx context.Context, query string, limit, offset int) ([]History, error) {
	sqlQuery := `
		SELECT id, url, title, content, summary,
			style_id, language, created_at
		FROM history
		WHERE title LIKE ? OR url LIKE ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	searchPattern := "%" + query + "%"
	rows, err := r.db.QueryContext(ctx, sqlQuery, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []History
	for rows.Next() {
		var h History
		err := rows.Scan(
			&h.ID, &h.URL, &h.Title, &h.Content,
			&h.Summary, &h.StyleID, &h.Language,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, h)
	}
	return histories, nil
}

func (r *historyRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM history"
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
