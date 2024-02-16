package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/anyuan-chen/hackathon/.gen/hackathon/public/model"
	. "github.com/anyuan-chen/hackathon/.gen/hackathon/public/table"
	"github.com/anyuan-chen/hackathon/src/auth"
	_ "github.com/lib/pq"
	// "github.com/anyuan-chen/hackathon/src/auth"
	// . "github.com/go-jet/jet/v2/postgres"
	// _ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

// Skill represents a skill with a name and rating
type Skill struct {
	Skill  string `json:"skill"`
	Rating int    `json:"rating"`
}

// Person represents a person with name, company, email, phone, and skills
type Person struct {
	Name    string  `json:"name"`
	Company string  `json:"company"`
	Email   string  `json:"email"`
	Phone   string  `json:"phone"`
	Skills  []Skill `json:"skills"`
}

func main() {
	//Open the JSON file
	file, err := os.Open("data.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	var people []Person
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&people)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	db, err := sql.Open("postgres", "postgresql://anyuan-chen:TiCM34meZERy@ep-round-river-a5two8jo.us-east-2.aws.neon.tech/hackathon?sslmode=require")
	if err != nil {
		fmt.Print(err, "db connect error")
		return
	}
	defer db.Close()

	skill_idx := 1
	for idx, person := range people {
		var user []model.Users
		password := "hi_eggy!"
		salt, hashed_pw := auth.GetHashedPassword(password)

		user_insert := Users.INSERT(Users.ID, Users.Name, Users.Company, Users.Email, Users.Phone, Users.Role, Users.Salt, Users.HashedSecret).VALUES(idx, person.Name, person.Company, person.Email, person.Phone, "hacker", salt, hashed_pw)

		err := user_insert.Query(db, &user)
		if err != nil {
			print("user", err.Error())
		}

		for _, skill := range person.Skills {
			skill_insert := Skills.INSERT(Skills.ID, Skills.Rating, Skills.Skill).VALUES(skill_idx, skill.Rating, skill.Skill)
			var skill []model.Skills
			err := skill_insert.Query(db, &skill)
			if err != nil {
				print("skill", err.Error())
				return
			}
			skill_idx++
		}
	}

	_, err = db.Exec("INSERT INTO test (id) VALUES (1);")
	if err != nil {
		print(err.Error())
	}

	// var t []model.Test
	// test := model.Test{
	// 	ID: 4,
	// }
	// test_insert := Test.INSERT(Test.ID).MODEL(test)
	// err = test_insert.Query(db, &t)
	// print(test_insert.Sql())
	// if err != nil {
	// 	print("test", test_insert.DebugSql(), err.Error())
	// }

	//generate an admin account for demo purposes
	password := "root"
	salt, hashed_pw := auth.GetHashedPassword(password)
	admin_insert := Users.INSERT(Users.ID, Users.Name, Users.Company, Users.Email, Users.Phone, Users.Role, Users.Salt, Users.HashedSecret).VALUES(666666, "Andrew Chen", "unemployed </3", "a22chen@uwaterloo.ca", "9059059055", "admin", salt, hashed_pw)
	var user []model.Users
	err = admin_insert.Query(db, &user)
	print(admin_insert.DebugSql())
	if err != nil {
		print("admin", err.Error())
		return
	}

}
