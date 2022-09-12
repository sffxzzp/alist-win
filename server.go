package cmd

import (
	"fmt"
	"net/http"

	"github.com/alist-org/alist/v3/cmd/flags"
	_ "github.com/alist-org/alist/v3/drivers"
	"github.com/alist-org/alist/v3/internal/bootstrap"
	"github.com/alist-org/alist/v3/internal/conf"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/alist-org/alist/v3/server"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sffxzzp/go-webview2"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server at the specified address",
	Long: `Start the server at the specified address
the address is defined in config file`,
	Run: func(cmd *cobra.Command, args []string) {
		Init()
		bootstrap.InitAria2()
		bootstrap.LoadStorages()
		if !flags.Debug && !flags.Dev {
			gin.SetMode(gin.ReleaseMode)
		}
		r := gin.New()
		r.Use(gin.LoggerWithWriter(log.StandardLogger().Out), gin.RecoveryWithWriter(log.StandardLogger().Out))
		server.Init(r)
		base := fmt.Sprintf("%s:%d", conf.Conf.Address, conf.Conf.Port)
		utils.Log.Infof("start server @ %s", base)
		srv := &http.Server{Addr: base, Handler: r}
		go func() {
			var err error
			if conf.Conf.Scheme.Https {
				//err = r.RunTLS(base, conf.Conf.Scheme.CertFile, conf.Conf.Scheme.KeyFile)
				err = srv.ListenAndServeTLS(conf.Conf.Scheme.CertFile, conf.Conf.Scheme.KeyFile)
			} else {
				err = srv.ListenAndServe()
			}
			if err != nil && err != http.ErrServerClosed {
				utils.Log.Fatalf("failed to start: %s", err.Error())
			}
		}()
		var prefix string
		if conf.Conf.Scheme.Https {
			prefix = "https"
		} else {
			prefix = "http"
		}
		w := webview2.NewWithOptions(webview2.WebViewOptions{
			Debug:    false,
			DataPath: "./data",
			WindowOptions: webview2.WindowOptions{
				Title:  "alist",
				Width:  1440,
				Height: 900,
				IconId: 2,
				Center: true,
			},
		})
		defer w.Destroy()
		w.Navigate(fmt.Sprintf("%s://127.0.0.1:%d", prefix, conf.Conf.Port))
		w.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
