package set

import "github.com/spf13/cobra"

var (
	baseCmd = &cobra.Command{
		Use:   "set",
		Short: "Set profile configurations",
		Long:  "Profile configurations ",
		RunE:  CobraRunE,
	}
	SetAll            bool
	SetComment        string
	SetVnfs           string
	SetKernel         string
	SetClusterName    string
	SetIpxe           string
	SetRuntimeOverlay string
	SetSystemOverlay  string
	SetIpmiNetmask    string
	SetIpmiUsername   string
	SetIpmiPassword   string
	SetNetDev         string
	SetIpaddr         string
	SetNetmask        string
	SetGateway        string
	SetHwaddr         string
	SetNetDevDel      bool
)

func init() {
	baseCmd.PersistentFlags().StringVar(&SetComment, "comment", "", "Set a comment for this node")
	baseCmd.PersistentFlags().StringVarP(&SetVnfs, "vnfs", "V", "", "Set node Virtual Node File System (VNFS)")
	baseCmd.PersistentFlags().StringVarP(&SetKernel, "kernel", "K", "", "Set Kernel version for nodes")
	baseCmd.PersistentFlags().StringVarP(&SetClusterName, "cluster", "C", "", "Set the node's cluster name")
	baseCmd.PersistentFlags().StringVarP(&SetIpxe, "ipxe", "P", "", "Set the node's iPXE template name")
	baseCmd.PersistentFlags().StringVarP(&SetRuntimeOverlay, "runtime", "R", "", "Set the node's runtime overlay")
	baseCmd.PersistentFlags().StringVarP(&SetSystemOverlay, "system", "S", "", "Set the node's system overlay")
	baseCmd.PersistentFlags().StringVar(&SetIpmiNetmask, "ipminetmask", "", "Set the node's IPMI netmask")
	baseCmd.PersistentFlags().StringVar(&SetIpmiUsername, "ipmiuser", "", "Set the node's IPMI username")
	baseCmd.PersistentFlags().StringVar(&SetIpmiPassword, "ipmipass", "", "Set the node's IPMI password")

	baseCmd.PersistentFlags().StringVarP(&SetNetDev, "netdev", "n", "", "Define the network device to configure")
	baseCmd.PersistentFlags().StringVarP(&SetIpaddr, "ipaddr", "I", "", "Set the node's network device IP address")
	baseCmd.PersistentFlags().StringVarP(&SetNetmask, "netmask", "M", "", "Set the node's network device netmask")
	baseCmd.PersistentFlags().StringVarP(&SetGateway, "gateway", "G", "", "Set the node's network device gateway")
	baseCmd.PersistentFlags().StringVarP(&SetHwaddr, "hwaddr", "H", "", "Set the node's network device HW address")
	baseCmd.PersistentFlags().BoolVar(&SetNetDevDel, "netdel", false, "Delete the node's network device")

	baseCmd.PersistentFlags().BoolVarP(&SetAll, "all", "a", false, "Set all profiles")
}

// GetRootCommand returns the root cobra.Command for the application.
func GetCommand() *cobra.Command {
	return baseCmd
}
