package index

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/elbombardi/siego/utils"
)

type HitField struct {
	Name               string
	Line               int
	PositionInDocument int
}

type Hit struct {
	DocId   int
	Fields  []HitField
	MinLine int
	MaxLine int
	MinPID  int
	MaxPID  int
}

type DocumentLanguage struct {
	Code string `json:"language"`
}

type Document struct {
	Id         int                `json:"id"`
	Name       string             `json:"name"`
	WordsCount int                `json:"words_count"`
	Languages  []DocumentLanguage `json:"languages"`
}

type Position struct {
	PositionInDoc int `json:"p"`
	LineNumber    int `json:"l"`
}

type Location struct {
	DocId     int        `json:"k"`
	Positions []Position `json:"p"`
}

type IndexEntry struct {
	Name      rune                `json:"n"`
	Locations map[int]Location    `json:"l"`
	Children  map[rune]IndexEntry `json:"c"`
	Parent    *IndexEntry         `json:"-"`
}

type SiegoIndex struct {
	DocumentsMap   []Document `json:"documents_map"`
	TotalWordCount int        `json:"total_word_count"`
	Target         string     `json:"target"`
	Root           IndexEntry `json:"root"`
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
	index.Root.propagateParent()
	return index, nil
}

func (ind *SiegoIndex) Index() {
	if ind.Target == "" {
		fmt.Println("Target directory is mandatory")
		os.Exit(1)
	}
	fmt.Println("Data path : ", ind.Target)
	ind.stepGenerateDocumentsMap()
	ind.stepGenerateEntries()
	ind.stepPurgeMostFrequent()
	fmt.Printf("Done!\n %d document(s) parsed, %d word(s) parsed, %d word(s) indexed\n",
		len(ind.DocumentsMap), ind.TotalWordCount, ind.CountEntries())
}

func (ind *SiegoIndex) Lookup(query string) (hits []Hit, found bool) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, false
	}
	query = utils.Normalise(query)
	queryFields := strings.Fields(query)
	entries := make(map[string]*IndexEntry)
	for _, queryField := range queryFields {
		entry := ind.Root.lookupWord([]rune(queryField))
		if entry == nil {
			//all fields are required,
			return nil, false
		}
		entries[queryField] = entry
	}
	//we look for them the words in the same order within a limited range (3 words)
	// for _, queryField := range queryFields {
	// 	entries[queryField].Locations[]
	// }
	return nil, false
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

func (ind *SiegoIndex) stepGenerateDocumentsMap() {
	ind.DocumentsMap = []Document{}
	docCounter := 0
	for _, filePath := range utils.BrowseDir(ind.Target) {
		if utils.IsTextFile(filePath) {
			document := Document{
				Id:   docCounter,
				Name: filePath,
			}
			ind.DocumentsMap = append(ind.DocumentsMap, document)
			docCounter++
		}
	}
}

func (ind *SiegoIndex) stepGenerateEntries() {
	totalDocNumber := len(ind.DocumentsMap)
	for i, doc := range ind.DocumentsMap {
		fmt.Printf("[%3.0f %%] Indexing '%s'...\n",
			(float32(i)/float32(totalDocNumber))*100, doc.Name)
		doc.WordsCount = ind.parseDocument(doc, i)
		ind.DocumentsMap[i] = doc
		ind.TotalWordCount += doc.WordsCount
	}
	ind.Root.propagateParent()
}

func (ind *SiegoIndex) stepPurgeMostFrequent() {
	// ind.Root.purgeFrequentEntries(len(ind.DocumentsMap), FREQUENCY_THRESHOLD)
}

func (ind *SiegoIndex) parseDocument(document Document, documentId int) (wordsCount int) {
	file, err := os.Open(document.Name)
	if err != nil {
		fmt.Println("Error while reading file.", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	wordsCount = 0
	linesCount := 1
	for scanner.Scan() {
		line := scanner.Text()
		line = utils.Normalise(line)
		for _, word := range strings.FieldsFunc(line, utils.IsNotLetter) {
			wordsCount++
			ind.indexWord(word, documentId, wordsCount, linesCount)
		}
		linesCount++
	}
	return wordsCount
}

func (ind *SiegoIndex) indexWord(word string, documentId, positionInDoc, lineNumber int) {
	word = strings.TrimSpace(word)
	// if len(word) < MIN_LENGTH_ELIGIBLE_WORD {
	// 	return
	// }
	word = strings.ToUpper(word)
	ind.Root.indexWord([]rune(word), documentId, positionInDoc, lineNumber)
}

func (ind *SiegoIndex) CountEntries() (count int) {
	return ind.Root.countEntries()
}

func (ind *SiegoIndex) PrintEntries() {
	ind.Root.printEntries("")
}

func (entry *IndexEntry) printEntries(indent string) {
	for k, child := range entry.Children {
		fmt.Printf("%s%s%s (parent : %s)\n", indent, "Name => ", string(k), string(child.Parent.Name))
		if len(child.Locations) != 0 {
			fmt.Printf("%s%s%v\n", indent, "Locations => ", child.Locations)
		}
		child.printEntries(indent + "\t")
	}
}

func (entry *IndexEntry) indexWord(word []rune, documentId, positionInDoc, lineNumber int) {
	if len(word) == 0 {
		return
	}
	header := word[0]
	if entry.Children == nil {
		entry.Children = make(map[rune]IndexEntry)
	}
	child, exist := entry.Children[header]
	if !exist {
		child = IndexEntry{
			Name: header,
		}
	}
	if len(word) == 1 {
		newPostion := Position{positionInDoc, lineNumber}
		if child.Locations == nil {
			child.Locations = make(map[int]Location)
		}
		loc, exists := child.Locations[documentId]
		if !exists {
			child.Locations[documentId] = Location{documentId, []Position{newPostion}}
		} else {
			loc.Positions = append(loc.Positions, newPostion)
			child.Locations[documentId] = loc
		}
	} else {
		child.indexWord(word[1:], documentId, positionInDoc, lineNumber)
	}
	entry.Children[header] = child
}

func (entry *IndexEntry) lookupWord(query []rune) *IndexEntry {
	if len(query) == 0 {
		return nil
	}
	if len(entry.Children) == 0 {
		return nil
	}
	header := query[0]
	child, exists := entry.Children[header]
	if !exists {
		return nil
	}
	if len(query) == 1 {
		if child.Locations == nil {
			return nil
		} else {
			return &child
		}
	} else {
		return child.lookupWord(query[1:])
	}
}

func (entry *IndexEntry) countEntries() (count int) {
	if len(entry.Locations) != 0 {
		count = 1
	}
	for _, child := range entry.Children {
		count += child.countEntries()
	}
	return count
}

func (entry *IndexEntry) propagateParent() {
	for key := range entry.Children {
		child := entry.Children[key]
		child.Parent = entry
		child.propagateParent()
		entry.Children[key] = child
	}
}

func (entry *IndexEntry) fullName() string {
	if entry.Parent == nil {
		return string(entry.Name)
	}
	return entry.Parent.fullName() + string(entry.Name)
}
