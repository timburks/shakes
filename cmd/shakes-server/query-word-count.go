package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/timburks/shakes/rpc"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const wordCountQuery = `SELECT * FROM bigquery-public-data.samples.shakespeare WHERE word = @word ORDER BY word_count DESC;`

type wordCountRow struct {
	Corpus     string `bigquery:"corpus"`
	CorpusDate int64  `bigquery:"corpus_date"`
	Word       string `bigquery:"word"`
	WordCount  int64  `bigquery:"word_count"`
}

func (queryServer) QueryWordCount(ctx context.Context, req *rpc.QueryWordCountRequest) (*rpc.QueryWordCountResponse, error) {
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

	var job *bigquery.Job
	if req.PageToken == "" {
		query := client.Query(wordCountQuery)
		query.QueryConfig.Parameters = append(query.QueryConfig.Parameters,
			bigquery.QueryParameter{
				Name:  "word",
				Value: req.Word,
			})
		job, err = query.Run(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error running query: %v", err)
		}
	} else {
		job, err = client.JobFromID(ctx, req.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "error getting job: %v", err)
		}
	}
	jobID := job.ID()
	fmt.Printf("The job ID is %s\n", jobID)

	iter, err := job.Read(ctx)
	rows := make([]*rpc.QueryWordCountRow, 0)
	for i := int32(0); i < req.PageSize; i++ {
		var row wordCountRow
		err := iter.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error iterating through results: %v", err)
		}
		rows = append(rows, &rpc.QueryWordCountRow{
			Corpus:     row.Corpus,
			Word:       row.Word,
			WordCount:  int32(row.WordCount),
			CorpusDate: int32(row.CorpusDate),
		})
	}
	return &rpc.QueryWordCountResponse{
		Rows:          rows,
		NextPageToken: job.ID(),
	}, nil
}
