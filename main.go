package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Torrent struct {
	ID      int     `json:"id" uri:"id"`
	Title   string  `json:"title" uri:"title"`
	Size    float64 `json:"size"`
	Seeder  int     `json:"seeder"`
	Leecher int     `json:"leecher"`
	Date    string  `json:"date"`
}

var torrents = []Torrent{
	{ID: 1, Title: "Riders.of.Justice.2020.DANISH.1080p.US.BluRay.AVC.DTS-HD.MA.5.1-FGT", Size: 22.05, Seeder: 23, Leecher: 62, Date: "2021-08-24 10:01:49"},
	{ID: 2, Title: "The.Suicide.Squad.2021.1080p.WEBRip.x264-RARBG", Size: 2.52, Seeder: 7587, Leecher: 694, Date: "2021-08-06 02:44:27"},
	{ID: 3, Title: "The.Green.Knight.2021.1080p.AMZN.WEBRip.DDP5.1.Atmos.x264-NOGRP", Size: 5.96, Seeder: 3425, Leecher: 389, Date: "2021-08-19 09:26:58"},
	{ID: 4, Title: "Jungle.Cruise.2021.720p.WEB.H264-TIMECUT", Size: 2.16, Seeder: 1, Leecher: 1, Date: "2020-05-11 01:22:28"},
	{ID: 5, Title: "Lady.Vengeance.2005.KOREAN.2160p.BluRay.REMUX.HEVC.DTS-HD.MA.5.1-FGT", Size: 72.56, Seeder: 57, Leecher: 58, Date: "2021-08-20 17:16:48"},
}

func getTorrents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, torrents)
}

func getFromStructById(id int) Torrent {
	var torrent Torrent

	for _, elem := range torrents {
		if elem.ID == id {
			torrent = elem
		}
	}

	return torrent
}

func getIndexFromStructById(id int) int {
	var selectedIndex int

	for index, elem := range torrents {
		if elem.ID == id {
			selectedIndex = index
		}
	}

	return selectedIndex
}

func deleteFromStructById(id int) bool {
	var selectedIndex int = getIndexFromStructById(id)

	torrents = torrents[:selectedIndex+copy(torrents[selectedIndex:], torrents[selectedIndex+1:])]

	return true
}

func getFromStructByTitle(searchParam string) []Torrent {
	var searchResults = []Torrent{}

	for _, elem := range torrents {
		if strings.Contains(strings.ToLower(elem.Title), strings.ToLower(searchParam)) {
			searchResults = append(searchResults, elem)
		}
	}

	return searchResults
}

func searchTorrent(c *gin.Context) {
	var torrents Torrent

	if err := c.ShouldBindUri(&torrents); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not valid search parameter"})

		return
	}

	var data = getFromStructByTitle(torrents.Title)

	c.IndentedJSON(http.StatusOK, data)
}

func getTorrent(c *gin.Context) {
	var torrent Torrent

	if err := c.ShouldBindUri(&torrent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not a valid ID"})

		return
	}

	var data = getFromStructById(torrent.ID)

	if data.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not a valid ID"})

		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

func updateTorrent(c *gin.Context) {
	var updatedTorrent Torrent

	if err := c.BindJSON(&updatedTorrent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})

		return
	}

	var data = getFromStructById(updatedTorrent.ID)

	if data.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not a valid ID"})

		return
	}

	var selectedIndex int = getIndexFromStructById(updatedTorrent.ID)

	p := &torrents[selectedIndex]

	p.Date = time.Now().Format("2006-01-06 15:04:05")

	if updatedTorrent.Size != 0 {
		p.Size = updatedTorrent.Size
	}

	if updatedTorrent.Seeder != 0 {
		p.Seeder = updatedTorrent.Seeder
	}

	if updatedTorrent.Leecher != 0 {
		p.Leecher = updatedTorrent.Leecher
	}

	if updatedTorrent.Title != "" {
		p.Title = updatedTorrent.Title
	}

	c.IndentedJSON(http.StatusCreated, torrents)
}

func deleteTorrent(c *gin.Context) {
	var torrent Torrent

	if err := c.ShouldBindUri(&torrent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not a valid ID"})

		return
	}

	var data = getFromStructById(torrent.ID)

	if data.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not a valid ID"})

		return
	}

	deleteFromStructById(data.ID)

	c.IndentedJSON(http.StatusOK, torrents)
}

func uploadTorrent(c *gin.Context) {
	var newTorrent Torrent

	newTorrent.ID = len(torrents) + 1
	newTorrent.Date = time.Now().Format("2006-01-06 15:04:05")

	if err := c.BindJSON(&newTorrent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})

		return
	}

	torrents = append(torrents, newTorrent)
	c.IndentedJSON(http.StatusCreated, newTorrent)
}

func getDocs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"/":            map[string]string{"method": "GET", "text": "Get api methods"},
		"/getTorrents": map[string]string{"method": "GET", "text": "Get all torrents"},
		"/getTorrent/:id": map[string]string{
			"method":   "GET",
			"text":     "Get specific torrent by id",
			"required": "id",
		},
		"/searchTorrent/:title": map[string]string{
			"method":   "GET",
			"text":     "Search torrents by name",
			"required": "title",
		},
		"/uploadTorrent": map[string]string{
			"method": "POST",
			"text":   "Upload torrent with JSON payload",
		},
		"/deleteTorrent/:id": map[string]string{
			"method":   "DELETE",
			"text":     "Delete torrent by id",
			"required": "id",
		},
		"/updateTorrent": map[string]string{
			"method":   "PUT",
			"text":     "Update torrent with JSON payload",
			"required": "id",
		},
	})
}

func main() {
	router := gin.Default()

	router.GET("/", getDocs)
	router.GET("/getTorrents", getTorrents)
	router.GET("/getTorrent/:id", getTorrent)
	router.GET("/searchTorrent/:title", searchTorrent)
	router.POST("/uploadTorrent", uploadTorrent)
	router.DELETE("/deleteTorrent/:id", deleteTorrent)
	router.PUT("/updateTorrent", updateTorrent)

	router.Run("localhost:8080")
}
