package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/magicLian/gostarter/cmd/gostarter-server/backgroundSvc"
	"github.com/magicLian/gostarter/pkg/setting"
	"gitlab.wise-paas.com/ifactory/patrol-models/log"
	"golang.org/x/sync/errgroup"
)

func NewPatrolServer(cfg *setting.Cfg, registry *backgroundSvc.BackgroundServiceRegistry) *PatrolServer {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)
	return &PatrolServer{
		context:       childCtx,
		shutdownFn:    shutdownFn,
		childRoutines: childRoutines,
		log:           log.New("server"),
		cfg:           cfg,
		services:      registry.GetServices(),
	}
}

type PatrolServer struct {
	context       context.Context
	shutdownFn    context.CancelFunc
	childRoutines *errgroup.Group
	log           log.Logger
	cfg           *setting.Cfg
	services      []backgroundSvc.BackgroundService
}

func (p *PatrolServer) Run() error {
	services := p.services
	for _, svc := range services {
		select {
		case <-p.context.Done():
			return p.context.Err()
		default:
		}
		service := svc
		serviceName := reflect.TypeOf(svc).String()
		p.childRoutines.Go(func() error {
			select {
			case <-p.context.Done():
				return p.context.Err()
			default:
			}
			p.log.Debug("Starting background service", "service", serviceName)
			err := service.Run(p.context)
			if !errors.Is(err, context.Canceled) && err != nil {
				p.log.Error("Stopped ", "reason", err)
				return fmt.Errorf("%s run error: %s", serviceName, err.Error())
			}
			return nil
		})
	}

	return p.childRoutines.Wait()
}

func (i *PatrolServer) Shutdown(reason string) {
	i.log.Info("Shutdown started", "reason", reason)
	// call cancel func on root context
	i.shutdownFn()

	// wait for child routines
	i.childRoutines.Wait()
}

func (i *PatrolServer) Exit(reason error) int {
	// default exit code is 1
	code := 1

	i.log.Error("Server shutdown", "reason", reason)
	return code
}
