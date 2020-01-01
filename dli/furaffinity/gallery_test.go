package fa

import "testing"

func TestGallery(t *testing.T) {
	var f = NewFurAffinity()

	subs, err := f.Posts("s-nina", 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(subs) <= 0 {
		t.Fatal("no subs returned")
	}

	//if subs[0].User().Name() != "s-nina" {
	//	t.Fatalf("username does not match input: %s, s-nina", subs[0].User().Name())
	//}

	nextBatch, err := f.Posts("s-nina", 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(nextBatch) <= 0 {
		t.Fatal("no subs in nextBatch")
	}

	if nextBatch[0].ID() == subs[0].ID() {
		t.Fatal("recieved the same batch")
	}
}

func TestGalleryScrapsRollover(t *testing.T) {
	var f = NewFurAffinity()

	firstBatch, err := f.Posts("s-nina", 2)
	if err != nil {
		t.Fatal(err)
	}

	if len(firstBatch) <= 0 {
		t.Fatal("recived no subs in firstBatch")
	}

	secondBatch, err := f.Posts("s-nina", 3)
	if err != nil {
		t.Fatal(err)
	}

	if len(secondBatch) <= 0 {
		t.Fatal("recived no subs in secondBatch")
	}

	sscrap := secondBatch[0]
	if _, err = sscrap.GetDetails(); err != nil {
		t.Fatal(err)
	}

	if sscrap.Folder() != "scraps" {
		t.Fatal("subs are not scraps, tests might be outdated! ", secondBatch[0].ID())
	}
}
