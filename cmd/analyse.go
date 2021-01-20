package cmd

import (
	"os"
	"path/filepath"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	sa "github.com/fabric8-analytics/cli-tools/analyses/stackanalyses"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var manifestFile string

// analyseCmd represents the analyse command
var analyseCmd = &cobra.Command{
	Use:     "analyse",
	Short:   "Performs full Stack Analyses on CRDA Platform.",
	Long:    `Performs full Stack Analyses on CRDA Platform. Supported ecosystems are Pypi (Python), Maven (Java), Npm (Node) and Golang (Go).`,
	Run:     runAnalyse,
	PostRun: destructor,
}

func init() {
	rootCmd.AddCommand(analyseCmd)
	analyseCmd.PersistentFlags().StringVarP(&manifestFile, "file", "f", "", "Manifest file absolute path.")
	analyseCmd.MarkPersistentFlagRequired("file")
}

// destructor deletes intermediary files used to have stack analyses
func destructor(cmd *cobra.Command, args []string) {
	log.Debug().Msgf("Running Destructor.\n")
	if debug {
		// Keep intermediary files, when on debug
		log.Debug().Msgf("Skipping file clearance on Debug Mode.\n")
		return
	}
	intermediaryFiles := []string{"generate_pylist.py", "pylist.json", "dependencies.txt", "golist.json", "npmlist.json"}
	for _, file := range intermediaryFiles {
		file = filepath.Join(os.TempDir(), file)
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				// If file doesn't exists, continue
				continue
			}
		}
		e := os.Remove(file)
		if e != nil {
			log.Fatal().Msgf("Error clearing files %s", file)
		}
	}
}

//runAnalyse is controller func for analyses cmd.
func runAnalyse(cmd *cobra.Command, args []string) {
	requestParams := driver.RequestType{
		UserID:          viper.GetString("crda-key"),
		ThreeScaleToken: viper.GetString("auth-token"),
		Host:            viper.GetString("host"),
		RawManifestFile: manifestFile,
	}
	analysesResult := sa.StackAnalyses(requestParams)
	if sa.ProcessResult(analysesResult) {
		// If Stack has vulnerability, exit with 2 code
		os.Exit(2)
	}
}