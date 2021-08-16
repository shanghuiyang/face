package face

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/shanghuiyang/oauth"
)

const (
	baiduFaceURL = "https://aip.baidubce.com/rest/2.0/face/v3/search"
)

// BaiduFace ...
type BaiduFace struct {
	auth    oauth.Oauth
	groupID string
}

type baiduResponse struct {
	ErrorCode int          `json:"error_code"`
	ErrorMsg  string       `json:"error_msg"`
	LogID     int64        `json:"log_id"`
	Timestamp int64        `json:"timestamp"`
	Cached    int64        `json:"cached"`
	Result    *baiduResult `json:"result"`
}

type baiduResult struct {
	FaceToken string  `json:"face_token"`
	UserList  []*user `json:"user_list"`
}

type user struct {
	GroupID  string  `json:"group_id"`
	UserID   string  `json:"user_id"`
	UserInfo string  `json:"user_info"`
	Score    float64 `json:"score"`
}

// New ...
func NewBaiduFace(auth oauth.Oauth, groupID string) *BaiduFace {
	return &BaiduFace{
		auth:    auth,
		groupID: groupID,
	}
}

// Recognize ...
func (f *BaiduFace) Recognize(image []byte) (string, error) {
	token, err := f.auth.Token()
	if err != nil {
		return "", err
	}

	b64img := base64.StdEncoding.EncodeToString(image)
	formData := url.Values{
		"access_token":  {token},
		"image":         {b64img},
		"image_type":    {"BASE64"},
		"group_id_list": {f.groupID},
	}
	resp, err := http.PostForm(baiduFaceURL, formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res baiduResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	if res.ErrorCode > 0 {
		return "", fmt.Errorf("error_code: %v, error_msg: %v", res.ErrorCode, res.ErrorMsg)
	}
	if res.Result == nil {
		return "", fmt.Errorf("not found")
	}

	maxScore := -1.
	result := ""
	for _, u := range res.Result.UserList {
		if u.Score > maxScore {
			maxScore = u.Score
			result = u.UserID
		}
	}
	return result, nil
}
