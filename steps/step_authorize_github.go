package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/mitchellh/multistep"
	"io/ioutil"
	"log"
	"log"
	"net/http"
	"os"
	"strings"
)

type stepAuthorizeGithub struct{}

type AuthorizeResponse struct {
	Token string `json:token`
}

// The authorization sequence, required for someone without a ~/.getconfig
func (*stepAuthorizeGithub) Run(state map[string]interface{}) multistep.StepAction {
	log.Println("Begin authorization sequence...")

	// Let the user know what will happen.
	fmt.Println(`Your username and password will be used once to obtain a unique
authorization token from GitHub's API, which will be stored in
~/.getconfig.
`)

	// Asks the user for credentials, and then makes a request to the
	// GitHub API to get an authorization token to store in ~/.getconfig
	log.Println("Asking user for credentials...")

	var username string

	fmt.Printf("Please enter your GitHub username: ")
	_, err := fmt.Scanf("%s", &username)

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Please enter your GitHub password: ")

	p := gopass.GetPasswd()
	fmt.Printf("\n")

	password := string(p)

	client := &http.Client{}

	reqBody := strings.NewReader(`
        {"scopes":["repo"],
        "note":"get command line client",
        "note_url": "https://github.com/pearkes/get"}`)

	req, err := http.NewRequest("POST",
		"https://api.github.com/authorizations", reqBody)

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	dec := json.NewDecoder(bytes.NewReader(body))

	log.Println(dec)

	auth := AuthorizeResponse{}

	dec.Decode(&auth)

	log.Println(resp.Status)

	if resp.StatusCode != 201 {
		fmt.Println(red + "Uh oh, there was an error authenticating with GitHub. Here's what we got back:\n" + clear)
		fmt.Println(string(body))
		return multistep.ActionHalt
	} else {
		fmt.Println(green + "Succesfully authenticated with Github." + clear)
	}

	log.Println(string(body))

	// Set the discovered credentials into the bag of state
	state["token"] = auth.Token
	state["username"] = username

	log.Println(username)

	return multistep.ActionContinue
}

func (*stepAuthorizeGithub) Cleanup(map[string]interface{}) {}
