package redis

import "fmt"

func ChannelID(videoID int64) string {
	return fmt.Sprintf("comment:%d", videoID)
}
