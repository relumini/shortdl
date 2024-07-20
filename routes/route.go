package routes

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/kkdai/youtube/v2"
	pb "github.com/relumini/shortdl/protos"
	syoutube "github.com/relumini/shortdl/services"
)

var (
	// serverAddr      = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	videoRegexpList = []*regexp.Regexp{
		regexp.MustCompile(`(?:v|embed|shorts|watch\?v)(?:=|/)([^"&?/=%]{11})`),
		regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
		regexp.MustCompile(`([^"&?/=%]{11})`),
	}
)

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
func InitRoute(router *gin.Engine, Client pb.DownloadShortClient) {
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
	router.GET("/tshort", func(ctx *gin.Context) {
		request := &pb.ParamsRequest{
			Url: ctx.Query("url"),
		}

		c, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Set a timeout for the request
		defer cancel()

		response, err := Client.DownTiktok(c, request)
		if err != nil {
			log.Fatalf("Failed to call DownTiktok: %v", err)
		}
		// req := &pb.ParamsRequest{Url: ctx.Query("url")}
		// c, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		// defer cancel()
		// metadata, err := grpcClient.DownTiktok(c, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Successfully downloaded TikTok",
			"data":    response.Status,
		})
	})
}
