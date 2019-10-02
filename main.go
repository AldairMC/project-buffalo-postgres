package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	guuid "github.com/google/uuid"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

func initDB(filepath string) *slq.DB {
	db, err := sql.Open("sqlite3", filepath)

	// Here we check for any db errors then exit
	if err != nil {
		panic(err)
	}

	// If we don't get any errors but somehow still don't get a db connection
	// we exit as well
	if db == nil {
		panic("db nil")
	}
	return db
}

//Post define model of player
type Post struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

//Posts define array of players
type Posts []Post

var posts Posts

func generateID() string {
	id := guuid.New()
	return id.String()
}

func getPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, posts)
}

func getPost(c echo.Context) error {
	id := c.Param("id")
	for _, post := range posts {
		if post.ID == id {
			c.JSON(http.StatusOK, post)
		}
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func putPlayer(c echo.Context) error {
	id := c.Param("id")
	for i, _ := range players {
		if players[i].ID == id {
			players[i].Online = true
			c.JSON(http.StatusOK, players)
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
	id, _ := c.Param("id")
	for i, _ := range posts {
		if posts[i].ID == id {
			posts = append(posts[:id], posts[i+1]...)
			return c.JSON(http.StatusOK, posts)
		}
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func main() {
	fmt.Println("Running....")
	e := echo.New()
	e.GET("/players", getPlayers)
	e.POST("/players", postPlayer)
	e.GET("/players/:id", getPlayer)
	e.PUT("/players/:id", putPlayer)
	e.Start(":8080")
}
