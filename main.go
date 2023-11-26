package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/elbombardi/bingo/internal/analyser"
	"github.com/elbombardi/bingo/internal/importer"
)

func main() {
	log.Println("Importing...")
	fi := &importer.FileImporter{}
	fip := &importer.FileImporterParams{
		Path:   "./datasets/shakespeares",
		Filter: "",
	}
	imported, err := fi.Import(fip)
	if err != nil {
		panic(err)
	}

	corpus := analyser.CorpusFromImport(imported)

	log.Println("Tokenizing...")
	tokenizer := analyser.NewTokenizer()
	corpus, err = tokenizer.Analyse(corpus)
	if err != nil {
		panic(err)
	}

	log.Println("Lemmatizing...")
	lemmatizer := analyser.NewLemmatizer()
	corpus, err = lemmatizer.Analyse(corpus)
	if err != nil {
		panic(err)
	}

	log.Println("Done!")
	json, err := json.MarshalIndent(corpus, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
}
