package cmd

import (
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"goblin/internal/config"
	"goblin/internal/paste"
	"goblin/internal/router"
)

var (
	verbose       bool
	enableMetrics bool
)

var rootCmd = &cobra.Command{
	Use:   "goblin",
	Short: "A description of your program",

	Run: func(cmd *cobra.Command, args []string) {
		// Initialize Logrus
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stderr)

		// Load configuration using the config package
		config, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		// Print config values if verbose flag is set
		if verbose {
			log.Infof("Config values:")
			log.Infof("Port: %s", config.Port)
			log.Infof("Private mode: %v", config.Private)
			log.Infof("Template directory: %s", config.TemplateDir)
			log.Infof("Log directory: %s", config.LogDir)
			log.Infof("Paste storage directory: %s", config.PasteDir)
			log.Infof("Expire: %s", strconv.Itoa(config.Expire))
		}

		// Set the port based on the configuration or the flag
		port := config.Port // Use the configuration value
		if cmd.Flags().Changed("port") {
			port, _ = cmd.Flags().GetString("port")
		}

		// Check if the private flag is set
		private := config.Private // Use the configuration value
		if cmd.Flags().Changed("private") {
			private, _ = cmd.Flags().GetBool("private")
		}

		// Log whether the private flag is set
		log.Infof("Private mode: %v", private)

		// Create the router and start the server using the router package
		router := router.NewRouter()
		http.Handle("/", router)

		// Add a new route for Prometheus metrics
		http.Handle("/metrics", promhttp.Handler()) // Use Prometheus handler

		log.Info("Server started")

		log.Infof("Listening on port %s\n", port)

		err = http.ListenAndServe(":"+port, nil) // Use nil handler for default routing
		if err != nil {
			log.Errorf("Error starting server: %v\n", err)
		}
	},
}

func initializeConfig() {
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := paste.InitializeTemplates(); err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}
}

func intializeFlags() {
	rootCmd.PersistentFlags().String("port", "", "Port to listen on")
	rootCmd.PersistentFlags().Bool("private", false, "Hide content from terminal")
	rootCmd.PersistentFlags().String("templateDirectory", "", "Template directory path")
	rootCmd.PersistentFlags().String("logDir", "", "Log file location")
	rootCmd.PersistentFlags().String("pasteDir", "", "Paste storage location")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Print config values and enable verbose logging")
	rootCmd.PersistentFlags().Int("expire", 24, "Paste expiration duration in hours")
	rootCmd.PersistentFlags().BoolVar(&enableMetrics, "metrics", false, "Enable Prometheus metrics")

	viper.BindPFlag("Port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("Private", rootCmd.PersistentFlags().Lookup("private"))
	viper.BindPFlag("TemplateDir", rootCmd.PersistentFlags().Lookup("templateDirectory"))
	viper.BindPFlag("LogDir", rootCmd.PersistentFlags().Lookup("logDir"))
	viper.BindPFlag("PasteDir", rootCmd.PersistentFlags().Lookup("pasteDir"))
	viper.BindPFlag("Expire", rootCmd.PersistentFlags().Lookup("expire"))
}

func init() {
	cobra.OnInitialize(initializeConfig)
	intializeFlags()
}
