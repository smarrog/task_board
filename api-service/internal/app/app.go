package app

import (
	"context"
	"errors"
	"net"
	"time"

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

	handler := http.NewHandler(a.log, a.cfg, a.coreConn) // вызываем gRPC

	v1 := app.Group("/v1")

	// заглушки
	v1.Post("/auth/register", http.NotImplemented("use Auth Service"))
	v1.Post("/auth/login", http.NotImplemented("use Auth Service"))

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

		timer := time.NewTimer(a.cfg.ShutdownTimeout)
		defer timer.Stop()

		select {
		case <-done:
			if !timer.Stop() {
				// если таймер уже успел сработать, вычистим сигнал, чтобы не утекал
				select {
				case <-timer.C:
				default:
				}
			}
			a.log.Info().Msg("HTTP gateway stopped gracefully")
		case <-timer.C:
			a.log.Error().Msg("HTTP gateway graceful shutdown timeout exceeded")
		}

		if a.coreConn != nil {
			_ = a.coreConn.Close()
		}

		return nil
	case err := <-errCh:
		if a.coreConn != nil {
			_ = a.coreConn.Close()
		}
		return err
	}
}
