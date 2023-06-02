package apps

import (
	"context"

	"twimgw/core"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (b *App) BeforeClose(ctx context.Context) (prevent bool) {
	dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   "Quit?",
		Message: "Are you sure you want to quit?",
		Buttons: []string{"ok", "cancel"},
	})

	if err != nil {
		return false
	}

	if dialog == "ok" {
		return false
	}
	return true
}

func (a *App) TwitterCore(tweetData *core.TwitterDownload) map[string]any {
	mediaData, mediaCounts, err := tweetData.CollectTweets()
	if err != nil {
		return JsonData(0, "collect error: "+err.Error(), []any{})
	}

	if mediaCounts == 0 {
		return JsonData(0, "media not found", []any{})
	}

	output, err := tweetData.Download(mediaData, tweetData.Socks5)
	if err != nil {
		return JsonData(0, "download error: "+err.Error(), []any{})
	}
	return JsonData(1, "twitter data", output)
}
