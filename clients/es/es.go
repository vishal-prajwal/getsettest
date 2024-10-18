package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/junglee_games/getsetgo/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ES struct {
	es *elasticsearch.Client
}

type ElasticConfig interface {
	GetElasticURL(ctx context.Context) string
	GetElasticToken(ctx context.Context) string
	GetElasticBuild(ctx context.Context) string
}

func New(ctx context.Context, cfg ElasticConfig) (*ES, error) {

	esConf := elasticsearch.Config{
		Addresses:           []string{cfg.GetElasticURL(ctx)},
		CompressRequestBody: true,
		MaxRetries:          3,
	}
	if cfg.GetElasticBuild(ctx) == "production" {
		esConf.Header = http.Header(map[string][]string{
			"Authorization": {"Basic " + cfg.GetElasticToken(ctx)},
		})
	}
	es, err := elasticsearch.NewClient(esConf)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
		return &ES{}, err
	}
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return &ES{}, err
	}

	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error in esInfo response: %s", res.String())
	}
	r := make(map[string]interface{})
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the esInfo response body: %s", err)
	}
	// Print client and server version numbers.
	log.Println(es.Info())
	return &ES{es}, nil
}

func (es ES) Post(ctx context.Context, data io.Reader, index string) error {
	// Set up the request object.
	req := esapi.IndexRequest{
		Index:   index,
		Body:    data,
		Refresh: "true",
		Pretty:  true,
		Human:   true,
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), es.es)
	if err != nil {
		log.Printf("Error getting response from elasticsearch POST: %s\n", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("error while reading es response in POST ", err)
		}
		bodyString := string(bodyBytes)
		err = fmt.Errorf("[%s] Error indexing document %s ID=%s", res.Status(), bodyString, req.DocumentID)
		log.Println(err)
		return err
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
			return err
		}
	}
	return nil
}

func (es ES) Patch(ctx context.Context, data io.Reader, index, docID string) error {
	// Set up the request object.
	req := esapi.UpdateRequest{
		Index:      index,
		Body:       data,
		DocumentID: docID,
		Refresh:    "true",
		Pretty:     true,
		Human:      true,
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), es.es)
	if err != nil {
		log.Printf("Error getting response from elasticsearch PATCH	: %s\n", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%s", res.String(), req.DocumentID)
		return fmt.Errorf("[%s] Error indexing document ID=%s", res.Status(), req.DocumentID)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
			return err
		}
	}
	return nil
}

func (es ES) SearchQuery(ctx context.Context, index string, query map[string]interface{}, size int, resGenerator func(req io.ReadCloser) error) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("Error encoding search query: %s", err)
		return err
	}
	// Perform the search request.
	res, err := es.es.Search(
		es.es.Search.WithContext(context.Background()),
		es.es.Search.WithIndex(index),
		es.es.Search.WithBody(&buf),
		es.es.Search.WithTrackTotalHits(true),
		es.es.Search.WithPretty(),
		es.es.Search.WithSize(size),
	)
	if err != nil {
		log.Printf("Error getting response from elasticsearch GET: %s\n", err)
		return nil
	}

	if res.IsError() {
		log.Printf("[%s] Error searching query=%v", res.String(), query)
		return fmt.Errorf("[%s] Error searching query=%v", res.Status(), query)
	}
	return resGenerator(res.Body)
}

func (es ES) DELETE(ctx context.Context, docID string, index string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: docID,
		Refresh:    "true",
		Pretty:     true,
		Human:      true,
	}
	res, err := req.Do(context.Background(), es.es)
	if err != nil {
		log.Printf("Error getting response from elasticsearch DELETE: %s\n", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error deleting document ID=%s", res.String(), req.DocumentID)
		return fmt.Errorf("[%s] Error deleting document ID=%s", res.Status(), req.DocumentID)
	}
	return nil
}

func (es ES) GetBatch(index string, batchSize int, scrollDuration time.Duration) (*esapi.Response, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"size": batchSize,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logger.Error(context.Background(), "Error encoding search query: %v", err)
	}
	res, err := es.es.Search(
		es.es.Search.WithContext(context.Background()),
		es.es.Search.WithIndex(index),
		es.es.Search.WithBody(&buf),
		es.es.Search.WithScroll(scrollDuration),
		es.es.Search.WithTrackTotalHits(true),
		es.es.Search.WithPretty(),
	)
	if err != nil {
		logger.Error(context.Background(), "Error getting response from elasticsearch GET: %v\n", err)
	}

	// Return scrollID and first batch of documents
	return res, nil
}

func (es ES) GetBatchWithQuery(index string, batchSize int, scrollDuration time.Duration, query interface{}) (*esapi.Response, error) {
	var buf bytes.Buffer
	qu := map[string]interface{}{
		"size":  batchSize,
		"query": query,
	}
	if err := json.NewEncoder(&buf).Encode(qu); err != nil {
		logger.Error(context.Background(), "Error encoding search query: %v", err)
	}
	res, err := es.es.Search(
		es.es.Search.WithContext(context.Background()),
		es.es.Search.WithIndex(index),
		es.es.Search.WithBody(&buf),
		es.es.Search.WithScroll(scrollDuration),
		es.es.Search.WithTrackTotalHits(true),
		es.es.Search.WithPretty(),
	)
	if err != nil {
		logger.Error(context.Background(), "Error getting response from elasticsearch GET: %v\n", err)
	}

	// Return scrollID and first batch of documents
	return res, nil
}

func (es ES) CreatePIT(index string, keepAlive string) (string, error) {
	res, err := es.es.OpenPointInTime(
		[]string{index},
		keepAlive, // Pass keepAlive dynamically
	)
	if err != nil {
		logger.Error(context.Background(), "Error creating PIT: %v", err)
		return "", err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	pitID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get PIT ID")
	}
	return pitID, nil
}

func (es ES) GetBatchWithPITQuery(index string, pitID string, batchSize int, query interface{}, searchAfter []interface{}, keepAlive string) (*esapi.Response, error) {
	var buf bytes.Buffer

	qu := map[string]interface{}{
		"size":  batchSize,
		"query": query,
		"sort": []map[string]interface{}{
			{
				"userID": map[string]string{"order": "asc"},
			},
		},
		"pit": map[string]interface{}{
			"id":         pitID,
			"keep_alive": keepAlive, // Keep the PIT alive for 1 minute give 1m
		},
	}

	// Add `search_after` if it’s not the first request
	if searchAfter != nil {
		qu["search_after"] = searchAfter
	}

	if err := json.NewEncoder(&buf).Encode(qu); err != nil {
		logger.Error(context.Background(), "Error encoding search query: %v", err)
		return nil, err
	}

	res, err := es.es.Search(
		es.es.Search.WithContext(context.Background()),
		es.es.Search.WithBody(&buf),
		es.es.Search.WithTrackTotalHits(true),
		es.es.Search.WithPretty(),
	)
	if err != nil {
		logger.Error(context.Background(), "Error executing PIT search: %v", err)
		return nil, err
	}

	return res, nil
}

func (es ES) ClosePIT(pitID string) error {
	res, err := es.es.ClosePointInTime(
		es.es.ClosePointInTime.WithBody(strings.NewReader(fmt.Sprintf(`{"id":"%s"}`, pitID))),
	)
	if err != nil {
		logger.Error(context.Background(), "Error closing PIT: %v", err)
		return err
	}
	defer res.Body.Close()

	return nil
}

func (es ES) Scroll(scrollID string, scrollDuration time.Duration) (*esapi.Response, error) {
	res, err := es.es.Scroll(
		es.es.Scroll.WithContext(context.Background()),
		es.es.Scroll.WithScrollID(scrollID),
		es.es.Scroll.WithScroll(scrollDuration),
	)
	if err != nil {
		logger.Error(context.Background(), "Error during scrolling: %v", err)
		return nil, err
	}

	return res, nil
}

func (es ES) ClearScroll(scrollID string) error {
	res, err := es.es.ClearScroll(
		es.es.ClearScroll.WithScrollID(scrollID),
	)
	if err != nil {
		logger.Error(context.Background(), "Error clearing scroll: %v", err)
		return err
	}
	defer res.Body.Close()
	return nil
}

func (es ES) BulkUpdate(index string, updatedDoc bytes.Buffer) error {
	// Perform the bulk update
	res, err := es.es.Bulk(bytes.NewReader(updatedDoc.Bytes()), es.es.Bulk.WithIndex("your-index"))
	if err != nil {
		logger.Error(context.Background(), "Error during bulk update: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error(context.Background(), "Error during bulk update: %s", res.String())
		return fmt.Errorf("Error during bulk update: %s", res.String())
	} else {
		fmt.Println("Bulk update successful")
	}
	return nil
}
