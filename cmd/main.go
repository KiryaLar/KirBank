package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"go_project/internal/config"
	"go_project/internal/handlers"
	"go_project/internal/middleware"
	"go_project/internal/repositories"
	"go_project/internal/services"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		logrus.Fatal("cannot load config:", err)
	}
	dbURL := cfg.Database.URL
	if dbURL == "" {
		logrus.Fatal("DATABASE_URL environment variable not set")
	}
	jwtSecret := cfg.Auth.JWTSecret
	if jwtSecret == "" {
		logrus.Fatal("JWT_SECRET not set")
	}
	encryptionKey := cfg.Auth.EncryptionKey
	if encryptionKey == "" {
		logrus.Fatal("CARD_ENC_KEY not set")
	}
	hmacSecret := cfg.Auth.HMACSecret
	if hmacSecret == "" {
		logrus.Fatal("HMAC_SECRET not set")
	}

	_ = cfg.SMTP.Host
	_ = cfg.SMTP.Port

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logrus.Fatal("Failed to connect to DB", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		logrus.Fatal("Database ping failed:", err)
	}

	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	cardRepo := repositories.NewCardRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	creditRepo := repositories.NewCreditRepository(db)
	scheduleRepo := repositories.NewPaymentScheduleRepository(db)

	authService := services.NewAuthService(userRepo, jwtSecret)
	accountService := services.NewAccountService(accountRepo)
	cardService := services.NewCardService(cardRepo, accountRepo, encryptionKey)
	transactionService := services.NewTransactionService(accountRepo, transactionRepo, hmacSecret)
	creditService := services.NewCreditService(creditRepo, accountRepo, scheduleRepo)

	creditService.StartOverduePayments()

	authHandler := handlers.NewAuthHandler(authService)
	accountHandler := handlers.NewAccountHandler(accountService)
	cardHandler := handlers.NewCardHandler(cardService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	creditHandler := handlers.NewCreditHandler(creditService)

	r := chi.NewRouter()

	// middleware для логирования запросов
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logrus.Infof("%s %s %v", r.Method, r.RequestURI, time.Since(start))
		})
	})

	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)

	r.Route("/", func(pr chi.Router) {
		pr.Use(middleware.JWTAuthMiddleware(jwtSecret))
		pr.Post("/accounts", accountHandler.CreateAccount)
		pr.Post("/cards", cardHandler.CreateCard)
		pr.Post("/transfer", transactionHandler.Transfer)
		pr.Get("/analytics", transactionHandler.Analytics)
		pr.Post("/credits", creditHandler.CreateCredit)
		pr.Get("/credits/{creditId}/schedule", creditHandler.GetPaymentSchedule)
		pr.Get("/accounts/{accountId}/balance", accountHandler.GetBalance)
		pr.Get("/accounts/{accountId}/predict", accountHandler.PredictBalance)
	})
	//	Запуск HTTP-сервера
	port := cfg.Server.Port
	logrus.Infof("Server started on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logrus.Fatal(err)
	}
}
