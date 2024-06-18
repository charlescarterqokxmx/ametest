import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// copyFile copies an object from one bucket to another.
func copyFile(w io.Writer, srcBucket, srcObject, dstBucket, dstObject string) error {
	// srcBucket := "bucket-1"
	// srcObject := "object-1"
	// dstBucket := "bucket-2"
	// dstObject := "object-2"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy is aborted if the
	// object's generation number does not match your precondition.
	attrs, err := src.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %v", err)
	}
	dst = dst.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Optional: set a metageneration-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy is aborted if the
	// object's metageneration does not match your precondition.
	dst = dst.If(storage.Conditions{MetagenerationMatch: attrs.Metageneration})

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %v", dstObject, srcObject, err)
	}
	fmt.Fprintf(w, "Blob %v in bucket %v copied to blob %v in bucket %v\n", srcObject, srcBucket, dstObject, dstBucket)
	return nil
}
  
