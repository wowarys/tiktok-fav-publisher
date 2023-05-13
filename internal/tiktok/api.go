package tiktok

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"

type ServiceApi struct {
	Log        *zap.Logger
	username   string
	count      int
	hdDownload int
	userAgent  string
}

type Videos []Video

type Video struct {
	ID             string `json:"video_id"`
	Title          string `json:"title"`
	AuthorUsername string
	DownloadLink   string
	ShareableLink  string
	Cover          string
}

type likedVideoResponse struct {
	Message string `json:"msg"`
	Data    struct {
		Videos `json:"videos"`
	} `json:"data"`
}

type videoInfoResponse struct {
	Message string `json:"msg"`
	Data    struct {
		Title       string `json:"title"`
		HDLink      string `json:"hdplay"`
		OriginCover string `json:"origin_cover"`
		RegularLink string `json:"play"`
		Author      struct {
			ID string `json:"unique_id"`
		} `json:"author"`
	}
}

func NewServiceApi(username string, count int, logger *zap.Logger) *ServiceApi {
	username = strings.TrimLeft(username, "@")
	return &ServiceApi{
		Log:        logger,
		username:   username,
		count:      count,
		userAgent:  userAgent,
		hdDownload: 1,
	}
}

func (sa *ServiceApi) GetLikedVideos() (Videos, error) {
	var (
		req      = &fasthttp.Request{}
		res      = &fasthttp.Response{}
		response likedVideoResponse
	)

	req.Header.SetMethod(http.MethodGet)
	req.SetRequestURI(fmt.Sprintf("https://www.tikwm.com/api/user/favorite?unique_id=%s&count=%d", sa.username, sa.count))
	req.Header.SetUserAgent(sa.userAgent)

	err := fasthttp.Do(req, res)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return nil, err
	}

	if response.Message != "success" {
		sa.Log.Error("not success message", zap.String("response", string(res.Body())))
		return nil, errors.New("not success message")
	}

	return response.Data.Videos, nil
}

func (sa *ServiceApi) SetVideoMetadata(video *Video) error {
	var (
		req    = &fasthttp.Request{}
		res    = &fasthttp.Response{}
		result videoInfoResponse
	)

	if video == nil {
		return errors.New("nil pointer")
	}

	req.Header.SetMethod(http.MethodGet)
	req.SetRequestURI(fmt.Sprintf("https://www.tikwm.com/api/?url=%s&hd=%d", video.ID, sa.hdDownload))
	req.Header.SetUserAgent(sa.userAgent)

	err := fasthttp.Do(req, res)
	if err != nil {
		return err
	}

	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		return err
	}

	if result.Message != "success" {
		sa.Log.Error("not success message", zap.String("response", string(res.Body())))
		return errors.New("not success message")
	}

	video.DownloadLink = result.Data.RegularLink
	video.AuthorUsername = result.Data.Author.ID
	video.Cover = result.Data.OriginCover
	video.Title = result.Data.Title
	video.ShareableLink = sa.constructShareLink(video.AuthorUsername, video.ID)
	return nil
}

func (sa *ServiceApi) constructShareLink(author, id string) string {
	return fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", author, id)
}
