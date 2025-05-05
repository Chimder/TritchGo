package db

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"tritchgo/internal/repository"

	"github.com/elastic/go-elasticsearch/v9"
)

type ElastiDB struct {
	client *elasticsearch.Client
}

func NewElasticDB() ElastiDB {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return ElastiDB{client: es}
}

func (el *ElastiDB) GetClient() *elasticsearch.Client {
	return el.client
}
func (el *ElastiDB) Mapping() error {
	indexName := "stream_stats"

	existsRes, err := el.client.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("error checking if index exists: %w", err)
	}
	defer existsRes.Body.Close()

	if existsRes.StatusCode == 200 {
		log.Println("index exists")
		return nil
	}

	mapping := `
	{
	  "mappings": {
	    "properties": {
	      "user_id": { "type": "keyword" },
	      "title": { "type": "text" }
	    }
	  }
	}`

	createRes, err := el.client.Indices.Create(indexName, el.client.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))))
	if err != nil {
		return fmt.Errorf("cannot create index: %w", err)
	}
	defer createRes.Body.Close()

	log.Println("index created")
	return nil
}

func (el *ElastiDB) FromPostToElastic(repo repository.Repository) error {

	stats, err := repo.Stats.GetAllStats()
	if err != nil {
		return fmt.Errorf("err create index to elastic from postgres %w", err)
	}

	slog.Info("RESP POST", "data", stats)

	for _, stat := range stats {
		doc := fmt.Sprintf(`
		{
			"user_id": "%s",
			"title": "%s"
		}`,
			stat.UserID,
			stat.Title,
		)

		req := bytes.NewReader([]byte(doc))

		res, err := el.client.Index(
			"stream_stats",
			req,
			el.client.Index.WithDocumentID(stat.ID.String()),
			el.client.Index.WithRefresh("true"),
		)
		if err != nil {
			log.Printf("error indexing document ID %s: %s", stat.ID, err)
			continue
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("error response indexing document ID %s: %s", stat.ID, res.String())
		}
	}

	return nil
}
