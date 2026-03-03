package misc

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/pkg/config"

	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

const (
	SHELL_BASH  = "bash"
	SHELL_ZSH   = "zsh"
	SHELL_FISH  = "fish"
	SHELL_POWER = "powershell"
)

// KubectlEnvCmd represents the kubectl-env command
var KubectlEnvCmd = &cobra.Command{
	Use:     "kubectl-env",
	Aliases: []string{"env", "ke"},
	Args:    cobra.NoArgs,
	Short:   "Generate a script that points the KUBECONFIG env var to the kubeconfig for the current kw session",
	Long: `Generate a script that points the KUBECONFIG env var to the kubeconfig for the current kw session.

The script can be called manually or added to your shell's profile to automatically set the KUBECONFIG env var when you start a new shell.`,
}

func init() {
	KubectlEnvCmd.AddCommand(getKubectlEnvCmdForShell(SHELL_BASH))
	KubectlEnvCmd.AddCommand(getKubectlEnvCmdForShell(SHELL_ZSH))
	KubectlEnvCmd.AddCommand(getKubectlEnvCmdForShell(SHELL_FISH))
	KubectlEnvCmd.AddCommand(getKubectlEnvCmdForShell(SHELL_POWER))
}

func getKubectlEnvCmdForShell(shell string) *cobra.Command {
	text := ""
	switch shell {
	case SHELL_BASH:
		text = `export KUBECONFIG='%s';

# Run this command to configure kubectl for your shell:
# eval $(kw kubectl-env bash)
`
	case SHELL_ZSH:
		text = `export KUBECONFIG='%s';

# Run this command to configure kubectl for your shell:
# eval $(kw kubectl-env zsh)
`
	case SHELL_FISH:
		text = `set -gx KUBECONFIG '%s';

# Run this command to configure kubectl for your shell:
# eval (kw kubectl-env fish)
`
	case SHELL_POWER:
		text = `$Env:KUBECONFIG = '%s';
# Run this command to configure kubectl for your shell:
# & kw kubectl-env powershell | Invoke-Expression
`
	default:
		libutils.Fatal(1, "unsupported shell: %s", shell)
	}
	return &cobra.Command{
		Use:   shell,
		Args:  cobra.NoArgs,
		Short: fmt.Sprintf("Generate a %s script that points the KUBECONFIG env var to the kubeconfig for the current kw session", shell),
		Long: fmt.Sprintf(`Generate a %s script that points the KUBECONFIG env var to the kubeconfig for the current kw session.

The script can be called manually or added to your shell's profile to automatically set the KUBECONFIG env var when you start a new shell.`, shell),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf(text, config.Runtime.KubeconfigPath())
		},
	}
}
