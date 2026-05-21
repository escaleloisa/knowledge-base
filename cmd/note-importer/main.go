package main

import (
	"log"

	noteimporter "github.com/escaleloisa/knowledge-base/internal/note-importer"
	"github.com/escaleloisa/knowledge-base/pkg/config"
)

func main() {
	cfg := config.Load()

	importer := noteimporter.New(cfg.NoteServiceURL, cfg.ImportPath)

	log.Println("Note importer starting...")
	if err := importer.Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("Note importer finished.")
}
