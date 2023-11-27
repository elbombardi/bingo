package analyser

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	str := " Hello!               beautiful;\n    :world 2"
	tokens := tokenize(str, 1)
	if len(tokens) != 4 {
		t.Errorf("Expected 4 tokens, got %d", len(tokens))
	}
	for i, token := range tokens {
		if str[token.Position.Start:token.Position.End] != token.Value {
			t.Errorf("Token %d: '%s' (%v) != '%s'", i, token.Value,
				token.Position, str[token.Position.Start:token.Position.End])
		}
		if token.DocId != 1 {
			t.Errorf("Token %d: DocId %d != 1", i, token.DocId)
		}
	}
}
