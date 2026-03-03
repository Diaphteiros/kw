package utils

import (
	"fmt"
	"path/filepath"

	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/Diaphteiros/kw/pkg/config"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

// GetId extracts the id from the given state directory.
// The id is either the id, if that file exists within the directory,
// or the apiserver url from the kubeconfig.
// Missing files do not cause an error.
// If neither file exists, it returns '<unknown>'.
func GetId(path string) (string, error) {
	idFilePath := filepath.Join(path, config.IdFileName)
	res := ""
	resRaw, err := vfs.ReadFile(fs.FS, idFilePath)
	if err == nil {
		res = string(resRaw)
	} else {
		if !vfs.IsNotExist(err) {
			return "", fmt.Errorf("error reading id file '%s': %w", idFilePath, err)
		}
		// fallback: try to extract apiserver url and use it instead of id
		kcfgFilePath := filepath.Join(path, config.KubeconfigFileName)
		kcfg, err := libutils.ParseKubeconfigFromFile(kcfgFilePath)
		if err != nil {
			if !vfs.IsNotExist(err) {
				return "", fmt.Errorf("error parsing kubeconfig from '%s': %w", kcfgFilePath, err)
			}
			res = "<unknown>"
		} else {
			host, err := libutils.GetCurrentApiserverHost(kcfg)
			if err != nil {
				return "", fmt.Errorf("error getting current apiserver host from kubeconfig: %w", err)
			}
			res = host
		}
	}
	return res, nil
}
