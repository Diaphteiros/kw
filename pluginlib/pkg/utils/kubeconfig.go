package utils

import (
	"fmt"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ParseKubeconfigFromFile reads the given file and parses it as kubeconfig.
func ParseKubeconfigFromFile(path string) (*clientcmdapi.Config, error) {
	data, err := vfs.ReadFile(fs.FS, path)
	if err != nil {
		if vfs.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("unable to read kubeconfig file: %w", err)
	}
	return ParseKubeconfig(data)
}

// ParseKubeconfigFromFileWithClient works like ParseKubeconfigFromFile, but it returns a client.Client in addition to the parsed kubeconfig.
func ParseKubeconfigFromFileWithClient(path string) (*clientcmdapi.Config, client.Client, error) {
	data, err := vfs.ReadFile(fs.FS, path)
	if err != nil {
		if vfs.IsNotExist(err) {
			return nil, nil, err
		}
		return nil, nil, fmt.Errorf("unable to read kubeconfig file: %w", err)
	}
	return ParseKubeconfigWithClient(data)
}

// ParseKubeconfig parses the given data as a kubeconfig and returns the parsed config.
func ParseKubeconfig(data []byte) (*clientcmdapi.Config, error) {
	return clientcmd.Load(data)
}

// ParseKubeconfigWithClient works like ParseKubeconfig, but it returns a client.Client in addition to the parsed kubeconfig.
func ParseKubeconfigWithClient(data []byte) (*clientcmdapi.Config, client.Client, error) {
	kcfg, err := ParseKubeconfig(data)
	if err != nil {
		return nil, nil, err
	}
	rest, err := clientcmd.RESTConfigFromKubeConfig(data)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating rest config from kubeconfig: %w", err)
	}
	c, err := client.New(rest, client.Options{})
	if err != nil {
		return nil, nil, fmt.Errorf("error creating client from kubeconfig: %w", err)
	}
	return kcfg, c, nil
}

// GetCurrentApiserverHost returns the current apiserver host from the given kubeconfig.
func GetCurrentApiserverHost(kcfg *clientcmdapi.Config) (string, error) {
	current := kcfg.CurrentContext
	curCon, ok := kcfg.Contexts[current]
	if !ok {
		return "", fmt.Errorf("current context '%s' not found in kubeconfig", current)
	}
	curCluster, ok := kcfg.Clusters[curCon.Cluster]
	if !ok {
		return "", fmt.Errorf("current cluster '%s' not found in kubeconfig", curCon.Cluster)
	}
	return curCluster.Server, nil
}

// MarshalKubeconfig returns the yaml representation of the given kubeconfig.
func MarshalKubeconfig(kcfg *clientcmdapi.Config) ([]byte, error) {
	return clientcmd.Write(*kcfg)
}
