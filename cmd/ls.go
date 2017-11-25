// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var bucketFlag string
var keyFlag string

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		bucket := aws.String(bucketFlag)
		key := aws.String(keyFlag)

		// Configure to use Minio Server
		s3Config := &aws.Config{
			Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
			Region:      aws.String(os.Getenv("AWS_REGION")),
		}
		newSession := session.New(s3Config)

		s3Client := s3.New(newSession)

		cparams := &s3.CreateBucketInput{
			Bucket: bucket, // Required
		}

		// Create a new bucket using the CreateBucket call.
		_, err := s3Client.CreateBucket(cparams)
		if err != nil {
			// Message from an error.
			fmt.Println(err.Error())
			return
		}

		// Upload a new object "testobject" with the string "Hello World!" to our "newbucket".
		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Body:   strings.NewReader("My Super File Contents"),
			Bucket: bucket,
			Key:    key,
		})
		if err != nil {
			fmt.Printf("Failed to upload data to %s/%s, %s\n", *bucket, *key, err.Error())
			return
		}
		fmt.Printf("Successfully created bucket %s and uploaded data with key %s\n", *bucket, *key)

		// Retrieve our "testobject" from our "newbucket" and store it locally in "testobject_local".
		file, err := os.Create("testobject_local")
		if err != nil {
			fmt.Println("Failed to create file", err)
			return
		}
		defer file.Close()

		downloader := s3manager.NewDownloader(newSession)
		numBytes, err := downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: bucket,
				Key:    key,
			})
		if err != nil {
			fmt.Println("Failed to download file", err)
			return
		}
		fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")

	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//lsCmd.PersistentFlags().String("bucket", "", "Bucket name")
	lsCmd.PersistentFlags().StringVar(&bucketFlag, "bucket", "", "Bucket name")
	lsCmd.PersistentFlags().StringVar(&keyFlag, "key", "", "Key name (the path inside the bucket)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
