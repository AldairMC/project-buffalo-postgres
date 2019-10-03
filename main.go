package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)

	if err != nil {
		panic(err)
	}

	if db == nil {
		panic("db nil")
	}
	return db
}

func migrate(db *sql.DB) {
	sql := `
		CREATE TABLE IF NOT EXISTS posts(
			id VARCHAR NOT NULL PRIMARY KEY,
			title VARCHAR NOT NULL,
			description VARCHAR NOT NULL,
			author VARCHAR NOT NULL,
			createAt DATETIME NOT NULL,
			updatedAt DATETIME NULL
		);
	`

	_, err := db.Exec(sql)

	if err != nil {
		panic(err)
	}
}

//Post define model of player
type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

//Posts define array of players
type Posts []Post

var posts Posts

func generateID() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(10000)
}

func getPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, posts)
}

func getPost(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, post := range posts {
		if post.ID == id {
			c.JSON(http.StatusOK, post)
		}
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func putPost(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, _ := range posts {
		if posts[i].ID == id {
			c.JSON(http.StatusOK, posts)
		}
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func postPost(c echo.Context) error {
	post := Post{}
	err := c.Bind(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}
	post.ID = generateID()
	post.CreatedAt = time.Now().String()
	posts = append(posts, post)
	return c.JSON(http.StatusCreated, posts)
}

func deletePost(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, _ := range posts {
		if posts[i].ID == id {
			posts = append(posts[:i], posts[i+1:]...)
			return c.JSON(http.StatusOK, posts)
		}
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func main() {
	db := initDB("storage.db")
	migrate(db)
	fmt.Println("Running....")
	e := echo.New()
	e.GET("/posts", getPosts)
	e.POST("/posts", postPost)
	e.GET("/posts/:id", getPost)
	e.PUT("/posts/:id", putPost)
	e.DELETE("/posts/:id", deletePost)
	e.Start(":8080")
}
