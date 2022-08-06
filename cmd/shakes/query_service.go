// Code generated. DO NOT EDIT.

package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	gapic "github.com/timburks/shakes/gapic"
)

var QueryConfig *viper.Viper
var QueryClient *gapic.QueryClient
var QuerySubCommands []string = []string{
	"list-word-counts",
}

func init() {
	rootCmd.AddCommand(QueryServiceCmd)

	QueryConfig = viper.New()
	QueryConfig.SetEnvPrefix("SHAKES_QUERY")
	QueryConfig.AutomaticEnv()

	QueryServiceCmd.PersistentFlags().Bool("insecure", false, "Make insecure client connection. Or use SHAKES_QUERY_INSECURE. Must be used with \"address\" option")
	QueryConfig.BindPFlag("insecure", QueryServiceCmd.PersistentFlags().Lookup("insecure"))
	QueryConfig.BindEnv("insecure")

	QueryServiceCmd.PersistentFlags().String("address", "", "Set API address used by client. Or use SHAKES_QUERY_ADDRESS.")
	QueryConfig.BindPFlag("address", QueryServiceCmd.PersistentFlags().Lookup("address"))
	QueryConfig.BindEnv("address")

	QueryServiceCmd.PersistentFlags().String("token", "", "Set Bearer token used by the client. Or use SHAKES_QUERY_TOKEN.")
	QueryConfig.BindPFlag("token", QueryServiceCmd.PersistentFlags().Lookup("token"))
	QueryConfig.BindEnv("token")

	QueryServiceCmd.PersistentFlags().String("api_key", "", "Set API Key used by the client. Or use SHAKES_QUERY_API_KEY.")
	QueryConfig.BindPFlag("api_key", QueryServiceCmd.PersistentFlags().Lookup("api_key"))
	QueryConfig.BindEnv("api_key")
}

var QueryServiceCmd = &cobra.Command{
	Use:       "query",
	Short:     "The Query service performs selected BigQuery...",
	Long:      "The Query service performs selected BigQuery queries.",
	ValidArgs: QuerySubCommands,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		var opts []option.ClientOption

		address := QueryConfig.GetString("address")
		if address != "" {
			opts = append(opts, option.WithEndpoint(address))
		}

		if QueryConfig.GetBool("insecure") {
			if address == "" {
				return fmt.Errorf("Missing address to use with insecure connection")
			}

			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}
			opts = append(opts, option.WithGRPCConn(conn))
		}

		if token := QueryConfig.GetString("token"); token != "" {
			opts = append(opts, option.WithTokenSource(oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: token,
					TokenType:   "Bearer",
				})))
		}

		if key := QueryConfig.GetString("api_key"); key != "" {
			opts = append(opts, option.WithAPIKey(key))
		}

		QueryClient, err = gapic.NewQueryClient(ctx, opts...)
		return
	},
}
