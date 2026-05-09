package main

import (
	"log"
	"net/http"
	"os"

	"casadavo/internal/handler"
	mw "casadavo/internal/middleware"
	"casadavo/internal/repository"
	"casadavo/internal/service"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./casadavo.db"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	adminEmail := os.Getenv("ADMIN_EMAIL")

	db, err := repository.Open(dbPath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	mw.SetSecret(jwtSecret)

	userRepo := repository.NewUserRepo(db)

	if adminEmail != "" {
		if err := userRepo.EnsureAdminRole(adminEmail); err != nil {
			log.Printf("warn: could not enforce admin role for %s: %v", adminEmail, err)
		}
	}
	tableRepo := repository.NewTableRepo(db)
	availRepo := repository.NewAvailabilityRepo(db)
	slotRepo := repository.NewTimeSlotRepo(db)
	resRepo := repository.NewReservationRepo(db)
	layoutRepo := repository.NewLayoutRepo(db)

	authSvc := service.NewAuthService(userRepo, jwtSecret, adminEmail)
	resSvc := service.NewReservationService(resRepo, slotRepo, availRepo, layoutRepo)

	authH := handler.NewAuthHandler(authSvc, userRepo)
	tableH := handler.NewTableHandler(tableRepo)
	availH := handler.NewAvailabilityHandler(availRepo, slotRepo)
	resH := handler.NewReservationHandler(resSvc, resRepo)
	layoutH := handler.NewLayoutHandler(layoutRepo)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", authH.Register)
		r.Post("/auth/login", authH.Login)
		r.With(mw.Auth).Get("/auth/me", authH.Me)

		r.Group(func(r chi.Router) {
			r.Use(mw.Auth)

			// Tables — list is open to all authenticated; mutations are manager-only
			r.Get("/tables", tableH.List)
			r.Group(func(r chi.Router) {
				r.Use(mw.RequireRole("gerente"))
				r.Post("/tables", tableH.Create)
				r.Put("/tables/{id}", tableH.Update)
				r.Delete("/tables/{id}", tableH.Delete)
			})

			// Availability — read open to all authenticated; write manager-only
			r.Get("/availability", availH.List)
			r.Get("/availability/{id}/slots", availH.ListSlots)
			r.Group(func(r chi.Router) {
				r.Use(mw.RequireRole("gerente"))
				r.Post("/availability", availH.Create)
				r.Put("/availability/{id}", availH.Update)
				r.Post("/availability/{id}/slots", availH.CreateSlot)
				r.Delete("/slots/{id}", availH.DeleteSlot)
			})

			// Table layout exclusions — read open to all authenticated; write manager-only
			r.Get("/layout", layoutH.ListExclusions)
			r.Group(func(r chi.Router) {
				r.Use(mw.RequireRole("gerente"))
				r.Post("/layout", layoutH.AddExclusion)
				r.Delete("/layout/{id}", layoutH.RemoveExclusion)
			})

			// Reservations
			r.Get("/reservations", resH.List)
			r.Post("/reservations", resH.Create)
			r.Put("/reservations/{id}", resH.Update)
			r.Delete("/reservations/{id}", resH.Delete)
			r.Group(func(r chi.Router) {
				r.Use(mw.RequireRole("gerente"))
				r.Put("/reservations/{id}/status", resH.UpdateStatus)
			})
		})
	})

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
