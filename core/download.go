package core

import (
	"fmt"
	URL "net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"twimgw/utils"
)

type TaskManager struct {
	Control sync.WaitGroup
	Workers chan map[string]any
}

type TwitterDownload struct {
	Username    string `json:"username"`
	MaxResults  string `json:"max_results"`
	StartID     string `json:"start_id"`
	Exclude     bool   `json:"exclude"`
	Socks5      string `json:"socks5"`
	StoragePath string `json:"storage_path"`
}

var MediaStorageFolder string

// setStorageFolder set storage folder
func (td *TwitterDownload) setStorageFolder(storageRoot, storageSub string) error {
	if storageRoot == "" {
		folder, err := utils.FileSuite.RootPath("")
		if err != nil {
			return err
		}
		storageRoot = folder
	}
	now := time.Now().Format("20060102150405")
	mediaStorageName := filepath.Join(storageRoot, storageSub+"_"+now)
	mediaStoragePath, err := utils.FileSuite.Mkdir(mediaStorageName)
	if err != nil {
		return err
	}
	MediaStorageFolder = mediaStoragePath
	return nil
}

// setFilename set filename
func (td *TwitterDownload) setFilename(mediaDate, mediaID, mediaURL string) (string, error) {
	uri, _ := URL.Parse(mediaURL)
	allPath := uri.Path
	allPaths := strings.Split(allPath, "/")
	mediaName := allPaths[len(allPaths)-1]
	mediaFile := fmt.Sprintf("%s_%s_%s", mediaDate, mediaID, mediaName)
	return filepath.Join(MediaStorageFolder, mediaFile), nil
}

// getMedia download media
func (td *TwitterDownload) getMedia(tasker *TaskManager, media map[string]any) error {
	defer tasker.Control.Done()
	mID := media["id"].(string)
	mDate := media["created_at"].(string)
	mURLs := media["media"].([]string)

	fmt.Printf("  -- %s Find %d media\n", mID, len(mURLs))
	for _, url := range mURLs {
		res, err := utils.Minireq.Get(url)
		if err != nil {
			return err
		}
		rawData, err := res.RawData()
		if err != nil {
			return err
		}
		mFile, err := td.setFilename(mDate, mID, url)
		if err != nil {
			return err
		}
		utils.FileSuite.WriteFile(mFile, rawData)
	}
	fmt.Printf("  -- %s Finished\n", mID)
	<-tasker.Workers
	return nil
}

// CollectTweets collect all data
func (td *TwitterDownload) CollectTweets() ([]map[string]any, int, error) {
	twitter := new(TwitterBasic)
	twitter.Username = td.Username
	twitter.MaxResults = td.MaxResults
	twitter.Exclude = td.Exclude
	twitter.StartID = td.StartID

	if td.Socks5 != "" {
		utils.Minireq.Socks5Address = td.Socks5
	}

	fmt.Println("> Setting token")
	err := twitter.SetToken()
	if err != nil {
		return nil, 0, err
	}

	fmt.Println("> Getting User ID")
	userID, err := twitter.GetUserID()
	if err != nil {
		return nil, 0, err
	}
	fmt.Printf("  -- User ID: %s\n", userID)

	allTweets := make([]map[string]any, 0)
	allCounts := 0
	pageToken := ""

	fmt.Println("> Getting media")
	for {
		tweets, counts, nextToken, err := twitter.GetTimelines(userID, pageToken)
		if err != nil {
			return nil, 0, err
		}

		fmt.Printf("  -- Fetching Media Counts: %d\n", counts)
		pageToken = nextToken
		allCounts = allCounts + counts
		allTweets = append(allTweets, tweets...)
		if pageToken == "" {
			break
		}
	}
	return allTweets, allCounts, nil
}

// Download download core
func (td *TwitterDownload) Download(mediaData []map[string]any, socks5 string) (string, error) {
	err := td.setStorageFolder(td.StoragePath, td.Username)
	if err != nil {
		return "", err
	}
	taskCounts := len(mediaData)

	tasker := new(TaskManager)
	tasker.Control.Add(taskCounts)
	tasker.Workers = make(chan map[string]any, 4)

	for i := 0; i < taskCounts; i++ {
		task := mediaData[i]
		tasker.Workers <- task
		go td.getMedia(tasker, task)
	}

	tasker.Control.Wait()
	if td.Socks5 != "" {
		utils.Minireq.Socks5Address = ""
	}
	return MediaStorageFolder, nil
}
