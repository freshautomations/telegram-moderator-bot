PACKAGES := $(shell go list)
BUILD_NUMBER ?= 0
BUILD_FLAGS = -ldflags "-s -extldflags -static -X github.com/freshautomations/telegram-moderator-bot/defaults.Release=$(BUILD_NUMBER)"
GOPATH ?= $(shell go env GOPATH)

# Set these before running targets. Example: ENVIRONMENT=prod make deploy
ENVIRONMENT ?= staging
LAMBDA_SECRET ?= staging
LAMBDA_TIMEOUT ?= 15
TELEGRAM_TOKEN ?= 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11

########################################
### Build

build:
	go build $(BUILD_FLAGS) -o build/tmb .

build-linux:
	GOOS=linux GOARCH=amd64 $(MAKE) build
#	docker run -it --rm -v $(GOPATH):/go -e BUILD_NUMBER=$(BUILD_NUMBER) golang:1.11.3 make -C /go/src/github.com/freshautomations/telegram-moderator-bot build

########################################
### Tools & dependencies

$(GOPATH)/bin/dep:
	@go get -u -v github.com/golang/dep/cmd/dep

get_vendor_deps: $(GOPATH)/bin/dep
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@$(GOPATH)/bin/dep ensure -v

########################################
### Testing

test:
	@go test -count 1 -p 1 $(PACKAGES)

########################################
### Localnet

localnet-start:
	build/tmb -webserver

localnet-lambda:
	# (Requirements: pip3 install aws-sam-cli)
	# Set up env.vars in template.yml since the --env-vars option doesn't seem to work
	sam local start-api

########################################
### Release management (set up requirements manually)

package:
	zip "build/tmb.zip" build/tmb

#sam-deploy:
#	sam deploy --template-file resources/template.yml --stack-name "tmb-staging" --capabilities CAPABILITY_IAM --region "us-east-1"

deploy:
	cd resources/terraform && terraform init && terraform apply -auto-approve -var ENVIRONMENT=$(ENVIRONMENT) -var LAMBDA_SECRET=$(LAMBDA_SECRET) -var LAMBDA_TIMEOUT=$(LAMBDA_TIMEOUT) -var TELEGRAM_TOKEN=$(TELEGRAM_TOKEN)

destroy:
	cd resources/terraform && terraform destroy -auto-approve -var ENVIRONMENT=$(ENVIRONMENT) -var LAMBDA_SECRET=$(LAMBDA_SECRET) -var LAMBDA_TIMEOUT=$(LAMBDA_TIMEOUT) -var TELEGRAM_TOKEN=$(TELEGRAM_TOKEN)

webhook:
	@curl https://api.telegram.org/bot$(TELEGRAM_TOKEN)/deleteWebhook
	@cd resources/terraform && export BASE_URL=`terraform output base_url` && curl -F "url=$${BASE_URL}" -F "allowed_updates=%5B\"message\"%5D" -F "max_connections=10" https://api.telegram.org/bot$(TELEGRAM_TOKEN)/setWebhook

webhook-info:
	@cd resources/terraform && export BASE_URL=`terraform output base_url` && curl https://api.telegram.org/bot$(TELEGRAM_TOKEN)/getWebhookinfo

getme:
	@cd resources/terraform && export BASE_URL=`terraform output base_url` && curl https://api.telegram.org/bot$(TELEGRAM_TOKEN)/getMe

list-lambda:
	aws lambda list-functions --region us-east-1

.PHONY: build build-linux get_vendor_deps test localnet-start localnet-lambda package deploy destroy webhook webhook-info getme list-lambda
