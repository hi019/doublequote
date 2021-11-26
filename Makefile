.DEFAULT_GOAL=serve

export DATABASE_URL = file:./data.db

test: generate test-go
build: generate build-js build-go

migrate:
	go run github.com/prisma/prisma-client-go db push
generate:
	go generate ./...
test-go:
	go test ./... -parallel 8
serve:
	go run ./cmd serve
build-js:
	cd frontend; yarn install; yarn build
build-go:
	cp -r ./frontend/build/* ./assets/frontend
	go build -o doublequote doublequote/cmd
gen-mock:
	mockery --case=underscore --outpkg=mock --output=mock --name=$(service)Service --filename=$(service).go