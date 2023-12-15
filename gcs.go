package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	speech "cloud.google.com/go/speech/apiv1p1beta1" //v1p1beta1
	"cloud.google.com/go/storage"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1" //v1p1beta1
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
	objectName string
}

// Get Public address, make sure the bocket's ACL is set to public-read.
func (c *ClientUploader) GetPulicAddress() string {
	if len(c.objectName) == 0 {
		return ""
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, c.uploadPath+c.objectName)
}

// Upload Image object
func (c *ClientUploader) UploadImage(file io.ReadCloser) error {
	c.objectName = buildFileName() + ".jpeg"
	return c.uploadFile(file)
}

// Upload video object
func (c *ClientUploader) UploadVideo(file io.ReadCloser) error {
	c.objectName = buildFileName() + ".mp4"
	return c.uploadFile(file)
}

// uploadFile uploads an object
func (c *ClientUploader) SpeachToText() (error, string) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	// Creates a client.
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return err, fmt.Sprintf("Failed to create client: %v", err)
	}
	// The path to the remote audio file to transcribe.
	fileURI := fmt.Sprintf("gs://%s/%s", c.bucketName, c.uploadPath+c.objectName)

	// Detects speech in the audio file.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_MP3,
			SampleRateHertz: 48000,
			LanguageCode:    "zh-TW",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: fileURI},
		},
	})
	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
		return err, fmt.Sprintf("failed to recognize: %v", err)
	}

	// Prints the results.
	var resultStr string
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			log.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
			resultStr = resultStr + alt.Transcript + " "
		}
	}

	return nil, resultStr
}

// uploadFile uploads an object
func (c *ClientUploader) uploadFile(file io.ReadCloser) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + c.objectName).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func buildFileName() string {
	return time.Now().Format("20060102150405")
}
