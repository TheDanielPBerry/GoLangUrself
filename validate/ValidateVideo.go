package validate

import (
	"strconv"
	"errors"
)

func VideoId(videoId string) error {
	_, err := strconv.ParseInt(videoId, 10, 64)
	if err != nil {
		return errors.New("Invalid Video ID")
	}
	return nil
}
