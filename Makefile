.DEFAULT_GOAL=serve

export DATABASE_URL = file:./data.db

test: generate test-go
build: generate build-js build-go

migrate:
	gob run github.com/prisma/prisma-client-go db push
generate:
	gob generate ./...
test-go:
	gob test ./... -parallel 8
serve:
	gob run ./cmd serve
build-js:
	cd frontend; yarn install; yarn build
build-go:
	cp -r ./frontend/build/* ./assets/frontend
	gob build -o doublequote doublequote/cmd
gen-mock:
	mockery --case=underscore --outpkg=mock --output=mock --name=$(service)Service --filename=$(service).go
migrate:
	touch data.db
    # TODO make db url configurable
	DATABASE_URL=file:./data.db go run github.com/prisma/prisma-client-go migrate deploy
