package main

import (
	"bytes"
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
)

type AwsS3 struct {
	sess *session.Session
	svc *s3.S3
	bucketName string
}

var AwsS3DefaultRegion = "us-west-2"

func ReadAwsS3ConnectionData() (string, string) {
	awsS3Region := flag.String("awsS3Region", os.Getenv("AwsS3Region"), "AWS S3 region")
	awsS3Bucket := flag.String("awsS3Bucket", os.Getenv("AwsS3Bucket"), "AWS S3 bucket name")

	flag.Parse()

	if len(*awsS3Region) == 0 {
		*awsS3Region = AwsS3DefaultRegion
		log.Printf("AWS S3 region not specified, using default region \"%s\"\n", AwsS3DefaultRegion)
	}

	if len(*awsS3Bucket) == 0 {
		log.Fatal("AWS S3 bucket name not specified")
	}

	return *awsS3Region, *awsS3Bucket
}

func ConnectToAwsS3(region, bucket string) *AwsS3 {
	sess, err := session.NewSession(aws.NewConfig().WithRegion(region))
	if err != nil {
		log.Fatalf("unable to start new AWS session: %v", err)
	}

	svc := s3.New(sess)

	log.Println("successfully connected to AWS S3")
	return &AwsS3{
		sess: sess,
		svc: svc,
		bucketName: bucket,
	}
}

func (awsS3 AwsS3) UploadClip(key string, file []byte) error {
	uploader := s3manager.NewUploader(awsS3.sess)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:	aws.String(awsS3.bucketName),
		Key:	aws.String(key),
		ACL:	aws.String("private"),
		Body: bytes.NewReader(file),
		ContentDisposition:	aws.String("attachment"),
	})

	if err != nil {
		return err
	}

	return err
}

func (awsS3 AwsS3) DownloadClip(key string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	dowloader := s3manager.NewDownloader(awsS3.sess)

	_, err := dowloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(awsS3.bucketName),
		Key:	aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (awsS3 AwsS3) DeleteClip(key string) error {
	_, err := awsS3.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(awsS3.bucketName),
		Key:	aws.String(key),
	})
	if err != nil {
		return err
	}

	err = awsS3.svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket:	aws.String(awsS3.bucketName),
		Key:	aws.String(key),
	})

	return err
}
