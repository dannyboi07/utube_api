package schema

import (
	"regexp"
	"strings"
)

var maxVideoFileSize uint = 2 * 1000 * 1000 // 2MB
var validVideoTypeRegEx *regexp.Regexp = regexp.MustCompile("^video/mp4$")

type CreateVideoRequest struct {
	Name *string `json:"name"`
	Size *uint   `json:"size"`
	// video/(mp4|webp|...)
	Type *string `json:"type"`
}

func (c *CreateVideoRequest) Validate() (bool, string) {
	if notValid := c.Name == nil || len(strings.Trim(*c.Name, " ")) < 3; notValid {
		return notValid, "Missing or empty file name"
	} else if notValid := c.Size == nil; notValid {
		return notValid, "Missing file size"
	} else if notValid := *c.Size > maxVideoFileSize; notValid {
		return notValid, "Video file too large"
	} else if notValid := c.Type == nil; notValid {
		return notValid, "Missing file type"
	} else if notValid := !validVideoTypeRegEx.MatchString(*c.Type); notValid {
		return notValid, "Invalid file type"
	}

	return false, ""
}

type CreateVideoResponse struct {
	Name string `json:"name"`
	Size uint   `json:"size"`
	Type string `json:"type"`
}

type UploadVideoChunkRequest struct {
	Name *string `json:"name"`
	// ChunkSize *uint            `json:"chunk_size"`
	Data *string `json:"data"`
}
