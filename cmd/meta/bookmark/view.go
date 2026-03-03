package bookmark

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Diaphteiros/kw/pkg/storage"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
	"sigs.k8s.io/yaml"
)

var output libutils.OutputFormat

var ViewCmd = &cobra.Command{
	Use:     "view",
	Aliases: []string{"list", "ls", "v"},
	Args:    cobra.NoArgs,
	Short:   "View the bookmarks",
	Long:    `View the bookmarks.`,
	Run: func(cmd *cobra.Command, args []string) {
		si, err := storage.RenderStorageInventory()
		if err != nil {
			libutils.Fatal(1, "error creating internal bookmark representation: %w\n", err)
		}

		var data []byte
		switch output {
		case libutils.OUTPUT_TEXT:
			data = []byte(si.String())
		case libutils.OUTPUT_JSON:
			data, err = json.MarshalIndent(si, "", "  ")
			if err != nil {
				libutils.Fatal(1, "error converting bookmark inventory to json: %w\n", err)
			}
		case libutils.OUTPUT_YAML:
			data, err = yaml.Marshal(si)
			if err != nil {
				libutils.Fatal(1, "error converting bookmark inventory to yaml: %w\n", err)
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
	libutils.AddOutputFlag(ViewCmd.Flags(), &output, libutils.OUTPUT_TEXT)
}
