// main package that executes the code
//
// telegram-moderator-bot is a web API that allows Telegram bots to moderate channels
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/freshautomations/telegram-moderator-bot/config"
	"github.com/freshautomations/telegram-moderator-bot/context"
	"github.com/freshautomations/telegram-moderator-bot/db"
	"github.com/freshautomations/telegram-moderator-bot/defaults"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"os"
	"time"
)

// lambdaInitialized is an indicator that tells if the AWS Lambda function is in the startup phase.
var lambdaInitialized = false

// Translates Gorilla Mux calls to AWS API Gateway calls
var lambdaProxy func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// LambdaHandler is the callback function when the application is set up as an AWS Lambda function.
func LambdaHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if !lambdaInitialized {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("[init] lambda start %s", defaults.Version)

		var err error
		ctx, err := Initialization(context.NewInitialContext())
		if err != nil {
			log.Printf("[init] initialization failed: %v", err)
			errbody, _ := json.Marshal(context.ErrorMessage{
				Message: "System could not be initialized, please contact the administrator.",
			})
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       string(errbody),
			}, nil
		}

		r := AddRoutes(ctx)
		muxLambda := gorillamux.New(r)
		lambdaProxy = muxLambda.Proxy

		lambdaInitialized = true
	}

	return lambdaProxy(req)

}

// WebserverHandler is the function that is called when the `--webserver` parameter is invoked.
// It sets up a local webserver for handling incoming requests.
func WebserverHandler(localCtx *context.InitialContext) {
	log.Printf("[init] webserver start %s", defaults.Version)

	var err error
	ctx, err := Initialization(localCtx)
	if err != nil {
		log.Fatalf("initialization failed: %v\n", err)
	}

	r := AddRoutes(ctx)

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", localCtx.WebserverIp, localCtx.WebserverPort),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Printf("[final] caught signal: %+v", sig)
		log.Print("[final] waiting 2 seconds to finish processing")
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Initialization creates and populates the context and sets up connectivity to the testnet.
func Initialization(initialContext *context.InitialContext) (ctx *context.Context, err error) {

	ctx = context.New()

	if initialContext.LocalExecution {
		log.Printf("[init] loading config file %s", initialContext.ConfigFile)
		ctx.Cfg, err = config.GetConfigFromFile(initialContext.ConfigFile)
		if err != nil {
			return
		}
	} else {
		log.Printf("[init] loading config from environment variables")
		ctx.Cfg, err = config.GetConfigFromENV()
		if err != nil {
			return
		}
	}

	db.Initialize(ctx)

	printCfg := *ctx.Cfg
	printCfg.TelegramToken = redact(printCfg.TelegramToken)
	log.Printf("[init] config loaded: %+v", printCfg)

	log.Print("[init] initialized context")

	return
}

// redact changes a string to XXXXXX - used to redact passwords when logging.
func redact(s string) string {
	if len(s) < 2 {
		return "RD"
	}
	return "REDACTED"
}

func main() {
	initialCtx := context.NewInitialContext()

	flag.BoolVar(&initialCtx.LocalExecution, "webserver", false, "run a local web-server instead of as an AWS Lambda function")
	flag.StringVar(&initialCtx.ConfigFile, "config", "tmb.conf", "read config from this local file")
	flag.StringVar(&initialCtx.WebserverIp, "ip", "127.0.0.1", "IP to listen on")
	flag.UintVar(&initialCtx.WebserverPort, "port", 3000, "Port to listen on")
	flag.Parse()

	//--webserver
	if initialCtx.LocalExecution {
		WebserverHandler(initialCtx)
	} else {
		//Lambda function on AWS
		lambda.Start(LambdaHandler)
	}
}
