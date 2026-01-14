package app

import (
	"context"
	"errors"
	"io"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/api-service/internal/config"
	"github.com/smarrog/task-board/api-service/internal/middleware"
	"github.com/smarrog/task-board/api-service/internal/transport/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg *config.Config
	log *zerolog.Logger

	coreConn *grpc.ClientConn
	authConn *grpc.ClientConn
	httpApp  *fiber.App
}

func New(cfg *config.Config, log *zerolog.Logger) *App {
	app := &App{
		cfg: cfg,
		log: log,
	}

	return app
}

func (a *App) Init() error {
	conn, err := grpc.NewClient(a.cfg.CoreGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	a.coreConn = conn

	authConn, err := grpc.NewClient(a.cfg.AuthGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		_ = a.coreConn.Close()
		return err
	}
	a.authConn = authConn

	app := fiber.New(fiber.Config{
		AppName:      a.cfg.FiberAppName,
		ErrorHandler: http.ErrorHandler(a.log),
		IdleTimeout:  a.cfg.FiberIdleTimeout,
		ReadTimeout:  a.cfg.FiberReadTimeout,
		WriteTimeout: a.cfg.FiberWriteTimeout,
	})

	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendString("ok") })

	app.Use(middleware.Recover(a.log))   // ловим панику, превращаем в 500 + лог
	app.Use(middleware.RequestID())      // проставляем request id
	app.Use(middleware.AccessLog(a.log)) // логируем запросы\ответы

	handler := http.NewHandler(a.log, a.cfg, a.coreConn)
	authHandler := http.NewAuthHandler(a.log, a.cfg, a.authConn)

	v1 := app.Group("/v1")

	v1.Post("/auth/register", authHandler.Register)
	v1.Post("/auth/login", authHandler.Login)

	protected := v1.Group("", middleware.JWT(a.cfg.JWTSecret))

	handler.Register(protected)

	a.httpApp = app
	return nil
}

func (a *App) Run(ctx context.Context) error {
	if a.httpApp == nil {
		return errors.New("app not initialized")
	}

	ln, err := net.Listen("tcp", a.cfg.HTTPAddr)
	if err != nil {
		return err
	}
	defer a.closeAll()

	errCh := make(chan error, 1)
	go func() {
		a.log.Info().Str("addr", a.cfg.HTTPAddr).Msg("HTTP gateway listening")
		errCh <- a.httpApp.Listener(ln)
	}()

	select {
	case <-ctx.Done():
		a.log.Info().Msg("Shutting down HTTP gateway")

		done := make(chan struct{})

		go func() {
			_ = a.httpApp.Shutdown()
			close(done)
		}()

		shCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
		defer cancel()

		select {
		case <-done:
			a.log.Info().Msg("HTTP gateway stopped gracefully")
		case <-shCtx.Done():
			a.log.Error().Msg("HTTP gateway graceful shutdown timeout exceeded")
		}

		return nil
	case err := <-errCh:
		if err == nil || errors.Is(err, io.EOF) {
			return nil
		}
		a.log.Error().Err(err).Msg("HTTP gateway stopped with error")
		return err
	}
}

func (a *App) closeAll() {
	if a.coreConn != nil {
		_ = a.coreConn.Close()
	}
	if a.authConn != nil {
		_ = a.authConn.Close()
	}
}
