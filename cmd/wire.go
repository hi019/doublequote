//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject
// +build wireinject

package main

import (
	dq "doublequote"
	"doublequote/asynq"
	"doublequote/crypto"
	"doublequote/http"
	"doublequote/redis"
	"doublequote/sql"
	"github.com/google/wire"
)

// Setup functions for services that require configuration.
// This file is used by wire (https://github.com/google/wire) for dependency injection.

func setupSQL(cfg *dq.Config) (*sql.SQL, func(), error) {
	d := sql.NewSQL(cfg.Database.URL)

	err := d.Open()
	if err != nil {
		return nil, nil, err
	}

	return d, func() {
		d.Close()
	}, nil
}

func setupCache(cfg *dq.Config) (*redis.CacheService, error) {
	d := redis.NewCache(cfg.Redis.URL)

	return d, nil
}

func setupEventService(cfg *dq.Config) (dq.EventService, func(), error) {
	s := asynq.NewEventService(cfg.Redis.URL)

	err := s.Open()
	if err != nil {
		return nil, nil, err
	}

	return s, func() {
		s.Close()
	}, nil
}

func setupCryptoService(cfg *dq.Config) dq.CryptoService {
	s := crypto.NewService(cfg.App.Secret)
	return s
}

func setupServer(
	cfg *dq.Config,
	userService dq.UserService,
	cryptoService dq.CryptoService,
	sessionService dq.SessionService,
	collectionService dq.CollectionService,
	feedService dq.FeedService,
) (*http.Server, func(), error) {
	s := http.NewServer()

	s.CryptoService = cryptoService
	s.UserService = userService
	s.SessionService = sessionService
	s.CollectionService = collectionService
	s.FeedService = feedService
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
) dq.FeedService {
	return sql.NewFeedService(s)
}

func newApplication(
	userService dq.UserService,
	cryptoService dq.CryptoService,
	sessionService dq.SessionService,
	server *http.Server,
) *application {
	return &application{
		userService:    &userService,
		cryptoService:  &cryptoService,
		sessionService: &sessionService,
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

		wire.Bind(new(dq.CacheService), new(*redis.CacheService)),
		setupCache,

		wire.Bind(new(dq.SessionService), new(*redis.SessionService)),
		sql.NewUserService,

		wire.Bind(new(dq.UserService), new(*sql.UserService)),
		redis.NewSessionService,

		wire.Bind(new(dq.CollectionService), new(*sql.CollectionService)),
		sql.NewCollectionService,
	)

	return &application{}, nil, nil
}
