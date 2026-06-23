package app

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Run initializes dependencies, starts the HTTP server, and blocks until shutdown.
func Run(ctx context.Context) error {
	app, err := newApplication(ctx)
	if err != nil {
		return err
	}
	defer app.Close()

	router := app.newRouter(ctx)
	addr := fmt.Sprintf(":%s", app.conf.AppPort)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		app.logger.Info().Str("addr", addr).Msg("Server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	<-ctx.Done()
	app.logger.Info().Err(ctx.Err()).Msg("Received shutdown signal. Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		app.logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	app.logger.Info().Msg("Server stopped gracefully")

	return nil
}
