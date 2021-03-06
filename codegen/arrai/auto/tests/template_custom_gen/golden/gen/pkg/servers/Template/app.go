// Code generated by sysl DO NOT EDIT.
package template

import (
	"context"
	"log"

	"github.com/anz-bank/sysl-go/config"
	"github.com/anz-bank/sysl-go/core"
	"github.com/anz-bank/sysl-go/handlerinitialiser"
)

// DownstreamClients for Template
type DownstreamClients struct {
}

// BuildDownstreamClients ...
func BuildDownstreamClients(ctx context.Context, hooks *core.Hooks, cfg *config.DefaultConfig) (*DownstreamClients, error) {
	downstreamConfig := cfg.GenCode.Downstream.(*DownstreamConfig)
	if downstreamConfig == nil {
		downstreamConfig = &DownstreamConfig{}
	}
	var err error
	return &DownstreamClients{}, err
}

// Serve starts the server.
//
// createService must be a function with the following signature:
//
//    func(ctx context.Context, config AppConfig) (*template.ServiceInterface, *core.Hooks, error)
//
// where AppConfig is a type defined by the application programmer to
// hold application-level configuration.
//
// For a quick start, copy main.go.sample from the same directory as this
// file into its own package directory, rename it to main.go, and run it.
func Serve(
	ctx context.Context,
	createService interface{},
) error {
	err := core.Serve(
		ctx,
		&DownstreamConfig{}, createService, &ServiceInterface{},
		func(
			ctx context.Context,
			cfg *config.DefaultConfig,
			serviceIntf interface{},
			hooks *core.Hooks,
		) (interface{}, error) {
			serviceInterface := serviceIntf.(*ServiceInterface)

			if hooks == nil {
				hooks = &core.Hooks{}
			}

			var downstream *DownstreamConfig
			var is bool
			if downstream, is = cfg.GenCode.Downstream.(*DownstreamConfig); !is {
				downstream = &DownstreamConfig{
					ContextTimeout: 30,
				}
			}

			genCallbacks := config.NewCallback(
				&cfg.GenCode,
				downstream.ContextTimeout,
				hooks.MapError,
			)

			serviceHandler, err := NewServiceHandler(
				ctx,
				cfg,
				hooks,
				genCallbacks,
				serviceInterface,
			)
			if err != nil {
				return nil, err
			}

			// Construct a HTTPManager to wrap our HandlerInitialiser AKA ServiceRouter.
			// TODO standardise terminology / generally refactor.
			handlerInitialiser := NewServiceRouter(genCallbacks, serviceHandler)

			libraryConfig := &(cfg.Library)
			var adminServerConfig *config.CommonHTTPServerConfig // TODO make configurable
			var publicServerConfig *config.CommonHTTPServerConfig = &(cfg.GenCode.Upstream.HTTP)
			// TODO make it possible to disable handlers through configuration.
			enabledHandlers := []handlerinitialiser.HandlerInitialiser{handlerInitialiser}
			manager := core.NewHTTPManagerShim(libraryConfig, adminServerConfig, publicServerConfig, enabledHandlers)
			return manager, nil
		},
	)
	if err != nil {
		log.Print(err)
	}
	return err
}
