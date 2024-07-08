package compcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// New returns a new compcmd Command.
func New() *cobra.Command {
	return &cobra.Command{
		Use:                   "completion [bash|zsh|fish]",
		Short:                 "Generate completion script",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			}
			return fmt.Errorf("unamanaged shell")
		},
	}
}
