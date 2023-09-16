package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func GetVersionCmd(v string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of fjira",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("fjira version: %s", v)
		},
	}
}
