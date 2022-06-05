package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/go-flac/flacpicture"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
)

const (
	regexString = `^https://hi-resmusic.sonyselect.kuke.com/page/album.html\?id=(\d+)$`
	apiBase     = "https://api.sonyselect.com.cn/streaming/"
	key         = "RF9q4w<X$dof3pFF"
	iv          = "SiGK&MvKm9Y+c6f@"
	urlKey      = "@VvqA6gLLeNhuCQnPzlY@MYzvihjGbo75txdGkTLpJyj#9v3t8ewKwFjHqp"
)

var client = &http.Client{Transport: &Transport{}}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Bad responses if not set.
	req.Header.Add(
		"User-Agent", "",
	)
	return http.DefaultTransport.RoundTrip(req)
}

func handleErr(errText string, err error, _panic bool) {
	errString := fmt.Sprintf("%s\n%s", errText, err)
	if _panic {
		panic(errString)
	}
	fmt.Println(errString)
}

func getScriptDir() (string, error) {
	var (
		ok    bool
		err   error
		fname string
	)
	if filepath.IsAbs(os.Args[0]) {
		_, fname, _, ok = runtime.Caller(0)
		if !ok {
			return "", errors.New("Failed to get script filename.")
		}
	} else {
		fname, err = os.Executable()
		if err != nil {
			return "", err
		}
	}
	scriptDir := filepath.Dir(fname)
	return scriptDir, nil
}

func fileExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		return !f.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func readTxtFile(path string) ([]string, error) {
	var lines []string
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return lines, nil
}

func contains(lines []string, value string) bool {
	for _, line := range lines {
		if strings.EqualFold(line, value) {
			return true
		}
	}
	return false
}

func processUrls(urls []string) ([]string, error) {
	var (
		processed []string
		txtPaths  []string
	)
	for _, _url := range urls {
		if strings.HasSuffix(_url, ".txt") && !contains(txtPaths, _url) {
			txtLines, err := readTxtFile(_url)
			if err != nil {
				return nil, err
			}
			for _, txtLine := range txtLines {
				if !contains(processed, txtLine) {
					processed = append(processed, txtLine)
				}
			}
			txtPaths = append(txtPaths, _url)
		} else {
			if !contains(processed, _url) {
				processed = append(processed, _url)
			}
		}
	}
	return processed, nil
}

func readConfig() (*Config, error) {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var obj Config
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func parseArgs() *Args {
	var args Args
	arg.MustParse(&args)
	return &args
}

func makeDirs(path string) error {
	err := os.MkdirAll(path, 0755)
	return err
}

func checkUrl(_url string) string {
	regex := regexp.MustCompile(regexString)
	match := regex.FindStringSubmatch(_url)
	if match == nil {
		return ""
	}
	return match[1]
}

func parseCfg() (*Config, error) {
	cfg, err := readConfig()
	if err != nil {
		return nil, err
	}
	args := parseArgs()
	cfg.Urls, err = processUrls(args.Urls)
	if err != nil {
		errString := fmt.Sprintf("Failed to process URLs.\n%s", err)
		return nil, errors.New(errString)
	}
	return cfg, nil
}

func getPostData(sonySelectId string) *PostData {
	var data PostData
	data.Header.AccessKey = "sonyhires"
	data.Header.ContentEncryption = false
	data.Header.Imei = "BFEBFBFF000306C3"
	data.Header.Model = "Microsoft Windows NT 10.0.19043.0"
	data.Header.SignEncryption = false
	data.Header.SonySelectID = sonySelectId
	data.Header.Version = "1.1.6.0"
	return &data
}

func pkcs5Trimming(data []byte) []byte {
	padding := data[len(data)-1]
	return data[:len(data)-int(padding)]
}

func decryptFileMeta(encResp string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encResp)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(iv))
	decrypted := make([]byte, len(decoded))
	ecb.CryptBlocks(decrypted, decoded)
	return pkcs5Trimming(decrypted), nil
}

func generateEpoch() int64 {
	now := time.Now()
	return now.UnixMilli()
}

func getAlbumMeta(albumId string, data *PostData) (*AlbumMeta, error) {
	epoch := generateEpoch()
	data.Content.AlbumID = albumId
	data.Header.Timestamp = epoch
	mData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, apiBase+"album/get_detail/v1/windows", bytes.NewBuffer(mData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(mData)))
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return nil, errors.New(do.Status)
	}
	var obj AlbumMeta
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func getFileMeta(indexId int, data *PostData) (*FileMeta, error) {
	epoch := generateEpoch()
	data.Content.IndexID = indexId
	data.Header.Timestamp = epoch
	mData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, apiBase+"play/get_segment_index/v2/windows", bytes.NewBuffer(mData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(mData)))
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return nil, errors.New(do.Status)
	}
	var obj FileMetaEnc
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}
	// Endpoint always forces encryption.
	decrypted, err := decryptFileMeta(obj.Content.Encrypcontent)
	if err != nil {
		return nil, err
	}
	var fileMeta FileMeta
	err = json.Unmarshal(decrypted, &fileMeta)
	if err != nil {
		return nil, err
	}
	return &fileMeta, nil
}

func getTrackMeta(musicId int, data *PostData) (*TrackMeta, error) {
	epoch := generateEpoch()
	data.Content.MusicId = musicId
	data.Header.Timestamp = epoch
	mData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, apiBase+"music/get_detail/v1/windows", bytes.NewBuffer(mData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(mData)))
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return nil, errors.New(do.Status)
	}
	var obj TrackMeta
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func parseAlbumMeta(_meta *AlbumMeta) map[string]string {
	meta := _meta.Content
	parsedMeta := map[string]string{
		"album":       meta.Name,
		"albumArtist": meta.Artist,
		"year":        meta.ReleaseTime[:4],
	}
	return parsedMeta
}

func parseTrackMeta(_meta *TrackMeta, albMeta map[string]string, trackNum, trackTotal int) map[string]string {
	meta := _meta.Content
	workName := meta.WorkName
	trackName := meta.MusicName
	if workName != "" {
		trackName = workName + ": " + trackName
	}
	albMeta["artist"] = meta.Artist
	albMeta["title"] = trackName
	albMeta["track"] = strconv.Itoa(trackNum)
	albMeta["trackPad"] = fmt.Sprintf("%02d", trackNum)
	albMeta["trackTotal"] = strconv.Itoa(trackTotal)
	return albMeta
}

func sanitize(filename string) string {
	regex := regexp.MustCompile(`[\/:*?"><|]`)
	sanitized := regex.ReplaceAllString(filename, "_")
	return sanitized
}

func generateUrl(_url string) (string, error) {
	date := time.Now().Format("200601021504")
	u, err := url.Parse(_url)
	if err != nil {
		return "", err
	}
	pathAndQuery := u.Path + u.RawQuery
	h := md5.New()
	_, err = h.Write([]byte(urlKey + date + pathAndQuery))
	if err != nil {
		return "", err
	}
	hash := hex.EncodeToString(h.Sum(nil))
	genUrl := u.Scheme + "://" + u.Host + "/" + date + "/" + hash + pathAndQuery
	return genUrl, nil
}

func downloadSegs(tmpPath string, meta *FileMeta) ([]string, error) {
	var segPaths []string
	base := meta.BaseURL
	segTotal := len(meta.Names)
	for segNum, fname := range meta.Names {
		segNum++
		segPath := filepath.Join(tmpPath, fname)
		f, err := os.OpenFile(segPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
		if err != nil {
			return nil, err
		}
		fmt.Printf("\rSegment %d of %d.", segNum, segTotal)
		url, err := generateUrl(base + fname)
		if err != nil {
			f.Close()
			return nil, err
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			f.Close()
			return nil, err
		}
		do, err := client.Do(req)
		if err != nil {
			f.Close()
			return nil, err
		}
		if do.StatusCode != http.StatusOK {
			do.Body.Close()
			f.Close()
			return nil, errors.New(do.Status)
		}
		_, err = io.Copy(f, do.Body)
		do.Body.Close()
		f.Close()
		if err != nil {
			return nil, err
		}
		segPaths = append(segPaths, segPath)
	}
	fmt.Println("")
	return segPaths, nil
}

func getTmpPath() (string, error) {
	return os.MkdirTemp(os.TempDir(), "")
}

func cleanup(segPaths []string, txtPath string) {
	exists, _ := fileExists(txtPath)
	if exists {
		os.Remove(txtPath)
	}
	for _, segPath := range segPaths {
		os.Remove(segPath)
	}
}

func mergeSegs(trackPath, tmpPath string, segPaths []string) error {
	txtPath := filepath.Join(tmpPath, "tmp.txt")
	defer cleanup(segPaths, txtPath)
	f, err := os.OpenFile(txtPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	for _, segPath := range segPaths {
		line := fmt.Sprintf("file '%s'\n", segPath)
		_, err := f.WriteString(line)
		if err != nil {
			f.Close()
			return err
		}
	}
	f.Close()
	var (
		errBuffer bytes.Buffer
		args      = []string{"-f", "concat", "-safe", "0", "-i", txtPath, "-c", "flac", "-y", trackPath}
	)
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stderr = &errBuffer
	err = cmd.Run()
	if err != nil {
		errString := fmt.Sprintf("%s\n%s", err, errBuffer.String())
		return errors.New(errString)
	}
	return nil
}

func getTracktotal(albumMeta *AlbumMeta) int {
	trackTotal := 0
	for _, cd := range albumMeta.Content.CdList {
		trackTotal += len(cd.Musiclist)
	}
	return trackTotal
}

func parseTemplate(templateText string, tags map[string]string) string {
	var buffer bytes.Buffer
	for {
		err := template.Must(template.New("").Parse(templateText)).Execute(&buffer, tags)
		if err == nil {
			break
		}
		fmt.Println("Failed to parse template. Default will be used instead.")
		templateText = "{{.trackPad}}. {{.title}}"
		buffer.Reset()
	}
	return html.UnescapeString(buffer.String())
}

func formatFreq(freq int) string {
	freqStr := strconv.Itoa(freq / 100)
	if strings.HasSuffix(freqStr, "0") {
		return strconv.Itoa(freq / 1000)
	} else {
		freqStrLen := len(freqStr)
		return freqStr[:freqStrLen-1] + "." + freqStr[freqStrLen-1:]
	}
}

func downloadCover(_url, path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	req, err := client.Get(_url)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return errors.New(req.Status)
	}
	_, err = io.Copy(f, req.Body)
	return err
}

func readCover(coverPath string) ([]byte, error) {
	imgData, err := ioutil.ReadFile(coverPath)
	if err != nil {
		return nil, err
	}
	return imgData, nil
}

func writeTags(trackPath, coverPath string, tags map[string]string) error {
	var (
		err     error
		imgData []byte
	)
	if coverPath != "" {
		imgData, err = readCover(coverPath)
		if err != nil {
			return err
		}
	}
	delete(tags, "trackPad")
	f, err := flac.ParseFile(trackPath)
	if err != nil {
		return err
	}
	tag := flacvorbis.New()
	for k, v := range tags {
		tag.Add(strings.ToUpper(k), v)
	}
	tagMeta := tag.Marshal()
	f.Meta = append(f.Meta, &tagMeta)
	if imgData != nil {
		picture, err := flacpicture.NewFromImageData(
			flacpicture.PictureTypeFrontCover, "", imgData, "image/jpeg",
		)
		if err != nil {
			return err
		}
		pictureMeta := picture.Marshal()
		f.Meta = append(f.Meta, &pictureMeta)
	}
	return f.Save(trackPath)
}

func init() {
	fmt.Println(`
 _____                _____     _          ____                _           _         
|   __|___ ___ _ _   |  |  |_ _| |_ ___   |    \ ___ _ _ _ ___| |___ ___ _| |___ ___ 
|__   | . |   | | |  |    -| | | '_| -_|  |  |  | . | | | |   | | . | .'| . | -_|  _|
|_____|___|_|_|_  |  |__|__|___|_,_|___|  |____/|___|_____|_|_|_|___|__,|___|___|_|  
              |___|                                                                  
`)
}

func main() {
	var albumFolder string
	scriptDir, err := getScriptDir()
	if err != nil {
		panic(err)
	}
	tmpPath, err := getTmpPath()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(scriptDir)
	if err != nil {
		panic(err)
	}
	cfg, err := parseCfg()
	if err != nil {
		handleErr("Failed to parse config file.", err, true)
	}
	err = makeDirs(cfg.OutPath)
	if err != nil {
		handleErr("Failed to make output folder.", err, true)
	}
	postData := getPostData(cfg.SonySelectID)
	albumTotal := len(cfg.Urls)
out:
	for albumNum, _url := range cfg.Urls {
		fmt.Printf("Album %d of %d:\n", albumNum+1, albumTotal)
		albumId := checkUrl(_url)
		if albumId == "" {
			fmt.Println("Invalid URL:", _url)
			continue
		}
		albumMeta, err := getAlbumMeta(albumId, postData)
		if err != nil {
			handleErr("Failed to fetch album metadata.", err, false)
			continue
		}
		parsedAlbMeta := parseAlbumMeta(albumMeta)
		if cfg.OmitArtists {
			albumFolder = parsedAlbMeta["album"]
		} else {
			albumFolder = parsedAlbMeta["albumArtist"] + " - " + parsedAlbMeta["album"]
		}
		fmt.Println(albumFolder)
		if len(albumFolder) > 120 {
			fmt.Println("Album folder was chopped as it exceeds 120 characters.")
			albumFolder = albumFolder[:120]
		}
		albumPath := filepath.Join(cfg.OutPath, sanitize(albumFolder))
		err = makeDirs(albumPath)
		if err != nil {
			handleErr("Failed to make album folder.", err, false)
			continue
		}
		coverPath := filepath.Join(albumPath, "cover.jpg")
		err = downloadCover(albumMeta.Content.LargeIcon, coverPath)
		if err != nil {
			handleErr("Failed to get cover.", err, false)
			coverPath = ""
		}
		trackNum := 0
		trackTotal := getTracktotal(albumMeta)
		for _, cd := range albumMeta.Content.CdList {
			for _, track := range cd.Musiclist {
				trackNum++
				trackMeta, err := getTrackMeta(track.MusicID, postData)
				if err != nil {
					handleErr("Failed to get track metadata.", err, false)
					continue
				}
				parsedMeta := parseTrackMeta(trackMeta, parsedAlbMeta, trackNum, trackTotal)
				if trackMeta.Content.PlayModels[0].Type != "streaming" {
					fmt.Println("Album is purchase-only.")
					continue out
				}
				trackFname := parseTemplate(cfg.TrackTemplate, parsedMeta)
				sanTrackFname := sanitize(trackFname)
				trackPath := filepath.Join(albumPath, sanTrackFname+".flac")
				exists, err := fileExists(trackPath)
				if err != nil {
					handleErr("Failed to check if track already exists locally.", err, false)
					continue
				}
				if exists {
					fmt.Println("Track already exists locally.")
					continue
				}
				fileMeta, err := getFileMeta(trackMeta.Content.PlayModels[0].IndexID, postData)
				if err != nil {
					handleErr("Failed to get file metadata.", err, false)
					continue
				}
				fmt.Printf(
					"Downloading track %d of %d: %s - %d-bit / %s kHz FLAC\n",
					trackNum, trackTotal, parsedMeta["title"], fileMeta.SampleBit, formatFreq(fileMeta.SampleRate),
				)
				segPaths, err := downloadSegs(tmpPath, fileMeta)
				if err != nil {
					handleErr("Failed to download segments.", err, false)
					continue
				}
				err = mergeSegs(trackPath, tmpPath, segPaths)
				if err != nil {
					handleErr("Failed to merge segments.", err, false)
					continue
				}
				err = writeTags(trackPath, coverPath, parsedMeta)
				if err != nil {
					fmt.Printf("Failed to write tags.\n%s", err)
					continue
				}
				if coverPath != "" && !cfg.KeepCover {
					err := os.Remove(coverPath)
					if err != nil {
						handleErr("Failed to delete cover.", err, false)
					}
				}
			}
		}
	}
}
