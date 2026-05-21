package noteimporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Importer struct {
	noteServiceURL string
	importPath     string
}

func New(noteServiceURL, importPath string) *Importer {
	return &Importer{noteServiceURL: noteServiceURL, importPath: importPath}
}

func (imp *Importer) Run() error {
	files, err := filepath.Glob(filepath.Join(imp.importPath, "*.md"))
	if err != nil {
		return fmt.Errorf("glob: %w", err)
	}

	if len(files) == 0 {
		log.Println("no .md files found in", imp.importPath)
		return nil
	}

	for _, file := range files {
		if err := imp.importFile(file); err != nil {
			log.Printf("failed to import %s: %v", file, err)
			continue
		}
		// Move processed file
		processed := file + ".imported"
		os.Rename(file, processed)
		log.Printf("imported: %s", filepath.Base(file))
	}
	return nil
}

func (imp *Importer) importFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	title := strings.TrimSuffix(filepath.Base(path), ".md")
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")

	body, _ := json.Marshal(map[string]any{
		"title":   title,
		"content": string(data),
		"tags":    []string{"imported"},
	})

	resp, err := http.Post(
		fmt.Sprintf("%s/api/notes", imp.noteServiceURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("note service returned %d", resp.StatusCode)
	}
	return nil
}
