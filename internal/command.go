package internal

import (
	"os"
	"fmt"
	"strings"
	"github.com/spf13/cobra"
)

var command = &cobra.Command {
	Use: "kube-switch",
	Version: "2.0.3",
	Short: "Switch between Kubernetes context & namespace using an interactive menu",
	Example: strings.Join ( [] string {
		"  kube-switch",
		"  kube-switch -c",
	}, "\n" ),
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
