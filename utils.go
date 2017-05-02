package main

import (
	"fmt"
	"net"

	"github.com/rancher/rancher-cni-bridge/macfinder"
	"github.com/rancher/rancher-cni-bridge/macfinder/metadata"
	"github.com/vishvananda/netlink"
)

func setInterfaceMacAddress(ifName, mac string) error {
	link, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("failed to lookup %q: %v", ifName, err)
	}

	hwaddr, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("failed to parse MAC address: %v", err)
	}
	err = netlink.LinkSetHardwareAddr(link, hwaddr)
	if err != nil {
		return fmt.Errorf("failed to set hw address of interface: %v", err)
	}

	return nil
}

func findMACAddressForContainer(containerID, rancherID string) (string, error) {
	var mf macfinder.MACFinder
	mf, err := metadata.NewMACFinderFromMetadata()
	if err != nil {
		return "", err
	}
	macString := mf.GetMACAddress(containerID, rancherID)
	if macString == "" {
		return "", fmt.Errorf("No MAC address found")
	}

	return macString, nil
}
