package index

import (
	"testing"
)

func TestIndex(t *testing.T) {
	index := SiegoIndex{
		Target: "./test_data",
	}
	index.Index()

	numberOfMappedFiles := len(index.DocumentsMap)
	if numberOfMappedFiles != 2 {
		t.Fatal("The number of mapped files expected 2, number found : ", numberOfMappedFiles)
	}

	locations, found := index.Lookup("errors")

	if !found {
		t.Fatal("The word 'erros' could not be found")
	}

	if len(locations) != 1 {
		t.Fatal("The number of files expected 1, number found ", len(locations))
	}

	var counter int
	for _, v := range locations {
		counter = len(v.Positions)
		break
	}
	if counter != 2 {
		t.Fatal("The number of occurences expected 2, number found ", counter)
	}
}
