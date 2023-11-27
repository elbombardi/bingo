package importer

// Document represents a document in a corpus.
// A document is a collection of tokens.
type Document struct {
	// The ID of the document.
	DocID int `json:"doc_id"`

	// URI of the document.
	URI string `json:"uri"`

	// Content of the document.
	Content string `json:"content"`
}

// Import represents a corpus of documents.
// A corpus is a collection of documents.
type Import struct {
	// The name of the corpus.
	Name string `json:"name"`

	// Description of the corpus.
	Description string `json:"description"`

	// Source of the corpus.
	Source string `json:"source"`

	// URI of the corpus.
	URI string `json:"uri"`

	// The documents in the corpus.
	Documents []*Document `json:"documents"`
}

// Importer represents an importer.
type Importer interface {
	// Import imports a corpus.
	Import(params any) (*Import, error)
}
