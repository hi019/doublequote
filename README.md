Doublequote: Slogan Pending
========================================

Doublequote is an RSS reader prioritizing excellent UI/UX.

## Tech

### Backend

* Written in Go
* Uses the standard library's net/http and [chi](https://github.com/go-chi/chi) for routing
* Structure based on @benbjohnson's ["Standard Package Layout"](https://www.gobeyond.dev/standard-package-layout/) and its accompanying example repo, [benbjohnson/wtf](https://github.com/benbjohnson/wtf)
* [google/wire](https://github.com/google/wire) for compile-time dependency injection

### Frontend

* Written in Typescript with React + create-react-app