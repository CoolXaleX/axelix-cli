package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	heapDumpLive bool
	heapDumpOut  string
)

var heapDumpCmd = &cobra.Command{
	Use:   "heap-dump",
	Short: "Download a heap dump",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.DownloadHeapDump(heapDumpLive)
		if err != nil {
			return err
		}
		outFile := heapDumpOut
		if outFile == "" {
			outFile = fmt.Sprintf("heapdump-%d.hprof", time.Now().Unix())
		}
		if err := os.WriteFile(outFile, data, 0644); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Heap dump saved to %s\n", outFile)
		return nil
	},
}

func init() {
	heapDumpCmd.Flags().BoolVar(&heapDumpLive, "live", false, "Dump only live objects (skip unreachable/GC-able objects)")
	heapDumpCmd.Flags().StringVar(&heapDumpOut, "out", "", "Output file path (default: heapdump-<timestamp>.hprof)")
	rootCmd.AddCommand(heapDumpCmd)
}
