// config/storj.go
package config

import "storj.io/uplink"

type StorjConfig struct {
	AccessGrant string
	BucketName  string
}

func NewStorjConfig() *StorjConfig {
	return &StorjConfig{
		AccessGrant: "15M6fjnUPHBJnVt3KofVPkEHb8P9AbE48hP6135x3k9yFnHJkno7zkXCD44poevuBaGjYUjYH8fXrzyYdZHhDCizj9iryNh76bub3Lvdd9tq6ptyb5hn8UhP5VTtCF5RLXxsKC3CjyTBAMvGhsPwZSgpXrLTqZNhYhBjR9kySmPfQAfNaBz7wPkn4Pv5Em1vuGdSkcbt2YftKPxRHShgZLMVuWFtCwLqcY1aYy8oxGPRWfceUWw19Ge3vEepXH2TFswjy7o3hWZVxn9HAMx7dtppzHzZ769tJ",
		BucketName:  "insurancepolicy",
	}
}

func (c *StorjConfig) GetAccess() (*uplink.Access, error) {
	return uplink.ParseAccess(c.AccessGrant)
}
