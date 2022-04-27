.DEFAULT_GOAL=serve

export DATABASE_URL = file:./data.db

test-all: generate test-go
build-all: generate build-js build-go

generate:
	go generate ./...
test-go:
	go test ./... -parallel 8
serve:
	go run ./cmd serve
build-js:
	cd frontend; pnpm install; pnpm build
build-go:
	cp -r ./frontend/build/* ./assets/frontend
	go build -o doublequote doublequote/cmd