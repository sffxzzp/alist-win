package cmd

import (
	"fmt"
	"net/http"
	"os"

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

var rootCmd = &cobra.Command{
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

		admin, _ := db.GetAdmin()
		os.WriteFile("password.txt", []byte("username: "+admin.Username+"\npassword: "+admin.Password), 0777)
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// let it directly run from explorer.exe
	cobra.MousetrapHelpText = ""
	rootCmd.PersistentFlags().StringVar(&flags.Config, "conf", "data/config.json", "config file")
	rootCmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "start with debug mode")
	rootCmd.PersistentFlags().BoolVar(&flags.NoPrefix, "no-prefix", false, "disable env prefix")
	rootCmd.PersistentFlags().BoolVar(&flags.Dev, "dev", false, "start with dev mode")
}