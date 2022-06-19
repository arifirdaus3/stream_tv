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
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbConnection, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	if dbConnection == 0 {
		dbConnection = 10
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbPort, dbUser, dbName, dbPassword)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
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

	fmt.Println("listen at port 80")
	http.ListenAndServe(":80", r)
}
