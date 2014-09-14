package gifsecrets

import (
	"testing"
)

const TestSecret = "test secret"

func TestGifSecrets(t *testing.T) {
	if err := Encode("coin.gif", TestSecret); err != nil {
		t.Fatal(err)
	}

	secret, err := Decode("out.gif")
	if err != nil {
		t.Fatal(err)
	}

	// intentionally weird behavior until I learn regex
	if secret != TestSecret {
		t.Errorf("Decoded secret was %s instead of %s", secret, TestSecret)
	}
}
