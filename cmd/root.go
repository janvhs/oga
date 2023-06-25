package cmd

import "github.com/spf13/cobra"

func New(appName, version, commitSHA string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:                   appName,
		Version:               version,
		DisableFlagsInUseLine: true,
		SilenceErrors:         true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	if len(commitSHA) >= 7 {
		vt := rootCmd.VersionTemplate()
		rootCmd.SetVersionTemplate(vt[:len(vt)-1] + " (" + commitSHA[0:7] + ")\n")
	}

	return rootCmd
}
