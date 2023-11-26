package analyser

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	str := " Hello!               beautiful;\n    :world 2"
	tokens := tokenize(str)
	if len(tokens) != 4 {
		t.Errorf("Expected 4 tokens, got %d", len(tokens))
	}
	for i, token := range tokens {
		if str[token.Position:token.Position+len(token.Value)] != token.Value {
			t.Errorf("Token %d: '%s' (%v) != '%s'", i, token.Value,
				token.Position, str[token.Position:len(token.Value)])
		}
	}
}
