package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"stock-tracker/internal/domain/usecases"
	"stock-tracker/internal/infrastructure/auth"
	"stock-tracker/internal/infrastructure/config"
	"stock-tracker/internal/infrastructure/database"
	infraMiddleware "stock-tracker/internal/infrastructure/middleware"
	"stock-tracker/internal/presentation/handlers"
	"stock-tracker/pkg/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Initialize logger
	log := logger.New(cfg.LogLevel)
	log.Info("Starting stock recommendation API server")

	// Initialize database connection pool
	dbPool, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Error("Failed to initialize database", "error", err)
		panic(err)
	}
	defer dbPool.Close()

	// Initialize repositories
	stockRepo := database.NewStockRepository(dbPool.GetPool(), log)
	brokerRepo := database.NewBrokerRepository(dbPool.GetPool())
	userRepo := database.NewUserRepository(dbPool.GetPool(), log)
	sessionRepo := database.NewSessionRepository(dbPool.GetPool())
	subscriptionRepo := database.NewSubscriptionRepository(dbPool.GetPool(), log)

	// Initialize JWT service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-in-production"
		log.Warn("Using default JWT secret - change in production!")
	}
	jwtService := auth.NewJWTService(jwtSecret)

	// Initialize use cases
	stockQueryUC := usecases.NewStockQueryUseCase(stockRepo, brokerRepo, log)
	userUC := usecases.NewUserUseCase(userRepo, subscriptionRepo, sessionRepo, jwtService, log)
	// subscriptionUC := usecases.NewSubscriptionUseCase(subscriptionRepo, userRepo, log) // TODO: Use when subscription handler is implemented

	// Initialize middleware
	authMiddleware := infraMiddleware.NewAuthMiddleware(jwtService, log)
	rateLimiter := infraMiddleware.NewRateLimiter(log)

	// Initialize handlers
	stockHandler := handlers.NewStockHandler(stockQueryUC, log)
	authHandler := handlers.NewAuthHandler(userUC, log)

	// Initialize router
	r := setupRouter(stockHandler, authHandler, authMiddleware, rateLimiter, log, dbPool)

	// Configure server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info("Starting HTTP server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Shutdown server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Failed to shutdown server", "error", err)
	}

	log.Info("Server stopped")
}

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Expose-Headers", "Link")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "300")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func setupRouter(
	stockHandler *handlers.StockHandler,
	authHandler *handlers.AuthHandler,
	authMiddleware *infraMiddleware.AuthMiddleware,
	rateLimiter *infraMiddleware.RateLimiter,
	log logger.Logger,
	dbPool *database.Connection,
) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(corsMiddleware)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database connectivity
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := dbPool.GetPool().Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unhealthy","database":"disconnected"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","database":"connected","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// Public authentication routes
			r.Route("/auth", func(r chi.Router) {
				r.Use(rateLimiter.RateLimit) // Rate limiting for auth
				r.Post("/register", authHandler.Register)
				r.Post("/login", authHandler.Login)
				r.Post("/refresh", authHandler.RefreshToken)
			})

			// Stock routes with optional authentication
			r.Route("/stocks", func(r chi.Router) {
				r.Use(authMiddleware.OptionalAuth) // Guest users can access with limitations
				r.Use(rateLimiter.RateLimit)       // Tier-based rate limiting
				r.Get("/", stockHandler.GetStocks)
				r.Get("/{id}", stockHandler.GetStockByID)
				r.Get("/{ticker}", stockHandler.GetStockByTicker)
				r.Get("/stats", stockHandler.GetStats)

				// Protected routes for authenticated users
				r.Group(func(r chi.Router) {
					r.Use(authMiddleware.RequireAuth)
					r.Post("/", stockHandler.CreateStock)
					r.Put("/{id}", stockHandler.UpdateStock)
					r.Delete("/{id}", stockHandler.DeleteStock)
				})
			})

			// Protected user routes
			r.Route("/user", func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(rateLimiter.RateLimit)
				// TODO: Add user profile endpoints
			})

			// Premium subscription routes
			r.Route("/subscriptions", func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(rateLimiter.RateLimit)
				// TODO: Add subscription endpoints when handler is implemented
			})

			// Premium features (AI chat, advanced analytics)
			r.Route("/premium", func(r chi.Router) {
				r.Use(authMiddleware.RequirePremium)
				r.Use(rateLimiter.RateLimit)
				// TODO: Add premium endpoints
			})
		})
	})

	return r
}
