// Collection of routes for the REST API
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Play the playlist on Sonos (Main Method from NFC Tags)
func playAlbum(c *gin.Context) {
	id := c.Param("id")
	v, _ := storage.getOne(id, false)
	if len(v.Links.SpotifyAlbumURI) > 0 {
		if e := sonosPlayer.PlaySpotifyAlbum(v.Links.SpotifyAlbumURI); e != nil {
			c.AbortWithError(500, e)
		}
	}
	// redirect to SPA app
	c.Redirect(http.StatusTemporaryRedirect, "/view/"+id)
}

// get full metadata for an ID
func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	a, e := storage.getOne(id, false)
	if e != nil {
		c.AbortWithError(http.StatusNotFound, e)
		return
	}

	c.IndentedJSON(http.StatusOK, a)
}

func createAlbum(ctx *gin.Context) {
	req := CreateRecordRequest{}
	if errA := ctx.BindJSON(&req); errA != nil {
		ctx.AbortWithError(400, errA)
		return
	}

	storage.Create(&req)
	ctx.IndentedJSON(200, req)
}

func searchAll(ctx *gin.Context) {
	r, _ := goldenRecord.CombinedSearch(ctx.Param("q"))
	ctx.IndentedJSON(http.StatusOK, r)
}

func searchDiscogs(ctx *gin.Context) {
	r, e := discogs.FindByQuery(ctx.Param("q"))
	if e != nil {
		ctx.AbortWithError(500, e)
		return
	}
	ctx.IndentedJSON(http.StatusOK, r)
}
