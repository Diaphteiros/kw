package history

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Diaphteiros/kw/pkg/config"
	"github.com/Diaphteiros/kw/pkg/storage"
	"sigs.k8s.io/yaml"

	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var (
	Global bool
	output libutils.OutputFormat
)

var HistoryViewCmd = &cobra.Command{
	Use:     "view",
	Aliases: []string{"list", "ls", "v"},
	Args:    cobra.NoArgs,
	Short:   "View the history",
	Long:    `View the history.`,
	Run: func(cmd *cobra.Command, args []string) {
		if config.Runtime.Config().Kubeswitcher.HistoryDepth == 0 {
			cmd.Println("History is disabled.")
		}
		hist, err := storage.RenderHistory(Global)
		if err != nil {
			libutils.Fatal(1, "error creating internal history representation: %w\n", err)
		}

		var data []byte
		switch output {
		case libutils.OUTPUT_TEXT:
			data = []byte(hist.String())
		case libutils.OUTPUT_JSON:
			data, err = json.MarshalIndent(hist, "", "  ")
			if err != nil {
				libutils.Fatal(1, "error converting history to json: %w\n", err)
			}
		case libutils.OUTPUT_YAML:
			data, err = yaml.Marshal(hist)
			if err != nil {
				libutils.Fatal(1, "error converting history to yaml: %w\n", err)
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
	libutils.AddOutputFlag(HistoryViewCmd.Flags(), &output, libutils.OUTPUT_TEXT)
}
