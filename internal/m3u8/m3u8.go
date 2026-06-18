package m3u8

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/grafov/m3u8"
	"github.com/injoyai/bar/internal/shell"
)

// defaultClient 用于下载 m3u8 索引文件，设置合理超时避免无限阻塞
var defaultClient = &http.Client{Timeout: 30 * time.Second}

func Decode(url string) ([]string, error) {

	// 校验 URL 必须含 "/"，否则后续 baseURL 拼接会出错
	idx := strings.LastIndex(url, "/")
	if idx < 0 {
		return nil, errors.New("无效的 m3u8 url: " + url)
	}

	// 下载 m3u8 文件
	resp, err := defaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("m3u8 下载失败, status=%d", resp.StatusCode)
	}

	playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		return nil, err
	}

	if listType != m3u8.MEDIA {
		return nil, errors.New("不是 MediaPlaylist（可能是 MasterPlaylist）")
	}

	media := playlist.(*m3u8.MediaPlaylist)

	baseURL := url[:idx+1]

	ls := make([]string, 0, len(media.Segments))

	for _, segment := range media.Segments {
		if segment == nil {
			continue
		}
		tsURL := segment.URI
		if !strings.HasPrefix(tsURL, "http") {
			tsURL = baseURL + tsURL
		}

		ls = append(ls, tsURL)
	}

	return ls, nil
}

func MergeByFFmpeg(dir, output string) error {
	lsFilename := filepath.Join(dir, "ts_list.txt")
	lsFilename = strings.ReplaceAll(lsFilename, "\\", "/")
	file, err := os.Create(lsFilename)
	if err != nil {
		return err
	}
	defer os.Remove(lsFilename)
	defer file.Close()

	es, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range es {
		info, err := e.Info()
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".ts") {
			_, err = file.WriteString("file '" + info.Name() + "'\r\n")
		}
	}

	cmd := fmt.Sprintf("ffmpeg -y -f concat -i %s -c copy %s", lsFilename, output)
	return shell.Run(cmd)
}
