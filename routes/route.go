package routes

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kkdai/youtube/v2"
	syoutube "github.com/relumini/shortdl/services"
)

var videoRegexpList = []*regexp.Regexp{
	regexp.MustCompile(`(?:v|embed|shorts|watch\?v)(?:=|/)([^"&?/=%]{11})`),
	regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
	regexp.MustCompile(`([^"&?/=%]{11})`),
}

func ExtractVideoID(videoID string) (string, error) {
	if strings.Contains(videoID, "youtu") || strings.ContainsAny(videoID, "\"?&/<%=") {
		for _, re := range videoRegexpList {
			if isMatch := re.MatchString(videoID); isMatch {
				subs := re.FindStringSubmatch(videoID)
				videoID = subs[1]
			}
		}
	}

	if strings.ContainsAny(videoID, "?&/<%=") {
		return "", youtube.ErrInvalidCharactersInVideoID
	}

	if len(videoID) < 10 {
		return "", youtube.ErrVideoIDMinLength
	}

	return videoID, nil
}
func InitRoute(router *gin.Engine) {
	router.GET("/yshort", func(ctx *gin.Context) {
		parseUrl, err := ExtractVideoID(ctx.Query("url"))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Error extracting video ID",
			})
		}
		metadata, err := syoutube.GetYoutubeShort(parseUrl)
		if err != nil {
			errStr := fmt.Sprintf("%v", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": errStr,
				"data":    metadata.Checksum,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Successfully downloaded youtube",
			"data":    metadata,
		})
	})
}
