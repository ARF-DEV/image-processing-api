package googlecloudstorage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/storage"
	"github.com/ARF-DEV/image-processing-api/configs"
	"github.com/ARF-DEV/image-processing-api/model"
	"google.golang.org/api/iterator"
)

type GoogleCloudStorageRepoImpl struct {
	client *storage.Client
	config *configs.Config
}

func New(ctx context.Context, cfg *configs.Config) GoogleCloudStorageRepo {
	s, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("error when creating gcs client: ", err)
	}
	isExists := false
	itr := s.Buckets(ctx, cfg.GOOGLE_PROJECT_ID)
	for {
		attr, err := itr.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			log.Fatal("error when iterating gcs bucket: ", err)
		}
		if attr.Name == cfg.GCS_BUCKET_NAME {
			isExists = true
			break
		}
	}

	bucket := s.Bucket(cfg.GCS_BUCKET_NAME)
	if !isExists {
		err = bucket.Create(ctx, cfg.GOOGLE_PROJECT_ID, &storage.BucketAttrs{
			Name:     cfg.GCS_BUCKET_NAME,
			Location: "ASIA-SOUTHEAST1",
		})
		if err != nil {
			log.Fatal("error when creating gcs bucket: ", err)
		}
		fmt.Println("bucket created")
	}

	policy, err := bucket.IAM().V3().Policy(ctx)
	if err != nil {
		log.Fatal("error when getting IAM policy: ", err)
	}

	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    "roles/storage.objectViewer",
		Members: []string{iam.AllUsers},
	})

	if err := bucket.IAM().V3().SetPolicy(ctx, policy); err != nil {
		log.Fatal("error when setting IAM policy: ", err)
	}

	return &GoogleCloudStorageRepoImpl{
		client: s,
		config: cfg,
	}
}

func (r *GoogleCloudStorageRepoImpl) CreateBucket(ctx context.Context) error {
	bucket := r.client.Bucket(r.config.GCS_BUCKET_NAME)
	err := bucket.Create(ctx, r.config.GOOGLE_PROJECT_ID, &storage.BucketAttrs{
		Name:     r.config.GCS_BUCKET_NAME,
		Location: "ASIA-SOUTHEAST1",
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *GoogleCloudStorageRepoImpl) UploadImage(ctx context.Context, req model.UploadImageRequest) (string, error) {
	o := r.client.Bucket(r.config.GCS_BUCKET_NAME).Object(req.Name)
	o = o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)
	if _, err := io.Copy(wc, req.Reader); err != nil {
		return "", fmt.Errorf("error when uploading to bucket: %v", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("error when closing writer: %v", err)
	}

	publicUrl := fmt.Sprintf("/%s/%s", r.config.GCS_BUCKET_NAME, req.Name)
	return publicUrl, nil
}

func (r *GoogleCloudStorageRepoImpl) Close() {
	r.client.Close()
}
