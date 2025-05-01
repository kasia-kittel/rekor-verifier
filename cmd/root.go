package cmd

import (
	"os"

	"github.com/kasia-kittel/rekor-verifier/pkg/log"
	"github.com/kasia-kittel/rekor-verifier/pkg/verifier"
	"github.com/spf13/cobra"
)

// TODO:
// -r custom Rekor instance (could be set up by Viper)

var Version string = "alpha" // set during build
var path string
var sha string

var rootCmd = &cobra.Command{
	Use:     "rekor-verifier [filename]",
	Short:   "rekor-verifier automates certificates verification for binary signatures stored in Rekor",
	Args:    cobra.ExactArgs(0),
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {

		res := false

		if path != "" {
			res = verifier.VerifyFile(path)
		}

		if sha != "" {
			res = verifier.VerifySha(sha)
		}

		if !res {
			log.StdLogger.Println("Verification unsuccessful")
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "path to the binary")
	rootCmd.PersistentFlags().StringVarP(&sha, "sha", "s", "", "shasum of the binary")
	rootCmd.MarkFlagsOneRequired("path", "sha")
	rootCmd.MarkFlagsMutuallyExclusive("path", "sha")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.StdLogger.Fatalf("Oops. An error while executing rekor-verifier '%s'\n", err)
		os.Exit(1)
	}
}
