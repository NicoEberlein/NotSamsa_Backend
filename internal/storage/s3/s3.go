package s3

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectToS3() *minio.Client {
	endpoint := "s3:9000"
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("notsamsa", "notsamsapw", ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}

	if err = minioClient.MakeBucket(context.Background(), "notsamsa", minio.MakeBucketOptions{}); err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), "notsamsa")
		if !(errBucketExists == nil && exists) {
			panic(err)
		}
	}

	return minioClient
}
