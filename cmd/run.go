package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Samyak2/guntainer/core"
)

var runCmd = &cobra.Command{
	Use:   "run [image-path] [command]",
	Short: "Run an existing container",
	Long:  `Run the archive of root FS as a container.

 image-path:	 is an archive (tar, zip, etc.) of a root filesystem which will be 'chroot'ed into.
 command:	is a program *inside* the container root FS that is run (note that PATH is not set unless you execute a shell).

Everything after the 'command' is passed as arguments directly to the command inside the container.`,
	DisableFlagsInUseLine: true,

	Args: cobra.MinimumNArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		core.Run(args[0], args[1], args[2:])
	},
}
