package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/anyuan-chen/hackathon/.gen/hackathon/public/model"
	"github.com/stretchr/testify/assert"
)

func waitUntilServiceReady(ctx context.Context) error {
	client := http.Client{}
	start := time.Now()
	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"http://localhost:8000/health",
			nil,
		)
		if err != nil {
			log.Println(fmt.Errorf("failed to create request: %w", err))
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(fmt.Errorf("failed to execute request %w", err))
			return err
		}

		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			log.Printf("the service took %d milliseconds to startup", time.Since(start).Milliseconds())
			return nil
		}
		resp.Body.Close()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(start) >= 5*time.Second {
				err := fmt.Errorf("timeout reached while waiting for endpoint")
				log.Println(err)
				return err
			}
			time.Sleep(500)
		}
	}
}

func startupService(ctx context.Context, t *testing.T) {
	go run(ctx)
	err := waitUntilServiceReady(ctx)
	if err != nil {
		log.Println(err.Error())
		t.FailNow()
	}
}

func printReqBody(r *http.Response) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("printing request body failed :(")
	}
	log.Println(string(body))
}

func LoginNonAdmin(ctx context.Context, t *testing.T) (person_id string, person_pw string, person_bearer_token string) {
	startupService(ctx, t)
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8000/login", nil)
	if err != nil {
		log.Println("did not create request successfully")
		t.FailNow()
	}
	id := "3"
	pass := "hi_eggy!"
	code := "Basic " + base64.StdEncoding.EncodeToString([]byte(id+":"+pass))
	req.Header.Add("Authorization", code)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("did not execute request successfully", err.Error())
		t.FailNow()
	}
	if resp.StatusCode != 200 {
		log.Println("unsuccessful req" + resp.Status)
		printReqBody(resp)
		t.FailNow()
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("malformed body", err.Error())
		t.FailNow()
	}
	// log.Println("raw token: ", body)
	type response struct {
		Access_token string `json:"access_token"`
	}
	var token response
	json.Unmarshal(body, &token)

	person_bearer_token = token.Access_token
	log.Println("person_bearer_token", person_bearer_token)
	return id, pass, person_bearer_token
}

func TestFetchUser(t *testing.T) {
	ctx := context.Background()
	client := http.Client{}
	id, _, bearer_token := LoginNonAdmin(ctx, t)
	assert.NotEqual(t, bearer_token, "")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/users/"+id, nil)
	if err != nil {
		log.Println("did not create request successfully")
		t.FailNow()
	}
	req.Header.Add("Authorization", "Bearer "+bearer_token)
	log.Println("bearer_token: ", bearer_token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("request not executed successfully", err.Error())
		t.FailNow()
	}
	var user model.Users
	body_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("body not read successfully")
		t.FailNow()
	}
	json.Unmarshal(body_bytes, &user)
	//assert user things here
	log.Println("bytes", string(body_bytes), user)
	assert.Equal(t, *user.Name, "Emily May")
	assert.Equal(t, *user.Company, "Graham Group")
	assert.Equal(t, *user.Email, "estradadana@example.org")
	assert.Equal(t, *user.Phone, "947.098.3138x493")
	assert.Equal(t, *user.Role, "hacker")
	assert.Nil(t, user.HashedSecret)
	assert.Nil(t, user.Salt)
}
