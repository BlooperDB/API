package storage

import (
	"strconv"

	"strings"

	"bytes"

	"github.com/minio/minio-go"
)

var BucketName = "blooper-blueprints"

var client *minio.Client

func Initialize(minioClient *minio.Client) {
	client = minioClient
}

func SaveRevision(revisionId uint, blueprintString string) {
	reader := strings.NewReader(blueprintString)
	client.PutObject(BucketName, RevisionToString(revisionId), reader, "text/plain")
}

func LoadRevision(revisionId uint) (string, error) {
	reader, err := client.GetObject(BucketName, RevisionToString(revisionId))

	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.String(), nil
}

func RevisionToString(revisionId uint) string {
	return "revision-blueprint-" + strconv.FormatUint(uint64(revisionId), 10)
}
