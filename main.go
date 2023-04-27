package main

import (
    "fmt"
	"log"
    "net/http"
	"crypto/hmac"
	"crypto/sha256"
	"os"
	"os/exec"
	"io"
)

func get_hash(secret string, payload []byte) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write(payload)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
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

	output, err := exec.Command("git -C ~/blog pull").Output()
	if err != nil {
		fmt.Println(err)
		return;
	}
	fmt.Println(output)
	output, err = exec.Command("~/blog/publish.sh").Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output)
	fmt.Println("updated the blog to the latest version")
}

func main() {
	fmt.Println("auto-publisher bot started...")
	http.HandleFunc("/blog_update", handler)
	log.Fatal(http.ListenAndServe(":27015", nil))
}