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

func TestWebPathRE(t *testing.T) {
	tcs := []struct {
		desc string
		body string // input to regexp
		repl string // replacement strin
		res  string // result we expect
	}{
		{
			desc: "Clasior",
			body: `../Royksopp/Senior/01 - ...And TheForest Began To Sing.mp3
../Royksopp/Senior/02 - Tricky Two.mp3
../Royksopp/Senior/03 - The Alcoholic.mp3
../Royksopp/Senior/04 - Senior Living.mp3
../Royksopp/Senior/05 - The Drug.mp3
../Royksopp/Senior/06 - Forsaken Cowboy.mp3
../Royksopp/Senior/07 - The Fear.mp3
../Royksopp/Senior/08 - Coming Home.mp3
../Ratatat/Classics/01 Montanita.mp3
../Ratatat/Classics/02 Lex.mp3
../Ratatat/Classics/03 Gettysburg.mp3
../Ratatat/Classics/04 Wildcat.mp3
../Ratatat/Classics/05 Tropicana.mp3
../Ratatat/Classics/06 Loud Pipes.mp3
../Ratatat/Classics/07 Kennedy.mp3
../Ratatat/Classics/08 Swisha.mp3
../Ratatat/Classics/09 Nostrand.mp3
../Ratatat/Classics/10 Tacobel Canon.mp3
../Royksopp/The Understanding/04 - Sombre Detune.ogg
`,
			repl: "audio/",
			res: `audio/Royksopp/Senior/01 - ...And TheForest Began To Sing.mp3
audio/Royksopp/Senior/02 - Tricky Two.mp3
audio/Royksopp/Senior/03 - The Alcoholic.mp3
audio/Royksopp/Senior/04 - Senior Living.mp3
audio/Royksopp/Senior/05 - The Drug.mp3
audio/Royksopp/Senior/06 - Forsaken Cowboy.mp3
audio/Royksopp/Senior/07 - The Fear.mp3
audio/Royksopp/Senior/08 - Coming Home.mp3
audio/Ratatat/Classics/01 Montanita.mp3
audio/Ratatat/Classics/02 Lex.mp3
audio/Ratatat/Classics/03 Gettysburg.mp3
audio/Ratatat/Classics/04 Wildcat.mp3
audio/Ratatat/Classics/05 Tropicana.mp3
audio/Ratatat/Classics/06 Loud Pipes.mp3
audio/Ratatat/Classics/07 Kennedy.mp3
audio/Ratatat/Classics/08 Swisha.mp3
audio/Ratatat/Classics/09 Nostrand.mp3
audio/Ratatat/Classics/10 Tacobel Canon.mp3
audio/Royksopp/The Understanding/04 - Sombre Detune.ogg
`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			if out := webPathRE.ReplaceAll([]byte(tc.body), []byte(tc.repl)); string(out) != tc.res {
				t.Fatalf("Mismatch:\n%s\n%s", tc.res, string(out))
			}
		})
	}
}
