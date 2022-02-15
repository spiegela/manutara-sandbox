package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/graphql-go/handler"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/spiegela/manutara/pkg/apis"
	"github.com/spiegela/manutara/pkg/service/schema"
)

func main() {
	flag.StringVar(&bindAddr, "bind-addr", ":8080",
		"The address to which the service endpoint binds.")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8081",
		"The address to which the metrics endpoint binds.")
	flag.Parse()
	logf.SetLogger(logf.ZapLogger(false))
	log := logf.Log.WithName("graphql-service")

	// Get a config to talk to the apiserver
	log.Info("setting up client for manager")
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "unable to set up client config")
		os.Exit(1)
	}

	// Create a new Cmd to provide shared dependencies and start components
	log.Info("setting up manager")
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: metricsAddr})
	if err != nil {
		log.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	log.Info("setting up scheme")
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "unable add APIs to scheme")
		os.Exit(1)
	}

	graphQLHandler := handler.New(&handler.Config{
		Schema:     schema.GenerateSchema(cfg),
		Pretty:     true,
		Playground: true,
	})

	//http.Handle("/graphql/login", loginHandler)
	http.Handle("/graphql", graphQLHandler)

	log.Info("GraphQL server starting", "bind-addr", bindAddr)
	err = http.ListenAndServe(bindAddr, nil)
	if err != nil {
		panic(err)
	}
}
