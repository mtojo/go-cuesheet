package cuesheet

import (
	"os"
	"reflect"
	"testing"
)

const cueFile = "test.cue"

func TestCuesheet(t *testing.T) {
	actual := Cuesheet{
		[]string{
			"GENRE \"Electronica\"",
			"DATE \"2015\"",
		},
		"1234567890123",
		"test.cdt",
		"Test Album Title",
		"Test Album Performer",
		"Test Album SongWriter",
		0,
		0,
		[]File{
			File{
				"test.wav",
				"WAVE",
				[]Track{
					Track{
						1,
						"AUDIO",
						None,
						"ABCDE1234567",
						"Test Title",
						"Test Performer",
						"Test SongWriter",
						0,
						0,
						[]TrackIndex{
							TrackIndex{
								1,
								0,
							},
						},
					},
					Track{
						2,
						"MODE1/2048",
						Dcp,
						"ABCDE1234567",
						"Test Title",
						"Test Performer",
						"Test SongWriter",
						0,
						0,
						[]TrackIndex{
							TrackIndex{
								1,
								2715,
							},
						},
					},
				},
			},
		},
	}

	{
		w, err := os.Create(cueFile)
		if err != nil {
			t.Error(err)
		}
		defer w.Close()

		if err := WriteFile(w, &actual); err != nil {
			t.Error(err)
		}
	}

	var expected *Cuesheet

	{
		r, err := os.Open(cueFile)
		if err != nil {
			t.Error(err)
		}
		defer r.Close()
		defer os.Remove(cueFile)

		expected, err = ReadFile(r)
		if err != nil {
			t.Error(err)
		}
	}

	if !reflect.DeepEqual(actual, *expected) {
		t.Errorf("wrong reading data")
	}

	genre := (*expected).Rem[0]
	if ReadString(&genre) != "GENRE" || ReadString(&genre) != "Electronica" {
		t.Errorf("wrong reading data")
	}

	date := (*expected).Rem[1]
	if ReadString(&date) != "DATE" || ReadUint(&date) != 2015 {
		t.Errorf("wrong reading data")
	}
}
