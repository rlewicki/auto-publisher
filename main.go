package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func get_hash(secret string, payload []byte) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write(payload)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received a new request...")
	defer r.Body.Close()
	body_content, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Failed to read request's body content")
		return
	}

	received_hash := r.Header.Get("X-Hub-Signature-256")[7:]
	secret_token := os.Getenv("AUTO_PUBLISHER_SECRET_TOKEN")
	local_hash := get_hash(secret_token, body_content)
	if local_hash != received_hash {
		fmt.Println("Hash mismatch, aborting request")
		return
	}

	fmt.Println("updating blog local repository...")
	blog_path := os.Getenv("BLOG_PATH")
	fmt.Println("path to blog local repository: ", blog_path)
	git_cmd := exec.Command("git", "pull")
	git_cmd.Dir = blog_path
	err = git_cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("publishing new version of the blog...")
	publish_cmd := exec.Command("./publish.sh")
	publish_cmd.Dir = blog_path
	err = publish_cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("new version of the blog released with success")
}

func main() {
	fmt.Println("auto-publisher bot started...")
	http.HandleFunc("/blog_update", handler)
	log.Fatal(http.ListenAndServe(":27015", nil))
}
