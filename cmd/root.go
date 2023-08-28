package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/strategodev/vuyo/router"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "my-calc",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: rootCmdRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// fmt.Println("flag")
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.my-calc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func rootCmdRun(cmd *cobra.Command, _ []string) {
	// Create a new HTTP server instance to handle inbound requests from the Panel
	// and external clients.
	s := &http.Server{
		Addr: "127.0.0.1:8080",
		// Addr:      api.Host + ":" + strconv.Itoa(api.Port),
		Handler: router.Configure(),
		// TLSConfig: config.DefaultTLSConfig,
	}

	// Check if main http server should run with TLS. Otherwise, reset the TLS
	// config on the server and then serve it over normal HTTP.
	// if api.Ssl.Enabled {
	// 	if err := s.ListenAndServeTLS(api.Ssl.CertificateFile, api.Ssl.KeyFile); err != nil {
	// 		log.WithFields(log.Fields{"auto_tls": false, "error": err}).Fatal("failed to configure HTTPS server")
	// 	}
	// 	return
	// }
	s.TLSConfig = nil
	if err := s.ListenAndServe(); err != nil {
		log.WithField("error", err).Fatal("failed to configure HTTP server")
	}
}
