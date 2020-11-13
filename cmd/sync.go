package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bomoko/lagoon-sync/synchers"
	"github.com/spf13/cobra"
)

var ProjectName string
var sourceEnvironmentName string
var targetEnvironmentName string
var configurationFile string

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func detectLocalEnvironment() bool {
	isLocal := false
	if domain, exists := os.LookupEnv("LAGOON_DOMAIN"); exists {
		// fmt.Println(domain)
		isLocal = strings.Contains(domain, "docker.amazee.io")
		return isLocal
	}
	return isLocal
}

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync [mariadb|files|mongodb|postgres|etc.]",
	Short: "Sync a resource type",
	Long:  `Use Lagoon-Sync to sync an external environments resources with the local environment`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		moduleName := args[0]

		//For now, let's just try write up a command that generates the strings ...
		//TODO: make the lagoonYamlPath a configuration value, and overridable
		//Perhaps we should refactor this into some generic thing ...
		lagoonConfigBytestream, err := LoadLagoonConfig("./.lagoon.yml")
		if err != nil {
			fmt.Println("Couldn't load lagoon config file")
			os.Exit(1)
		}

		//TODO: we need some standard way of extracting the project name
		// For now, let's just pull it straight from the .lagoon.yml

		configRoot, err := synchers.UnmarshallLagoonYamlToLagoonSyncStructure(lagoonConfigBytestream)
		if err != nil {
			log.Printf("There was an issue unmarshalling the sync configuration: %v", err)
			return
		}

		defaultProjectName := fmt.Sprintf("%s", strings.Replace(getEnv("LAGOON_PROJECT", configRoot.Project), "_", "-", -1))

		sourceEnvironment := synchers.RemoteEnvironment{
			ProjectName:            ProjectName,
			DefaultProjectName:     defaultProjectName,
			EnvironmentName:        sourceEnvironmentName,
			DefaultEnvironmentName: getEnv("LAGOON_GIT_BRANCH", ""),
		}

		var lagoonSyncer synchers.Syncer
		//TODO: perhaps there's a more dynamic way of doing this match?
		switch moduleName {
		case "mariadb":
			lagoonSyncer = configRoot.LagoonSync.Mariadb.PrepareSyncer()
			break
		case "postgres":
			lagoonSyncer = configRoot.LagoonSync.Postgres.PrepareSyncer()
			break
		case "drupalconfig":
			lagoonSyncer = configRoot.LagoonSync.Drupalconfig.PrepareSyncer()
			break
		default:
			log.Print("Could not match type : %v", moduleName)
			return
			break
		}

		err = synchers.RunSyncProcess(sourceEnvironment, lagoonSyncer)
		fmt.Println(lagoonSyncer)
		if err != nil {
			log.Printf("There was an error running the sync process: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")
	syncCmd.PersistentFlags().StringVar(&ProjectName, "project-name", "", "The Lagoon project name of the remote system")
	// syncCmd.MarkPersistentFlagRequired("project-name")
	syncCmd.PersistentFlags().StringVar(&sourceEnvironmentName, "source-environment-name", "", "The Lagoon environment name of the source system")

	isLocal := detectLocalEnvironment()
	if isLocal {
		fmt.Print("(Local) ")
		syncCmd.MarkPersistentFlagRequired("source-environment-name")
	}

	syncCmd.PersistentFlags().StringVar(&targetEnvironmentName, "target-environment-name", "", "The Lagoon environment name of the source system (defaults to local)")
	syncCmd.PersistentFlags().StringVar(&configurationFile, "configuration-file", "", "File containing sync configuration. Defaults to ./.lagoon.yml")
	syncCmd.MarkPersistentFlagRequired("remote-environment-name")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
