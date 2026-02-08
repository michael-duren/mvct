package main

import (
	"fmt"
	"os"

	"github.com/michael-duren/mvct/internal/scaffold"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "mvct",
		Short: "MVCT is a Model-View-Controller framework for Bubble Tea",
		Long:  `A framework for building terminal user interfaces with Bubble Tea using an MVC architecture.`,
	}

	var withDB bool

	var scaffoldCmd = &cobra.Command{
		Use:   "scaffold [project-name]",
		Short: "Scaffold a new MVCT project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			config := scaffold.ProjectConfiguration{
				Name:   projectName,
				Path:   projectName,
				WithDB: withDB,
			}
			scaffold.ScaffoldProject(config)
		},
	}

	scaffoldCmd.Flags().BoolVar(&withDB, "db", false, "Include SQLite database with sqlc setup")

	rootCmd.AddCommand(scaffoldCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
