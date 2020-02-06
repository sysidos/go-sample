package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	viper.SetConfigFile(`.env`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {

	// Database Connections
	db, err := StoreInit().connect()
	if err != nil {
		log.Fatal(err)
	}
	println(db)

	connectRedis()

	// Echo instance
	e := echo.New()

	// Echo Flags
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	port := fmt.Sprintf(":%s", viper.GetString("port"))
	e.Logger.Fatal(e.Start(port))
}

type Store struct{}

func StoreInit() *Store {
	return &Store{}
}

func (store *Store) connect() (*gorm.DB, error) {
	host := viper.GetString("POSTGRES_HOST")
	port := viper.GetString("POSTGRES_PORT")
	user := viper.GetString("POSTGRES_USER")
	password := viper.GetString("POSTGRES_PASSWORD")
	dbname := viper.GetString("POSTGRES_DB")
	mode := viper.GetString("POSTGRES_SSL")

	args := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, mode)

	db, err := gorm.Open("postgres", args)
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
	}

	// sets the maximum number of connections in the idle connection pool.
	db.DB().SetMaxIdleConns(10)

	// sets the maximum number of open connections to the database.
	db.DB().SetMaxOpenConns(100)

	// sets the maximum amount of time a connection may be reused.
	db.DB().SetConnMaxLifetime(time.Hour)

	err = db.DB().Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Close Database Connection on Error
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	return db, err
}

func connectRedis()  {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
		log.Fatal(err)
		os.Exit(1)
	}

	defer func() {
		client.Close()
		fmt.Println("closed")
	}()
	//fmt.Println(pong, err)
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
