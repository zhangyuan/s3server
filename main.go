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
	if len(os.Args) != 2 {
		log.Fatalf("invalid arguments.")
	}

	bucket := os.Args[1]

	if err := serve(bucket); err != nil {
		log.Fatalf("error occurs, %v", err)
	}
}

func serve(bucket string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	svc := s3.NewFromConfig(cfg)

	r := gin.Default()

	r.Use(func(ctx *gin.Context) {
		output, err := GetObjectByPath(svc, bucket, ctx.Request.URL.Path)

		if err != nil {
			var nsk *types.NoSuchKey

			if errors.As(err, &nsk) {
				ctx.AbortWithStatus(404)
			} else {
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

func GetObjectByPath(svc *s3.Client, bucket string, path string) (*s3.GetObjectOutput, error) {
	path = strings.TrimLeft(path, "/")

	if path == "" || strings.HasSuffix(path, "/") {
		path = fmt.Sprintf("%sindex.html", strings.TrimLeft(path, "/"))
		return GetObject(svc, bucket, path)
	}

	output, err := GetObject(svc, bucket, path)

	var nsk *types.NoSuchKey

	if errors.As(err, &nsk) {
		path = fmt.Sprintf("%s/index.html", path)
		return GetObject(svc, bucket, path)
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
