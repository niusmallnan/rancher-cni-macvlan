package macfinder

//MACFinder is used to get MAC address given a container ID.
type MACFinder interface {
	GetMACAddress(cid, rancherid string) string
}
