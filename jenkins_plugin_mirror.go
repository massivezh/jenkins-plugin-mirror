package main

import (
	"fmt"
	"net/http"
	"net/url"
	"io"
	"github.com/bitly/go-simplejson"
	"os"
	"path"
	"log"
	"crypto/sha1"
	"path/filepath"
	"encoding/base64"
)
func chk_sha1(filePath string, sha1_sum string) bool {
	//Initialize variable returnMD5String now in case an error has to be returned
	
	//Open the filepath passed by the argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	
	//Tell the program to call the following function when the current function returns
	defer file.Close()
	
	//Open a new SHA1 hash interface to write to
	hash := sha1.New()
	
	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return false
	}
	
	hashInBytes := hash.Sum(nil)

	base64str := base64.URLEncoding.EncodeToString(hashInBytes)

	if base64str == sha1_sum {
		return true
	} else {
		return false
	}

}

func downloadFromUrl(url string, filename string) {
	fmt.Println("Downloading", url, "to", filename)

	output, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error while creating", filename, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}


func downloader(prefix string, download_url string, sha1_sum string) {
	if prefix == "" {
		prefix = "."
	}
	u, err := url.Parse(download_url)
	if err != nil {
		log.Fatal(err)
	}
	
	dirname := path.Join(prefix, filepath.Dir(u.Path))
	filename := filepath.Base(u.Path)
	fullname := path.Join(dirname , filename)


	if _, err := os.Stat(dirname); os.IsNotExist(err) {
	    os.MkdirAll(dirname, os.ModePerm)
	}
	
	if !chk_sha1(fullname, sha1_sum) {
		os.Remove(fullname)
		downloadFromUrl(download_url, fullname)
	}
}	
func make_new_url(host string, prefix string, download_url string) string {
	u, err := url.Parse(download_url)
	if err != nil {
		log.Fatal(err)
	}
	
	u.Path = path.Join(prefix, u.Path)
	u.Host = host
	return u.String()

}

func main() {

    // then config file settings
	prefix := "jenkins"
	host := "repo.dev.netis.com.cn"

    configFile, err := os.Open("json")
    if err != nil {
        print("opening config file", err.Error())
    }

    js, err := simplejson.NewFromReader(configFile)
    js.Set("connectionCheckUrl","http://www.baidu.com")

    js.Get("core").Set("url","http://repo.dev.netis.com.cn/download/war/2.17/jenkins.war")


	m, err := js.Get("plugins").Map()
	for k := range m {
		download_url, _ := js.Get("plugins").Get(k).Get("url").String()
		sha1_sum, _ := js.Get("plugins").Get(k).Get("sha1").String()
		downloader(prefix, download_url, sha1_sum)

		new_url := make_new_url(host, prefix, download_url)
		js.Get("plugins").Get(k).Set("url",new_url)
	}

    w, err := os.Create("./result.json")
    defer w.Close()

    o, _ := js.EncodePretty()
    w.Write(o)


}
