package importer

import (
	"fmt"
	"log"
)

// MemoryImporter import corpus from memory.
type MemoryImporter struct {
}

// File Importer params.
type MemoryImporterParams []string

// Import imports a corpus from a directory.
func (mi *MemoryImporter) Import(params any) (*Import, error) {
	mip, ok := params.(*MemoryImporterParams)
	if !ok {
		return nil, fmt.Errorf("invalid params")
	}
	log.Printf("Importing corpus from memory %v...\n", mip)

	corpus := &Import{}
	corpus.Name = "Memory"
	corpus.Description = "Corpus imported from memory"
	corpus.Source = "Memory"
	corpus.URI = "-"
	for _, content := range *mip {
		doc := &Document{}
		doc.URI = "-"
		doc.Content = content
		corpus.Documents = append(corpus.Documents, doc)
	}
	return corpus, nil
}
