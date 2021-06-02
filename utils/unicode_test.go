package utils

import (
	"testing"
)

func TestNormalise(t *testing.T) {
	source := "žůžo"
	trans := Normalise(source)
	if trans != "ZUZO" {
		t.Fatalf("'%s' not zuzo\n", trans)
	}
}
