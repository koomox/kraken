// +build !linux,!darwin,!windows,!freebsd,!dragonfly,!netbsd,!openbsd

package mem

func sysTotalMemory() uint64 {
	return 0
}
func sysFreeMemory() uint64 {
	return 0
}
