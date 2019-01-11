// Context package defines primitives for generic context handling during execution.
package context

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/freshautomations/telegram-moderator-bot/config"
	"github.com/freshautomations/telegram-moderator-bot/defaults"
	"log"
	"net/http"
)

// Context holds current execution details.
type Context struct {
	AWSSession *session.Session

	DDBSession *dynamodb.DynamoDB

	DBUserTable string

	DBWarnTable string

	// Application configuration
	Cfg *config.Config
}

// InitialContext holds the input parameter details at the start of execution.
type InitialContext struct {
	// --ip IP address of local webserver
	WebserverIp string

	// --port Port number of local webserver
	WebserverPort uint

	// --config Config file for local execution
	ConfigFile string

	// --webserver was set
	LocalExecution bool
}

// New creates a fresh Context.
func New() *Context {
	return &Context{}
}

// NewInitialContext creates a fresh InitialContext.
func NewInitialContext() *InitialContext {
	return &InitialContext{}
}

// ErrorMessage defines the message structure returned when an error happens.
type ErrorMessage struct {
	Message string `json:"message"`
}

// Handler is an abstraction layer to standardize web API returns, if an error happens.
type Handler struct {
	C *Context
	H func(*Context, http.ResponseWriter, *http.Request) (int, error)
}

// ServeHTTP is a wrapper around web API calls, that adds a default Content-Type and formats outgoing error messages.
func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", defaults.ContentType)
	if status, err := fn.H(fn.C, w, r); err != nil {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(ErrorMessage{err.Error()})
		log.Printf("%d %s", status, err.Error())
	}
}
