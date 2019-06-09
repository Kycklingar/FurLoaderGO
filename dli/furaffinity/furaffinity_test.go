package fa

import "testing"

func TestLogin(t *testing.T) {
	var fa furaffinity
	if err := fa.Login("", ""); err != nil {
		//t.Fatal(err)
	}
}
