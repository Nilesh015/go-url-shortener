package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testStoreService = &StorageService{}

func init() {
	testStoreService = InitializeStore()
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testStoreService.redisClient != nil)
}

func TestInsertionAndRetrieval(t *testing.T) {
	initialLink := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortURL := "du_myURL"

	// Delete any existing entry for "du_myURL".
	err := DeleteShortUrlEntry(shortURL)
	assert.True(t, err)

	// Persist data mapping
	assert.True(t, SaveUrlMapping(shortURL, initialLink))

	// Retrieve initial URL
	retrievedUrl := RetrieveLongUrl(shortURL)

	assert.Equal(t, initialLink, retrievedUrl)
}
