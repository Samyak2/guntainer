package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Samyak2/guntainer/core"
)

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "Build a new container from a Gunfile",
	Long:  "Executes the Gunfile to build a new container image.",
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			core.Build(args[1], args[0])
		} else {
			core.Build("Gunfile", args[0])
		}
	},
}

