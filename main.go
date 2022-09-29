package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
)

func main() {
	log.Println("server started")
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func handleWebhook(w http.ResponseWriter, r *http.Request) {
// 	fmt.Printf("headers: %v\n", r.Header)

// 	_, err := io.Copy(os.Stdout, r.Body)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// }

// func handleWebhook(w http.ResponseWriter, r *http.Request) {
// 	webhookData := make(map[string]interface{})
// 	err := json.NewDecoder(r.Body).Decode(&webhookData)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	fmt.Println("got webhook payload: ")
// 	for k, v := range webhookData {
// 		fmt.Printf("%s : %v\n", k, v)
// 	}
// }

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// payload, err := github.ValidatePayload(r, []byte("my-secret-key"))
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {
	case *github.PushEvent:
		// this is a commit push, do something with it
		fmt.Println("Push Event Triggered")
	case *github.PullRequestEvent:
		// this is a pull request, do something with it
		fmt.Println("Pull Request Event Triggered")
	case *github.WatchEvent:
		// https://developer.github.com/v3/activity/events/types/#watchevent
		// someone starred our repository
		fmt.Println("Watch Event Triggered")
		if e.Action != nil && *e.Action == "starred" {
			fmt.Printf("%s starred repository %s\n",
				*e.Sender.Login, *e.Repo.FullName)
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}
