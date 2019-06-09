package fa

import (
	"testing"
)

func TestSubmission(t *testing.T) {
	var f = NewFurAffinity()

	var sub submission
	sub.fa = f
	sub.id = 31698875
	if _, err := sub.GetDetails(); err != nil {
		t.Fatal(err)
		return
	}

	if sub.user.id != "s-nina" {
		t.Error("sub id is: ", sub.user.id, "want: s-nina")
	}
	if sub.user.name != "S-Nina" {
		t.Error("sub name is: ", sub.user.name, "want: S-Nina")
	}

	if sub.fileURL != "https://d.facdn.net/art/s-nina/1559057227/1559057220.s-nina_mihari.png" {
		t.Error("sub fileURL is: ", sub.fileURL, "want: https://d.facdn.net/art/s-nina/1559057227/1559057220.s-nina_mihari.png")
	}

	if sub.scraps {
		t.Error("sub is a scrap when it should be main gallery")
	}

	if _, err := sub.Download(); err != nil {
		t.Fatal(err)
		return
	}
}

func TestScrap(t *testing.T) {
	var f = NewFurAffinity()
	var sub submission
	sub.fa = f
	sub.id = 7893454

	if _, err := sub.GetDetails(); err != nil {
		t.Fatal(err)
		return
	}

	if !sub.scraps {
		t.Error("sub was expected to be scrap, it was not")
	}
}
