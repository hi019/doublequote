package main

import (
	"doublequote/pkg/asynq"
	"doublequote/pkg/blob"
	"doublequote/pkg/crypto"
	"doublequote/pkg/domain"
	"doublequote/pkg/http"
	"doublequote/pkg/ingest"
	"doublequote/pkg/job"
	"doublequote/pkg/listener"
	"doublequote/pkg/redis"
	"doublequote/pkg/smtp"
	"doublequote/pkg/sql"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the application",
	RunE:  RunServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func RunServe(_ *cobra.Command, _ []string) error {
	app := initApp()

	app.Run()

	return nil
}

func initApp() *fx.App {
	app := fx.New(
		fx.Provide(
			func() domain.Config { return cfg },
			sql.NewSQL,
			job.NewIngestJob,

			fx.Annotate(asynq.NewEventService, fx.As(new(domain.EventService))),
			fx.Annotate(ingest.NewService, fx.As(new(domain.IngestService))),
			fx.Annotate(smtp.NewEmailService, fx.As(new(domain.EmailService))),
			fx.Annotate(blob.NewStorageService, fx.As(new(domain.StorageService))),
			fx.Annotate(crypto.NewService, fx.As(new(domain.CryptoService))),
			fx.Annotate(redis.NewCache, fx.As(new(domain.CacheService))),
			fx.Annotate(redis.NewSessionService, fx.As(new(domain.SessionService))),
			fx.Annotate(sql.NewCollectionService, fx.As(new(domain.CollectionService))),
			fx.Annotate(sql.NewUserService, fx.As(new(domain.UserService))),
			fx.Annotate(sql.NewEntryService, fx.As(new(domain.EntryService))),
			fx.Annotate(sql.NewFeedService, fx.As(new(domain.FeedService))),
		),

		fx.Invoke(listener.NewService),
		fx.Invoke(http.NewServer),
	)

	return app
}
