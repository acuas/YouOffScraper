package v1

import (
	"github.com/gofiber/fiber"
	. "gitlab.com/c0b/go-ordered-json"
)

func V1(api *fiber.Group) {
	_ = api.Group("/v1", func(c *fiber.Ctx) {
		c.JSON(NewOrderedMapFromKVPairs([]*KVPair{
			{Key: "ok", Value: 1},
			{Key: "data", Value: NewOrderedMapFromKVPairs([]*KVPair{
				{Key: "name", Value: "scrape"},
				{Key: "path", Value: "/scrape"},
				{Key: "info", Value: "Scrape an entire channel or video from YouTube"},
				{Key: "types", Value: NewOrderedMapFromKVPairs([]*KVPair{
					{Key: "channel", Value: fiber.Map{}},
					{Key: "video", Value: fiber.Map{}},
				})},
			})},
		}))
	})
}