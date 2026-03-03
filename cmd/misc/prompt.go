package misc

import (
	"fmt"

	"github.com/spf13/cobra"

	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var PromptCmd = &cobra.Command{
	Use:   "prompt",
	Args:  cobra.NoArgs,
	Short: "Generate a script that generates a prompt to display in the shell",
	Long: `Generate a script that generates a prompt to display in the shell.

The script can be called manually or added to your shell's profile.`,
}

func init() {
	PromptCmd.AddCommand(getPromptCmdForShell(SHELL_BASH))
	PromptCmd.AddCommand(getPromptCmdForShell(SHELL_ZSH))
}

func getPromptCmdForShell(shell string) *cobra.Command {
	text := ""
	switch shell {
	case SHELL_BASH:
		text = `_kw_prompt() { if [[ -f "$KUBECONFIG" ]]; then; id_path="$(dirname "$KUBECONFIG")/id"; if [[ -f "$id_path" ]]; then; ns=$(kubectl config view --minify -o jsonpath='{.contexts[0].context.namespace}'); if [[ ${ns:-"default"} == "default" ]]; then; ns=""; else; ns=" <$ns>"; fi; echo -n "$(cat $id_path)$ns"; fi; fi; }

# Run this command to enable the _kw_prompt function for your shell:
# eval $(kw prompt bash)
`
	case SHELL_ZSH:
		text = `_kw_prompt() { if [[ -f "$KUBECONFIG" ]]; then; id_path="$(dirname "$KUBECONFIG")/id"; if [[ -f "$id_path" ]]; then; ns=$(kubectl config view --minify -o jsonpath='{.contexts[0].context.namespace}'); if [[ ${ns:-"default"} == "default" ]]; then; ns=""; else; ns=" <$ns>"; fi; echo -n "$(cat $id_path)$ns"; fi; fi; }

# Run this command to enable the _kw_prompt function for your shell:
# eval $(kw prompt zsh)
`
	default:
		libutils.Fatal(1, "unsupported shell: %s", shell)
	}
	return &cobra.Command{
		Use:   shell,
		Args:  cobra.NoArgs,
		Short: fmt.Sprintf("Generate a %s script that generates a prompt to display in the shell", shell),
		Long: fmt.Sprintf(`Generate a %s script that generates a prompt to display in the shell.

The script can be called manually or added to your shell's profile.`, shell),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Print(text)
		},
	}
}
