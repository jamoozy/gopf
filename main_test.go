package main

import (
	"bytes"
	"testing"

	"html/template"
)

func TestTemplate(t *testing.T) {
	tmpl, err := template.ParseFiles("index.tmpl.html")
	if err != nil {
		t.Fatal(err)
	}

	tcs := []struct {
		desc string
		data map[string]interface{}
	}{
		{
			desc: "simple values",
			data: map[string]interface{}{
				"Title": "title",
				"PlaylistItems": &PlaylistData{
					Playlist: "",
					Fnames:   []string{"foo"},
				},
				"MediaItems": &MediaData{},
			},
		},

		{
			desc: "nil values",
			data: map[string]interface{}{
				"Title":         "title",
				"PlaylistItems": nil,
				"MediaItems":    nil,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, tc.data); err != nil {
				t.Fatal(err)
			}
		})
	}
}
