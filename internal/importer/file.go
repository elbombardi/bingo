package importer

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// This is a file importer
// It will implement the Importer interface
// It will import a corpus from a directory of files
type FileImporter struct {
}

// File Importer params.
type FileImporterParams struct {
	// Path of the corpus.
	Path string `json:"path"`

	// Filter of the corpus.
	Filter string `json:"filter"`
}

// Import imports a corpus from a directory.
func (fi *FileImporter) Import(params any) (*Import, error) {
	fip, ok := params.(*FileImporterParams)
	if !ok {
		return nil, fmt.Errorf("invalid params")
	}
	log.Printf("Importing corpus from directory %v...\n", fip.Path)

	corpus := &Import{}
	corpus.Name = filepath.Base(fip.Path)
	corpus.Description = fmt.Sprintf("Corpus imported from %s", fip.Path)
	corpus.Source = "File"
	corpus.URI = fip.Path
	err := filepath.Walk(fip.Path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if fip.Filter != "" && filepath.Ext(path) != fip.Filter {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		doc := &Document{}
		doc.URI, _ = strings.CutPrefix(path, fip.Path)
		content, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		doc.Content = string(content)
		corpus.Documents = append(corpus.Documents, doc)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return corpus, nil
}
