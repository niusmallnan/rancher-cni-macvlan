package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/containernetworking/cni/pkg/ip"
	"github.com/containernetworking/cni/pkg/ns"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
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

func checkIfContainerInterfaceExists(args *skel.CmdArgs) bool {
	err := ns.WithNetNSPath(args.Netns, func(_ ns.NetNS) error {
		_, err := netlink.LinkByName(args.IfName)
		if err != nil {
			return fmt.Errorf("failed to lookup %q: %v", args.IfName, err)
		}
		return nil
	})

	if err == nil {
		return true
	}
	return false
}

func setInterfaceDown(args *skel.CmdArgs) error {
	err := ns.WithNetNSPath(args.Netns, func(_ ns.NetNS) error {
		link, err := netlink.LinkByName(args.IfName)
		if err != nil {
			return fmt.Errorf("failed to lookup %q: %v", args.IfName, err)
		}
		err = netlink.LinkSetDown(link)
		if err != nil {
			return fmt.Errorf("failed to setdown %q: %v", args.IfName, err)
		}
		return nil
	})
	return err

}

func configureInterface(ifName string, res *types.Result) error {
	link, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("failed to lookup %q: %v", ifName, err)
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("failed to set %q UP: %v", ifName, err)
	}

	// TODO(eyakubovich): IPv6
	addr := &netlink.Addr{IPNet: &res.IP4.IP, Label: ""}
	if err = netlink.AddrAdd(link, addr); err != nil {
		if err.Error() == "file exists" {
			logrus.Infof("rancher-cni-macvlan: Interface %q already has IP address: %v, no worries", ifName, addr)
		} else {
			return fmt.Errorf("failed to add IP addr to %q: %v", ifName, err)
		}
	}

	for _, r := range res.IP4.Routes {
		gw := r.GW
		if gw == nil {
			gw = res.IP4.Gateway
		}
		if err = ip.AddRoute(&r.Dst, gw, link); err != nil {
			// we skip over duplicate routes as we assume the first one wins
			if !os.IsExist(err) {
				return fmt.Errorf("failed to add route '%v via %v dev %v': %v", r.Dst, gw, ifName, err)
			}
		}
	}

	return nil
}
