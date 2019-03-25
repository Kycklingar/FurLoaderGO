package inkbunny

import "github.com/kycklingar/FurLoaderGO/dli"

func init() {
	var inkbunny InkBunny
	dli.Logins["inkbunny"] = &inkbunny
	dli.Watchers["inkbunny"] = &inkbunny
	dli.Galleries["inkbunny"] = &inkbunny
}
