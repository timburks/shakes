package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type ShakespeareRow struct {
	Corpus     string `bigquery:"corpus"`
	CorpusDate int64  `bigquery:"corpus_date"`
	Word       string `bigquery:"word"`
	WordCount  int64  `bigquery:"word_count"`
}

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}

	ctx := context.Background()
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
	rows, err := query.Read(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var row ShakespeareRow
		err := rows.Next(&row)
		if err == iterator.Done {
			return
		}
		if err != nil {
			log.Fatal(fmt.Errorf("error iterating through results: %v", err))
		}
		fmt.Fprintf(os.Stdout, "%s (%d): %d\n", row.Corpus, row.CorpusDate, row.WordCount)
	}
}
