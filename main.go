package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/jackc/pgx/v5"
)

// var db *pgx.Conn
var conn *pgx.Conn

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func main() {

	// access .env as default
	errEnv := godotenv.Load()
	if errEnv != nil {
		fmt.Println("Error occurred while loading .env file:", errEnv)
		os.Exit(1)
	}

	// read .env file
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASS")

	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)

	var errSql error
	conn, errSql = pgx.Connect(context.Background(), connStr)
	if errSql != nil {
		fmt.Println("Error connecting to the database:", errSql)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	fmt.Println("Connected!")

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)                      //optional to not get warning
	router.SetTrustedProxies([]string{"192.168.1.2"}) //to trust only a specific value

	router.GET("/albums", func(c *gin.Context) {
		albums, err := getAllAlbums(conn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.JSON(http.StatusOK, albums)
	})

	router.Run("localhost:8080")

}

func getAllAlbums(db *pgx.Conn) ([]Album, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM album")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An albums slice to hold data from returned rows.
	var albums []Album
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var a Album
		if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, nil
}
