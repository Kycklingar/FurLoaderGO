package inkbunny

import (
	"testing"
)

func TestInkbunnyLogin(t *testing.T) {
	var i InkBunny

	if err := i.Login("guest", ""); err != nil {
		t.Error(err)
	}

	sid := i.GetCookies()

	i = InkBunny{}

	if err := i.Login("faileduserlogin", ""); err == nil {
		t.Error("Expecting an error")
	}

	if err := i.LoginCookies(sid); err != nil {
		t.Error(err)
	}
}

func TestInkbunnySearch(t *testing.T) {
	var i InkBunny

	if err := i.Login("guest", ""); err != nil {
		t.Fatal(err)
	}

	subs, err := i.Feed(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(subs) <= 0 {
		t.Fatal("sub length <= 0")
	}
}
