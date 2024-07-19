package main

import (
	"github.com/matteopic/hashdir/internal"
	"github.com/spf13/cobra"
)

func init() {
	LoadStatsCmd.Flags().String("index", "checksum.txt", "Specify the filename of the index file. This index file contains the checksum data that the application processes.")
}

var LoadStatsCmd = &cobra.Command{
	Use:  "stats [--index indexfile]",
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, err := cmd.Flags().GetString("index")
		if err != nil {
			return err
		}
		s := internal.NewScanner(internal.WithIndexFilename(filename))
		if err := s.LoadIndex(); err != nil {
			return err
		}

		s.PrintStats()
		return nil
	},
}
