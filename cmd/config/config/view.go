package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"sigs.k8s.io/yaml"

	"github.com/Diaphteiros/kw/pkg/config"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var (
	core   bool
	plugin string
	raw    bool
	output libutils.OutputFormat
)

var ViewCmd = &cobra.Command{
	Use:     "view",
	Aliases: []string{"v"},
	Args:    cobra.NoArgs,
	Short:   "View the kubeswitcher configuration file",
	Long: `View the kubeswitcher configuration file.

Use the '--core' flag to exclude plugin configuration from the output.
To only view a specific plugin's configuration, use the '--plugin <plugin_name>' flag.

Usually, the configuration file is parsed and before being printed. This means that comments don't appear in the output and some fields might have changed due to defaulting.
To view the raw configuration file, use the '--raw' flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		validateArgs()

		cpath := filepath.Join(config.Runtime.ConfigDirectory(), config.ConfigFileName)
		var data []byte
		var err error
		if raw {
			rawData, err := os.ReadFile(cpath)
			if err != nil {
				libutils.Fatal(1, "error reading configuration file: %w\n", err)
			}
			cmd.Print(string(rawData))
			return
		} else {
			cfg := config.Runtime.Config()
			var toPrint any
			if core {
				cfg.Plugins = nil
				toPrint = cfg
			} else if plugin != "" {
				var pCfg *config.PluginConfig
				for _, p := range cfg.Plugins {
					if p.Name == plugin {
						pCfg = p
						break
					}
				}
				if pCfg != nil {
					toPrint = pCfg
				} else {
					libutils.Fatal(1, "plugin '%s' not found in configuration\n", plugin)
				}
			} else {
				toPrint = cfg
			}

			switch output {
			case libutils.OUTPUT_JSON:
				data, err = json.MarshalIndent(toPrint, "", "  ")
				if err != nil {
					libutils.Fatal(1, "error converting configuration to json: %w\n", err)
				}
			case libutils.OUTPUT_YAML:
				data, err = yaml.Marshal(toPrint)
				if err != nil {
					libutils.Fatal(1, "error converting configuration to yaml: %w\n", err)
				}
			}
		}

		sData := string(data)
		if strings.HasSuffix(sData, "\n") {
			cmd.Print(sData)
		} else {
			cmd.Println(sData)
		}
	},
}

func init() {
	ViewCmd.Flags().BoolVarP(&core, "core", "c", false, "Print only the core configuration, excluding plugin configurations.")
	ViewCmd.Flags().StringVarP(&plugin, "plugin", "p", "", "Print only the configuration for the specified plugin.")
	ViewCmd.Flags().BoolVarP(&raw, "raw", "r", false, "Print the raw configuration file without parsing it. Cannot be combined with any arguments that filter the output or change its format.")
	libutils.AddOutputFlag(ViewCmd.Flags(), &output, libutils.OUTPUT_YAML, libutils.OUTPUT_JSON, libutils.OUTPUT_YAML)
}

func validateArgs() {
	if core && plugin != "" {
		libutils.Fatal(1, "'--core' and '--plugin' are mutually exclusive\n")
	}

	if raw && (core || plugin != "" || output != libutils.OUTPUT_YAML) {
		libutils.Fatal(1, "'--raw' cannot be combined with any arguments that filter the output or change its format\n")
	}
}
