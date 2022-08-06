package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/timburks/shakes/pkg/shakespeare"
	"github.com/timburks/shakes/rpc"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type queryServer struct {
	rpc.UnimplementedQueryServer
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

	query := client.Query(`SELECT * FROM bigquery-public-data.samples.shakespeare WHERE word = @word ORDER BY word_count DESC;`)
	query.QueryConfig.Parameters = append(query.QueryConfig.Parameters,
		bigquery.QueryParameter{
			Name:  "word",
			Value: req.Word,
		})
	iter, err := query.Read(ctx)

	rows := make([]*rpc.QueryWordCountRow, 0)
	for {
		var row shakespeare.WordCountRow
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
		Rows: rows,
	}, nil
}

func newServer() *queryServer {
	s := &queryServer{}
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterQueryServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
