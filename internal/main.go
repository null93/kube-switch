package internal

import (
	"fmt"
	"os"

	"github.com/null93/aws-knox/sdk/picker"
	"github.com/null93/kube-switch/sdk/k8s"
)

func FatalError(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func getInitialSelectionIndex(haystack []string, needle string) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}
	return 0
}

func PickContextAndNamespace(useCurrentContext bool) {
	contexts, currentContext, error := k8s.GetContexts()
	if error != nil {
		FatalError("Error: Could not load configured contexts.")
	}
	if !useCurrentContext {
		p := picker.NewPicker()
		p.WithMaxHeight(10)
		p.WithFilterStrategy("fuzzy")
		p.WithEmptyMessage("No Kubernetes Contexts Found")
		p.WithTitle("Pick Kubernetes Context")
		p.WithHeaders("CONTEXT")
		p.WithInitialIndex(getInitialSelectionIndex(contexts, currentContext))
		for _, name := range contexts {
			p.AddOption(name, name)
		}
		selection, _ := p.Pick("")
		if selection == nil {
			return
		}
		currentContext = selection.Value.(string)
		if error := k8s.ChangeContext(currentContext); error != nil {
			FatalError("Error: Failed to write current context to config.")
		}
	}
	namespaces, currentNamespace, error := k8s.GetNamespaces()
	if error != nil {
		FatalError("Error: Failed to load list of namespaces in current context.")
	}
	p := picker.NewPicker()
	p.WithMaxHeight(10)
	p.WithFilterStrategy("fuzzy")
	p.WithEmptyMessage("No Kubernetes Namespaces Found")
	p.WithTitle("Pick Kubernetes Namespace")
	p.WithHeaders("NAMESPACE")
	p.WithInitialIndex(getInitialSelectionIndex(namespaces, currentNamespace))
	for _, name := range namespaces {
		p.AddOption(name, name)
	}
	selection, _ := p.Pick("")
	if selection == nil {
		return
	}
	currentNamespace = selection.Value.(string)
	if error := k8s.ChangeNamespace(currentContext, currentNamespace); error != nil {
		FatalError("Error: Failed to write current namespace to config.")
	}
}
