package importcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/manipulate/manipcli"
)

// MakeImportCmd returns the import sub command.
func MakeImportCmd(mmaker manipcli.ManipulatorMaker) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "import",
		Short:            "Manage import files",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			fFile := viper.GetString("file")
			fSet := viper.GetStringSlice("set")
			fValues := viper.GetString("values")
			fDelete := viper.GetBool("delete")
			fRender := viper.GetBool("render")

			_ = fSet
			_ = fFile
			_ = fValues
			_ = fDelete
			_ = fRender

			fmt.Println("Not implemented yet. You can use api create import for now")

			return nil
		},
	}

	cmd.Flags().StringP("file", "F", "", "Path or URL to the import file.")
	cmd.Flags().StringSliceP("set", "S", nil, "Set the value for one key in the template.")
	cmd.Flags().StringP("values", "V", "", "Path to a values file.")
	cmd.Flags().BoolP("delete", "D", false, "Delete the previously created data declared in the import file.")
	cmd.Flags().BoolP("render", "R", false, "Only renders the templated information locally.")

	return cmd
}
