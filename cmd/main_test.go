package main

import (
	"bytes"
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

func init() {
	go run(context.Background())
	err := waitUntilServiceReady(context.Background())
	if err != nil {
		log.Println("uh oh this is bad bad bad", err.Error())
	}
}

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

func LoginAs(ctx context.Context, id string, pass string, t *testing.T) (person_id string, person_pw string, person_bearer_token string) {
	// startupService(ctx, t)
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8000/login", nil)
	if err != nil {
		log.Println("did not create request successfully")
		t.FailNow()
	}
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
	type response struct {
		Access_token string `json:"access_token"`
	}
	var token response
	json.Unmarshal(body, &token)

	person_bearer_token = token.Access_token
	return id, pass, person_bearer_token
}

type UserWithSkills struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Email   string         `json:"email"`
	Company string         `json:"company"`
	Phone   string         `json:"phone"`
	Skills  []model.Skills `json:"skills"`
}

func TestGetSelfUser(t *testing.T) {
	ctx := context.Background()
	client := http.Client{}
	//make sure the user exists
	id := "3"
	pass := "hi_eggy!"
	_, _, bearer_token := LoginAs(ctx, id, pass, t)
	assert.NotEqual(t, bearer_token, "")

	//request for user
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/users/3", nil)
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
	//read response
	var user UserWithSkills
	body_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("body not read successfully")
		t.FailNow()
	}
	json.Unmarshal(body_bytes, &user)
	//assert user things here
	// log.Println("bytes", string(body_bytes), user)
	assert.Equal(t, user.Name, "Emily May")
	assert.Equal(t, user.Company, "Graham Group")
	assert.Equal(t, user.Email, "estradadana@example.org")
	assert.Equal(t, user.Phone, "947.098.3138x493")
	assert.NotEqual(t, len(user.Skills), 0)
}
func TestNoPermissionsFetchUser(t *testing.T) {
	ctx := context.Background()
	client := http.Client{}
	//make sure the user exists
	_, _, bearer_token := LoginAs(ctx, "3", "hi_eggy!", t)
	assert.NotEqual(t, bearer_token, "")

	//request for user
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/users/5", nil)
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
	//read response
	var user model.Users
	body_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("body not read successfully")
		t.FailNow()
	}
	json.Unmarshal(body_bytes, &user)
	assert.Equal(t, string(body_bytes), "only admins or owner of account may call this")
}

func TestUpdateUsers(t *testing.T) {
	newName := "James Su"
	newSkillName := "league of legends"
	newRating := int32(10)
	newSkills := make([]model.Skills, 0)
	newSkills = append(newSkills, model.Skills{
		Skill:  &newSkillName,
		Rating: &newRating,
	})

	ctx := context.Background()
	client := http.Client{}
	_, _, bearer_token := LoginAs(ctx, "3", "hi_eggy!", t)
	type request struct {
		Name    *string         `json:"name"`
		Email   *string         `json:"email"`
		Company *string         `json:"company"`
		Phone   *string         `json:"phone"`
		Skills  *[]model.Skills `json:"skills"`
	}
	req_body := request{
		Name:   &newName,
		Skills: &newSkills,
	}
	log.Println("should be league", *newSkills[0].Skill)

	req_body_json, err := json.Marshal(req_body)
	log.Println(req_body_json)
	assert.Nil(t, err)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "http://localhost:8000/users/3", bytes.NewBuffer(req_body_json))
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+bearer_token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("request not executed successfully", err.Error())
		t.FailNow()
	}
	//read response
	var user UserWithSkills
	body_bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	json.Unmarshal(body_bytes, &user)
	assert.Equal(t, user.Name, "James Su")
	assert.Equal(t, len(user.Skills), 1)
	assert.Equal(t, *user.Skills[0].Rating, int32(10))
	assert.Equal(t, *user.Skills[0].Skill, "league of legends")

	newName2 := "Emily May"
	newRating2 := int32(3)
	newSkillName2 := "Julia"
	newSkill := model.Skills{
		Rating: &newRating2,
		Skill:  &newSkillName2,
	}
	newSkills = make([]model.Skills, 0)
	newSkills = append(newSkills, newSkill)
	log.Println("should be juilia", *newSkills[0].Skill)
	req_body = request{
		Name:   &newName2,
		Skills: &newSkills,
	}
	req_body_json, err = json.Marshal(req_body)
	assert.Nil(t, err)
	req, err = http.NewRequestWithContext(ctx, http.MethodPut, "http://localhost:8000/users/3", bytes.NewBuffer(req_body_json))
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+bearer_token)
	resp, err = client.Do(req)
	if err != nil {
		log.Println("request not executed successfully", err.Error())
		t.FailNow()
	}
	//read response
	body_bytes, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)
	json.Unmarshal(body_bytes, &user)
	assert.Equal(t, user.Name, "Emily May")
	assert.Equal(t, len(user.Skills), 1)
	assert.Equal(t, *user.Skills[0].Rating, int32(3))
	assert.Equal(t, *user.Skills[0].Skill, "Julia")
}

func TestFetchAllUsers(t *testing.T) {
	ctx := context.Background()
	client := http.Client{}
	_, _, access_token := LoginAs(ctx, "666666", "root", t)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/users", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+access_token)
	resp, err := client.Do(req)
	assert.Equal(t, resp.Status, "200 OK")
	assert.Nil(t, err)
	resp_body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	type response struct {
		Users []UserWithSkills `json:"users"`
	}
	var users response
	json.Unmarshal(resp_body, &users)
	var t1, t2 bool
	for _, user := range users.Users {
		if user.ID == 3 {
			assert.Equal(t, user.Name, "Emily May")
			assert.Equal(t, user.Company, "Graham Group")
			assert.Equal(t, user.Email, "estradadana@example.org")
			assert.Equal(t, user.Phone, "947.098.3138x493")
			assert.NotEqual(t, len(user.Skills), 0)

			t1 = true
		} else if user.ID == 666666 {
			assert.Equal(t, user.Name, "Andrew Chen")
			assert.Equal(t, user.Company, "unemployed </3")
			assert.Equal(t, user.Email, "a22chen@uwaterloo.ca")
			assert.Equal(t, user.Phone, "9059059055")
			assert.Equal(t, len(user.Skills), 0)
			t2 = true
		}
	}
	assert.Greater(t, len(users.Users), 0)
	assert.True(t, t1)
	assert.True(t, t2)
}

func TestNoPermissionsFetchAllUsers(t *testing.T) {
	ctx := context.Background()
	client := http.Client{}
	_, _, access_token := LoginAs(ctx, "3", "hi_eggy!", t)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/users", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+access_token)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	resp_body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(resp_body), "only admins may call this")
}

func TestFetchAllSkills(t *testing.T) {
	type groupedSkill struct {
		Skill       string `json:"skill"`
		Skill_count int    `json:"skill_count"`
	}
	ctx := context.Background()
	client := http.Client{}
	_, _, access_token := LoginAs(ctx, "666666", "root", t)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/skills", nil)
	req.Header.Add("Authorization", "Bearer "+access_token)
	assert.Nil(t, err)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	var skills []groupedSkill
	resp_body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	json.Unmarshal(resp_body, &skills)
	assert.Equal(t, len(skills), 106)
	skillExists := make([]bool, 4)
	for _, skill := range skills {
		if skill.Skill == "Materialize" {
			assert.Equal(t, skill.Skill_count, 26)
			skillExists[0] = true
		}
		if skill.Skill == "Tailwind" {
			assert.Equal(t, skill.Skill_count, 32)
			skillExists[1] = true
		}
		if skill.Skill == "TypeScript" {
			assert.Equal(t, skill.Skill_count, 24)
			skillExists[2] = true
		}
		if skill.Skill == "Rust" {
			assert.Equal(t, skill.Skill_count, 34)
			skillExists[3] = true
		}
	}
	for _, exists := range skillExists {
		assert.True(t, exists)
	}
}

func TestFetchAllSkillsWithMinFreq(t *testing.T) {
	type groupedSkill struct {
		Skill       string `json:"skill"`
		Skill_count int    `json:"skill_count"`
	}
	ctx := context.Background()
	client := http.Client{}
	_, _, access_token := LoginAs(ctx, "666666", "root", t)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/skills?min_freq=30", nil)
	req.Header.Add("Authorization", "Bearer "+access_token)
	assert.Nil(t, err)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	var skills []groupedSkill
	resp_body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	json.Unmarshal(resp_body, &skills)
	skillExists := make([]bool, 2)
	for _, skill := range skills {
		if skill.Skill == "Materialize" {
			t.FailNow()
		}
		if skill.Skill == "Tailwind" {
			assert.Equal(t, skill.Skill_count, 32)
			skillExists[0] = true
		}
		if skill.Skill == "TypeScript" {
			t.FailNow()
		}
		if skill.Skill == "Rust" {
			assert.Equal(t, skill.Skill_count, 34)
			skillExists[1] = true
		}
	}
	for _, exists := range skillExists {
		assert.True(t, exists)
	}
}

func TestFetchAllSkillsWithMaxFreq(t *testing.T) {
	type groupedSkill struct {
		Skill       string `json:"skill"`
		Skill_count int    `json:"skill_count"`
	}
	ctx := context.Background()
	client := http.Client{}
	_, _, access_token := LoginAs(ctx, "666666", "root", t)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/skills?max_freq=30", nil)
	req.Header.Add("Authorization", "Bearer "+access_token)
	assert.Nil(t, err)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	var skills []groupedSkill
	resp_body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	json.Unmarshal(resp_body, &skills)
	skillExists := make([]bool, 2)
	for _, skill := range skills {
		if skill.Skill == "Materialize" {
			assert.Equal(t, skill.Skill_count, 26)
			skillExists[0] = true
		}
		if skill.Skill == "Tailwind" {
			t.FailNow()
		}
		if skill.Skill == "TypeScript" {
			assert.Equal(t, skill.Skill_count, 24)
			skillExists[1] = true
		}
		if skill.Skill == "Rust" {
			t.FailNow()
		}
	}
	for _, exists := range skillExists {
		assert.True(t, exists)
	}
}
