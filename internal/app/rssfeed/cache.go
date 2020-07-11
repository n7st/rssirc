// Package rssfeed runs in parallel to rssbot. It polls an RSS feed URL and
// sends formatted output to the IRC bot.
//
// This file provides a simple, ordered cache with a limited length for RSS feed
// items.
package rssfeed

import "github.com/mmcdole/gofeed"

// ItemMap is a map containing items by their title.
type ItemMap map[string]*gofeed.Item

// Cache contains saved feed items we don't want to display again.
type Cache struct {
	Items  ItemMap
	Order  []string
	Length int
}

// NewCache creates a new, empty cache.
func NewCache(length int) *Cache {
	return &Cache{
		Items:  ItemMap(make(ItemMap)),
		Length: length,
	}
}

// Save saves items to the cache.
func (c *Cache) Save(item *gofeed.Item) {
	c.Order = append(c.Order, item.Title)
	c.Items[item.Title] = item

	// Remove old items from the cache on every save
	c.Clean()
}

// Exists checks whether an item exists in the cache.
func (c *Cache) Exists(title string) bool {
	return c.Items[title] != nil
}

// Clean removes old items from the cache.
func (c *Cache) Clean() {
	if len(c.Order) <= c.Length {
		return
	}

	newCacheOrder := c.Order[len(c.Order)-c.Length:]

	items := ItemMap(make(ItemMap))

	for _, key := range newCacheOrder {
		items[key] = c.Items[key]
	}

	c.Order = newCacheOrder
	c.Items = items
}
