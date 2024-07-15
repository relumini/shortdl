package syoutube

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/kkdai/youtube/v2"
	"github.com/relumini/shortdl/database"
	handler_error "github.com/relumini/shortdl/handler"
	"github.com/relumini/shortdl/models"
	"github.com/relumini/shortdl/utils"
)

type Metadata struct {
	Description string
	Caption     youtube.CaptionTrack
	Title       string
	Transcript  string
}

func GetYoutubeShort(idv string) (Metadata, error) {
	client := youtube.Client{}
	var metadata Metadata
	var msgerr string
	var wg sync.WaitGroup

	videoChan := make(chan *youtube.Video)
	errChan := make(chan error)

	// Get video metadata
	go func() {
		defer close(videoChan)
		defer close(errChan)
		vid, err := client.GetVideo(idv)
		fmt.Print(err)
		if err != nil {

			errChan <- errors.New("cannot get video")
			return
		}
		videoChan <- vid
	}()
	vid, err := <-videoChan, <-errChan
	if err != nil {
		return Metadata{}, err
	}

	format := vid.Formats.WithAudioChannels()
	sfile := "download/" + utils.ComputeChecksum(strings.ReplaceAll(vid.Title, " ", "_")) + ".mp4"

	wg.Add(2)
	if vid.CaptionTracks == nil {
		return Metadata{}, errors.New(string(handler_error.ErrSensitiveContent))
	}
	// Get Transcript
	go func() {
		defer wg.Done()
		transcript, err := client.GetTranscript(vid, vid.CaptionTracks[0].LanguageCode)
		if err != nil {
			msgerr = "cannot get transcript\n"
			return
		}
		transcript_fil := regexp.MustCompile(`\d{1,2}:\d{2}\s*-\s*`).ReplaceAllString(transcript.String(), "")
		metadata.Transcript = transcript_fil
	}()
	if msgerr != "" {
		return Metadata{}, errors.New(msgerr)
	}
	checkSum := utils.ComputeChecksum(metadata.Transcript)

	var existingChecksum models.ChecksumData
	if err := database.DB.Where("checksum_value = ?", checkSum).First(&existingChecksum).Error; err == nil {
		msg := "already downloaded the video"

		return Metadata{}, errors.New(msg)
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return Metadata{}, err
		// }
		// if err != nil {

		// }
	}

	newChecksum := models.ChecksumData{ChecksumValue: checkSum}
	if err := database.DB.Create(&newChecksum).Error; err != nil {
		msg := "duplicated video"

		return Metadata{}, errors.New(msg)
	}

	go func() {
		defer wg.Done()
		stream, _, err := client.GetStream(vid, &format[0])
		if err != nil {
			msgerr += "cannot get stream\n"
			return
		}
		defer stream.Close()

		file, err := os.Create(sfile)
		if err != nil {
			msgerr += "cannot create file\n" + err.Error()
			return
		}
		defer file.Close()

		_, err = io.Copy(file, stream)
		if err != nil {
			msgerr += "cannot copy file\n"
			return
		}
	}()

	// Set metadata fields
	metadata.Description = vid.Description
	metadata.Caption = vid.CaptionTracks[0]
	metadata.Title = vid.Title

	// Wait for all goroutines to finish
	wg.Wait()

	if msgerr != "" {
		return Metadata{}, errors.New(msgerr)
	}

	return metadata, nil
}
