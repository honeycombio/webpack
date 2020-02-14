package webpack

import (
	"testing"
)

func TestManifestAssetHelper(t *testing.T) {
	assets := map[string][]string{
		"main.js":   []string{"main.1.js", "main.2.js"},
		"image.png": []string{"image.3.png"},
	}

	tagHelper := createAssetTagHelper(&Config{
		Plugin: "manifest",
	}, assets)

	urlHelper := createAssetURLHelper(&Config{
		Plugin: "manifest",
	}, assets)

	html, err := tagHelper("main.js")
	if err != nil {
		t.Fatalf("error %v returned from asset tag helper for valid asset", err)
	}
	expectedHTML :=
		`<script type="text/javascript" src="main.1.js"></script>
<script type="text/javascript" src="main.2.js"></script>`

	if string(html) != expectedHTML {
		t.Fatalf("unexpected <script> tags\nexpected:\n%s\nactual:\n%s", expectedHTML, html)
	}

	html, err = tagHelper("image.png")
	if err != nil {
		t.Fatalf("error %v returned from asset tag helper for valid asset", err)
	}
	expectedHTML =
		`<img src="image.3.png"></img>`

	if string(html) != expectedHTML {
		t.Fatalf("unexpected <script> tags\nexpected:\n%s\nactual:\n%s", expectedHTML, html)
	}

	url, err := urlHelper("main.js")
	if err != nil {
		t.Fatalf("error %v returned from asset url helper for valid asset", err)
	}

	expectedURL := "main.1.js,main.2.js"
	if url != expectedURL {
		t.Fatalf("unexpected url\nexpected:\n%s\nactual:\n%s", expectedURL, url)
	}

	// IgnoreMissing = false
	_, err = tagHelper("maiin.js")
	if err == nil {
		t.Fatalf("error nil when it shouldn't have been")
	}

	_, err = urlHelper("maiin.js")
	if err == nil {
		t.Fatalf("error nil when it shouldn't have been")
	}
}

func TestManifestAssetHelperWithAssetHost(t *testing.T) {
	assets := map[string][]string{
		"main.js":   []string{"main.1.js"},
		"image.png": []string{"image.3.png"},
	}

	tagHelper := createAssetTagHelper(&Config{
		Plugin:    "manifest",
		AssetHost: "//cdn.com/prefix/",
	}, assets)

	urlHelper := createAssetURLHelper(&Config{
		Plugin:    "manifest",
		AssetHost: "//cdn.com/prefix/",
	}, assets)

	html, err := tagHelper("main.js")
	if err != nil {
		t.Fatalf("error %v returned from asset tag helper for valid asset", err)
	}
	expectedHTML :=
		`<script type="text/javascript" src="//cdn.com/prefix/main.1.js"></script>`

	if string(html) != expectedHTML {
		t.Fatalf("unexpected <script> tag\nexpected:\n%s\nactual:\n%s", expectedHTML, html)
	}

	html, err = tagHelper("image.png")
	if err != nil {
		t.Fatalf("error %v returned from asset tag helper for valid asset", err)
	}
	expectedHTML =
		`<img src="//cdn.com/prefix/image.3.png"></img>`

	if string(html) != expectedHTML {
		t.Fatalf("unexpected <script> tags\nexpected:\n%s\nactual:\n%s", expectedHTML, html)
	}

	url, err := urlHelper("main.js")
	if err != nil {
		t.Fatalf("error %v returned from asset url helper for valid asset", err)
	}

	expectedURL := "//cdn.com/prefix/main.1.js"
	if url != expectedURL {
		t.Fatalf("unexpected url\nexpected:\n%s\nactual:\n%s", expectedURL, url)
	}

	// IgnoreMissing = false
	_, err = tagHelper("maiin.js")
	if err == nil {
		t.Fatalf("error nil when it shouldn't have been")
	}

	_, err = urlHelper("maiin.js")
	if err == nil {
		t.Fatalf("error nil when it shouldn't have been")
	}
}
