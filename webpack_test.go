package webpack

import (
	"testing"
)

func TestManifestAssetHelper(t *testing.T) {
	assets := map[string][]string{
		"main.js": []string{"main.1.js", "main.2.js"},
	}

	helper := createAssetHelper(&Config{
		Plugin: "manifest",
	}, assets)

	html, err := helper("main.js")
	if err != nil {
		t.Fatalf("error %v returned from asset helper for valid asset", err)
	}
	expectedHTML :=
		`<script type="text/javascript" src="main.1.js"></script>
<script type="text/javascript" src="main.2.js"></script>`

	if string(html) != expectedHTML {
		t.Fatalf("unexpected <script> tags\nexpected:\n%s\nactual:\n%s", expectedHTML, html)
	}

	// IgnoreMissing = false
	_, err = helper("maiin.js")
	if err == nil {
		t.Fatalf("error nil when it shouldn't have been")
	}
}

func TestManifestAssetHelperWithAssetHost(t *testing.T) {
	assets := map[string][]string{
		"main.js": []string{"main.1.js", "main.2.js"},
	}

	helper := createAssetHelper(&Config{
		Plugin:    "manifest",
		AssetHost: "//cdn.com/prefix/",
	}, assets)

	html, err := helper("main.js")
	if err != nil {
		t.Fatalf("error %v returned from asset helper for valid asset", err)
	}
	expectedHTML :=
		`<script type="text/javascript" src="//cdn.com/prefix/main.1.js"></script>
<script type="text/javascript" src="//cdn.com/prefix/main.2.js"></script>`

	if string(html) != expectedHTML {
		t.Fatalf("unexpected <script> tags\nexpected:\n%s\nactual:\n%s", expectedHTML, html)
	}

	// IgnoreMissing = false
	_, err = helper("maiin.js")
	if err == nil {
		t.Fatalf("error nil when it shouldn't have been")
	}
}
