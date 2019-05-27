package cmd

import (
	"context"
	"fmt"
	"github.com/orestrepov/metadatahost/model"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/orestrepov/metadatahost/api"
	"github.com/orestrepov/metadatahost/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the api",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := app.New()
		if err != nil {
			return err
		}
		defer a.Close()

		api, err := api.New(a)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
			logrus.Info("signal caught. shutting down...")
			cancel()
		}()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			serveAPI(ctx, api)
		}()

		// start new
		a.Database.AutoMigrate(&model.Host{}, &model.Host{})
		a.Database.AutoMigrate(&model.Server{}, &model.Server{})
		// end new

		wg.Wait()
		return nil
	},
}

func serveAPI(ctx context.Context, api *api.API) {

	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "HEAD", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,                             // Log API request calls
		middleware.DefaultCompress,                    // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
		middleware.Recoverer,                          // Recover from panics without crashing server
		cors.Handler,
	)

	api.Init(router.Route("/api", nil))

	s := &http.Server{
		Addr:        fmt.Sprintf(":%d", api.Config.Port),
		Handler:     router,
		ReadTimeout: 2 * time.Minute,
	}

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			logrus.Error(err)
		}
		close(done)
	}()

	logrus.Infof("serving api at http://127.0.0.1:%d", api.Config.Port)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Error(err)
	}
	<-done
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
