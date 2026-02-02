package upload

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/url"
	"time"
)

func Minio(file string, name string) *url.URL {
	endpoint := "115.190.154.22:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	filePath := file
	bucketName := "2308a"
	objectName := name

	ctx := context.Background()
	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	if errBucketExists != nil {
		log.Fatalln(errBucketExists)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Successfully created %s\n", bucketName)
	}
	contentType := "image/png"
	_, err = minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	//_, err = minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, time.Hour*24*7, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(presignedURL)
	fmt.Printf("Successfully uploaded %s to %s/%s\n", filePath, bucketName, objectName)
	return presignedURL
}

/*
// FileUploader.go MinIO example
package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx := context.Background()
	endpoint := "play.min.io"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called testbucket.
	bucketName := "testbucket"
	location := "us-east-1"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload the test file
	// Change the value of filePath if the file is in another location
	objectName := "testdata"
	filePath := "/tmp/testdata"
	contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}
*/
