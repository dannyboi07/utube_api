package redis

import (
	"strconv"
	"time"
)

func createCreateVideoReqKey(userId uint64, videoName string) (string, time.Duration) {
	return "create_video:" + strconv.FormatUint(userId, 10) + ":" + videoName, time.Hour * 24
}
