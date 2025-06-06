package web

import (
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"bitbucket.org/junglee_games/getsetgo/validate"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

// Create a seperate struct for storing keyValue pair of headers and values and then use that ctx.Values().Set(KeyValues, &v)

type WebValues struct {
	PlatformName PlatformName `header:"X-Platform-Name"`
	UserID       string       `header:"X-User-Id"`
	RequestID    string       `header:"X-Request-Id"`
}

type PlatformName string

const (
	PS_CASHAPP PlatformName = "psrmg"
	CASHAPP    PlatformName = "apk"
	PSAPP      PlatformName = "psapp"
	IPA        PlatformName = "ipa"
)

const (
	webKeyValues = "web.service.ctx"
)

// KeyValues is how request values are stored/retrieved.
const KeyValues string = "dms.service.ctx"

// Values represent state for each request.
type Values struct {
	RequestID  string
	Now        time.Time
	StatusCode int
	Err        error
}

// A Handler is a type that handles an http request within our own little mini
// framework.
type Handler func(ctx iris.Context) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	*iris.Application
	shutdown chan os.Signal
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...iris.Handler) *App {
	fw := iris.New()
	iris.WithFireMethodNotAllowed(fw)
	log := fw.Logger()
	log.SetLevel("debug")
	log.SetPrefix(string(log.Prefix))
	log.SetTimeFormat("2006/01/02 15:04:05")
	log.Debug(`Log level set to "debug"`)
	log.Debugf(`Setting %d middlewares`, len(mw))
	app := App{
		Application: fw,
		shutdown:    shutdown,
	}
	app.wrapMiddleware(mw)
	// Catch a specific error code.
	app.Application.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		// Set the context with the required values to
		// process the request.
		v := Values{
			RequestID: uuid.NewString(),
			Now:       time.Now(),
		}
		ctx.Values().Set(KeyValues, &v)
		er := validate.ErrorResponse{
			Err: "Path not supported",
		}
		Respond(ctx, er, http.StatusMethodNotAllowed)
	})
	return &app
}

// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
// pair, this makes for really easy, convenient routing.
func (a *App) Handle(method string, path string, handler Handler, mw ...iris.Handler) {
	a.wrapMiddleware(mw)
	// The function to execute for each request.
	h := func(ctx iris.Context) {

		// Call the wrapped handler functions.
		if err := handler(ctx); err != nil {
			v := ctx.Values().Get(KeyValues).(*Values)
			v.Err = err
			// If we receive the shutdown err we need to return it
			// back to the base handler to shutdown the service.
			if ok := IsShutdown(err); ok {
				ctx.StopWithError(http.StatusInternalServerError, errors.New("something went wrong"))
				a.SignalShutdown()
			}
		}
	}

	// Add this handler for the specified verb and route.
	a.Application.Handle(method, path, h)
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
