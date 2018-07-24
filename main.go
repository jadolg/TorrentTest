package main

import (
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"gopkg.in/cheggaaa/pb.v1"
)

func main() {
	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.Debug = false
	clientConfig.Seed = false
	clientConfig.DefaultStorage = storage.NewMMap("")

	client, err := torrent.NewClient(clientConfig)

	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}

	defer client.Close()

	if len(os.Args) < 2 {
		log.Fatal("no magnet link specified")
	}

	t, err := client.AddMagnet(os.Args[1])

	if err != nil {
		log.Fatalf("error adding magnet: %s", err)
	}

	<-t.GotInfo()

	log.Printf("Downloading %s", t.Info().Name)

	go func() {
		t.DownloadAll()
	}()

	bar := pb.StartNew(int(t.Info().TotalLength()))
	diff := int64(0)

	for t.BytesCompleted() != t.Info().TotalLength() {
		time.Sleep(time.Second)
		bar.Add64(t.BytesCompleted() - diff)
		diff = t.BytesCompleted()
	}
}
