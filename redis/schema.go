package redis

import (
	"encoding/json"
)

type CreateVideoCached struct {
	Name      string
	Size      uint
	Uploaded  uint
	UserId    uint64
	UserEmail string
}

func (c CreateVideoCached) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CreateVideoCached) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
