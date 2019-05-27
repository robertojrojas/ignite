package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"io"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdLogs gets the logs for a Firecracker VM
func NewCmdLogs(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [id]",
		Short: "Gets the logs for a Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunLogs(out, cmd, args)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunLogs(out io.Writer, cmd *cobra.Command, args []string) error {
	var md *vmmd.VMMetadata

	// Match a single VM using the VMFilter
	if matches, err := filter.NewFilterer(vmmd.NewVMFilter(args[0]), metadata.VM.Path(), vmmd.LoadVMMetadata); err == nil {
		if filterable, err := matches.Single(); err == nil {
			if md, err = vmmd.ToVMMetadata(filterable); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	// Check if the VM is running
	if !md.Running() {
		return fmt.Errorf("%s is not running", md.ID)
	}

	dockerArgs := []string{
		"logs",
		md.ID,
	}

	// Fetch the VM logs from docker
	output, err := util.ExecuteCommand("docker", dockerArgs...)
	if err != nil {
		return fmt.Errorf("failed to get logs for VM %q: %v", md.ID, err)
	}

	// Print the ID and the VM logs
	fmt.Println(md.ID)
	fmt.Println(output)
	return nil
}
