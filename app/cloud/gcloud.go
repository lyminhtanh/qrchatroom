package cloud

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/revel/revel"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type GoodleCloud struct {
	cloudOnce sync.Once
	client *storage.Client
	bucket string
}

func (gcloud *GoodleCloud) Read(object string) ([]byte, error) {
	// [START download_file]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := gcloud.client.Bucket(gcloud.bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
	// [END download_file]
}

func (gcloud *GoodleCloud) MakePublic(object string) (string, error) {
	// [START public]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	acl := gcloud.client.Bucket(gcloud.bucket).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}
	// [END public]
	fileUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s", gcloud.bucket, object)
	return fileUrl, nil
}

func (gcloud *GoodleCloud) Delete(object string) error {
	// [START delete_file]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := gcloud.client.Bucket(gcloud.bucket).Object(object)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	// [END delete_file]
	return nil
}

func (gcloud *GoodleCloud) Write(object, filePath string) error {
	ctx := context.Background()
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := gcloud.client.Bucket(gcloud.bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END upload_file]
	return nil
}

// lazy initialize db
func GCloudLazyInit() Cloud{
	g := GoodleCloud{}
	g.cloudOnce.Do(func() {
		client, bucket, err := getGCloudInfo()
		if err != nil {
			panic(err)
		}
		g.client = client
		g.bucket = bucket
	})
	return &g
}


func getGCloudInfo() (client *storage.Client, bucket string, err error) {
	//projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	projectID, ok := revel.Config.String("google.projectid")
	if !ok {
		fmt.Println("projectID not found")
		return
	}
	fmt.Println("projectID")
	fmt.Println(projectID)

	bucket, bok := revel.Config.String("google.bucket")
	if !bok {
		fmt.Println("bucket not found")
	}

	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx, storage.ScopeFullControl)
	if err != nil {
		fmt.Println("credential not found")
		return client,  bucket, err
	}
	client, errc := storage.NewClient(ctx, option.WithCredentials(creds))
	if errc != nil {
		log.Fatal(errc)
		return client,  bucket, errc
	}
	return client,  bucket, nil
}
