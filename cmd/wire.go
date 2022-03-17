//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject
// +build wireinject

package main

import (
	"doublequote/pkg/asynq"
	"doublequote/pkg/blob"
	dq "doublequote/pkg/config"
	"doublequote/pkg/crypto"
	"doublequote/pkg/domain"
	"doublequote/pkg/http"
	"doublequote/pkg/ingest"
	redis2 "doublequote/pkg/redis"
	sql2 "doublequote/pkg/sql"
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

func setupSQL(cfg *dq.Config) (*sql2.SQL, func(), error) {
	d := sql2.NewSQL(cfg.Database.URL)

	err := d.Open()
	if err != nil {
		return nil, nil, err
	}

	return d, func() {
		d.Close()
	}, nil
}

func setupCache(cfg *dq.Config) (*redis2.CacheService, error) {
	d := redis2.NewCache(cfg.Redis.URL)

	return d, nil
}

func setupEventService(cfg *dq.Config) (domain.EventService, func(), error) {
	s := asynq.NewEventService(cfg.Redis.URL)

	err := s.Open()
	if err != nil {
		return nil, nil, err
	}

	return s, func() {
		s.Close()
	}, nil
}

func setupCryptoService(cfg *dq.Config) domain.CryptoService {
	s := crypto.NewService(cfg.App.Secret)
	return s
}

func setupServer(
	cfg *dq.Config,
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
	s *sql2.SQL,
) domain.FeedService {
	return sql2.NewFeedService(s)
}

func setupEntryService(
	s *sql2.SQL,
) domain.EntryService {
	return sql2.NewEntryService(s)
}

func setupIngestService(
	feedService domain.FeedService,
	entryService domain.EntryService,
	storageService domain.StorageService,
) domain.IngestService {
	return ingest.NewService(feedService, entryService, storageService)
}

func setupStorageService(
	cfg *dq.Config,
) (domain.StorageService, func(), error) {
	a, b, c := blob.NewStorageService(cfg.App.BucketName)
	return a, func() {
		b()
	}, c
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

func initializeApplication(cfg *dq.Config) (*application, func(), error) {
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

		wire.Bind(new(domain.CacheService), new(*redis2.CacheService)),
		setupCache,

		wire.Bind(new(domain.SessionService), new(*redis2.SessionService)),
		sql2.NewUserService,

		wire.Bind(new(domain.UserService), new(*sql2.UserService)),
		redis2.NewSessionService,

		wire.Bind(new(domain.CollectionService), new(*sql2.CollectionService)),
		sql2.NewCollectionService,
	)

	return &application{}, nil, nil
}
