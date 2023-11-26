package analyser

import (
	"strings"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
)

type Lemmatizer struct {
	Language string
}

func NewLemmatizer() Analyser {
	return &Lemmatizer{}
}

func (l *Lemmatizer) Analyse(corpus *Corpus) (*Corpus, error) {
	lemmatizer, err := golem.New(en.New())
	if err != nil {
		return nil, err
	}
	for i := range corpus.Documents {
		corpus.Documents[i].Roots = lemmatize(lemmatizer, corpus.Documents[i].Tokens)
	}
	return corpus, nil
}

func lemmatize(lemmatizer *golem.Lemmatizer, tokens []*Token) map[string]*Root {
	roots := map[string]*Root{}
	for _, token := range tokens {
		value := strings.ToLower(token.Value)
		value = strings.TrimSpace(value)
		value = lemmatizer.Lemma(value)
		_, ok := roots[value]
		if !ok {
			roots[value] = &Root{
				Value:  value,
				Tokens: []*Token{},
				TF:     0,
			}
		}
		root := roots[value]
		root.Tokens = append(root.Tokens, token)
		root.TF++

	}
	return roots
}
