package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Category string

const (
	CategoryCareer       Category = "career"
	CategoryArchitecture Category = "architecture"
	CategoryLifestyle    Category = "lifestyle"
	CategoryProjects     Category = "projects"
)

type Metadata struct {
	Timeline string `json:"timeline,omitempty"`
	Priority int    `json:"priority"`
}

type KnowledgeChunk struct {
	ID        string    `json:"id"`
	Category  Category  `json:"category"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Metadata  Metadata  `json:"metadata"`
	Embedding []float32 `json:"embedding,omitempty"`
}

type KnowledgeLoader interface {
	LoadDirectory(ctx context.Context, dirPath string) ([]KnowledgeChunk, error)
	ValidateSchema(ctx context.Context, chunk KnowledgeChunk) error
}

type FileKnowledgeLoader struct{}

func (loader FileKnowledgeLoader) LoadDirectory(ctx context.Context, dirPath string) ([]KnowledgeChunk, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("read knowledge directory %q: %w", dirPath, err)
	}

	chunks := make([]KnowledgeChunk, 0)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		path := filepath.Join(dirPath, entry.Name())
		fileChunks, err := loadChunks(path)
		if err != nil {
			return nil, err
		}
		for index, chunk := range fileChunks {
			if err := loader.ValidateSchema(ctx, chunk); err != nil {
				return nil, fmt.Errorf("validate %s chunk %d: %w", path, index, err)
			}
			chunks = append(chunks, chunk)
		}
	}

	return chunks, nil
}

func (FileKnowledgeLoader) ValidateSchema(ctx context.Context, chunk KnowledgeChunk) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	switch {
	case chunk.ID == "":
		return errors.New("id is required")
	case !validCategory(chunk.Category):
		return fmt.Errorf("category %q is not supported", chunk.Category)
	case chunk.Title == "":
		return errors.New("title is required")
	case chunk.Content == "":
		return errors.New("content is required")
	case chunk.Tags == nil:
		return errors.New("tags are required")
	}

	for index, tag := range chunk.Tags {
		if tag == "" {
			return fmt.Errorf("tags[%d] must not be empty", index)
		}
	}
	return nil
}

func validCategory(category Category) bool {
	switch category {
	case CategoryCareer, CategoryArchitecture, CategoryLifestyle, CategoryProjects:
		return true
	default:
		return false
	}
}

func loadChunks(path string) ([]KnowledgeChunk, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open knowledge file %q: %w", path, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	var chunks []KnowledgeChunk
	if err := decoder.Decode(&chunks); err != nil {
		return nil, fmt.Errorf("decode knowledge file %q: %w", path, err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return nil, fmt.Errorf("decode knowledge file %q: multiple JSON values", path)
		}
		return nil, fmt.Errorf("decode knowledge file %q: %w", path, err)
	}
	if chunks == nil {
		return nil, fmt.Errorf("decode knowledge file %q: expected a JSON array", path)
	}
	return chunks, nil
}
