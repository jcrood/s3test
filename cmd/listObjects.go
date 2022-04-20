package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"time"
)

const delimiter = "/"
const defaultObjectPrefix = "transient_BatchManagerCache/"
const defaultTimeout = 50
const defaultLimit = 1000

var listObjectsCmd = &cobra.Command{
	Use:   "listObjects",
	Short: "execute s3 listObjects",
	Long:  `Executes the problematic listObjects API request that times out.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("listObjects")

		s3key := viper.GetString("key")
		s3secret := viper.GetString("secret")
		s3endpoint := viper.GetString("endpoint")
		s3bucket := viper.GetString("bucket")
		requestTimeout := viper.GetInt("timeout")
		objectPrefix := viper.GetString("prefix")
		debug := viper.GetBool("debug")
		limit := viper.GetInt("limit")

		if s3key == "" || s3secret == "" || s3endpoint == "" || s3bucket == "" {
			log.Fatalf("missing params, abort!")
		}

		client, err := getClient(s3key, s3secret, s3endpoint, debug)
		if err != nil {
			log.Fatalf("failed to get s3 client: %v", err)
		}

		noFiles, err := listObjects(client, s3bucket, objectPrefix, time.Duration(requestTimeout)*time.Second, limit)
		if err != nil {
			log.Fatalf("failed to listObjects: %v", err)
		}

		log.Printf("done. found %d files\n", noFiles)
	},
}

func init() {
	rootCmd.AddCommand(listObjectsCmd)

	listObjectsCmd.Flags().IntP("timeout", "t", defaultTimeout, "request timeout in seconds")
	listObjectsCmd.Flags().StringP("prefix", "p", defaultObjectPrefix, "object prefix to use")
	listObjectsCmd.Flags().Bool("debug", false, "enable debug mode")
	listObjectsCmd.Flags().IntP("limit", "l", defaultLimit, "limit max-keys")
	_ = viper.BindPFlag("timeout", listObjectsCmd.Flags().Lookup("timeout"))
	_ = viper.BindPFlag("prefix", listObjectsCmd.Flags().Lookup("prefix"))
	_ = viper.BindPFlag("debug", listObjectsCmd.Flags().Lookup("debug"))
	_ = viper.BindPFlag("limit", listObjectsCmd.Flags().Lookup("limit"))
	viper.SetDefault("timeout", defaultTimeout)
	viper.SetDefault("prefix", defaultObjectPrefix)
	viper.SetDefault("limit", defaultLimit)
}

const defaultRegion = "eu-west-3"

func getClient(s3key, s3secret, endpoint string, debugMode bool) (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           endpoint,
			SigningRegion: defaultRegion,
		}, nil
	})

	creds := credentials.NewStaticCredentialsProvider(s3key, s3secret, "")
	configOptions := []func(o *config.LoadOptions) error{
		config.WithCredentialsProvider(creds),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRetryMaxAttempts(0),
	}
	if debugMode {
		configOptions = append(configOptions, config.WithClientLogMode(aws.LogRequest))
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		configOptions...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return client, nil
}

func listObjects(client *s3.Client, bucketName, prefix string, timeout time.Duration, limit int) (int, error) {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	out, err := client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket:       aws.String(bucketName),
		Prefix:       aws.String(prefix),
		Delimiter:    aws.String(delimiter),
		EncodingType: types.EncodingTypeUrl,
		MaxKeys:      int32(limit),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to list objects: %w", err)
	}

	return len(out.Contents), nil
}
