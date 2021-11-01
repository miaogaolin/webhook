package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)


type WatchItem struct {
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	Script string `json:"script"`
	Secret string `json:"secret"`
}

type Config struct {
	BindHost string      `json:"bind"`
	Items    []WatchItem `json:"items"`
}


type Repository struct {
	Url         string `json:"url"`
	AbsoluteUrl string `json:"absolute_url"`
}

type Commit struct {
	Branch string `json:"branch"`
}

type Payload struct {
	Ref      string     `json:"ref"` // "refs/heads/develop"
	Repo     Repository `json:"repository"`
	CanonUrl string     `json:"canon_url"`
	Commits  []Commit   `json:"commits"`
	Secret string `json:"secret"`
	body []byte
}


var cfg Config



func runScript(item *WatchItem) (err error) {
	script := "./" + item.Script
	out, err := exec.Command("bash", "-c", script).Output()
	if err != nil {
		log.Printf("Exec command failed: %s\n", err)
	}

	log.Printf("Run %s output: %s\n", script, string(out))
	return
}

func handleGithub(event Payload, cfg *Config) (err error) {
	
	for _, item := range cfg.Items {
		if item.Secret != "" &&
			event.Secret != githubSecret(event.body, []byte(item.Secret)) {
			log.Println("validate secret failed")
			break
		}

		if event.Repo.Url == item.Repo && strings.Contains(event.Ref, item.Branch) {
			err = runScript(&item)
			if err != nil {
				log.Printf("run script error: %s\n", err)
			}
			break
		}
	}
	return
}

func githubSecret(data, secret []byte) string {
	return "sha256=" + NewSha256(data, secret)
}


func handle(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("read request failed:%s\n", err)
		return
	}

	var event Payload
	err = json.Unmarshal(body, &event)
	if err != nil {
		log.Printf("payload json decode failed: %s\n", err)
		return
	}

	event.Secret = req.Header.Get("X-Hub-Signature-256")
	event.body = body

	handleGithub(event, &cfg)
}


func main() {

	if len(os.Args) < 2 {
		log.Println("Usage: webhook <ConfigFile>")
		return
	}

	cfgbuf, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Println("Read config file failed:", err)
		return
	}

	err = json.Unmarshal(cfgbuf, &cfg)
	if err != nil {
		log.Println("Unmarshal config failed:", err)
		return
	}

	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(cfg.BindHost, nil))
}

