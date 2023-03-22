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

func (es ES) SearchQuery(ctx context.Context, index string, query map[string]interface{}, resGenerator func(req io.ReadCloser) error) error {
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
		es.es.Search.WithSize(100),
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
