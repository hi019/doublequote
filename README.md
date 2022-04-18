# Doublequote: Slogan Pending

Doublequote is an RSS reader prioritizing excellent UI/UX.

## Status

Doublequote is not yet at MVP status, but it's under active development. Here's what's implemented so far:

- File & environment based configuration
- User login, signup, email verification and authorization
- Feed CRUD (backend only)
- Collection (groups of Feeds) CRUD (backend only)
- RSS Entry ingestion (half done, backend only)

## Building

Make sure the following dependencies are installed:

- `make`
- `pnpm`
- `go`
- `redis` (on the default port, 6379)

Then, after cloning the repository, run:

1. `make build`
2. `make migrate`
3. `./doublequote serve`

And Doublequote should be running on port 8080.

## Tech

### Backend

- Written in Go
- Uses the standard library's net/http and [chi](https://github.com/go-chi/chi) for routing
- Structure based on @benbjohnson's ["Standard Package Layout"](https://www.gobeyond.dev/standard-package-layout/) and its accompanying example repo, [benbjohnson/wtf](https://github.com/benbjohnson/wtf)
- [google/wire](https://github.com/google/wire) for compile-time dependency injection

### Frontend

- Written in Typescript with React + create-react-app
