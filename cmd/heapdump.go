package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
)

func newHeapDumpCmd(cl *client.Client) *cobra.Command {
	var live bool
	var outFile string

	cmd := &cobra.Command{
		Use:   "heap-dump",
		Short: "Download a heap dump",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cl.DownloadHeapDump(live)
			if err != nil {
				return err
			}
			path := outFile
			if path == "" {
				path = fmt.Sprintf("heapdump-%d.hprof", time.Now().Unix())
			}
			if err := os.WriteFile(path, data, 0644); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "✓ heap dump saved to %s\n", path)
			return nil
		},
	}
	cmd.Flags().BoolVar(&live, "live", false, "Dump only live objects (skip unreachable/GC-able objects)")
	cmd.Flags().StringVar(&outFile, "out", "", "Output file path (default: heapdump-<timestamp>.hprof)")
	return cmd
}
