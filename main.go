package main

import (
	"fmt"
	"net/http"

	guuid "github.com/google/uuid"
	"github.com/labstack/echo"
)

//Player define model of player
type Player struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Online   bool   `json:"online"`
}

//Players define array of players
type Players []Player

var players Players

func generateID() string {
	id := guuid.New()
	return id.String()
}

func getPlayers(c echo.Context) error {
	return c.JSON(http.StatusOK, players)
}

func getPlayer(c echo.Context) error {
	id := c.Param("id")
	for _, player := range players {
		if player.ID == id {
			c.JSON(http.StatusOK, player)
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

func postPlayer(c echo.Context) error {
	player := Player{}
	err := c.Bind(&player)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}
	player.ID = generateID()
	players = append(players, player)
	return c.JSON(http.StatusCreated, players)
}

func main() {
	fmt.Println("Running....")
	e := echo.New()
	e.GET("/players", getPlayers)
	e.POST("/players", postPlayer)
	e.GET("/players/:id", getPlayer)
	e.PUT("/players/:id", putPlayer)
	e.Start(":9000")
}
