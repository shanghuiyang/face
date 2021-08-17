package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/shanghuiyang/face"
	"github.com/shanghuiyang/oauth"
)

const (
	apiKey    = "your_baidu_app_key"
	secretKey = "your_baidu_secret_key"

	groupID = "mygroup"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("error: invalid args")
		fmt.Println("usage: image-recognizer test.jpg")
		os.Exit(1)
	}
	imgf := os.Args[1]
	img, err := ioutil.ReadFile(imgf)
	if err != nil {
		log.Printf("failed to read image file: %v, error: %v\n", imgf, err)
		os.Exit(1)
	}

	auth := oauth.NewBaiduOauth(apiKey, secretKey, oauth.NewCacheImp())
	f := face.NewBaiduFace(auth, groupID)
	name, err := f.Recognize(img)
	if err != nil {
		log.Printf("failed to recognize the image, error: %v", err)
		os.Exit(1)
	}
	fmt.Println(name)
}
