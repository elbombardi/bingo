package analyser

import "strings"

const (
	Separators  = " \t\n\r"
	Puntuations = ".,;:!?()[]{}\"#@&*+-_\\/'`~%^<>|="
)

type Tokenizer struct {
}

func NewTokenizer() Analyser {
	return &Tokenizer{}
}

func (t *Tokenizer) Analyse(corpus *Corpus) (*Corpus, error) {
	for i := range corpus.Documents {
		corpus.Documents[i].Tokens = tokenize(corpus.Documents[i].Content)
		corpus.Statistics.NumTokens += len(corpus.Documents[i].Tokens)
	}
	return corpus, nil
}

func tokenize(content string) []*Token {
	c := replaceAll(content, Puntuations)
	c = replaceAll(c, Separators)
	a := []*Token{}
	s := c
	p := 0
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
			Value:    s[:m],
			Position: p,
		})
		s = s[m+1:]
		p = p + m + 1
	}
	if s != "" {
		a = append(a, &Token{
			Value:    s,
			Position: p,
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
