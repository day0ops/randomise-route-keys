package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	defaultAddress      = ":8081"
	defaultShutdownWait = 5 * time.Second
)

type Server struct {
	ctx     context.Context
	log     *zap.Logger
	address string
	httpSrv *http.Server
}

type RouteStrList struct {
	Routes []string `json:"route-keys"`
}

type RouteDecision struct {
	Decision string `json:"decision"`
}

type Option func(*Server)

func New(ctx context.Context, log *zap.Logger, handler *gin.Engine, opts ...Option) (*Server, error) {
	srv := &Server{
		ctx: ctx,
		log: log,
	}

	for _, opt := range opts {
		opt(srv)
	}

	if srv.address == "" {
		srv.address = defaultAddress
	}

	// register the main routes
	RegisterRoutes(handler)

	srv.httpSrv = &http.Server{
		Addr:    srv.address,
		Handler: handler.Handler(),
	}

	err := ReadRouteListFile()
	if err != nil {
		srv.log.Error("error reading route list", zap.Error(err))
		return nil, err
	}

	return srv, nil
}

func (s *Server) Serve() error {
	if s.ctx == nil {
		s.ctx = context.TODO()
	}

	errCh := make(chan error, 1)

	go func() {
		s.log.Info("starting server", zap.String("address", s.address))
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("cannot listen: %w", err)
			return
		}
		errCh <- nil
	}()

	select {
	case <-s.ctx.Done():
		return s.Stop()
	case err := <-errCh:
		return err
	}
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if s.httpSrv != nil {
		s.log.Info("stopping grpc server")
		if err := s.httpSrv.Shutdown(ctx); err != nil {
			s.log.Error("server shutdown", zap.Error(err))
		}
	}
	time.Sleep(defaultShutdownWait)
	return nil
}

func WithServerAddress(address string) Option {
	return func(s *Server) {
		s.address = fmt.Sprintf(":%s", address)
	}
}

func HealthHandler(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func RootHandler(c *gin.Context) {
	rd, err := getRandomRoute()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, rd)
}

func getRandomRoute() (RouteDecision, error) {
	routerList := GetCachedRouteList()
	if len(routerList.Routes) > 0 {
		random := routerList.Routes[rand.Intn(len(routerList.Routes))]
		return RouteDecision{Decision: random}, nil
	}

	return RouteDecision{}, fmt.Errorf("no route list found")
}
