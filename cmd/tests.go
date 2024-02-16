package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
			"/health",
			nil,
		)
		if err != nil {
			print(fmt.Errorf("failed to create request: %w", err))
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			print(fmt.Errorf("failed to execute request %w", err))
			return err
		}

		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			print(fmt.Sprintf("the service took %d milliseconds to startup", time.Since(start).Milliseconds()))
			return nil
		}
		resp.Body.Close()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(start) >= 5*time.Second {
				err := fmt.Errorf("timeout reached while waiting for endpoint")
				print(err)
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
		t.FailNow()
	}
}

func LoginNonAdmin(ctx context.Context, t *testing.T) (person_id string, person_pw string, person_bearer_token string) {
	startupService(ctx, t)
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/login", nil)
	if err != nil {
		fmt.Printf("did not create request successfully")
		t.FailNow()
	}
	id := ""
	pass := ""
	code := "Basic " + base64.StdEncoding.EncodeToString([]byte(id+":"+pass))
	req.Header.Add("Authorization", code)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("did not execute request successfully")
		t.FailNow()
	}
	if resp.StatusCode != 200 {
		fmt.Printf("unsuccessful req" + resp.Status)
		t.FailNow()
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("malformed body")
		t.FailNow()
	}
	type response struct {
		Access_token string `json:"access_token"`
	}
	var token response
	json.Unmarshal(body, &token)

	person_bearer_token = token.Access_token
	return id, pass, person_bearer_token
}

func TestFetchUser(t *testing.T) {
	ctx := context.Background()
	client := http.Client{}
	id, _, bearer_token := LoginNonAdmin(ctx, t)
	assert.NotEqual(t, bearer_token, "")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/users/"+id, nil)
	if err != nil {
		fmt.Printf("did not create request successfully")
		t.FailNow()
	}
	req.Header.Add("Authorization", "Bearer "+bearer_token)
	resp, err := client.Do(req)
	if err != nil {

	}
	var user model.Users
	body_bytes, err := io.ReadAll(resp.Body)
	if err != nil {

	}
	json.Unmarshal(body_bytes, &user)
	//assert user things here
}
