// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	"google.golang.org/api/iterator"

	"os"

	rpcpb "github.com/timburks/shakes/rpc"
)

var ListWordCountsInput rpcpb.ListWordCountsRequest

var ListWordCountsFromFile string

func init() {
	QueryServiceCmd.AddCommand(ListWordCountsCmd)

	ListWordCountsCmd.Flags().StringVar(&ListWordCountsInput.Word, "word", "", "The word to match.")

	ListWordCountsCmd.Flags().Int32Var(&ListWordCountsInput.PageSize, "page_size", 10, "Default is 10. The maximum number of responses to return.")

	ListWordCountsCmd.Flags().StringVar(&ListWordCountsInput.PageToken, "page_token", "", "A token to use for paginated requests.")

	ListWordCountsCmd.Flags().StringVar(&ListWordCountsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var ListWordCountsCmd = &cobra.Command{
	Use:   "list-word-counts",
	Short: "ListWordCounts returns counts for a specified...",
	Long:  "ListWordCounts returns counts for a specified word.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if ListWordCountsFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if ListWordCountsFromFile != "" {
			in, err = os.Open(ListWordCountsFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &ListWordCountsInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Query", "ListWordCounts", &ListWordCountsInput)
		}
		iter := QueryClient.ListWordCounts(ctx, &ListWordCountsInput)

		// populate iterator with a page
		_, err = iter.Next()
		if err != nil && err != iterator.Done {
			return err
		}

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(iter.Response)

		return err
	},
}
