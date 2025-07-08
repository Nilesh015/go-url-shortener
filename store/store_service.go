package store

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Define the struct wrapper around raw db and cache clients.
type StorageService struct {
	// Redis client for caching.
	redisClient *redis.Client

	// MongoDB client for persistent storage.
	mongoClient *mongo.Client
}

// Top level declarations for the storeService and context.
var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

// Define constants to be used.
const kRedisServerUrl = "localhost:6379"
const kMongoServerUrl = "mongodb://localhost:27017"
const kMongoDbName = "short_url_db"
const kMongoCollectionName = "urls"
const kShortUrlKey = "shortUrl"
const kLongUrlKey = "longUrl"

// Initializing the store service and return a store pointer
func InitializeStore() *StorageService {
	// Initialize the Redis client.
	redisClient := redis.NewClient(&redis.Options{
		Addr:     kRedisServerUrl,
		Password: "",
		DB:       0,
	})
	pong, redisErr := redisClient.Ping(ctx).Result()
	if redisErr != nil {
		panic(fmt.Sprintf("Error init Redis: %v\n", redisErr))
	}
	fmt.Printf("\nConnected to Redis: pong message = {%s}\n", pong)

	// Initialize the MongoDB client.
	mongoClient, mongoErr := mongo.Connect(options.Client().ApplyURI(kMongoServerUrl))
	if mongoErr != nil {
		panic(fmt.Sprintf("Error connecting to MongoDB server: %v\n", mongoErr))
	}
	mongoErr = mongoClient.Ping(ctx, nil)
	if mongoErr != nil {
		panic(fmt.Sprintf("Could not ping MongoDB server: %v\n", mongoErr))
	}
	fmt.Println("Connected to MongoDB server.")

	// Create the database and collection if they do not exist.
	db := mongoClient.Database(kMongoDbName)
	collection := db.Collection(kMongoCollectionName)
	if collection == nil {
		// Create the collection if it does not exist
		mongoErr = db.CreateCollection(ctx, kMongoCollectionName)
		if mongoErr != nil {
			panic(fmt.Sprintf("Failed to create collection: %v\n", mongoErr))
		}
		fmt.Printf("Collection %s created successfully\n", kMongoCollectionName)
	}

	// Index urls collection by shortUrl key sorted in ascending order.
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{kShortUrlKey, 1}},
		Options: options.Index().SetUnique(true),
	}

	var index_name string
	index_name, mongoErr = collection.Indexes().CreateOne(ctx, indexModel)
	if mongoErr != nil {
		panic(fmt.Sprintf("Failed to create index on key: %s for collection: %s due to error: %v\n", kShortUrlKey, kMongoCollectionName, mongoErr))
	}
	fmt.Printf("Index created in '%s' collection on key: %s\n", kMongoCollectionName, index_name)

	// Initialize the store service with Redis and MongoDB clients.
	storeService = &StorageService{}
	storeService.redisClient = redisClient
	storeService.mongoClient = mongoClient
	return storeService
}

// Save a mapping from short URL to original long URL in MongoDB and Redis cache.
func SaveUrlMapping(shortUrl string, longUrl string) bool {
	// Insert shortURL mapping into MongoDB.
	urls := storeService.mongoClient.Database(kMongoDbName).Collection(kMongoCollectionName)
	doc := map[string]interface{}{kShortUrlKey: shortUrl, kLongUrlKey: longUrl}
	_, mongoErr := urls.InsertOne(ctx, doc)
	if mongoErr != nil {
		fmt.Printf("Could not insert URL into MongoDB %s due to error: %v\n", shortUrl, mongoErr)
		return false
	}

	// Insert shortURL mapping into Redis cache.
	redisErr := storeService.redisClient.Set(ctx, shortUrl, longUrl, 0 /* No TTL set */).Err()
	if redisErr != nil {
		fmt.Printf("Failed saving k,v: { %s: %s, %s: %s} due to error: %v\n", kShortUrlKey, shortUrl, kLongUrlKey, longUrl, redisErr)
		return false
	}

	return true
}

// Given a short URL, retrieve the original long URL from Redis cache. If not found, it is retreived from MongoDB.
func RetrieveLongUrl(shortUrl string) string {
	redisResult, redisGetErr := storeService.redisClient.Get(ctx, shortUrl).Result()
	if redisGetErr != nil {
		// If not found in Redis, try to fetch from MongoDB.
		var mongoResult string
		filter := bson.D{{Key: kShortUrlKey, Value: shortUrl}}
		mongoErr := storeService.mongoClient.Database(kMongoDbName).Collection(kMongoCollectionName).FindOne(ctx, filter).Decode(&mongoResult)
		if mongoErr != nil {
			return ""
		}

		// If found in MongoDB, save it to Redis cache for future requests.
		redisSetErr := storeService.redisClient.Set(ctx, shortUrl, mongoResult, 0 /* No TTL set */).Err()
		if redisSetErr != nil {
			return ""
		}

		return mongoResult
	}
	return redisResult
}

// Utility method to delete a URL entry from both Redis cache and MongoDB.
// This is not triggered via user workflow and is defined solely to support internal cleanups.
func DeleteShortUrlEntry(shortUrl string) bool {
	// Delete from Redis cache.
	redisError := storeService.redisClient.Del(ctx, shortUrl).Err()
	if redisError != nil {
		fmt.Printf("Failed to delete key %s from Redis cache due to error: %v\n", shortUrl, redisError)
		return false
	}

	// Delete from MongoDB.
	urls := storeService.mongoClient.Database(kMongoDbName).Collection(kMongoCollectionName)
	filter := bson.D{{Key: kShortUrlKey, Value: shortUrl}}
	_, mongoErr := urls.DeleteOne(ctx, filter)
	if mongoErr != nil {
		fmt.Printf("Failed to delete key %s from MongoDB due to error: %v\n", shortUrl, mongoErr)
		return false
	}

	return true
}
