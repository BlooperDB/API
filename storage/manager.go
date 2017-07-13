package storage

import (
	"strconv"

	"strings"

	"fmt"
	"os"

	"net/http"

	"github.com/BlooperDB/API/utils"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/policy"
)

var PublicURL string

var BlueprintStringBucket = "blooper-blueprints"
var BlueprintRenderBucket = "blooper-blueprint-renders"

var client *minio.Client

func Initialize(minioClient *minio.Client, url string) {
	client = minioClient
	PublicURL = url

	MakeBucket(BlueprintStringBucket)
	MakeBucket(BlueprintRenderBucket)

	client.SetBucketPolicy(BlueprintStringBucket, "", policy.BucketPolicyReadOnly)
	client.SetBucketPolicy(BlueprintRenderBucket, "", policy.BucketPolicyReadOnly)
}

func MakeBucket(name string) {
	err := client.MakeBucket(name, "")
	if err != nil {
		exists, err := client.BucketExists(name)
		if err == nil && !exists {
			fmt.Println("What?")
			os.Exit(1)
		} else if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func SaveRevision(revisionId uint, blueprintString string) {
	reader := strings.NewReader(blueprintString)
	client.PutObject(BlueprintStringBucket, RevisionToString(revisionId), reader, "text/plain")
}

func RevisionToString(revisionId uint) string {
	return "revision-blueprint-" + strconv.FormatUint(uint64(revisionId), 10)
}

func RenderAndSaveBlueprint(blueprintString string) {
	sha265 := utils.SHA265(blueprintString)

	// Normal
	reader := strings.NewReader(blueprintString)
	resp, _ := http.Post(os.Getenv("RENDERER_URL")+"/", "text/plain", reader)
	client.PutObject(BlueprintRenderBucket, sha265+".png", resp.Body, "image/png")

	// Square
	reader = strings.NewReader(blueprintString)
	resp, _ = http.Post(os.Getenv("RENDERER_URL")+"/?square", "text/plain", reader)
	client.PutObject(BlueprintRenderBucket, sha265+"-square.png", resp.Body, "image/png")

	// Thumbnail
	reader = strings.NewReader(blueprintString)
	resp, _ = http.Post(os.Getenv("RENDERER_URL")+"/?squarethumb", "text/plain", reader)
	client.PutObject(BlueprintRenderBucket, sha265+"-thumbnail.png", resp.Body, "image/png")
}
