package clipstore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// presignTTL is how long a generated audio URL stays valid.
//
// Short on purpose. The site sits behind a password, and a presigned URL is the one thing
// that leaves that protection — anyone holding it can fetch the clip without logging in.
// Long enough to start playing and to retry, short enough that a copied URL is not a
// lasting hole. Each request mints a fresh one, so nothing breaks when it lapses.
const presignTTL = 15 * time.Minute

// R2 serves clips from a Cloudflare R2 bucket.
//
// Requests are answered with a redirect to a presigned URL rather than by streaming the
// bytes through this process. That is the whole point of using object storage: audio is
// the only bandwidth-heavy thing the site serves, and proxying it would put every byte
// back through the container it was moved out of.
type R2 struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
}

// R2Config is what the Cloudflare dashboard gives you when creating a bucket and an
// Object Read & Write API token.
type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
}

// Configured reports whether enough is set to talk to R2. Anything missing means the
// local store is used instead, which is what a fresh checkout wants.
func (c R2Config) Configured() bool {
	return c.AccountID != "" && c.AccessKeyID != "" && c.SecretAccessKey != "" && c.Bucket != ""
}

// NewR2 builds a store against the bucket described by cfg.
func NewR2(cfg R2Config) *R2 {
	client := s3.New(s3.Options{
		// R2 has no regions in the AWS sense, but SigV4 requires one in the signature.
		// "auto" is what Cloudflare documents for its S3-compatible endpoint.
		Region:       "auto",
		BaseEndpoint: aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)),
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID, cfg.SecretAccessKey, "",
		),
	})

	return &R2{client: client, presign: s3.NewPresignClient(client), bucket: cfg.Bucket}
}

func (r *R2) Describe() string { return "Cloudflare R2 bucket " + r.bucket }

// Serve redirects the caller to a short-lived URL for the object.
//
// downloadName travels as the presigned request's response-content-disposition parameter
// rather than as a header on this response, because a redirect hands the browser off to
// R2 and any header set here is discarded. Signing it in is what keeps a saved file named
// after the clip's display name instead of its storage filename.
func (r *R2) Serve(w http.ResponseWriter, req *http.Request, filename, downloadName string) {
	in := &s3.GetObjectInput{Bucket: aws.String(r.bucket), Key: aws.String(filename)}
	if downloadName != "" {
		in.ResponseContentDisposition = aws.String(downloadName)
	}

	signed, err := r.presign.PresignGetObject(req.Context(), in, s3.WithPresignExpires(presignTTL))
	if err != nil {
		http.Error(w, "failed to prepare audio url", http.StatusInternalServerError)
		return
	}

	// 302 rather than 307: this is a GET being pointed elsewhere, and the URL is
	// per-request and expiring, so it must never be cached as a permanent location.
	w.Header().Set("Cache-Control", "no-store")
	http.Redirect(w, req, signed.URL, http.StatusFound)
}

func (r *R2) Put(ctx context.Context, filename string, data io.Reader, size int64) error {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(r.bucket),
		Key:           aws.String(filename),
		Body:          data,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String("audio/mpeg"),
	})
	if err != nil {
		return fmt.Errorf("upload %s to R2: %w", filename, err)
	}
	return nil
}

func (r *R2) List(ctx context.Context) ([]string, error) {
	var names []string
	pages := s3.NewListObjectsV2Paginator(r.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
	})
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list R2 bucket: %w", err)
		}
		for _, obj := range page.Contents {
			if obj.Key != nil {
				names = append(names, *obj.Key)
			}
		}
	}
	sort.Strings(names)
	return names, nil
}

func (r *R2) Delete(ctx context.Context, filename string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return fmt.Errorf("delete %s from R2: %w", filename, err)
	}
	return nil
}
