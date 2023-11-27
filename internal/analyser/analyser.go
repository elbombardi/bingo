package analyser

import (
	"strings"
	"unicode"

	"github.com/elbombardi/bingo/internal/importer"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type TokenPosition struct {
	Start int `json:"s"`
	End   int `json:"e"`
}

// Token represents a word in a document.
type Token struct {
	// The index of the word in the document.
	Index int `json:"i"`

	// The original form of the word in the document,
	// before any preprocessing.
	Value string `json:"v"`

	// The ID of the document that the word belongs to.
	DocId int `json:"d"`

	// The position of the word in the document.
	Position TokenPosition `json:"p"`
}

// DocTerm represents a term in a document.
// It is the result of processing a token,
// And it keeps track of the tokens that it is derived from.
type DocTerm struct {
	// The value of the root.
	Value string `json:"v"`

	// The tokens that the root is derived from.
	Tokens []*Token `json:"ts"`

	// The TF of the term (calculated by tfidf analyser)
	TF float64 `json:"tf"`

	// The TFIDF of the term (calculated by tfidf analyser)
	TFIDF float64 `json:"tfidf"`
}

type Term struct {
	// The value of the root.
	Value string `json:"v"`

	// The tokens that the root is derived from.
	Tokens []*Token `json:"ts"`

	// The documents that the root is derived from.
	Documents []int `json:"ds"`

	// The DF of the term (calculated by tfidf analyser)
	DF float64 `json:"df"`

	// The IDF of the term (calculated by tfidf analyser)
	IDF float64 `json:"idf"`
}

// Document represents a document in a corpus.
// A document is a collection of tokens.
type Document struct {
	// The ID of the document.
	DocID int `json:"doc_id"`

	// URI of the document.
	URI string `json:"uri"`

	// Content of the document.
	Content string `json:"content"`

	// Tokens of the document.
	Tokens []*Token `json:"tokens"`

	// The roots of the document.
	Terms map[string]*DocTerm `json:"terms"`
}

// CorpusStatistics represents statistics of a corpus.
type CorpusStatistics struct {
	// The number of documents in the corpus.
	NumDocs int `json:"num_docs"`

	// The number of tokens in the corpus.
	NumTokens int `json:"num_tokens"`

	// The number of terms in the corpus.
	NumDocTerms int `json:"num_doc_terms"`

	// The number of unique terms in the corpus.
	NumUniqueTerms int `json:"num_unique_terms"`
}

// Corpus represents a corpus of documents.
// A corpus is a collection of documents.
type Corpus struct {
	// The name of the corpus.
	Name string `json:"name"`

	// Description of the corpus.
	Description string `json:"description"`

	// Source of the corpus.
	Source string `json:"source"`

	// URI of the corpus.
	URI string `json:"uri"`

	// Statistics of the corpus.
	Statistics CorpusStatistics `json:"statistics"`

	// The documents in the corpus.
	Documents []*Document `json:"documents"`

	// The unique terms in the corpus.
	Terms map[string]*Term `json:"terms"`
}

// Analyser is an interface to be implemented by analysers.
type Analyser interface {
	// Analyse analyses an imported corpus and returns an analysed corpus.
	Analyse(corpus *Corpus) (*Corpus, error)
}

// CorpusFromImport converts an import to a corpus.
func CorpusFromImport(imported *importer.Import) *Corpus {
	c := &Corpus{}
	c.Name = imported.Name
	c.Description = imported.Description
	c.Source = imported.Source
	c.URI = imported.URI
	c.Documents = make([]*Document, len(imported.Documents))
	for i, doc := range imported.Documents {
		c.Documents[i] = &Document{
			DocID:   i,
			URI:     doc.URI,
			Content: doc.Content,
		}
	}
	return c
}

func normalize(s string) string {
	s = strings.ToLower(s)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ = transform.String(t, s)
	return s
}
