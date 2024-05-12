/*
Copyright © 2023 Daniel Wu <wxc@wxccs.org>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/routes"
	"github.com/iceking2nd/rustdesk-api-server/docs"
	"github.com/iceking2nd/rustdesk-api-server/global"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginglog "github.com/szuecs/gin-glog"
	"github.com/toorop/gin-logrus"
)

var (
	cfgFile string
	logFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rustdesk-api-server",
	Short: "API service for RustDesk",
	/*Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,*/
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		//gin.SetMode(gin.ReleaseMode)
		apiEngine := gin.New()
		docs.SwaggerInfo.BasePath = "/api"
		corsConfig := cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
			AllowAllOrigins:  true,
		}
		apiEngine.Use(cors.New(corsConfig))
		apiServer := &http.Server{
			Addr:         fmt.Sprintf("%s:%d", viper.GetString("API.Host"), viper.GetInt("API.Port")),
			Handler:      apiEngine,
			ReadTimeout:  120 * time.Second,
			WriteTimeout: 120 * time.Second,
		}
		apiEngine.NoRoute(func(c *gin.Context) {
			global.Log.WithField("request", c.Request).Debugln("received 404 requests")
			c.JSON(http.StatusNotFound, Controllers.ResponseError{Error: "server_not_support"})
			return
		})
		root := apiEngine.Group("/")
		root.Use(ginglog.Logger(3 * time.Second))
		root.Use(Database.SetContext())
		root.Use(ginlogrus.Logger(global.Log), gin.Recovery())
		routes.SetupRouter(root)
		go func() {
			var err error
			if err = apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("HTTP server listen: %s\n", err.Error())
			}
		}()

		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
		ticker := time.NewTicker(time.Millisecond)
		for {
			select {
			case sig := <-signalChan:
				log.Println("Get Signal:", sig)
				log.Println("Shutdown Server ...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := apiServer.Shutdown(ctx); err != nil {
					log.Fatal("Closing web service error: ", err)
				}
				log.Println("Server exiting")
				os.Exit(0)
			case <-ticker.C:
				// do sth every tick
			}
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/rustdesk-api-server.yaml)")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "logging file")
	rootCmd.PersistentFlags().Uint32Var(&global.LogLevel, "log-level", 3, "log level (0 - 6, 3 = warn , 5 = debug)")

	rootCmd.SetVersionTemplate(fmt.Sprintf(`{{with .Name}}{{printf "%%s version information: " .}}{{end}}
   {{printf "Version:    %%s" .Version}}
   Build Time:		%s
   Git Revision:	%s
   Go version:		%s
   OS/Arch:			%s/%s
`, global.BuildTime, global.GitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".rustdesk-api-server" (without extension).
		viper.AddConfigPath(home + string(os.PathSeparator) + ".config")
		viper.AddConfigPath("/etc/rustdesk-api-server")
		viper.AddConfigPath("/usr/local/etc/rustdesk-api-server")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("rustdesk-api-server")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err.Error())
		viper.SetDefault("API.Host", "127.0.0.1")
		viper.SetDefault("API.Port", 21114)
		viper.SetDefault("API.PublicURL", "http://127.0.0.1:21114")
		viper.SetDefault("MySQL.Host", "localhost")
		viper.SetDefault("MySQL.Port", 3306)
		viper.SetDefault("MySQL.User", "rustdesk-api-server")
		viper.SetDefault("MySQL.Pass", "rustdesk-api-server")
		viper.SetDefault("MySQL.DB", "rustdesk-api-server")
		viper.SetDefault("SMTP.From", "mail@example.com")
		viper.SetDefault("SMTP.Name", "Rustdesk")
		viper.SetDefault("SMTP.Host", "127.0.0.1")
		viper.SetDefault("SMTP.Port", 587)
		viper.SetDefault("SMTP.Username", "mail@example.com")
		viper.SetDefault("SMTP.Password", "")

		if len(cfgFile) > 0 {
			_ = viper.WriteConfigAs(cfgFile)
		} else {
			_ = viper.WriteConfigAs("rustdesk-api-server.yaml")
		}
		fmt.Println("Configuration file created.")
		os.Exit(0)
	}

	global.Log = logrus.New()
	var logWriter io.Writer
	if logFile == "" {
		logWriter = os.Stdout
	} else {
		logFileHandle, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err.Error())
		}
		logWriter = io.MultiWriter(os.Stdout, logFileHandle)
	}
	global.Log.SetOutput(logWriter)
	global.Log.SetLevel(logrus.Level(global.LogLevel))
}
