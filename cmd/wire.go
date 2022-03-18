//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject
// +build wireinject

package main

import (
	"doublequote/pkg/asynq"
	"doublequote/pkg/blob"
	"doublequote/pkg/crypto"
	"doublequote/pkg/domain"
	"doublequote/pkg/http"
	"doublequote/pkg/ingest"
	"doublequote/pkg/listener"
	"doublequote/pkg/reaper"
	"doublequote/pkg/redis"
	"doublequote/pkg/sql"
	"github.com/google/wire"
)

type application struct {
	userService       domain.UserService
	cryptoService     domain.CryptoService
	sessionService    domain.SessionService
	collectionService domain.CollectionService
	entryService      domain.EntryService
	storageService    domain.StorageService
	ingestService     domain.IngestService

	httpServer *http.Server
}

// Setup functions for services that require configuration.
// This file is used by wire (https://github.com/google/wire) for dependency injection.

func setupSQL(cfg *domain.Config) (*sql.SQL, func(), error) {
	d := sql.NewSQL(cfg.Database.URL)

	err := d.Open()
	if err != nil {
		return nil, nil, err
	}

	return d, func() {
		d.Close()
	}, nil
}

func setupCache(cfg *domain.Config) (*redis.CacheService, error) {
	d := redis.NewCache(cfg.Redis.URL)

	return d, nil
}

func setupEventService(cfg *domain.Config) (domain.EventService, func(), error) {
	s := asynq.NewEventService(cfg.Redis.URL)

	err := s.Open()
	if err != nil {
		return nil, nil, err
	}

	return s, func() {
		s.Close()
	}, nil
}

func setupCryptoService(cfg *domain.Config) domain.CryptoService {
	return crypto.NewService(cfg.App.Secret)
}

func setupReaperService(entryService *domain.EntryService) domain.ReaperService {
	return reaper.New(entryService)
}

func setupServer(
	cfg *domain.Config,
	userService domain.UserService,
	cryptoService domain.CryptoService,
	sessionService domain.SessionService,
	storageService domain.StorageService,
	ingestService domain.IngestService,
	collectionService domain.CollectionService,
	feedService domain.FeedService,
) (*http.Server, func(), error) {
	s := http.NewServer()

	s.CryptoService = cryptoService
	s.UserService = userService
	s.SessionService = sessionService
	s.CollectionService = collectionService
	s.FeedService = feedService
	s.StorageService = storageService
	s.IngestService = ingestService
	s.Config = *cfg

	err := s.Open()
	if err != nil {
		return nil, nil, err
	}

	return s, func() {
		s.Close()
	}, nil
}

func setupFeedService(
	s *sql.SQL,
) domain.FeedService {
	return sql.NewFeedService(s)
}

func setupEntryService(
	s *sql.SQL,
) domain.EntryService {
	return sql.NewEntryService(s)
}

func setupIngestService(
	feedService domain.FeedService,
	entryService domain.EntryService,
	storageService domain.StorageService,
) domain.IngestService {
	return ingest.NewService(feedService, entryService, storageService)
}

func setupStorageService(
	cfg *domain.Config,
) (domain.StorageService, func(), error) {
	a, b, c := blob.NewStorageService(cfg.App.DataFolder)
	return a, func() {
		b()
	}, c
}

func setupListenerService(
	eventService domain.EventService,
	emailService domain.EmailService,
	cryptoService domain.CryptoService,
	reaperService domain.ReaperService,
	cfg domain.Config,
) {
	s := listener.NewService(eventService, emailService, cryptoService, reaperService, cfg)
	s.Start()
}

func newApplication(
	userService domain.UserService,
	cryptoService domain.CryptoService,
	sessionService domain.SessionService,
	entryService domain.EntryService,
	storageService domain.StorageService,
	server *http.Server,
) *application {
	return &application{
		userService:    userService,
		cryptoService:  cryptoService,
		sessionService: sessionService,
		entryService:   entryService,
		storageService: storageService,
		httpServer:     server,
	}
}

func initializeApplication(cfg *domain.Config) (*application, func(), error) {
	wire.Build(
		newApplication,
		setupServer,

		setupSQL,
		setupEventService,
		setupCryptoService,
		setupFeedService,
		setupEntryService,
		setupIngestService,
		setupStorageService,
		setupReaperService,
		setupListenerService,

		wire.Bind(new(domain.CacheService), new(*redis.CacheService)),
		setupCache,

		wire.Bind(new(domain.SessionService), new(*redis.SessionService)),
		sql.NewUserService,

		wire.Bind(new(domain.UserService), new(*sql.UserService)),
		redis.NewSessionService,

		wire.Bind(new(domain.CollectionService), new(*sql.CollectionService)),
		sql.NewCollectionService,
	)

	return &application{}, nil, nil
}
