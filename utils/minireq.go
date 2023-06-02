package utils

import (
	"github.com/qmaru/minireq/v2"
)

// UserAgent Global UA
var UserAgent = " Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 Edg/108.0.1462.54"

var Minireq *minireq.HttpClient

// MiniHeaders Headers
type MiniHeaders = minireq.Headers

// MiniParams Params
type MiniParams = minireq.Params

// MiniJSONData application/json
type MiniJSONData = minireq.JSONData

// MiniFormData multipart/form-data
type MiniFormData = minireq.FormData

// MiniFormKV application/x-www-from-urlencoded
type MiniFormKV = minireq.FormKV

// MiniAuth HTTP Basic Auth
type MiniAuth = minireq.Auth

func init() {
	Minireq = minireq.NewClient()
}
