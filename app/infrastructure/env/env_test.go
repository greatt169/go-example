package env

import (
	"os"
	"testing"
)

func TestGetter(t *testing.T) {
	firstTest, secondTest := "first test", "second test"
	env := Getter("ENV_FOR_TEST", firstTest)
	if env != firstTest {
		t.Fatal()
	}
	if err := os.Setenv("ENV_FOR_TEST", secondTest); err == nil {
		env = Getter("ENV_FOR_TEST", firstTest)
		if env == firstTest {
			t.Fatal()
		}
	}
}
