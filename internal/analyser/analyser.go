package analyser

import "github.com/elbombardi/bingo/internal/importer"

// Token represents a word in a document.
type Token struct {
	// The original form of the word in the document,
	// before any preprocessing.
	Value string `json:"v"`

	// The position of the word in the document.
	Position int `json:"p"`
}

// Root represents a word root in a document.
// It is the result of processing a token,
// And it keeps track of the tokens that it is derived from.
type Root struct {
	// The value of the root.
	Value string `json:"v"`

	// The tokens that the root is derived from.
	Tokens []*Token `json:"ts"`

	// The term frequency of the root in the document.
	TF float64 `json:"tf"`

	// The document frequency of the root in the corpus.
	DF int `json:"df"`

	// The inverse document frequency of the root in the corpus.
	IDF float64 `json:"idf"`

	// The term frequency-inverse document frequency of the root in the corpus.
	TFIDF float64 `json:"tfidf"`
}

// Document represents a document in a corpus.
// A document is a collection of tokens.
type Document struct {
	// The ID of the document.
	DocID string `json:"doc_id"`

	// URI of the document.
	URI string `json:"uri"`

	// Content of the document.
	Content string `json:"content"`

	// Tokens of the document.
	Tokens []*Token `json:"tokens"`

	// The roots of the document.
	Roots map[string]*Root `json:"roots"`
}

// CorpusStatistics represents statistics of a corpus.
type CorpusStatistics struct {
	// The number of documents in the corpus.
	NumDocs int `json:"num_docs"`

	// The number of tokens in the corpus.
	NumTokens int `json:"num_tokens"`
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
}

// Analyser is an interface to be implemented by analysers.
type Analyser interface {
	// Analyse analyses an imported corpus and returns an analysed corpus.
	Analyse(corpus *Corpus) (*Corpus, error)
}

// CorpusFromImport converts an imported corpus to a corpus.
func CorpusFromImport(imported *importer.Import) *Corpus {
	c := &Corpus{}
	c.Name = imported.Name
	c.Description = imported.Description
	c.Source = imported.Source
	c.URI = imported.URI
	c.Documents = make([]*Document, len(imported.Documents))
	for i, doc := range imported.Documents {
		c.Documents[i] = &Document{
			DocID:   doc.DocID,
			URI:     doc.URI,
			Content: doc.Content,
		}
	}
	return c
}
