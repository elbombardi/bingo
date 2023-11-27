package analyser

import (
	"errors"
	"log"
	"strings"
)

const (
	Separators  = " \t\n\r"
	Puntuations = ".,;:!?()[]{}\"#@&*+-_\\/'`’~%^<>|="
)

type Tokenizer struct {
}

func NewTokenizer() Analyser {
	return &Tokenizer{}
}

func (t *Tokenizer) Analyse(corpus *Corpus) (*Corpus, error) {
	log.Printf("Tokenizing corpus '%s'...\n", corpus.Name)
	if corpus.Statistics.NumDocs == 0 {
		return nil, errors.New("no documents to tokenize")
	}
	corpus.Statistics.NumTokens = 0
	for i := range corpus.Documents {
		corpus.Documents[i].Tokens = tokenize(corpus.Documents[i].Content, i)
		corpus.Statistics.NumTokens += len(corpus.Documents[i].Tokens)
		corpus.Documents[i].Content = ""
	}
	return corpus, nil
}

func tokenize(content string, docId int) []*Token {
	c := replaceAll(content, Puntuations)
	c = replaceAll(c, Separators)
	a := []*Token{}
	s := c
	p := 0
	i := 1
	for s != "" {
		m := strings.IndexByte(s, ' ')
		if m < 0 {
			break
		}
		if m == 0 {
			s = s[1:]
			p++
			continue
		}
		a = append(a, &Token{
			Index: i,
			Value: s[:m],
			Position: TokenPosition{
				Start: p,
				End:   p + m,
			},
			DocId: docId,
		})
		s = s[m+1:]
		p = p + m + 1
		i++
	}
	if s != "" {
		a = append(a, &Token{
			Index: i,
			Value: s,
			Position: TokenPosition{
				Start: p,
				End:   p + len(s),
			},
			DocId: docId,
		})
	}
	return a
}

func replaceAll(content string, symbols string) string {
	for _, p := range symbols {
		content = strings.ReplaceAll(content, string(p), " ")
	}
	return content
}
