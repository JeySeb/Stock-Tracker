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

	authUsecases "stock-tracker/internal/domain/authentication/usecase"
	marketDataUsecase "stock-tracker/internal/domain/marketdata/usecase"
	recommendationUsecases "stock-tracker/internal/domain/recommendation/usecase"
	recommendationValidation "stock-tracker/internal/domain/recommendation/validation"
	stockUsecases "stock-tracker/internal/domain/stocks/usecase"
	subscriptionUsecases "stock-tracker/internal/domain/subscription/usecase"
	"stock-tracker/internal/infrastructure/auth"
	"stock-tracker/internal/infrastructure/cache"
	"stock-tracker/internal/infrastructure/config"
	"stock-tracker/internal/infrastructure/database"
	"stock-tracker/internal/infrastructure/external"
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
	defer func() {
		if err := dbPool.Close(); err != nil {
			log.Error("Failed to close database connection", "error", err)
		}
	}()

	// Initialize repositories
	stockRepo := database.NewStockRepository(dbPool.GetPool(), log)
	userRepo := database.NewUserRepository(dbPool.GetPool(), log)
	brokerRepo := database.NewBrokerRepository(dbPool.GetPool())
	subscriptionRepo := database.NewSubscriptionRepository(dbPool.GetPool(), log)
	sessionRepo := database.NewSessionRepository(dbPool.GetPool())
	marketDataAnalysisRepo := database.NewMarketDataAnalysisRepository(dbPool.GetPool(), log)
	// Initialize JWT service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-in-production"
		log.Warn("Using default JWT secret - change in production!")
	}
	jwtService := auth.NewJWTService(jwtSecret)

	// Initialize cache
	cache := cache.NewInMemoryCache()

	// Initialize external API clients
	yahooClient := external.NewYahooFinanceClient(log)
	alphaVantageClient := external.NewAlphaVantageClient("", log) // Empty API key for now

	// Initialize recommendation components
	basicCalculator := recommendationValidation.NewBasicScoringCalculator(stockRepo, log)
	externalEnricher := recommendationValidation.NewExternalDataEnricher(yahooClient, alphaVantageClient, cache, log)

	// Initialize use cases
	stockQueryUC := stockUsecases.NewStockQueryUseCase(stockRepo, brokerRepo, log)
	brokerQueryUC := stockUsecases.NewBrokerQueryUseCase(stockRepo, brokerRepo, log)
	userUC := authUsecases.NewUserUseCase(userRepo, subscriptionRepo, sessionRepo, jwtService, log)
	recommendationUC := recommendationUsecases.NewTieredRecommendationUseCase(stockRepo, basicCalculator, externalEnricher, cache, log)
	// subscriptionUC := usecases.NewSubscriptionUseCase(subscriptionRepo, userRepo, log) // TODO: Use when subscription handler is implemented
	marketDataUC := marketDataUsecase.NewMarketDataAnalysisUseCase(marketDataAnalysisRepo, log)

	// Initialize middleware
	authMiddleware := infraMiddleware.NewAuthMiddleware(jwtService, log)
	rateLimiter := infraMiddleware.NewRateLimiter(log)

	// Initialize subscription use case (was commented out)
	subscriptionUC := subscriptionUsecases.NewSubscriptionUseCase(subscriptionRepo, userRepo, log)

	// Initialize handlers
	stockHandler := handlers.NewStockHandler(stockQueryUC, log)
	brokerHandler := handlers.NewBrokerHandler(brokerQueryUC, log)
	authHandler := handlers.NewAuthHandler(userUC, log)
	recommendationHandler := handlers.NewRecommendationHandler(recommendationUC, log)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionUC, log)
	marketDataHandler := handlers.NewMarketDataHandler(marketDataUC, log)
	// Initialize router
	r := setupRouter(stockHandler, brokerHandler, authHandler, recommendationHandler, subscriptionHandler, marketDataHandler, authMiddleware, rateLimiter, log, dbPool)

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
	brokerHandler *handlers.BrokerHandler,
	authHandler *handlers.AuthHandler,
	recommendationHandler *handlers.RecommendationHandler,
	subscriptionHandler *handlers.SubscriptionHandler,
	marketDataHandler *handlers.MarketDataHandler,
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

	// NEW: Additional production middleware
	r.Use(middleware.Compress(5))        // Response compression
	r.Use(middleware.Heartbeat("/ping")) // Simple health endpoint
	r.Use(middleware.NoCache)            // Disable caching for API responses
	r.Use(middleware.RedirectSlashes)    // Handle trailing slashes
	r.Use(middleware.CleanPath)          // Clean URL paths

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database connectivity
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := dbPool.GetPool().Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			if _, writeErr := w.Write([]byte(`{"status":"unhealthy","database":"disconnected"}`)); writeErr != nil {
				log.Error("Failed to write health check response", "error", writeErr)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, writeErr := w.Write([]byte(`{"status":"healthy","database":"connected","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`)); writeErr != nil {
			log.Error("Failed to write health check response", "error", writeErr)
		}
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
				r.Get("/enhanced", stockHandler.GetStocksWithEnhancedFilters)
				r.Get("/tickers", stockHandler.GetUniqueTickers)
				r.Get("/companies", stockHandler.GetUniqueCompanies)
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

			// Broker routes with optional authentication
			r.Route("/brokers", func(r chi.Router) {
				r.Use(authMiddleware.OptionalAuth) // Guest users can access with limitations
				r.Use(rateLimiter.RateLimit)       // Tier-based rate limiting
				r.Get("/scores", brokerHandler.GetBrokersWithScores)
			})

			// Protected user routes
			r.Route("/user", func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(rateLimiter.RateLimit)
				// TODO: Add user profile endpoints
			})

			// Subscription routes (IMPLEMENTED functionality)
			r.Route("/subscriptions", func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(rateLimiter.RateLimit)

				// ✅ IMPLEMENTED: Create new subscription
				r.Post("/", subscriptionHandler.CreateSubscription)

				// ✅ IMPLEMENTED: Simulate payment for subscription
				r.Post("/{id}/payment", subscriptionHandler.SimulatePayment)

				// TODO: Additional routes that can be implemented with existing logic:
				r.Get("/{id}", subscriptionHandler.GetSubscriptionByID)     // Uses: GetByID
				r.Get("/active", subscriptionHandler.GetActiveSubscription) // Uses: GetActiveByUserID
				r.Get("/plans", subscriptionHandler.GetSubscriptionPlans)   // Static data
			})

			// Recommendation routes with optional authentication
			r.Route("/recommendations", func(r chi.Router) {
				r.Use(authMiddleware.OptionalAuth) // Guest users can access with limitations
				r.Use(rateLimiter.RateLimit)       // Tier-based rate limiting
				r.Get("/", recommendationHandler.GetRecommendations)
				r.Get("/{ticker}", recommendationHandler.GetRecommendationByTicker)

				// ACTUALLY IMPLEMENTED: Preview endpoint from the handler
				r.Get("/preview/{ticker}", recommendationHandler.GetRecommendationPreview)
			})

			// Market data routes with optional authentication
			r.Route("/market-data", func(r chi.Router) {
				r.Use(authMiddleware.OptionalAuth) // Guest users can access with limitations
				r.Use(rateLimiter.RateLimit)       // Tier-based rate limiting
				r.Get("/analysis/{ticker}", marketDataHandler.GetMarketDataAnalysis)
				r.Get("/trend/{ticker}", marketDataHandler.GetMarketDataTrend)
				r.Get("/summary", marketDataHandler.GetMarketDataSummary)
				r.Get("/top-performers", marketDataHandler.GetTopPerformers)
				r.Get("/worst-performers", marketDataHandler.GetWorstPerformers)
				r.Get("/latest/{ticker}", marketDataHandler.GetLatestMarketData)
				//r.Get("/compare", marketDataHandler.GetMarketDataComparison)
				//r.Get("/volatile", marketDataHandler.GetMostVolatile)
				//r.Get("/active", marketDataHandler.GetMostActive)
				//r.Get("/high-risk", marketDataHandler.GetHighRiskTickers)
				//r.Get("/low-risk", marketDataHandler.GetLowRiskTickers)
				//r.Get("/analysis-combined/{ticker}", marketDataHandler.GetMarketDataWithStockAnalysis)
				//r.Get("/broker-correlation/{ticker}", marketDataHandler.GetCorrelationWithBrokerActions)
				//r.Get("/impact/{ticker}", marketDataHandler.GetMarketDataImpactOnRecommendations)
				r.Get("/{ticker}", marketDataHandler.GetMarketDataByTicker)
			})
			// Premium features (basic framework ready)
			r.Route("/premium", func(r chi.Router) {
				r.Use(authMiddleware.RequirePremium)
				r.Use(rateLimiter.RateLimit)
				// TODO: Add premium endpoints when AI features are implemented
			})
		})
	})

	return r
}
