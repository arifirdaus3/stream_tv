package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/arifirdaus3/stream_tv/model"
	"github.com/arifirdaus3/stream_tv/populate"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/patrickmn/go-cache"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbDSN := os.Getenv("DB_DSN")
	dbConnection, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	if dbConnection == 0 {
		dbConnection = 10
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dbDSN,
	}), &gorm.Config{})

	sql, _ := db.DB()
	sql.SetMaxOpenConns(dbConnection)
	sql.SetMaxIdleConns(dbConnection / 2)

	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Category{}, &model.Language{}, &model.Country{}, &model.Region{}, &model.Subdivision{}, &model.Channel{}, &model.Guide{})
	if err != nil {
		log.Fatal(err)
	}
	cache := cache.New(24*time.Hour, 10*time.Minute)
	routeHandler := handler{
		db:    db,
		cache: cache,
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Syncronize data
	r.Get("/sync", func(w http.ResponseWriter, r *http.Request) {
		err := populate.All(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/category", routeHandler.handleCategory)
	r.Get("/language", routeHandler.handleLanguage)
	r.Get("/country", routeHandler.handleCountry)
	r.Get("/region", routeHandler.handleRegion)
	r.Get("/subdivision", routeHandler.handleSubDivision)
	r.Get("/channel", routeHandler.handleChannel)
	r.Get("/guide", routeHandler.handleGuide)

	port := os.Getenv("PORT")
	fmt.Println("listen at port ", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
