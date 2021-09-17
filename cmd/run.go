package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Samyak2/guntainer/core"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an existing container",
	Long:  "Run the archive of root FS as a container.",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		core.Run(args[0], args[1], args[2:])
	},
}
