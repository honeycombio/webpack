package webpack

import (
	"errors"
	"html/template"
	"log"
	"strings"

	"github.com/honeycombio/webpack/helper"
	"github.com/honeycombio/webpack/reader"
)

// DevHost webpack-dev-server host:port
var DevHost = "localhost:3808"

// FsPath filesystem path to public webpack dir
var FsPath = "./public/webpack"

// WebPath http path to public webpack dir
var WebPath = "webpack"

// Plugin webpack plugin to use, can be stats or manifest
var Plugin = "deprecated-stats"

// IgnoreMissing ignore assets missing on manifest or fail on them
var IgnoreMissing = true

// Verbose error messages to console (even if error is ignored)
var Verbose = true

// provide an optional asset host to support a CDN
var AssetHost = ""

type Config struct {
	// DevHost webpack-dev-server host:port
	DevHost string
	// FsPath filesystem path to public webpack dir
	FsPath string
	// WebPath http path to public webpack dir
	WebPath string
	// Plugin webpack plugin to use, can be stats or manifest
	Plugin string
	// IgnoreMissing ignore assets missing on manifest or fail on them
	IgnoreMissing bool
	// Verbose - show more info
	Verbose bool
	// IsDev - true to use webpack-serve or webpack-dev-server, false to use filesystem and manifest.json
	IsDev bool
	// AssetHost - optionally provide an asset host to support a CDN
	AssetHost string
}

type AssetTagHelperFunc func(string) (template.HTML, error)
type AssetURLHelperFunc func(string) (string, error)

var AssetTagHelper AssetTagHelperFunc
var AssetURLHelper AssetURLHelperFunc

// Init Set current environment and preload manifest
func Init(dev bool) error {
	if Plugin == "deprecated-stats" {
		Plugin = "stats"
		log.Println("go-webpack: default plugin will be changed to manifest instead of stats-plugin")
		log.Println("go-webpack: to continue using stats-plugin, please set webpack.Plugin = 'stats' explicitly")
	}

	var err error
	AssetTagHelper, AssetURLHelper, err = GetAssetHelpers(&Config{
		DevHost:       DevHost,
		FsPath:        FsPath,
		WebPath:       WebPath,
		Plugin:        Plugin,
		IgnoreMissing: IgnoreMissing,
		Verbose:       Verbose,
		IsDev:         dev,
		AssetHost:     AssetHost,
	})
	return err
}

func BasicConfig(host, path, webPath string) *Config {
	return &Config{
		DevHost:       host,
		FsPath:        path,
		WebPath:       webPath,
		Plugin:        "manifest",
		IgnoreMissing: true,
		Verbose:       true,
		IsDev:         false,
		AssetHost:     AssetHost,
	}
}

// AssetHelper renders asset tag with url from webpack manifest to the page

func readManifest(conf *Config) (map[string][]string, error) {
	return reader.Read(conf.Plugin, conf.DevHost, conf.FsPath, conf.WebPath, conf.IsDev)
}

func GetAssetHelpers(conf *Config) (AssetTagHelperFunc, AssetURLHelperFunc, error) {
	preloadedAssets := map[string][]string{}

	var err error
	if conf.IsDev {
		// Try to preload manifest, so we can show an error if webpack-dev-server is not running
		_, err = readManifest(conf)
		if err != nil {
			log.Println(err)
		}
	} else {
		preloadedAssets, err = readManifest(conf)
		// we won't ever re-check assets in this case.  this should be a hard error.
		if err != nil {
			return nil, nil, err
		}
	}

	return createAssetTagHelper(conf, preloadedAssets), createAssetURLHelper(conf, preloadedAssets), nil
}

func getValues(key, kind string, conf *Config, preloadedAssets map[string][]string) ([]string, error) {
	var err error

	var assets map[string][]string
	if conf.IsDev {
		assets, err = readManifest(conf)
		if err != nil {
			return nil, err
		}
	} else {
		assets = preloadedAssets
	}

	v, ok := assets[key]
	if !ok {
		message := "go-webpack: Asset file '" + key + "' not found in manifest"
		if conf.Verbose {
			log.Printf("%s. Manifest contents:", message)
			for k, a := range assets {
				log.Printf("%s: %s", k, a)
			}
		}
		if conf.IgnoreMissing {
			return nil, nil
		}
		return nil, errors.New(message)
	}

	values := []string{}
	for _, s := range v {
		if strings.HasSuffix(s, "."+kind) {
			url := s
			if len(conf.AssetHost) > 0 {
				url = conf.AssetHost + url
			}
			values = append(values, url)
		} else {
			log.Println("skip asset", s, ": bad type")
		}
	}

	return values, nil
}

func createAssetTagHelper(conf *Config, preloadedAssets map[string][]string) AssetTagHelperFunc {
	return func(key string) (template.HTML, error) {
		parts := strings.Split(key, ".")
		kind := parts[len(parts)-1]

		urls, err := getValues(key, kind, conf, preloadedAssets)
		if err != nil {
			return template.HTML(""), err
		}

		buf := []string{}
		for _, url := range urls {
			buf = append(buf, helper.AssetTag(kind, url))
		}
		return template.HTML(strings.Join(buf, "\n")), nil
	}
}

func createAssetURLHelper(conf *Config, preloadedAssets map[string][]string) AssetURLHelperFunc {
	return func(key string) (string, error) {
		parts := strings.Split(key, ".")
		kind := parts[len(parts)-1]

		urls, err := getValues(key, kind, conf, preloadedAssets)
		if err != nil {
			return "", err
		}

		buf := []string{}
		for _, url := range urls {
			buf = append(buf, url)
		}
		// this seems very wrong for multiple urls, but we don't have that problem right now :grimacing:
		return strings.Join(buf, ","), nil
	}
}
