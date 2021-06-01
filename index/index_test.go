package index

import (
	"testing"
)

func TestIndex(t *testing.T) {
	index := SiegoIndex{
		Target: "./small",
	}
	index.Index()

	numberOfMappedFiles := len(index.LocationsMap)
	if numberOfMappedFiles != 2 {
		t.Fatal("The number of mapped files expected 2, number found : ", numberOfMappedFiles)
	}

	locations := index.Lookup("errors")
	if len(locations) != 1 {
		t.Fatal("The number of files expected 1, number found ", len(locations))
	}

	var counter int
	for _, v := range locations {
		counter = v.Counter
		break
	}
	if counter != 2 {
		t.Fatal("The number of occurences expected 2, number found ", counter)
	}
}
