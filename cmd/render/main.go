package main

import (
	"fmt"

	"flag"

	"github.com/BlooperDB/API"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/storage"
)

func main() {
	var postgresHost string
	var minioHost string

	flag.StringVar(&postgresHost, "postgres-host", "postgres", "sets the postgres host to connect to")
	flag.StringVar(&minioHost, "minio-host", "minio", "sets the minio host to connect to")
	flag.Parse()

	blooper.InitializeDB(postgresHost)

	blooper.InitializeStorage(minioHost)

	revisions := db.FindUnrenderedRevisions()

	for _, rev := range revisions {
		revision := storage.GetRevision(rev.ID)
		if revision != nil {
			fmt.Println("Rendering " + string(rev.ID))
			storage.RenderAndSaveAndUpdateBlueprint(*revision, rev)
		}
	}
}
