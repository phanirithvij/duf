//go:build freebsd
// +build freebsd

package duf

func isFuseFs(m Mount) bool {
	//FIXME: implement
	return false
}

func isNetworkFs(m Mount) bool {
	fs := []string{"nfs", "smbfs"}

	for _, v := range fs {
		if m.Fstype == v {
			return true
		}
	}

	return false
}

func isSpecialFs(m Mount) bool {
	fs := []string{"devfs", "tmpfs", "linprocfs", "linsysfs", "fdescfs", "procfs"}

	for _, v := range fs {
		if m.Fstype == v {
			return true
		}
	}

	return false
}

func IsHiddenFs(m Mount) bool {
	return false
}
