package kubernetes

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"log"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
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

// GetCurrentNamespace --
func GetCurrentNamespace() (string, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeContainer, err := shouldUseContainerMode()
		if err != nil {
			return "", err
		}
		if kubeContainer {
			return getNamespaceFromKubernetesContainer()
		}
	}
	if kubeconfig == "" {
		var err error
		kubeconfig, err = getDefaultKubeConfigFile()
		if err != nil {
			fmt.Printf("ERROR: cannot get information about current user: %v", err)
		}
	}
	if kubeconfig == "" {
		return "default", nil
	}

	data, err := ioutil.ReadFile(kubeconfig)
	if err != nil {
		return "", err
	}
	conf := clientcmdapi.NewConfig()
	if len(data) == 0 {
		return "", errors.New("kubernetes config file is empty")
	}

	decoded, _, err := clientcmdlatest.Codec.Decode(data, &schema.GroupVersionKind{Version: clientcmdlatest.Version, Kind: "Config"}, conf)
	if err != nil {
		return "", err
	}

	clientcmdconfig := decoded.(*clientcmdapi.Config)

	cc := clientcmd.NewDefaultClientConfig(*clientcmdconfig, &clientcmd.ConfigOverrides{})
	ns, _, err := cc.Namespace()
	return ns, err
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

func getNamespaceFromKubernetesContainer() (string, error) {
	var nsba []byte
	var err error
	if nsba, err = ioutil.ReadFile(inContainerNamespaceFile); err != nil {
		return "", err
	}
	return string(nsba), nil
}
