package main

import (
	"context"
	"os"
	"strconv"

	"github.com/timburks/shakes/rpc"
	"golang.org/x/oauth2/google"
	bq "google.golang.org/api/bigquery/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const wordCountQuery = `SELECT * FROM bigquery-public-data.samples.shakespeare WHERE word = @word ORDER BY word_count DESC;`

func wordCountParameters(req *rpc.ListWordCountsRequest) []*bq.QueryParameter {
	return []*bq.QueryParameter{
		{
			Name:           "word",
			ParameterValue: &bq.QueryParameterValue{Value: req.Word},
			ParameterType:  &bq.QueryParameterType{Type: "STRING"},
		},
	}
}

func wordCountRow(schema []*bq.TableFieldSchema, row *bq.TableRow) (*rpc.WordCount, error) {
	result := &rpc.WordCount{}
	for i, f := range row.F {
		value, ok := f.V.(string)
		if !ok {
			return nil, status.Errorf(codes.Internal, "unexpected type in response: %T", f.V)
		}

		switch schema[i].Name {
		case "word":
			result.Word = value
		case "word_count":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "%s", err)
			}
			result.WordCount = int32(v)
		case "corpus":
			result.Corpus = value
		case "corpus_date":
			v, err := strconv.Atoi(value)
			result.CorpusDate = int32(v)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "%s", err)
			}
		}
	}
	return result, nil
}

func wordCountRows(response *bq.GetQueryResultsResponse) ([]*rpc.WordCount, error) {
	rows := make([]*rpc.WordCount, 0)
	for i := range response.Rows {
		row, err := wordCountRow(response.Schema.Fields, response.Rows[i])
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%s", err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (queryServer) ListWordCounts(ctx context.Context, req *rpc.ListWordCountsRequest) (*rpc.ListWordCountsResponse, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		credentials, err := google.FindDefaultCredentials(context.Background())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%s", err)
		}
		projectID = credentials.ProjectID
	}
	if projectID == "" {
		return nil, status.Errorf(codes.Internal, "GOOGLE_CLOUD_PROJECT environment variable must be set.")
	}
	bqService, err := bq.NewService(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
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
		return nil, status.Errorf(codes.Internal, "error creating job %s", err)
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
		return nil, status.Errorf(codes.Internal, "error making query %s", err)
	}
	rows, err := wordCountRows(response)
	resp := &rpc.ListWordCountsResponse{
		WordCounts:    rows,
		NextPageToken: response.PageToken,
	}
	return resp, nil
}
