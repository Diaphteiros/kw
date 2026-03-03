package meta

import (
	"encoding/json"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/pkg/cmdgroups"
	"github.com/Diaphteiros/kw/pkg/config"

	"sigs.k8s.io/yaml"

	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var output libutils.OutputFormat

var InfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Args:    cobra.NoArgs,
	GroupID: cmdgroups.Meta,
	Short:   "Shows information about the current configuration",
	Long: `Shows information about the current configuration.

Supports json and yaml output, with the latter one being the default.`,
	Run: func(cmd *cobra.Command, args []string) {
		libutils.ValidateOutputFormat(output, libutils.OUTPUT_YAML, libutils.OUTPUT_JSON)
		state := config.Runtime.State()
		rawData := map[string]any{}
		// id
		data, err := vfs.ReadFile(fs.FS, config.Runtime.IdPath())
		if err != nil {
			if !vfs.IsNotExist(err) {
				libutils.Fatal(1, "error reading id backup file '%s': %w\n", config.Runtime.IdPath(), err)
			}
			data = []byte("<unknown>")
		}
		rawData["id"] = string(data)
		// generic state
		rawData["genericState"] = state.GenericState
		// plugin state
		if len(state.RawPluginState) > 0 {
			// unmarshal plugin state into map[string]any to properly convert it to json/yaml
			ps := map[string]any{}
			if err := json.Unmarshal(state.RawPluginState, &ps); err != nil {
				libutils.Fatal(1, "error unmarshalling plugin state: %w\n", err)
			}
			rawData["pluginState"] = ps
		}
		switch output {
		case libutils.OUTPUT_YAML:
			data, err = yaml.Marshal(rawData)
			if err != nil {
				libutils.Fatal(1, "error marshalling state to yaml: %w\n", err)
			}
		case libutils.OUTPUT_JSON:
			data, err = json.Marshal(rawData)
			if err != nil {
				libutils.Fatal(1, "error marshalling state to json: %w\n", err)
			}
		default:
			libutils.Fatal(1, "unsupported output format: %s\n", output)
		}
		cmd.Println(string(data))
	},
}

func init() {
	libutils.AddOutputFlag(InfoCmd.Flags(), &output, libutils.OUTPUT_YAML, libutils.OUTPUT_YAML, libutils.OUTPUT_JSON)
}
