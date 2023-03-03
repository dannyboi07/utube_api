package redis

import (
	"utube/utils"
)

func SetCreateVideoReq(createVideoReqCached CreateVideoCached) error {
	key, duration := createCreateVideoReqKey(createVideoReqCached.UserId, createVideoReqCached.Name)

	var err error
	if err = redisClient.Set(ctx, key, createVideoReqCached, duration).Err(); err != nil {
		utils.Log.Println("Failed to cache create video request for userId:", createVideoReqCached.UserId)
	}

	return err
}
