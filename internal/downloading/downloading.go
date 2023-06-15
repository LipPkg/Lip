package downloading

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

type ProgressBarStyleType int

const (
	StyleDefault ProgressBarStyleType = iota
	StylePercentageOnly
	StyleNone
)

// DownloadFile downloads a file from a url and saves it to a local path.
// Note that if the style is not StyleNone, the progress bar will be shown on the terminal.
func DownloadFile(url string, filePath string, progressBarStyle ProgressBarStyleType) error {
	var err error

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("cannot create HTTP request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot download file (HTTP %s): %s", resp.Status, url)
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer file.Close()

	switch progressBarStyle {
	case StyleNone:
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return fmt.Errorf("cannot download file from %s: %w", url, err)
		}
		return nil

	case StylePercentageOnly:
		bar := progressbar.NewOptions64(
			resp.ContentLength,
			progressbar.OptionClearOnFinish(),
			progressbar.OptionSetElapsedTime(false),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionSetWidth(0),
		)
		_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
		if err != nil {
			return fmt.Errorf("cannot download file from %s: %w", url, err)
		}

		return nil

	case StyleDefault:
		bar := progressbar.NewOptions64(
			resp.ContentLength,
			progressbar.OptionClearOnFinish(),
			progressbar.OptionShowBytes(true),
			progressbar.OptionShowCount(),
		)
		_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
		if err != nil {
			return fmt.Errorf("cannot download file from %s: %w", url, err)
		}

		return nil
	}

	// Never reached.
	panic("unreachable")
}
