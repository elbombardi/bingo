package index

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/elbombardi/siego/utils"
)

// TODO add exact position in text
type Location struct {
	Key     int `json:"key"`
	Counter int `json:"counter"`
}

type IndexEntry struct {
	Name      rune                `json:"name"`
	Locations map[int]Location    `json:"locations"`
	Children  map[rune]IndexEntry `json:"children"`
}

type SiegoIndex struct {
	LocationsMap []string            `json:"locations_map"`
	Target       string              `json:"target"`
	Entries      map[rune]IndexEntry `json:"entries"`
}

func (ind *SiegoIndex) Index() {
	if ind.Target == "" {
		fmt.Println("Target directory is mandatory")
		os.Exit(1)
	}
	fmt.Println("Data path : ", ind.Target)
	ind.stepGenerateMap()
	ind.stepGenerateEntries()
	fmt.Println("Done!")
}

func (ind *SiegoIndex) Lookup(query string) []Location {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}
	query = strings.ToUpper(query)
	query = strings.Fields(query)[0] //for this first version, we only look for the first word

	return lookupWordInEntries(ind.Entries, []rune(query))
}

func (ind *SiegoIndex) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	err = enc.Encode(ind)
	if err != nil {
		return err
	}
	return nil
}

func Load(filename string) (*SiegoIndex, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	index := &SiegoIndex{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(index)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func (ind *SiegoIndex) indexWord(word string, locationKey int) {
	word = strings.TrimSpace(word)
	if word == "" {
		return
	}
	word = strings.ToUpper(word)
	if ind.Entries == nil {
		ind.Entries = make(map[rune]IndexEntry)
	}
	addWordToEntries(ind.Entries, []rune(word), locationKey)
}

func (ind *SiegoIndex) stepGenerateMap() {
	ind.LocationsMap = []string{}
	for _, filePath := range utils.BrowseDir(ind.Target) {
		if utils.IsTextFile(filePath) {
			ind.LocationsMap = append(ind.LocationsMap, filePath)
		}
	}
}

func (ind *SiegoIndex) generateEntries(location string, locationKey int) {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error while reading file.", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		line = utils.RemoveNonWordCharacters(line)
		for _, word := range strings.Fields(line) {
			ind.indexWord(word, locationKey)
		}
	}
}

func (ind *SiegoIndex) stepGenerateEntries() {
	for i, filename := range ind.LocationsMap {
		fmt.Printf("Indexing '%s'...\n", filename)
		ind.generateEntries(filename, i)
	}
}

func (ind *SiegoIndex) PrintEntries() {
	printEntries(" ", ind.Entries)
}

func printEntries(indent string, entries map[rune]IndexEntry) {
	for k, v := range entries {
		fmt.Printf("%s%s%s\n", indent, "Name => ", string(k))
		if len(v.Locations) != 0 {
			fmt.Printf("%s%s%v\n", indent, "Locations => ", v.Locations)
		}
		if v.Children != nil {
			printEntries(indent+"\t", v.Children)
		}
	}
}

func addWordToEntries(entries map[rune]IndexEntry, word []rune, locKey int) {
	if len(word) == 0 {
		return
	}
	header := word[0]
	entry, exist := entries[header]
	if !exist {
		entry = IndexEntry{
			Name: header,
		}
	}
	if len(word) == 1 {
		if entry.Locations == nil {
			entry.Locations = make(map[int]Location)
		}
		loc, exists := entry.Locations[locKey]
		if !exists {
			entry.Locations[locKey] = Location{locKey, 1}
		} else {
			loc.Counter++
			entry.Locations[locKey] = loc
		}
	} else {
		if entry.Children == nil {
			entry.Children = make(map[rune]IndexEntry)
		}
		addWordToEntries(entry.Children, word[1:], locKey)
	}
	entries[header] = entry
}

func lookupWordInEntries(entries map[rune]IndexEntry, query []rune) []Location {
	if len(query) == 0 {
		return nil
	}
	header := query[0]
	entry, exists := entries[header]
	if !exists {
		return nil
	}
	if len(query) == 1 {
		if entry.Locations == nil {
			return nil
		} else {
			return valuesFromLocationMap(entry.Locations)
		}
	} else {
		if entry.Children == nil {
			return nil
		} else {
			return lookupWordInEntries(entry.Children, query[1:])
		}
	}
}

func valuesFromLocationMap(locations map[int]Location) []Location {
	result := []Location{}
	for _, location := range locations {
		result = append(result, location)
	}
	return result
}
