package analyser

import (
	"errors"
	"log"

	"github.com/kljensen/snowball"
)

type Stemmer struct {
	Language string
}

func NewStemmer(language string) Analyser {
	return &Stemmer{
		Language: language,
	}
}

func (l *Stemmer) Analyse(corpus *Corpus) (*Corpus, error) {
	log.Printf("Stemming corpus '%s'...\n", corpus.Name)

	if corpus.Statistics.NumTokens == 0 {
		return nil, errors.New("no tokens to stem")
	}

	var err error
	corpus.Statistics.NumDocTerms = 0
	for i := range corpus.Documents {
		corpus.Documents[i].Terms, err = stem(l.Language, corpus.Documents[i].Tokens)
		if err != nil {
			return nil, err
		}
		corpus.Statistics.NumDocTerms += len(corpus.Documents[i].Terms)
		// corpus.Documents[i].Tokens = nil
	}
	return corpus, nil
}

func stem(lang string, tokens []*Token) (map[string]*DocTerm, error) {
	terms := map[string]*DocTerm{}
	for _, token := range tokens {
		v, err := snowball.Stem(token.Value, lang, true)
		if err != nil {
			return nil, err
		}
		v = normalize(v)
		_, ok := terms[v]
		if !ok {
			terms[v] = &DocTerm{
				Value: v,
			}
		}
		term := terms[v]
		term.Tokens = append(term.Tokens, token)
	}
	return terms, nil
}
