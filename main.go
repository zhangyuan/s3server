package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
)

func main() {
	var args Args
	if len(os.Args) == 1 {
		fmt.Println("Usage: s3server bucketName rootDirectory")
	} else if len(os.Args) == 2 {
		args.Bucket = os.Args[1]
		args.Prefix = ""
	} else if len(os.Args) == 3 {
		args.Bucket = os.Args[1]
		args.Prefix = strings.TrimLeft(os.Args[2], "/")
	} else {
		log.Fatalf("invalid arguments.")
	}

	if err := serve(&args); err != nil {
		log.Fatalf("error occurs, %v", err)
	}
}

type Args struct {
	Bucket string
	Prefix string
}

func serve(args *Args) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	svc := s3.NewFromConfig(cfg)

	r := gin.Default()

	r.Use(func(ctx *gin.Context) {
		output, err := GetObjectByPath(svc, args.Bucket, args.Prefix, ctx.Request.URL.Path)

		if err != nil {
			var nsk *types.NoSuchKey

			if errors.As(err, &nsk) {
				ctx.AbortWithStatus(404)
			} else {
				fmt.Printf("%v", err)
				ctx.AbortWithStatus(500)
			}
			return
		}

		b := make([]byte, 128)
		reader := output.Body
		writer := ctx.Writer

		for {
			n, err := reader.Read(b)

			writer.Write(b[:n])

			if err == io.EOF {
				writer.Flush()
				break
			}
		}

		ctx.AbortWithStatus(200)
	})

	r.Run(":5050")
	return nil
}

func GetObjectByPath(svc *s3.Client, bucket string, prefix string, path string) (*s3.GetObjectOutput, error) {
	objectKey := strings.Trim(strings.Join([]string{strings.TrimRight(prefix, "/"), strings.TrimLeft(path, "/")}, "/"), "/")

	if objectKey == "" {
		objectKey = fmt.Sprintf("%sindex.html", strings.TrimLeft(objectKey, "/"))
		return GetObject(svc, bucket, objectKey)
	}

	output, err := GetObject(svc, bucket, objectKey)

	var nsk *types.NoSuchKey

	if errors.As(err, &nsk) {
		objectKey = fmt.Sprintf("%s/index.html", objectKey)
		return GetObject(svc, bucket, objectKey)
	} else {
		return output, err
	}
}

func GetObject(svc *s3.Client, bucket string, key string) (*s3.GetObjectOutput, error) {
	return svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}
