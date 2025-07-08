package shortener

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShortLinkGenerator(t *testing.T) {
	initialLink_1 := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortLink_1 := GenerateShortLink(initialLink_1)
	fmt.Printf("Short url for %s is %s\n", initialLink_1, shortLink_1)
	assert.Equal(t, len(shortLink_1), 8, "Short link should be 8 characters long")

	initialLink_2 := "https://www.eddywm.com/lets-build-a-url-shortener-in-go-with-redis-part-2-storage-layer/"
	shortLink_2 := GenerateShortLink(initialLink_2)
	fmt.Printf("Short url for %s is %s\n", initialLink_2, shortLink_2)
	assert.Equal(t, len(shortLink_2), 8, "Short link should be 8 characters long")

	initialLink_3 := "https://spectrum.ieee.org/automaton/robotics/home-robots/hello-robots-stretch-mobile-manipulator"
	shortLink_3 := GenerateShortLink(initialLink_3)
	fmt.Printf("Short url for %s is %s\n", initialLink_3, shortLink_3)
	assert.Equal(t, len(shortLink_3), 8, "Short link should be 8 characters long")
}
