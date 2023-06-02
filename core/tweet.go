package core

import (
	"errors"
	"fmt"
	"time"

	"twimgw/configs"
	"twimgw/utils"
)

// TwitterBasic twitter basic info
type TwitterBasic struct {
	Username     string
	MaxResults   string
	StartID      string
	Exclude      bool
	token        string
	commonHeader utils.MiniHeaders
}

// dateFormat RFC3339 to your layout
func dateFormat(layout string, t string) (string, error) {
	cnLocal, err := time.LoadLocation("PRC")
	if err != nil {
		return "", err
	}
	// time.RFC3339
	rawTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return "", err
	}
	uTime := rawTime.In(cnLocal).Format(layout)
	return uTime, nil
}

// SetToken set twitter token
func (tb *TwitterBasic) SetToken() error {
	secrets, err := configs.Secrets()
	if err != nil {
		return err
	}

	twitterKey := secrets["twitter_key"].(string)
	twitterSecret := secrets["twitter_secret"].(string)
	twitterToken := secrets["token"].(string)
	if twitterToken != "" {
		tb.token = twitterToken
		tb.commonHeader = utils.MiniHeaders{
			"Authorization": "Bearer " + twitterToken,
			"User-Agent":    "v2TwitterGolang",
		}
		return nil
	}

	if twitterKey == "" || twitterSecret == "" {
		return errors.New("please configure configs/secrets.json")
	}

	tokenURL := "https://api.twitter.com/oauth2/token"
	data := utils.MiniFormKV{"grant_type": "client_credentials"}
	auth := utils.MiniAuth{twitterKey, twitterSecret}

	res, err := utils.Minireq.Post(tokenURL, data, auth)
	if err != nil {
		return err
	}

	if res.Response.StatusCode != 200 {
		return errors.New("get token failed")
	}

	resRaw, err := res.RawJSON()
	if err != nil {
		return err
	}

	resJSON := resRaw.(map[string]any)
	if value, ok := resJSON["access_token"]; ok {
		token := value.(string)
		if token != "" {
			tb.token = twitterToken
			tb.commonHeader = utils.MiniHeaders{
				"Authorization": "Bearer " + twitterToken,
			}
			return nil
		}
		return errors.New("token is empty")
	}
	return errors.New("access_token not found")
}

// GetToken get twitter token
func (tb *TwitterBasic) GetToken() string {
	return tb.token
}

// GetUserID get user_id by username
func (tb *TwitterBasic) GetUserID() (string, error) {
	// https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-by-username-username
	userLookup := fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s", tb.Username)

	res, err := utils.Minireq.Get(userLookup, tb.commonHeader)
	if err != nil {
		return "", err
	}

	resRaw, err := res.RawJSON()
	if err != nil {
		return "", err
	}

	if resRaw != nil {
		data := resRaw.(map[string]any)
		if data, ok := data["data"]; ok {
			userData := data.(map[string]any)
			userID := userData["id"].(string)
			return userID, nil
		}
	}
	return "", errors.New("user_id not found")
}

// GetTimelines get user timelines by user_id
func (tb *TwitterBasic) GetTimelines(userID, nextToken string) ([]map[string]any, int, string, error) {
	// https://developer.twitter.com/en/docs/twitter-api/tweets/timelines/api-reference/get-users-id-tweets
	userTimeline := fmt.Sprintf("https://api.twitter.com/2/users/%s/tweets", userID)

	paras := utils.MiniParams{
		"max_results":  tb.MaxResults,
		"expansions":   "attachments.media_keys",
		"tweet.fields": "created_at",
		"media.fields": "type,url,height,width,variants",
	}

	if nextToken != "" {
		paras["pagination_token"] = nextToken
	}

	if tb.Exclude {
		paras["exclude"] = "retweets"
	}

	if tb.StartID != "" {
		paras["since_id"] = tb.StartID
	}

	res, err := utils.Minireq.Get(userTimeline, tb.commonHeader, paras)
	if err != nil {
		return nil, 0, "", err
	}

	resRaw, err := res.RawJSON()
	if err != nil {
		return nil, 0, "", err
	}

	var pageToken string
	var tweetCounts int
	userTweets := make([]map[string]any, 0)

	if resRaw != nil {
		data := resRaw.(map[string]any)
		tweetMediaData := make(map[string]string)

		if meta, ok := data["meta"]; ok {
			metaData := meta.(map[string]any)
			if nToken, ok := metaData["next_token"]; ok {
				pageToken = nToken.(string)
			}
		}

		if includes, ok := data["includes"]; ok {
			includesData := includes.(map[string]any)
			userMedia := includesData["media"].([]any)
			for _, uMedia := range userMedia {
				mediaData := uMedia.(map[string]any)
				mediaKey := mediaData["media_key"].(string)
				mediaType := mediaData["type"].(string)
				var mediaBody string
				if mediaType == "photo" {
					mediaBody = mediaData["url"].(string) + "?format=jpg&name=orig"
				} else if mediaType == "video" {
					var mediaMaxBitrate float64
					mediaVariants := mediaData["variants"].([]any)
					for _, medimediaVariant := range mediaVariants {
						mediaDetails := medimediaVariant.(map[string]any)
						mediaContentType := mediaDetails["content_type"].(string)
						if mediaContentType == "video/mp4" {
							mediaBitrate := mediaDetails["bit_rate"].(float64)
							if mediaBitrate > mediaMaxBitrate {
								mediaMaxBitrate = mediaBitrate
								mediaURL := mediaDetails["url"].(string)
								mediaBody = mediaURL
							}
						}
					}
				}
				tweetMediaData[mediaKey] = mediaBody
			}
		}

		if data, ok := data["data"]; ok {
			userData := data.([]any)
			for _, uData := range userData {
				tmp := make(map[string]any)
				tweetBody := uData.(map[string]any)

				tweetCreateAd, err := dateFormat("20060102150405", tweetBody["created_at"].(string))
				if err != nil {
					return nil, 0, "", err
				}

				if tweetAttachmentsData, ok := tweetBody["attachments"]; ok {
					tweetAttachments := tweetAttachmentsData.(map[string]any)
					tweetAttachmentsMediaKeys := tweetAttachments["media_keys"].([]any)
					var tweetMedia []string
					for _, mediaKey := range tweetAttachmentsMediaKeys {
						key := mediaKey.(string)
						mediaURL := tweetMediaData[key]
						tweetMedia = append(tweetMedia, mediaURL)
					}

					tmp["id"] = tweetBody["id"].(string)
					tmp["created_at"] = tweetCreateAd
					tmp["media"] = tweetMedia
					tweetCounts = tweetCounts + len(tweetMedia)
					userTweets = append(userTweets, tmp)
				}
			}
		}
	}
	return userTweets, tweetCounts, pageToken, nil
}
