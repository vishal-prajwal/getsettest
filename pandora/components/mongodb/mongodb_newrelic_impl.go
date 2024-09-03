package mongodb

import (
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	mgo "gopkg.in/mgo.v2"
)

type MongoDriverNewRelic struct {
	*MongoDriver
	txn newrelic.Transaction
}

type segmentDetails struct {
	Dbname      string
	Collection  string
	Operation   string
	QueryParams map[string]interface{}
}

func getDataStoreSegment(txn newrelic.Transaction, data segmentDetails) *newrelic.DatastoreSegment {
	ds := &newrelic.DatastoreSegment{
		StartTime:       txn.StartSegmentNow(),
		Product:         newrelic.DatastoreMongoDB,
		Collection:      data.Collection,
		Operation:       data.Operation,
		QueryParameters: data.QueryParams,
		DatabaseName:    data.Dbname,
	}
	return ds
}

// FindOne returns first matching item
func (this *MongoDriverNewRelic) FindOne(collection string, query map[string]interface{}) (interface{}, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:      this.MongoDriver.conn.Name,
		Collection:  collection,
		Operation:   "FindOne",
		QueryParams: query,
	})
	defer ds.End()
	return this.MongoDriver.FindOne(collection, query)
}

func (this *MongoDriverNewRelic) FindOnenLoad(coll string, q map[string]interface{}, payload interface{}) *MDBError {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:      this.MongoDriver.conn.Name,
		Collection:  coll,
		Operation:   "FindOnenLoad",
		QueryParams: q,
	})
	defer ds.End()

	return this.MongoDriver.FindOnenLoad(coll, q, payload)
}

func (this *MongoDriverNewRelic) FindSortnLoad(coll string, q map[string]interface{}, selField map[string]interface{}, order string, page int, limit int, payload interface{}) (int, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:      this.MongoDriver.conn.Name,
		Collection:  coll,
		Operation:   "FindSortnLoad",
		QueryParams: q,
	})
	defer ds.End()
	return this.MongoDriver.FindSortnLoad(coll, q, selField, order, page, limit, payload)
}

func (this *MongoDriverNewRelic) FindOneUsingSession(sess *MSession, collection string, query map[string]interface{}) (interface{}, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:      this.MongoDriver.conf.DbName,
		Collection:  collection,
		Operation:   "FindOneUsingSession",
		QueryParams: query,
	})
	defer ds.End()
	return this.MongoDriver.FindOneUsingSession(sess, collection, query)
}

// FindAll returns all matching items
func (this *MongoDriverNewRelic) FindAll(collection string, query map[string]interface{}) ([]interface{}, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:      this.MongoDriver.conn.Name,
		Collection:  collection,
		Operation:   "FindAll",
		QueryParams: query,
	})
	defer ds.End()

	return this.MongoDriver.FindAll(collection, query)
}

func (this *MongoDriverNewRelic) FindAllUsingSession(sess *MSession, collection string, query map[string]interface{}) ([]interface{}, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:      this.MongoDriver.conf.DbName,
		Collection:  collection,
		Operation:   "FindAllUsingSession",
		QueryParams: query,
	})
	defer ds.End()
	return this.MongoDriver.FindAllUsingSession(sess, collection, query)
}

// FindWithLimit returns all matching items with limit
func (this *MongoDriverNewRelic) FindWithLimit(collection string, query map[string]interface{}, skip int, limit int) ([]interface{}, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:      this.MongoDriver.conn.Name,
		Collection:  collection,
		Operation:   "FindWithLimit",
		QueryParams: query,
	})
	defer ds.End()
	return this.MongoDriver.FindWithLimit(collection, query, skip, limit)
}

// FindWithField used to find data and returns only given fields
func (this *MongoDriverNewRelic) FindWithField(collection string, query map[string]interface{}, field map[string]interface{}) ([]interface{}, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "FindWithField",
		QueryParams: map[string]interface{}{
			"query": query,
			"field": field,
		},
	})
	defer ds.End()
	return this.MongoDriver.FindWithField(collection, query, field)
}

func (this *MongoDriverNewRelic) GetNewID(collection, idName, counterName string) (interface{}, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "GetNewID",
		QueryParams: map[string]interface{}{
			"counter": idName,
		},
	})
	defer ds.End()
	return this.MongoDriver.GetNewID(collection, idName, counterName)
}

// Insert add one item
func (this *MongoDriverNewRelic) Insert(collection string, value interface{}) *MDBError {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "Insert",
		QueryParams: map[string]interface{}{
			"value": value,
		},
	})
	defer ds.End()
	return this.MongoDriver.Insert(collection, value)
}

// Insert add one or more items
func (this *MongoDriverNewRelic) BulkInsert(collection string, value []interface{}) *MDBError {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "BulkInsert",
		QueryParams: map[string]interface{}{
			"value": value,
		},
	})
	defer ds.End()
	return this.MongoDriver.BulkInsert(collection, value)
}

// Update modify existing item
func (this *MongoDriverNewRelic) Update(collection string, query map[string]interface{}, value interface{}) *MDBError {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "Update",
		QueryParams: map[string]interface{}{
			"query": query,
			"value": value,
		},
	})
	defer ds.End()
	return this.MongoDriver.Update(collection, query, value)
}

func (this *MongoDriverNewRelic) UpdateAll(collection string, query map[string]interface{}, value interface{}) (*mgo.ChangeInfo, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "UpdateAll",
		QueryParams: map[string]interface{}{
			"query": query,
			"value": value,
		},
	})
	defer ds.End()
	return this.MongoDriver.UpdateAll(collection, query, value)
}

// Update modify existing item or insert new item if does not exist
func (this *MongoDriverNewRelic) Upsert(collection string, query map[string]interface{}, value interface{}) *MDBError {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "Upsert",
		QueryParams: map[string]interface{}{
			"query": query,
			"value": value,
		},
	})
	defer ds.End()
	return this.MongoDriver.Upsert(collection, query, value)
}

// Upsert one or more items where each item will be query and value
func (this *MongoDriverNewRelic) BulkUpsert(collection string, query []map[string]interface{}, value []interface{}) (BulkResult, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "BulkUpsert",
		QueryParams: map[string]interface{}{
			"query": query,
			"value": value,
		},
	})
	defer ds.End()
	return this.MongoDriver.BulkUpsert(collection, query, value)
}

// Update one or more items where each item will be query and value
func (this *MongoDriverNewRelic) BulkUpdate(collection string, query []map[string]interface{}, value []interface{}) *MDBError {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "BulkUpdate",
		QueryParams: map[string]interface{}{
			"query": query,
			"value": value,
		},
	})
	defer ds.End()
	return this.MongoDriver.BulkUpdate(collection, query, value)
}

// Remove delete existing item
func (this *MongoDriverNewRelic) Remove(collection string, query map[string]interface{}) *MDBError {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:      this.MongoDriver.conn.Name,
		Collection:  collection,
		Operation:   "Remove",
		QueryParams: query,
	})
	defer ds.End()
	return this.MongoDriver.Remove(collection, query)
}

func (this *MongoDriverNewRelic) RemoveAll(collection string, query map[string]interface{}) (*mgo.ChangeInfo, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:      this.MongoDriver.conn.Name,
		Collection:  collection,
		Operation:   "RemoveAll",
		QueryParams: query,
	})
	defer ds.End()
	return this.MongoDriver.RemoveAll(collection, query)
}

// Count
func (this *MongoDriverNewRelic) Count(collection string, query map[string]interface{}) (int, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:      this.MongoDriver.conn.Name,
		Collection:  collection,
		Operation:   "Count",
		QueryParams: query,
	})
	defer ds.End()
	return this.MongoDriver.Count(collection, query)
}

// Find -- collection,selectedField,limit,skip
func (this *MongoDriverNewRelic) Find(collection string, query map[string]interface{},
	selectField map[string]interface{}, skip int, limit int) ([]interface{}, int, *MDBError) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "Find",
		QueryParams: map[string]interface{}{
			"query":       query,
			"selectField": selectField,
			"skip":        skip,
			"limit":       limit,
		},
	})
	defer ds.End()

	return this.MongoDriver.Find(collection, query, selectField, skip, limit)
}

// FindAll queries the mongo DB and returns the sorted results with limit
func (this *MongoDriverNewRelic) FindAndSort(collection string, query map[string]interface{},
	selectField map[string]interface{}, sortField string, skip int, limit int) ([]interface{}, int, *MDBError) {

	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "FindAndSort",
		QueryParams: map[string]interface{}{
			"query":       query,
			"selectField": selectField,
			"sortField":   sortField,
			"skip":        skip,
			"limit":       limit,
		},
	})
	defer ds.End()
	return this.MongoDriver.FindAndSort(collection, query, selectField, sortField, skip, limit)
}

// Pipe is used for aggregation
func (this *MongoDriverNewRelic) Pipe(collection string, query interface{}) (*mgo.Pipe, error) {
	ds := getDataStoreSegment(this.txn, segmentDetails{
		Dbname:     this.MongoDriver.conn.Name,
		Collection: collection,
		Operation:  "Pipe",
		QueryParams: map[string]interface{}{
			"query": query,
		},
	})
	defer ds.End()
	return this.MongoDriver.Pipe(collection, query)
}
