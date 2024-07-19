package main

import (
	"fmt"
	"os"

	"github.com/matteopic/hashdir/internal"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.Flags().String("index", "checksum.txt", "Output file name")

	RootCmd.AddCommand(LoadStatsCmd)
}

var RootCmd = &cobra.Command{
	Use:   "hashdir",
	Short: "hashdir is a very fast file duplicate finder",
	Long: `A Fast hashing file tool that can be used to find duplicate files performing a checksum of each file and indexing them.
Bring to you with love by matteopic and friends in Go.
Complete documentation is available at https://github.com/matteopic/hashdir`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, err := cmd.Flags().GetString("index")
		if err != nil {
			return err
		}

		s := internal.NewScanner(internal.WithIndexFilename(filename))
		if err := s.Scan(args); err != nil {
			return err
		}

		s.PrintStats()
		return nil
	},
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
