package routes

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iawia002/lux/downloader"
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/extractors/tiktok"
	"github.com/iawia002/lux/request"
	"github.com/kkdai/youtube/v2"
	pb "github.com/relumini/shortdl/protos"
	syoutube "github.com/relumini/shortdl/services"
)

var (
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
	request.SetOptions(request.Options{
		Cookie: "tt_csrf_token=YmksDB6a-h4cT2fF7JpORI2O9UBMCWjsntIc; ttwid=1%7C0FVb9fFc-sjDG2UdJwdC1AirqYozQ0xfbAS4N72vN2Y%7C1713886256%7C78a9d83445b82b73ca8d4e0cf024ea6cdf1329b7f3866c826b0a69a300ebce46; ak_bmsc=51B1D53481A3A4E4D0CEFF2BCF622DA2~000000000000000000000000000000~YAAQ7uIsF6c4j+SOAQAAANmUCxfRGVXZ4D9xnO97l1yDw0OWyomnVkNY7IUKaggUja0kQzFQ+WG4xaxBcPt0AN0n26KeHXGGKgHYpHPUMUBHGHQGDtE4RLyy7U+LPbSJCqVaSDiPuzxHht0YUIbWogvrFmBfkP4ohcmjkZxWtEI9qQ4Whaobb2CFHGdKNt0zlVNBjJQ3uYRAvUe12zSBynQB18y6QhE8goneRkCEw9VIeft2pFIwNQ8tkWWEjDt6wHNaqeND7eASg5WLzYskWbTt6bPAOhSNRLJ38HZrOB5QNg+xxN5uuCSYmjMXCl8SkvQr91pInmOng+V898FLLBQtefs95whvbpfE0mKwBk5Cz2TkkHcUJa/IoC0CLmNqoEk3AtKxpw/J; tt_chain_token=46Xkv2ukMzyJ2e7XU7y0AQ==; bm_sv=A2E67B998DE8E6A4F1C2C02485467446~YAAQ7uIsF6g4j+SOAQAABdqUCxf1J/K4dYG0k7bbw2m5rFujdlSqMoCKDubu4R602nFvbY6zWC5puJczBv3IXwJJRpQxxR03wDCMVlKTCqjQvgDs8BoCuoNQxfY2fdS+F3bKut2lxXPQ2qctqz4kHBrgspJArHn/zu/IuKCIeSzmV4KcyxW6Zvw3/xMRA0MeHgyuHsTRBS+VrFk8Ju2NbJWWC8uSHbLCM/dhFT7/ktw8RE30r24XpQmhLpVTsUSC~1; tiktok_webapp_theme=light; msToken=ySXERzKCE0QUG0cCg6TWLw3wfEB-6kh6kAfuzhzjcQvmV1jBFloSgIsT9xk-QXFVdI99U1Fqb9mhUpIOldoDkjdZwskB8rvt66MHZaHnvBRZRtOKtTYsWT8osDyQXDVZWdPkvyE598h9; passport_csrf_token=1a47d95ebf68fc3648b0018ee75afc9f; passport_csrf_token_default=1a47d95ebf68fc3648b0018ee75afc9f; perf_feed_cache={%22expireTimestamp%22:1714057200000%2C%22itemIds%22:[%227346425092966206766%22%2C%227353812964207594795%22%2C%227343343741916171563%22]}; msToken=yWwG-ITrCnjJbx5ltBa9FTXdCImOJrl-wtQJSQH3afeEumWZcbo_qcrF6F7-NjYcrG6JVxtJiOU208REZeCSgXEZrrs5_65K741fQ7PSzCGOhz6vUyycq3Xvj4Mu-S0kJ6SqyltHnpJp",
	})
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
		url := ctx.Query("url")
		request := &pb.ParamsRequest{
			Url: url,
		}

		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// add handler here
		response, err := Client.DownTiktok(c, request)

		if err != nil {
			fmt.Printf("Failed to call DownTiktok: %v \n", err)
			data, err := tiktok.New().Extract(url, extractors.Options{})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			for _, v := range data {
				v.FillUpStreamsData()
			}
			fmt.Print(data)
			for _, item := range data {
				dt := downloader.New(downloader.Options{InfoOnly: true, OutputPath: "download"}).Download(item)
				if dt != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": dt,
					})
					return
				}
			}
			ctx.JSON(http.StatusOK, gin.H{
				"message": "downloaded",
			})
			return

		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Successfully downloaded TikTok",
			"data":    response.Status,
		})
	})
}
