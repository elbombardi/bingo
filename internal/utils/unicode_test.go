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

	// t.Fatal(strings.Join(strings.FieldsFunc( 	("1|6|اهْدِنَا الصِّرَاطَ الْمُسْتَقِيمَ"), IsNotLetter), "/"))
	//t.Fatal()
}
