package k8s

import (
	"context"
	"path"
	"sort"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetContexts() ([]string, string, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	if config, error := kubeConfig.RawConfig(); error == nil {
		contexts := []string{}
		for name, _ := range config.Contexts {
			contexts = append(contexts, name)
		}
		sort.Strings(contexts)
		return contexts, config.CurrentContext, nil
	} else {
		return []string{}, "", error
	}
}

func GetNamespaces() ([]string, string, error) {
	currentNamespace := "default"
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	if config, error := kubeConfig.RawConfig(); error == nil {
		if config.Contexts[config.CurrentContext] != nil {
			currentNamespace = config.Contexts[config.CurrentContext].Namespace
		}
	} else {
		return []string{}, "", error
	}
	kubeConfigPath := path.Join(homedir.HomeDir(), ".kube", "config")
	if config, error := clientcmd.BuildConfigFromFlags("", kubeConfigPath); error == nil {
		config.Timeout = time.Millisecond * 5000
		if clientset, error := kubernetes.NewForConfig(config); error == nil {
			if response, error := clientset.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{}); error == nil {
				items := []string{}
				for _, namespace := range response.Items {
					items = append(items, namespace.ObjectMeta.Name)
				}
				return items, currentNamespace, nil
			} else {
				return []string{}, "", error
			}
		} else {
			return []string{}, "", error
		}
	} else {
		return []string{}, "", error
	}
}

func ChangeContext(context string) error {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	if config, error := kubeConfig.RawConfig(); error == nil {
		configAccess := kubeConfig.ConfigAccess()
		config.CurrentContext = context
		return clientcmd.ModifyConfig(configAccess, config, true)
	} else {
		return error
	}
}

func ChangeNamespace(context, namespace string) error {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	if config, error := kubeConfig.RawConfig(); error == nil {
		configAccess := kubeConfig.ConfigAccess()
		if config.Contexts[context] != nil {
			config.Contexts[context].Namespace = namespace
		}
		return clientcmd.ModifyConfig(configAccess, config, true)
	} else {
		return error
	}
}
