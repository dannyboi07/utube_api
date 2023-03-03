package redis

func GetVideoUploadStatus(userId uint64, videoName string) (string, error) {
	key, _ := createCreateVideoReqKey(userId, videoName)

	var (
		result string
		err    error
	)
	result, err = redisClient.Get(ctx, key).Result()

	return result, err
}
