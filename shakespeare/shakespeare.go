package shakespeare

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
)

type WordCountRow struct {
	Corpus     string `bigquery:"corpus"`
	CorpusDate int64  `bigquery:"corpus_date"`
	Word       string `bigquery:"word"`
	WordCount  int64  `bigquery:"word_count"`
}

type QueryWordCountRequest struct {
	Word string
}

func QueryWordCount(ctx context.Context, req *QueryWordCountRequest) (*bigquery.RowIterator, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	query := client.Query(
		`SELECT
			*
		FROM ` + "`bigquery-public-data.samples.shakespeare`" + `
		WHERE word = 'Caesar'
		ORDER BY word_count DESC
		LIMIT 3;`)
	return query.Read(ctx)

}
