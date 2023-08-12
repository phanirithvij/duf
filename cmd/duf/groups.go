package main

import (
	. "github.com/muesli/duf"
	"strings"
)

// FilterOptions contains all filters.
type FilterOptions struct {
	HiddenDevices map[string]struct{}
	OnlyDevices   map[string]struct{}

	HiddenFilesystems map[string]struct{}
	OnlyFilesystems   map[string]struct{}

	HiddenMountPoints map[string]struct{}
	OnlyMountPoints   map[string]struct{}
}

// renderTables renders all tables.
func renderTables(m []Mount, filters FilterOptions, opts TableOptions) {
	deviceMounts := make(map[string][]Mount)
	hasOnlyDevices := len(filters.OnlyDevices) != 0

	_, hideLocal := filters.HiddenDevices[LocalDevice]
	_, hideNetwork := filters.HiddenDevices[NetworkDevice]
	_, hideFuse := filters.HiddenDevices[FuseDevice]
	_, hideSpecial := filters.HiddenDevices[SpecialDevice]
	_, hideLoops := filters.HiddenDevices[LoopsDevice]
	_, hideBinds := filters.HiddenDevices[BindsMount]

	_, onlyLocal := filters.OnlyDevices[LocalDevice]
	_, onlyNetwork := filters.OnlyDevices[NetworkDevice]
	_, onlyFuse := filters.OnlyDevices[FuseDevice]
	_, onlySpecial := filters.OnlyDevices[SpecialDevice]
	_, onlyLoops := filters.OnlyDevices[LoopsDevice]
	_, onlyBinds := filters.OnlyDevices[BindsMount]

	// sort/filter devices
	for _, v := range m {
		if len(filters.OnlyFilesystems) != 0 {
			// skip not onlyFs
			if _, ok := filters.OnlyFilesystems[strings.ToLower(v.Fstype)]; !ok {
				continue
			}
		} else {
			// skip hideFs
			if _, ok := filters.HiddenFilesystems[strings.ToLower(v.Fstype)]; ok {
				continue
			}
		}

		// skip hidden devices
		if IsHiddenFs(v) && !*all {
			continue
		}

		// skip bind-mounts
		if strings.Contains(v.Opts, "bind") {
			if (hasOnlyDevices && !onlyBinds) || (hideBinds && !*all) {
				continue
			}
		}

		// skip loop devices
		if strings.HasPrefix(v.Device, "/dev/loop") {
			if (hasOnlyDevices && !onlyLoops) || (hideLoops && !*all) {
				continue
			}
		}

		// skip special devices
		if v.Blocks == 0 && !*all {
			continue
		}

		// skip zero size devices
		if v.BlockSize == 0 && !*all {
			continue
		}

		// skip not only mount point
		if len(filters.OnlyMountPoints) != 0 {
			if !findInKey(v.Mountpoint, filters.OnlyMountPoints) {
				continue
			}
		}

		// skip hidden mount point
		if len(filters.HiddenMountPoints) != 0 {
			if findInKey(v.Mountpoint, filters.HiddenMountPoints) {
				continue
			}
		}

		t := DeviceType(v)
		deviceMounts[t] = append(deviceMounts[t], v)
	}

	// print tables
	for _, devType := range groups {
		mounts := deviceMounts[devType]

		shouldPrint := *all
		if !shouldPrint {
			switch devType {
			case LocalDevice:
				shouldPrint = (hasOnlyDevices && onlyLocal) || (!hasOnlyDevices && !hideLocal)
			case NetworkDevice:
				shouldPrint = (hasOnlyDevices && onlyNetwork) || (!hasOnlyDevices && !hideNetwork)
			case FuseDevice:
				shouldPrint = (hasOnlyDevices && onlyFuse) || (!hasOnlyDevices && !hideFuse)
			case SpecialDevice:
				shouldPrint = (hasOnlyDevices && onlySpecial) || (!hasOnlyDevices && !hideSpecial)
			}
		}

		if shouldPrint {
			printTable(devType, mounts, opts)
		}
	}
}
