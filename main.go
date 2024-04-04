package main

import (
	delugeclient "github.com/gdm85/go-libdeluge"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

var (
	Hostname string
	Port     int
	Username string
	Password string
	KeepFor  time.Duration
	RunEvery time.Duration
	DryRun   bool
	err      error
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)

	Hostname = os.Getenv("DELUGE_HOSTNAME")
	if Hostname == "" {
		logrus.Warning("No hostname provided, using localhost")
		Hostname = "localhost"
	}

	Port, err = strconv.Atoi(os.Getenv("DELUGE_PORT"))
	if err != nil || Port < 1 || Port > 65535 {
		logrus.Warning("Invalid port provided, using default port 58846")
		Port = 58846
	}

	Username = os.Getenv("DELUGE_USERNAME")
	if Username == "" {
		logrus.Warning("No username provided, using default username localclient")
		Username = "localclient"
	}

	Password = os.Getenv("DELUGE_PASSWORD")
	if Password == "" {
		logrus.Fatal("No password provided")
	}

	KeepFor, err = time.ParseDuration(os.Getenv("KEEP_FOR"))
	if err != nil {
		logrus.Warning("Invalid duration provided, using default duration 720h")
		KeepFor, _ = time.ParseDuration("720h")
	}

	RunEvery, err = time.ParseDuration(os.Getenv("RUN_EVERY"))
	if err != nil {
		logrus.Warning("Invalid duration provided, using default duration 24h")
		RunEvery, _ = time.ParseDuration("24h")
	}

	DryRun, err = strconv.ParseBool(os.Getenv("DRY_RUN"))
	if err != nil {
		DryRun = false
	}
}

func main() {
	ticker := time.NewTicker(RunEvery)

	go func() {
		for range ticker.C {
			CleanTorrents(Hostname, uint(Port), Username, Password, KeepFor, DryRun)
		}
	}()

	select {}
}

func CleanTorrents(hostname string, port uint, username string, password string, keepFor time.Duration, dryRun bool) {
	logrus.Infof("Running cleanup at %s", time.Now().Format(time.RFC3339))

	if dryRun {
		logrus.Info("Running in dry-run mode")
	}

	var settings = delugeclient.Settings{
		Hostname: hostname,
		Port:     port,
		Login:    username,
		Password: password,
	}
	var deluge = delugeclient.NewV2(settings)

	err = deluge.Connect()
	if err != nil {
		logrus.Error(err)
		return
	}
	defer deluge.Close()

	torrents, err := deluge.TorrentsStatus(delugeclient.StateUnspecified, nil)
	if err != nil {
		logrus.Error(err)
	}

	var torrentsToDelete []string

	for id, torrent := range torrents {
		var completedDate = time.Unix(torrent.CompletedTime, 0)
		if torrent.IsSeed && torrent.IsFinished && time.Now().Sub(completedDate) > keepFor {
			logrus.Infof("Torrent: %s was completed on %s and is ready to be deleted", torrent.Name, completedDate)
			torrentsToDelete = append(torrentsToDelete, id)
		}
	}

	logrus.Infof("Found %d torrents to delete", len(torrentsToDelete))

	if !dryRun {
		failedTorrents, err := deluge.RemoveTorrents(torrentsToDelete, true)
		if err != nil {
			logrus.Error("Failed to remove torrents with ids: ", failedTorrents, " error: ", err)
		}
	}

}
