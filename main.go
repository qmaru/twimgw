package main

import (
	"embed"

	"twimgw/apps"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := apps.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "twimg",
		Width:         420,
		Height:        740,
		MinWidth:      420,
		MinHeight:     700,
		MaxWidth:      600,
		MaxHeight:     800,
		OnBeforeClose: app.BeforeClose,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
