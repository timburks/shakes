package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/timburks/shakes/rpc"
	bq "google.golang.org/api/bigquery/v2"
)

const wordCountQuery = `SELECT * FROM bigquery-public-data.samples.shakespeare WHERE word = @word ORDER BY word_count DESC;`

func wordCountParameters(req *rpc.QueryWordCountRequest) []*bq.QueryParameter {
	return []*bq.QueryParameter{
		{
			Name:           "word",
			ParameterValue: &bq.QueryParameterValue{Value: req.Word},
			ParameterType:  &bq.QueryParameterType{Type: "STRING"},
		},
	}
}

func wordCountRow(schema []*bq.TableFieldSchema, row *bq.TableRow) *rpc.QueryWordCountRow {
	result := &rpc.QueryWordCountRow{}
	for i, f := range row.F {
		log.Printf("field %d %T %+v %s %s", i, f.V, f.V, schema[i].Name, schema[i].Type)
		switch schema[i].Name {
		case "word":
			result.Word = f.V.(string)
		case "word_count":
			v, _ := strconv.Atoi(f.V.(string))
			result.WordCount = int32(v)
		case "corpus":
			result.Corpus = f.V.(string)
		case "corpus_date":
			v, _ := strconv.Atoi(f.V.(string))
			result.CorpusDate = int32(v)
		}
	}
	return result
}

func (queryServer) QueryWordCount(ctx context.Context, req *rpc.QueryWordCountRequest) (*rpc.QueryWordCountResponse, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}

	bqService, err := bq.NewService(ctx)
	if err != nil {
		log.Fatalf("bigquery.NewService: %v", err)
	}
	jobsService := bq.NewJobsService(bqService)

	useLegacySql := false
	job := &bq.Job{
		Configuration: &bq.JobConfiguration{
			Query: &bq.JobConfigurationQuery{
				Query:           wordCountQuery,
				QueryParameters: wordCountParameters(req),
				UseLegacySql:    &useLegacySql,
			},
		},
	}
	insertCall := jobsService.Insert(projectID, job)
	job, err = insertCall.Do()
	if err != nil {
		log.Fatalf("error creating job: %v", err)
	}

	queryCall := jobsService.GetQueryResults(projectID, job.JobReference.JobId).Context(ctx)
	if req.PageSize != 0 {
		queryCall = queryCall.MaxResults(int64(req.PageSize))
	}
	if req.PageToken != "" {
		queryCall = queryCall.PageToken(req.PageToken)
	}
	response, err := queryCall.Do()
	if err != nil {
		log.Fatalf("error making query: %v", err)
	}

	rows := make([]*rpc.QueryWordCountRow, 0)
	for i := range response.Rows {
		rows = append(rows, wordCountRow(response.Schema.Fields, response.Rows[i]))
	}

	resp := &rpc.QueryWordCountResponse{
		Rows:          rows,
		NextPageToken: response.PageToken,
	}
	return resp, nil
}
