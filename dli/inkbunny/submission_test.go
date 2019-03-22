package inkbunny

import "testing"

func TestSubmissionDownload(t *testing.T) {
	var sub ibSub
	sub.fileURL = "https://nl.ib.metapix.net/files/full/2648/2648955_Kagemusha_page228.jpg"

	file, err := sub.Download()
	if err != nil {
		t.Fatal(err)
	}

	if file == nil {
		t.Fatal("file is nil")
	}

	defer file.Close()
}
