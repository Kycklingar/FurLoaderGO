package fa

import "testing"

func TestWatchlist(t *testing.T) {
	var fa = NewFurAffinity()
	users, err := fa.Watchlist("s-nina")
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(users) <= 0 {
		t.Fatal("users returned was <= 0")
	}
}
