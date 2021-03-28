package internal

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
)

var command = &cobra.Command {
	Use: "kube-switch",
	Version: "2.0.0",
	Short: "Example",
	Run: func ( cmd * cobra.Command, args [] string ) {
		current, _ := cmd.Flags ().GetBool ("current-context")
		PickContextAndNamespace ( current )
	},
}

func Execute () {
	if error := command.Execute (); error != nil {
		fmt.Println ( error )
		os.Exit ( 1 )
	}
}

func init () {
	command.Flags ().BoolP ( "current-context", "c", false, "use current context, only pick namespace" )
}
