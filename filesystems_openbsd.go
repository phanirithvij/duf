//go:build openbsd
// +build openbsd

package duf

func isFuseFs(m Mount) bool {
	//FIXME: implement
	return false
}

func isNetworkFs(m Mount) bool {
	//FIXME: implement
	return false
}

func isSpecialFs(m Mount) bool {
	return m.Fstype == "devfs"
}

func IsHiddenFs(m Mount) bool {
	return false
}
