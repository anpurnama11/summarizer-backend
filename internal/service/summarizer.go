package service

import "context"

type Summarizer interface {
	Summarize(ctx context.Context, content string, styleName string) (string, error)
}
