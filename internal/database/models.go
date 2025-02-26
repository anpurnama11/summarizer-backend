package database

import "time"

type History struct {
    ID        int64     `db:"id"`
    URL       string    `db:"url"`
    Title     string    `db:"title"`
    Content   string    `db:"content"`
    Summary   string    `db:"summary"`
    StyleID   int64     `db:"style_id"`
    Language  string    `db:"language"`
    CreatedAt time.Time `db:"created_at"`
}

type SummarizationStyle struct {
    ID             int64     `db:"id"`
    Name           string    `db:"name"`
    Description    string    `db:"description"`
    PromptTemplate string    `db:"prompt_template"`
    CreatedAt      time.Time `db:"created_at"`
}