package kubernetes

import (
	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const inContainerNamespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

func GetKubernetesConfigOrDie() *rest.Config {
	initialize(os.Getenv("KUBECONFIG"))
	return config.GetConfigOrDie()
}

func GetKubernetesConfig() (*rest.Config, error) {
	initialize(os.Getenv("KUBECONFIG"))
	return config.GetConfig()
}

// init initialize the k8s client for usage outside the cluster
func initialize(kubeconfig string) {
	if kubeconfig == "" {
		// skip out-of-cluster initialization if inside the container
		if kc, err := shouldUseContainerMode(); kc && err == nil {
			return
		} else if err != nil {
			log.Printf("could not determine if running in a container: %v", err)
		}
		var err error
		kubeconfig, err = getDefaultKubeConfigFile()
		if err != nil {
			panic(err)
		}
	}
	os.Setenv("KUBECONFIG", kubeconfig)
}

func getDefaultKubeConfigFile() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".kube", "config"), nil
}

func shouldUseContainerMode() (bool, error) {
	// When kube config is set, container mode is not used
	if os.Getenv("KUBECONFIG") != "" {
		return false, nil
	}
	// Use container mode only when the kubeConfigFile does not exist and the container namespace file is present
	configFile, err := getDefaultKubeConfigFile()
	if err != nil {
		return false, err
	}
	configFilePresent := true
	_, err = os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		configFilePresent = false
	} else if err != nil {
		return false, err
	}
	if !configFilePresent {
		_, err := os.Stat(inContainerNamespaceFile)
		if os.IsNotExist(err) {
			return false, nil
		}
		return true, err
	}
	return false, nil
}
