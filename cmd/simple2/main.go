package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/timburks/shakes/pkg/shakespeare"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()
	rows, err := shakespeare.QueryWordCount(ctx, &shakespeare.QueryWordCountRequest{
		Word: "Caesar",
	})
	if err != nil {
		log.Fatal(fmt.Errorf("error performing query: %v", err))
	}
	for {
		var row shakespeare.WordCountRow
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
