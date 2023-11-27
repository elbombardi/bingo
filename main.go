package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/elbombardi/bingo/internal/analyser"
	"github.com/elbombardi/bingo/internal/importer"
)

func main() {
	// Importing
	fi := &importer.FileImporter{}
	fip := &importer.FileImporterParams{
		Path:   "./datasets/shakespeares",
		Filter: "",
	}
	imported, err := fi.Import(fip)
	if err != nil {
		log.Fatal("Import failed", err)
	}

	corpus := analyser.CorpusFromImport(imported)
	corpus.Statistics.NumDocs = len(corpus.Documents)

	// Tokenizing
	tokenizer := analyser.NewTokenizer()
	corpus, err = tokenizer.Analyse(corpus)
	if err != nil {
		log.Fatal("Tokenization failed", err)
	}

	// Lemmatizing
	// lemmatizer := analyser.NewLemmatizer("fr")
	// corpus, err = lemmatizer.Analyse(corpus)
	// if err != nil {
	// 	panic(err)
	// }

	// Stemming
	stemmer := analyser.NewStemmer("english")
	corpus, err = stemmer.Analyse(corpus)
	if err != nil {
		log.Fatal("Stemming failed", err)
	}

	// TF-IDF
	tfidf := analyser.NewTFIDF()
	corpus, err = tfidf.Analyse(corpus)
	if err != nil {
		log.Fatal("TF-IDF failed", err)
	}

	log.Println("Done!")
	json, err := json.MarshalIndent(corpus, "", "  ")
	if err != nil {
		log.Fatal("Marshaling failed", err)
	}
	fmt.Println(string(json))
}
