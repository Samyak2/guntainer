package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Samyak2/guntainer/core"
)

var buildCommand = &cobra.Command{
	Use:   "build [image-path] (Gunfile-path)",
	Short: "Build a new container from a Gunfile",
	Long:  `Executes the Gunfile to build a new container image.

 image-path:	is where the resulting container image will be saved.
 Gunfile-path:	(optional) is the path to the Gunfile to use for building the image.`,
	DisableFlagsInUseLine: true,

	Args: cobra.RangeArgs(1, 2),
	Run: func(_ *cobra.Command, args []string) {
		if len(args) > 1 {
			core.Build(args[1], args[0])
		} else {
			core.Build("Gunfile", args[0])
		}
	},
}

