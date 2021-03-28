package internal

import (
	"os"
	"fmt"
	"github.com/null93/kube-switch/sdk/k8s"
	"github.com/null93/kube-switch/sdk/prompt"
)

func FatalError ( message string ) {
	fmt.Println ( message )
	os.Exit ( 1 )
}

func PickContextAndNamespace ( useCurrentContext bool ) {
	contexts, currentContext, error := k8s.GetContexts ()
	if error != nil {
		FatalError ("Error: Could not load configured contexts.")
	}
	if !useCurrentContext {
		pickedContext, exited := prompt.Pick ( prompt.Options {
			Choices: contexts,
			DefaultValue: currentContext,
			Header: "Pick Kubernetes Context",
			SelectionPrefix: "",
		})
		if exited {
			return
		}
		currentContext = pickedContext
		if error := k8s.ChangeContext ( pickedContext ); error != nil {
			FatalError ("Error: Failed to write current context to config.")
		}
	}
	namespaces, currentNamespace, error := k8s.GetNamespaces ()
	if error != nil {
		FatalError ("Error: Failed to load list of namespaces in current context.")
	}
	pickedNamespace, _ := prompt.Pick ( prompt.Options {
		Choices: namespaces,
		DefaultValue: currentNamespace,
		Header: "Pick Kubernetes Context",
		SelectionPrefix: currentContext + " > ",
	})
	if error := k8s.ChangeNamespace ( currentContext, pickedNamespace ); error != nil {
		FatalError ("Error: Failed to write current namespace to config.")
	}
}
