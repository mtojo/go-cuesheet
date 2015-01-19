package cuesheet

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

const (
	delims          = "\t\n\r "
	eol             = "\n"
	framesPerSecond = 75
)

type Frame uint64
type Flags int

const (
	None    Flags = iota
	Dcp           = 1 << iota
	Four_ch       = 1 << iota
	Pre           = 1 << iota
	Scms          = 1 << iota
)

type TrackIndex struct {
	Number uint
	Frame  Frame
}

type Track struct {
	TrackNumber   uint
	TrackDataType string
	Flags         Flags
	Isrc          string
	Title         string
	Performer     string
	SongWriter    string
	Pregap        Frame
	Postgap       Frame
	Index         []TrackIndex
}

type File struct {
	FileName string
	FileType string
	Tracks   []Track
}

type Cuesheet struct {
	Rem        []string
	Catalog    string
	CdTextFile string
	Title      string
	Performer  string
	SongWriter string
	Pregap     Frame
	Postgap    Frame
	File       []File
}

func ReadFile(r io.Reader) (*Cuesheet, error) {
	b := bufio.NewReader(r)
	cuesheet := &Cuesheet{}

	for {
		line, err := (*b).ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		line = strings.Trim(line, delims)
		command := ReadString(&line)

		switch command {
		case "REM":
			(*cuesheet).Rem = append((*cuesheet).Rem, line)
			break
		case "CATALOG":
			(*cuesheet).Catalog = line
			break
		case "CDTEXTFILE":
			(*cuesheet).CdTextFile = ReadString(&line)
			break
		case "TITLE":
			(*cuesheet).Title = ReadString(&line)
			break
		case "PERFORMER":
			(*cuesheet).Performer = ReadString(&line)
			break
		case "SONGWRITER":
			(*cuesheet).SongWriter = ReadString(&line)
			break
		case "PREGAP":
			(*cuesheet).Pregap = ReadFrame(&line)
			break
		case "POSTGAP":
			(*cuesheet).Postgap = ReadFrame(&line)
			break
		case "FILE":
			fname := ReadString(&line)
			ftype := ReadString(&line)
			tracks, err := readTracks(b)
			if err != nil {
				return nil, err
			}
			(*cuesheet).File = append((*cuesheet).File, File{fname, ftype, *tracks})
			break
		default:
			break
		}
	}

	return cuesheet, nil
}

func WriteFile(w io.Writer, cuesheet *Cuesheet) error {
	ws := bufio.NewWriter(w)

	for i := 0; i < len((*cuesheet).Rem); i++ {
		ws.WriteString("REM " + (*cuesheet).Rem[i] + eol)
	}

	if len((*cuesheet).Catalog) > 0 {
		ws.WriteString("CATALOG " + (*cuesheet).Catalog + eol)
	}

	if len((*cuesheet).CdTextFile) > 0 {
		ws.WriteString("CDTEXTFILE " + FormatString((*cuesheet).CdTextFile) + eol)
	}

	if len((*cuesheet).Title) > 0 {
		ws.WriteString("TITLE " + FormatString((*cuesheet).Title) + eol)
	}

	if len((*cuesheet).Performer) > 0 {
		ws.WriteString("PERFORMER " + FormatString((*cuesheet).Performer) + eol)
	}

	if len((*cuesheet).SongWriter) > 0 {
		ws.WriteString("SONGWRITER " + FormatString((*cuesheet).SongWriter) + eol)
	}

	if (*cuesheet).Pregap > 0 {
		ws.WriteString("PREGAP " + FormatFrame((*cuesheet).Pregap) + eol)
	}

	if (*cuesheet).Postgap > 0 {
		ws.WriteString("POSTGAP " + FormatFrame((*cuesheet).Postgap) + eol)
	}

	for i := 0; i < len((*cuesheet).File); i++ {
		file := (*cuesheet).File[i]
		ws.WriteString("FILE " + FormatString(file.FileName) +
			" " + file.FileType + eol)

		for i := 0; i < len(file.Tracks); i++ {
			track := file.Tracks[i]

			ws.WriteString("  TRACK " + FormatTrackNumber(track.TrackNumber) +
				" " + track.TrackDataType + eol)

			if track.Flags != None {
				ws.WriteString("    FLAGS")
				if (track.Flags & Dcp) != 0 {
					ws.WriteString(" DCP")
				}
				if (track.Flags & Four_ch) != 0 {
					ws.WriteString(" 4CH")
				}
				if (track.Flags & Pre) != 0 {
					ws.WriteString(" PRE")
				}
				if (track.Flags & Scms) != 0 {
					ws.WriteString(" SCMS")
				}
				ws.WriteString(eol)
			}

			if len(track.Isrc) > 0 {
				ws.WriteString("    ISRC " + track.Isrc + eol)
			}

			if len(track.Title) > 0 {
				ws.WriteString("    TITLE " + FormatString(track.Title) + eol)
			}

			if len(track.Performer) > 0 {
				ws.WriteString("    PERFORMER " + FormatString(track.Performer) + eol)
			}

			if len(track.SongWriter) > 0 {
				ws.WriteString("    SONGWRITER " + FormatString(track.SongWriter) + eol)
			}

			if track.Pregap > 0 {
				ws.WriteString("    PREGAP " + FormatFrame(track.Pregap) + eol)
			}

			if track.Postgap > 0 {
				ws.WriteString("    POSTGAP " + FormatFrame(track.Postgap) + eol)
			}

			for i := 0; i < len(track.Index); i++ {
				index := track.Index[i]
				ws.WriteString("    INDEX " + FormatTrackNumber(index.Number) +
					" " + FormatFrame(index.Frame) + eol)
			}
		}
	}

	ws.Flush()

	return nil
}

func ReadString(s *string) string {
	*s = strings.TrimLeft(*s, delims)
	if isQuoted(*s) {
		v := unquote(*s)
		*s = (*s)[len(v)+2:]
		return v
	}
	for i := 0; i < len(*s); i++ {
		if (*s)[i] == ' ' {
			v := (*s)[0:i]
			*s = (*s)[i+1:]
			return v
		}
	}
	v := *s
	*s = ""
	return v
}

func ReadInt(s *string) int {
	v := ReadString(s)
	if n, err := strconv.ParseInt(v, 10, 32); err == nil {
		return int(n)
	}
	return 0
}

func ReadUint(s *string) uint {
	v := ReadString(s)
	if n, err := strconv.ParseUint(v, 10, 32); err == nil {
		return uint(n)
	}
	return 0
}

func ReadFrame(s *string) Frame {
	v := strings.Split(ReadString(s), ":")
	if len(v) == 3 {
		mm, _ := strconv.ParseUint(v[0], 10, 32)
		ss, _ := strconv.ParseUint(v[1], 10, 32)
		ff, _ := strconv.ParseUint(v[2], 10, 32)
		return Frame((mm*60+ss)*framesPerSecond + ff)
	}
	return 0
}

func FormatString(s string) string {
	if strings.ContainsAny(s, delims) {
		return quote(s, '"')
	}
	return s
}

func FormatTrackNumber(n uint) string {
	return leftPad(strconv.FormatUint(uint64(n), 10), "0", 2)
}

func FormatFrame(frame Frame) string {
	n := frame / framesPerSecond
	mm := n / 60
	ss := n % 60
	ff := frame % framesPerSecond
	return leftPad(strconv.FormatUint(uint64(mm), 10), "0", 2) + ":" +
		leftPad(strconv.FormatUint(uint64(ss), 10), "0", 2) + ":" +
		leftPad(strconv.FormatUint(uint64(ff), 10), "0", 2)
}

func isQuoted(s string) bool {
	return s[0] == '"' || s[0] == '\''
}

func quote(s string, quote byte) string {
	buf := make([]byte, 0, 3*len(s)/2)
	buf = append(buf, quote)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == quote || c == '\\' {
			buf = append(buf, '\\')
			buf = append(buf, byte(c))
		} else {
			buf = append(buf, byte(c))
		}
	}
	buf = append(buf, quote)
	return string(buf)
}

func unquote(s string) string {
	quote := s[0]
	i := 1
	for ; i < len(s); i++ {
		if s[i] == quote {
			break
		}
		if s[i] == '\\' {
			i++
		}
	}
	return s[1:i]
}

func readTrack(b *bufio.Reader, track *Track) error {
L:
	for {
		before := *b
		line, err := (*b).ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if !strings.HasPrefix(line, "    ") {
			*b = before
			break
		}
		line = strings.Trim(line, delims)
		command := ReadString(&line)

		switch command {
		case "FLAGS":
			(*track).Flags = None
			for len(line) > 0 {
				switch ReadString(&line) {
				case "DCP":
					(*track).Flags |= Dcp
					break
				case "4CH":
					(*track).Flags |= Four_ch
					break
				case "PRE":
					(*track).Flags |= Pre
					break
				case "SCMS":
					(*track).Flags |= Scms
					break
				default:
					break
				}
			}
			break
		case "ISRC":
			(*track).Isrc = line
			break
		case "TITLE":
			(*track).Title = unquote(line)
			break
		case "PERFORMER":
			(*track).Performer = unquote(line)
			break
		case "SONGWRITER":
			(*track).SongWriter = unquote(line)
			break
		case "PREGAP":
			line = line
			(*track).Pregap = ReadFrame(&line)
			break
		case "POSTGAP":
			line = line
			(*track).Postgap = ReadFrame(&line)
			break
		case "INDEX":
			index := TrackIndex{}
			index.Number = ReadUint(&line)
			index.Frame = ReadFrame(&line)
			(*track).Index = append((*track).Index, index)
			break
		default:
			break L
		}
	}

	return nil
}

func readTracks(b *bufio.Reader) (*[]Track, error) {
	tracks := &[]Track{}

L:
	for {
		before := *b
		line, err := (*b).ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(line, "  ") {
			*b = before
			break
		}
		line = strings.Trim(line, delims)
		command := ReadString(&line)

		switch command {
		case "TRACK":
			track := Track{}
			track.TrackNumber = ReadUint(&line)
			track.TrackDataType = ReadString(&line)
			if err := readTrack(b, &track); err != nil {
				return nil, err
			}
			*tracks = append(*tracks, track)
			break
		default:
			break L
		}
	}

	return tracks, nil
}

func leftPad(s, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}
