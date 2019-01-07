# Telegram-Moderator-Bot

## Overview
Telegram Moderator Bot written in Go as an AWS Lambda function (and a stand-alone web service).

Please refer to the [User's Guide](freshautomations/telegram-moderator-bot/blob/master/GUIDE.md)
for additional information on how to use the bot.

## Setup
### Prerequisites
1. Set up [Go](https://golang.org).
1. Set up [Terraform](https://terraform.io).
1. Familiarize yourself with [@BotFather](https://core.telegram.org/bots#3-how-do-i-create-a-bot) on Telegram.
1. Set up AWS credentials for Terraform and GOPATH variable for Go.
1. Run the below to download and prepare the code:
```bash
cd $GOPATH
mkdir -p src/github.com/freshautomations
cd src/github.com/freshautomations
git clone https://github.com/freshautomations/telegram-moderator-bot
cd telegram-moderator-bot

make get_vendor_deps
```
### Build
```bash
BUILD_NUMBER=10 make build
```
This creates the `tmb` binary in the `build` folder. It can be run locally, with the `-webserver` option or deployed as a Lambda function in AWS.
### Run local webserver
Copy the `tmb.conf.template` file over to `tmb.conf`. Edit the file and fill in the details. Then run:
```bash
make localnet-start
```
### Deploy infrastructure
The `resources/terraform` folder contains Terraform scripts to deploy the compiled binary to AWS and set it up with Telegram.
It is a working example, however it is worth checking exactly what it does when deploying the bot to production.

To make life easier, let's `export` some variables that the build process and the deployment process will use:
```bash
export ENVIRONMENT=prod
export LAMBDA_SECRET=mylonglambdasecret
export TELEGRAM_TOKEN=abcdef 
```
Descriptions for these variables can be found in the [Environment variables](#environment-variables) section.

You also need to have access to AWS Lambda, Gateway, AWF and IAM. The easiest way is to request an `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` token
and export them as environment variables. You could also run the deployment from an AWS EC2 instance with proper IAM roles set up.
Setting up AWS access is beyond the scope of this documentation.

Run the below to build a Linux binary, deploy it to Lambda using Terraform and deploy the associated API Gateway:
```bash
BUILD_NUMBER=11 make build-linux package deploy
```

Run the below command to tell Telegram where to find the bot: (only if you used Terraform to deploy)
```bash
make webhook
```
And the below will tell you the details of the webhook from Telegram:
```bash
make webhook-info
```

## How to use it

Please refer to the [User's Guide](freshautomations/telegram-moderator-bot/blob/master/GUIDE.md)
for additional information on how to use the bot.

## Environment variables
```
BUILD_NUMBER = 0-dev
```
Required for build.

The release or patch version of the compiled code (semantic versioning).
For example: BUILD_NUMBER=23 - code version will be: 0.1.**23**.
```
LAMBDA_SECRET = <unset_by_default>
```
Required for deploy.

The URL path prefix that the Lambda function will be listening on.
Even if someone finds the Lambda function on the Internet, they will need this secret to be able to issue commands.

```
TELEGRAM_TOKEN = <unset_by_default>
```
Required for deploy and webhook.

The token received from @BotFather on Telegram.

```
ENVIRONMENT = staging
```
Optional for deploy.

Tag the deployed environment with a name. Default value: staging.

```
LAMBDA_TIMEOUT = 15
```
Optional for build and deploy.

Timeout in seconds for the Lambda function. Default value: 15 seconds.
