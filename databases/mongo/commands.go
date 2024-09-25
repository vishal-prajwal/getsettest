package mongo

import (
	"context"
	"errors"
	"fmt"

	"bitbucket.org/junglee_games/getsetgo/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("not found db")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrForbidden             = errors.New("attempted action is not allowed")
)

type MongoDB struct {
	client *mongo.Client
}

func Open(cfg *Config, build string) (MongoDB, error) {
	logger.Info(context.Background(), fmt.Sprint(cfg))
	client, err := GetMongoClient(cfg)
	return MongoDB{client: client}, err
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func (db MongoDB) StatusCheck(ctx context.Context) error {
	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	_, err := db.client.ListDatabaseNames(ctx, bson.D{})
	return err
}

func (db MongoDB) Disconnect(ctx context.Context) {
	db.client.Disconnect(ctx)
}

func (db MongoDB) Upsert(ctx context.Context, dbName string, collection string, filter interface{}, data interface{}) (bool, error) {
	var database *mongo.Database
	// Set the upsert option to true
	upsert := true

	// Define the options for the update operation
	updateOptions := options.Update().SetUpsert(upsert)
	database = db.client.Database(dbName)
	_, err := database.Collection(collection).UpdateOne(ctx, filter, data, updateOptions)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db MongoDB) Insert(ctx context.Context, dbName string, collection string, data interface{}, usePrimaryDB bool) (string, error) {
	var database *mongo.Database
	if !usePrimaryDB {
		database = secondaryReadPrefDB(dbName, db.client)
	} else {
		database = db.client.Database(dbName)
	}
	result, err := database.Collection(collection).InsertOne(ctx, data)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db MongoDB) InsertMany(ctx context.Context, dbName string, collection string, data []interface{}, usePrimaryDB bool) ([]interface{}, error) {
	var database *mongo.Database
	if !usePrimaryDB {
		database = secondaryReadPrefDB(dbName, db.client)
	} else {
		database = db.client.Database(dbName)
	}
	result, err := database.Collection(collection).InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil
}

func (db MongoDB) Update(ctx context.Context, dbName string, collection string, filter interface{}, data interface{}, usePrimaryDB bool) (bool, error) {
	var database *mongo.Database
	if !usePrimaryDB {
		database = secondaryReadPrefDB(dbName, db.client)
	} else {
		database = db.client.Database(dbName)
	}
	result, err := database.Collection(collection).UpdateOne(ctx, filter, data)
	if err != nil {
		return false, err
	}
	if result.MatchedCount == 0 {
		return false, ErrNotFound
	}
	return result.ModifiedCount > 0, nil
}

func (db MongoDB) Find(ctx context.Context, dbName, collection string, filter interface{}, results interface{}, usePrimaryDB bool) error {
	logger.Debug(ctx, "%s %s %v", dbName, collection, filter)
	var database *mongo.Database
	if !usePrimaryDB {
		database = secondaryReadPrefDB(dbName, db.client)
	} else {
		database = db.client.Database(dbName)
	}
	cur, err := database.Collection(collection).Find(ctx, filter)
	if err != nil {
		return err
	}
	if err = cur.All(ctx, results); err != nil {
		return err
	}
	return nil
}

func (db MongoDB) FindWithOpts(ctx context.Context, dbName, collection string, filter interface{}, results interface{}, opts interface{}, usePrimaryDB bool) error {
	opts2 := opts.(*options.FindOptions)
	var database *mongo.Database
	if !usePrimaryDB {
		database = secondaryReadPrefDB(dbName, db.client)
	} else {
		database = db.client.Database(dbName)
	}
	cur, err := database.Collection(collection).Find(ctx, filter, opts2)
	if err != nil {
		return err
	}
	if err = cur.All(ctx, results); err != nil {
		return err
	}
	return nil
}

func (db MongoDB) GetAllDistinct(ctx context.Context, dbName, collection, fieldName string, filter interface{}, usePrimaryDB bool) ([]interface{}, error) {
	var database *mongo.Database
	if !usePrimaryDB {
		database = secondaryReadPrefDB(dbName, db.client)
	} else {
		database = db.client.Database(dbName)
	}
	results, err := database.Collection(collection).Distinct(ctx, fieldName, filter)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (db MongoDB) GetByIDs(ctx context.Context, dbName, collection string, ids []string, results interface{}, usePrimaryDB bool) error {
	filter, err := matchIDs(ids)
	if err != nil {
		return err
	}
	var database *mongo.Database
	if !usePrimaryDB {
		database = secondaryReadPrefDB(dbName, db.client)
	} else {
		database = db.client.Database(dbName)
	}
	cur, err := database.Collection(collection).Find(ctx, filter)
	if err != nil {
		return err
	}
	if err = cur.All(ctx, results); err != nil {
		return err
	}
	return nil
}

func (db MongoDB) UpdateByID(ctx context.Context, dbName, collection, id string, elem interface{}, usePrimaryDB bool) (bool, error) {
	filter, err := matchID(id)
	if err != nil {
		return false, err
	}
	update := bson.D{
		{Key: "$set", Value: elem},
	}
	var database *mongo.Database
	if !usePrimaryDB {
		database = secondaryReadPrefDB(dbName, db.client)
	} else {
		database = db.client.Database(dbName)
	}
	res, err := database.Collection(collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

func secondaryReadPrefDB(dbname string, client *mongo.Client) *mongo.Database {
	secondary := readpref.Secondary()
	dbOpts := options.Database().SetReadPreference(secondary)
	db := client.Database(dbname, dbOpts)
	return db
}

func matchIDs(ids []string) (bson.D, error) {
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		err := CheckObjectIDFromHex(id)
		if err != nil {
			return nil, ErrInvalidID
		}
		objectID, _ := primitive.ObjectIDFromHex(id)
		objectIDs = append(objectIDs, objectID)
	}

	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: objectIDs}}}}
	return filter, nil
}

func matchID(id string) (bson.D, error) {
	err := CheckObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objectID}}
	return filter, nil
}

func CheckObjectIDFromHex(id string) error {
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}
	return nil
}
