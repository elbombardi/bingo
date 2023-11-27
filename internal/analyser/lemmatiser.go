package analyser

import (
	"errors"
	"log"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/aaaton/golem/v4/dicts/fr"
)

type Lemmatizer struct {
	Language string
}

func NewLemmatizer(language string) Analyser {
	return &Lemmatizer{
		Language: language,
	}
}

func languageToGolemDict(language string) golem.LanguagePack {
	switch language {
	case "en":
		return en.New()
	case "fr":
		return fr.New()
	default:
		return nil
	}
}
func (l *Lemmatizer) Analyse(corpus *Corpus) (*Corpus, error) {
	log.Printf("Lemmatizing corpus '%s'...\n", corpus.Name)
	if corpus.Statistics.NumTokens == 0 {
		return nil, errors.New("no tokens to lemmatize")
	}
	lp := languageToGolemDict(l.Language)
	if lp == nil {
		return nil, errors.New("unsupported language")
	}
	lemmatizer, err := golem.New(lp)
	if err != nil {
		return nil, err
	}
	corpus.Statistics.NumDocTerms = 0
	for i := range corpus.Documents {
		corpus.Documents[i].Terms = lemmatize(lemmatizer, corpus.Documents[i].Tokens)
		corpus.Statistics.NumDocTerms += len(corpus.Documents[i].Terms)
	}
	return corpus, nil
}

func lemmatize(lemmatizer *golem.Lemmatizer, tokens []*Token) map[string]*DocTerm {
	terms := map[string]*DocTerm{}
	for _, token := range tokens {
		v := lemmatizer.Lemma(token.Value)
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
	return terms
}
