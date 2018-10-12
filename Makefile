PACKAGE_NAME := github.com/yusufsyaifudin/tax-calculator-example
PROJECT_DIR := $(PWD)
CURRENT_TIME := `date +%s`

test:
	go test -v -cover -race ./...

# make create-migration NAME="create_users_table"
create-migration:
	@[ ! -z ${NAME} ] && echo assets/migrations/$(CURRENT_TIME)_${NAME}.sql
	@touch assets/migrations/$(CURRENT_TIME)_${NAME}.sql

install-dep:
	go get -v -u github.com/golang/dep/cmd/dep
	go get -v -u github.com/swaggo/swag/cmd/swag
	dep ensure -v

generate-doc:
	swag init -g internal/app/restapi/v1_server.go

build:
	rm -f out/tax-calculator-server
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o out/tax-calculator-server $(PACKAGE_NAME)/cmd/tax-calculator-server

install: install-dep generate-doc build

