package domain

import (
	"context"
	"testing"
)

func TestLoadDirectoryLoadsAndValidatesKnowledgeFiles(t *testing.T) {
	loader := FileKnowledgeLoader{}
	chunks, err := loader.LoadDirectory(context.Background(), "../../knowledge/raw")
	if err != nil {
		t.Fatalf("LoadDirectory() error = %v", err)
	}
	if len(chunks) != 4 {
		t.Fatalf("LoadDirectory() loaded %d chunks, want 4", len(chunks))
	}
}

func TestValidateSchemaRejectsInvalidChunk(t *testing.T) {
	loader := FileKnowledgeLoader{}
	chunk := KnowledgeChunk{
		Category: CategoryCareer,
		Title:    "Geçersiz kayıt",
		Content:  "İçerik",
		Tags:     []string{"Go"},
	}

	if err := loader.ValidateSchema(context.Background(), chunk); err == nil {
		t.Fatal("ValidateSchema() error = nil, want missing id error")
	}
}
