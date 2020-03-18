package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"0proxy.io/core/common"
	"0proxy.io/core/config"
	"0proxy.io/core/logging"
	zc "0proxy.io/zproxycore/common"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func initializeConfig() {
	config.Configuration.ChainID = viper.GetString("server_chain.id")
	config.Configuration.SignatureScheme = viper.GetString("server_chain.signature_scheme")
	config.Configuration.Port = viper.GetInt("port")
	config.Configuration.DataShards = viper.GetInt("allocation_default.data_shards")
	config.Configuration.ParityShards = viper.GetInt("allocation_default.parity_shards")
	config.Configuration.MinSubmit = viper.GetInt("allocation_default.min_submit")
	config.Configuration.MinConfirmation = viper.GetInt("allocation_default.min_confirmation")
	config.Configuration.Miners = viper.GetStringSlice("miners")
	config.Configuration.Sharders = viper.GetStringSlice("sharders")
}
func initHandlers(r *mux.Router) {
	r.HandleFunc("/", HomePageHandler)
	zc.SetupHandlers(r)
}

var startTime time.Time

func main() {
	deploymentMode := flag.Int("deployment_mode", 2, "deployment_mode")
	flag.Parse()

	config.Configuration.DeploymentMode = byte(*deploymentMode)
	config.SetupDefaultConfig()
	config.SetupConfig()

	if config.Development() {
		logging.InitLogging("development")
	} else {
		logging.InitLogging("production")
	}
	initializeConfig()

	common.SetupRootContext(context.Background())

	address := fmt.Sprintf(":%v", config.Configuration.Port)

	var server *http.Server
	r := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	rHandler := handlers.CORS(originsOk, headersOk, methodsOk)(r)
	if config.Development() {
		server = &http.Server{
			Addr:           address,
			ReadTimeout:    30 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        rHandler,
		}
	} else {
		server = &http.Server{
			Addr:           address,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        rHandler,
		}
	}
	common.HandleShutdown(server)

	initHandlers(r)
	startTime = time.Now().UTC()

	fmt.Println("Listening on port", address)
	log.Fatal(server.ListenAndServe())
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<div>Running since %v ...\n", startTime)
	fmt.Fprintf(w, "<div>Working on the chain: %v</div>\n", config.Configuration.ChainID)
	fmt.Fprintf(w, "<div>I am 0proxy with <ul><li>miners:%v</li><li>sharders:%v</li></ul></div>\n",
		config.Configuration.Miners, config.Configuration.Sharders)
}
