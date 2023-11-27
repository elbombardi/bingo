package analyser

import (
	"errors"
	"math"
)

type TFIDF struct {
}

func NewTFIDF() Analyser {
	return &TFIDF{}
}

func (t *TFIDF) Analyse(corpus *Corpus) (*Corpus, error) {
	if corpus.Statistics.NumDocTerms == 0 {
		return nil, errors.New("no terms to calculate TF-IDF")
	}

	// Grouping terms.
	corpus.Terms = map[string]*Term{}
	for i := range corpus.Documents {
		for _, term := range corpus.Documents[i].Terms {
			_, ok := corpus.Terms[term.Value]
			if !ok {
				corpus.Terms[term.Value] = &Term{
					Value: term.Value,
				}
			}
			corpus.Terms[term.Value].Documents = append(corpus.Terms[term.Value].Documents, i)
			corpus.Terms[term.Value].Tokens = append(corpus.Terms[term.Value].Tokens, term.Tokens...)
		}
	}
	corpus.Statistics.NumUniqueTerms = len(corpus.Terms)

	// Calculating TF-IDF.
	for i := range corpus.Documents {
		for _, docTerm := range corpus.Documents[i].Terms {
			term := corpus.Terms[docTerm.Value]
			docTerm.TF = float64(len(docTerm.Tokens)) / float64(len(corpus.Documents[i].Tokens))
			term.DF = float64(len(corpus.Terms[term.Value].Documents)) / float64(len(corpus.Documents))
			term.IDF = math.Log(1/term.DF) + 1
			docTerm.TFIDF = docTerm.TF * term.IDF
		}
	}
	return corpus, nil
}
