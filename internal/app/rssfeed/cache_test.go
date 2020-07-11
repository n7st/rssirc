package rssfeed

import (
	"reflect"
	"testing"

	"github.com/mmcdole/gofeed"
)

var (
	first  = &gofeed.Item{Title: "First"}
	second = &gofeed.Item{Title: "Second"}
	third  = &gofeed.Item{Title: "Third"}
	fourth = &gofeed.Item{Title: "Fourth"}
	fifth  = &gofeed.Item{Title: "Fifth"}
	sixth  = &gofeed.Item{Title: "Sixth"}
)

func TestCacheClean(t *testing.T) {
	type (
		want []*gofeed.Item
		args []*gofeed.Item
	)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty cache",
			args: args([]*gofeed.Item{}),
			want: want([]*gofeed.Item{}),
		},
		{
			name: "single item should keep single item",
			args: args([]*gofeed.Item{first}),
			want: want([]*gofeed.Item{first}),
		},
		{
			name: "three items should keep three items",
			args: args([]*gofeed.Item{first, second, third}),
			want: want([]*gofeed.Item{first, second, third}),
		},
		{
			name: "four items should keep last three",
			args: args([]*gofeed.Item{first, second, third, fourth}),
			want: want([]*gofeed.Item{second, third, fourth}),
		},
		{
			name: "six items should keep last three",
			args: args([]*gofeed.Item{first, second, third, fourth, fifth, sixth}),
			want: want([]*gofeed.Item{fourth, fifth, sixth}),
		},
	}

	for _, tt := range tests {
		c := NewCache(3)

		t.Run(tt.name, func(t *testing.T) {
			expected := NewCache(3)

			for _, item := range tt.args {
				c.Save(item)
			}

			for _, item := range tt.want {
				// Bypass Save method to avoid popping items out of the
				// "expected" cache
				expected.Items[item.Title] = item
			}

			c.Clean()

			if !reflect.DeepEqual(c.Items, expected.Items) {
				t.Errorf("Clean() = %v, want %v", c.Items, expected.Items)
			}
		})
	}
}
