package meta

import (
	"slices"
	"strings"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"

	"github.com/Diaphteiros/kw/pkg/cmdgroups"
	"github.com/Diaphteiros/kw/pkg/config"

	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	"github.com/Diaphteiros/kw/pluginlib/pkg/selector"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var NamespaceCmd = &cobra.Command{
	Use:     "namespace [<namespace>]",
	Aliases: []string{"ns"},
	Args:    cobra.RangeArgs(0, 1),
	GroupID: cmdgroups.Meta,
	Short:   "Change the default namespace in the current context of the kubeconfig",
	Long: `Change the default namespace in the current context of the kubeconfig.

This is basically the same thing as running 'kubectl config set-context --current --namespace=<namespace>'.
If called without any argument, the command fetches the namespaces from the currently selected cluster and prompts for a selection.

Note that this command does change the kubeconfig file, but doesn't create a new kubeswitcher history entry.`,
	Run: func(cmd *cobra.Command, args []string) {
		kcfg, c, err := libutils.ParseKubeconfigFromFileWithClient(config.Runtime.KubeconfigPath())
		if err != nil {
			if vfs.IsNotExist(err) {
				libutils.Fatal(1, "kubeconfig not set, run another kw command to switch to a kubeconfig first\n")
			}
			libutils.Fatal(1, "error parsing kubeconfig: %w\n", err)
		}
		curCtx, ok := kcfg.Contexts[kcfg.CurrentContext]
		if !ok {
			libutils.Fatal(1, "invalid kubeconfig: current context '%s' not found\n", kcfg.CurrentContext)
		}
		namespace := ""
		if len(args) == 0 {
			// fetch namespaces
			nsl := &corev1.NamespaceList{}
			if err := c.List(cmd.Context(), nsl); err != nil {
				libutils.Fatal(1, "error fetching namespaces: %w\n", err)
			}
			namespaces := make([]string, len(nsl.Items))
			for i, ns := range nsl.Items {
				namespaces[i] = ns.Name
			}
			slices.SortFunc(namespaces, func(a, b string) int {
				return -strings.Compare(a, b)
			})
			namespaces = append([]string{""}, namespaces...) // add empty namespace to let user select default namespace

			// let user select namespace
			_, namespace, err = selector.New[string]().
				WithPrompt("Select Namespace: ").
				WithFatalOnAbort("No namespace selected.").
				WithFatalOnError("error selecting namespace: %w").
				From(namespaces, func(elem string) string { return elem }).
				Select()
			if err != nil {
				libutils.Fatal(1, err.Error())
			}
		} else {
			namespace = args[0]
		}
		curCtx.Namespace = namespace
		kcfgData, err := libutils.MarshalKubeconfig(kcfg)
		if err != nil {
			libutils.Fatal(1, "error marshaling kubeconfig: %w\n", err)
		}
		if err := vfs.WriteFile(fs.FS, config.Runtime.KubeconfigPath(), kcfgData, vfs.ModePerm); err != nil {
			libutils.Fatal(1, "error writing kubeconfig: %w\n", err)
		}
	},
}
