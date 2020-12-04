package list

import (
	"fmt"
	"github.com/hpcng/warewulf/internal/pkg/node"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strings"
)

func CobraRunE(cmd *cobra.Command, args []string) error {
	var err error
	var nodes []node.NodeInfo

	n, err := node.New()
	if err != nil {
		wwlog.Printf(wwlog.ERROR, "Could not open node configuration: %s\n", err)
		os.Exit(1)
	}

	if len(args) > 0 {
		nodes, err = n.SearchByNameList(args)
	} else {
		nodes, err = n.FindAllNodes()
	}
	if err != nil {
		wwlog.Printf(wwlog.ERROR, "Cloud not get nodeList: %s\n", err)
		os.Exit(1)
	}

	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].ClusterName.Get() < nodes[j].ClusterName.Get() {
			return true
		} else if nodes[i].ClusterName.Get() == nodes[j].ClusterName.Get() {
			if nodes[i].Id.Get() < nodes[j].Id.Get() {
				return true
			}
		}
		return false
	})

	if ShowAll == true {
		for _, node := range nodes {
			fmt.Printf("################################################################################\n")
			fmt.Printf("%-20s %-18s %-12s %s\n", "NODE", "FIELD", "PROFILE", "VALUE")
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Id", node.Id.Source(), node.Id.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Comment", node.Comment.Source(), node.Comment.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "ClusterName", node.ClusterName.Source(), node.ClusterName.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Profiles", "--", strings.Join(node.Profiles, ","))

			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Vnfs", node.Vnfs.Source(), node.Vnfs.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "KernelVersion", node.KernelVersion.Source(), node.KernelVersion.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "KernelArgs", node.KernelArgs.Source(), node.KernelArgs.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "RuntimeOverlay", node.RuntimeOverlay.Source(), node.RuntimeOverlay.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "SystemOverlay", node.SystemOverlay.Source(), node.SystemOverlay.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Ipxe", node.Ipxe.Source(), node.Ipxe.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "IpmiIpaddr", node.IpmiIpaddr.Source(), node.IpmiIpaddr.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "IpmiNetmask", node.IpmiNetmask.Source(), node.IpmiNetmask.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "IpmiUserName", node.IpmiUserName.Source(), node.IpmiUserName.Print())

			for name, netdev := range node.NetDevs {
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":HWADDR", netdev.Hwaddr.Source(), netdev.Hwaddr.Print())
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":IPADDR", netdev.Ipaddr.Source(), netdev.Ipaddr.Print())
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":NETMASK", netdev.Netmask.Source(), netdev.Netmask.Print())
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":GATEWAY", netdev.Gateway.Source(), netdev.Gateway.Print())
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":TYPE", netdev.Type.Source(), netdev.Type.Print())
			}

			//			v := reflect.ValueOf(node)
			//			typeOfS := v.Type()
			//			for i := 0; i< v.NumField(); i++ {
			//				//TODO: Fix for NetDevs and Interface should print Fprint() method
			//				fmt.Printf("%-25s %s = %#v\n", node.Fqdn.Get(), typeOfS.Field(i).Name, v.Field(i).Interface())
			//			}
		}

	} else if ShowNet == true {
		fmt.Printf("%-22s %-6s %-18s %-15s %-15s\n", "NODE NAME", "DEVICE", "HWADDR", "IPADDR", "GATEWAY")
		fmt.Println(strings.Repeat("=", 80))

		for _, node := range nodes {
			if len(node.NetDevs) > 0 {
				for name, dev := range node.NetDevs {
					fmt.Printf("%-22s %-6s %-18s %-15s %-15s\n", node.Id.Get(), name, dev.Hwaddr, dev.Ipaddr, dev.Gateway)
				}
			} else {
				fmt.Printf("%-22s %-6s %-18s %-15s %-15s\n", node.Id.Get(), "--", "--", "--", "--")
			}
		}

	} else if ShowIpmi == true {
		fmt.Printf("%-22s %-16s %-20s %-20s\n", "NODE NAME", "IPMI IPADDR", "IPMI USERNAME", "IPMI PASSWORD")
		fmt.Println(strings.Repeat("=", 80))

		for _, node := range nodes {
			fmt.Printf("%-22s %-16s %-20s %-20s\n", node.Id.Get(), node.IpmiIpaddr.Print(), node.IpmiUserName.Print(), node.IpmiPassword.Print())
		}

	} else if ShowLong == true {
		fmt.Printf("%-22s %-26s %-35s %s\n", "NODE NAME", "KERNEL VERSION", "VNFS IMAGE", "OVERLAYS (S/R)")
		fmt.Println(strings.Repeat("=", 120))

		for _, node := range nodes {
			fmt.Printf("%-22s %-26s %-35s %s\n", node.Id.Get(), node.KernelVersion.Print(), node.Vnfs.Print(), node.SystemOverlay.Print()+"/"+node.RuntimeOverlay.Print())
		}

	} else {
		fmt.Printf("%-22s %-26s %s\n", "NODE NAME", "PROFILES", "NETWORK")

		for _, node := range nodes {
			var netdevs []string
			if len(node.NetDevs) > 0 {
				for name, dev := range node.NetDevs {
					netdevs = append(netdevs, fmt.Sprintf("%s:%s", name, dev.Ipaddr.Get()))
				}
			}
			sort.Strings(netdevs)

			fmt.Printf("%-22s %-26s %s\n", node.Id.Get(), strings.Join(node.Profiles, ","), strings.Join(netdevs, ", "))
		}

	}

	return nil
}
