package hackers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/anyuan-chen/hackathon/.gen/hackathon/public/model"
	. "github.com/anyuan-chen/hackathon/.gen/hackathon/public/table"
	"github.com/anyuan-chen/hackathon/src/auth"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/gorilla/mux"
)

// Skill represents a skill with a name and rating
type InputSkill struct {
	Skill  *string `json:"skill"`
	Rating *int    `json:"rating"`
}

// Person represents a person with name, company, email, phone, and skills
type InputPerson struct {
	Name    *string       `json:"name"`
	Company *string       `json:"company"`
	Email   *string       `json:"email"`
	Phone   *string       `json:"phone"`
	Skills  *[]InputSkill `json:"skills"`
}

type UserWithSkills struct {
	Id      int            `json:"id"`
	Name    string         `json:"name"`
	Email   string         `json:"email"`
	Company string         `json:"company"`
	Phone   string         `json:"phone"`
	Skills  []model.Skills `json:"skills"`
}

func writeBody(r *http.Request, w http.ResponseWriter, resp interface{}) {
	json_resp, err := json.Marshal(resp)
	if err != nil {
		panic("coudln't serialize body")
	}
	w.Write(json_resp)
}

func databaseError(dest Statement, w http.ResponseWriter) {
	w.Write([]byte(dest.DebugSql() + " query failed"))
	w.WriteHeader(http.StatusInternalServerError)
}

func handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func handleLogin(db *sql.DB) http.HandlerFunc {
	type response struct {
		Access_token string `json:"access_token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id, pass, err := auth.GetIdPasswordFromRequest(r)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		stmt := Users.SELECT(Users.HashedSecret, Users.Salt).WHERE(Users.ID.EQ(Int32(id)))
		var users []model.Users
		err = stmt.Query(db, &users)
		if err != nil || len(users) != 1 {
			databaseError(stmt, w)
			return
		}
		user := users[0]
		hashed_scrt := auth.VerifyHashedPassword(pass, *user.Salt)
		if *user.HashedSecret != hashed_scrt {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("user hash: ", *user.HashedSecret)
			log.Println("provided hash ", hashed_scrt)

			w.Write([]byte("incorrect pw for user"))
			return
		}
		token, err := auth.GetBearerTokenByUserId(id, db)
		if err != nil {
			w.Write([]byte("no auth provided"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		/* do work here */
		resp := response{
			Access_token: *token.BearerToken,
		}
		log.Println("handleLogin", *token.BearerToken)
		var json_resp []byte
		json_resp, err = json.Marshal(resp)
		if err != nil {
			panic("bruh")
		}
		w.Write(json_resp)
	}
}

func handleGetAllUsers(db *sql.DB) http.HandlerFunc {
	type response struct {
		Users []UserWithSkills `json:"users"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var users []model.Users
		stmt := Users.SELECT(Users.ID, Users.Name, Users.Email, Users.Company, Users.Role, Users.Role, Users.Phone)
		err := stmt.Query(db, &users)
		if err != nil {
			databaseError(stmt, w)
			return
		}
		users_with_skills := make([]UserWithSkills, 0, len(users))
		skillsMap := make(map[int32]*UserWithSkills)
		for _, user := range users {
			userWithSkill := UserWithSkills{
				Id:      int(user.ID),
				Name:    *user.Name,
				Email:   *user.Email,
				Company: *user.Company,
				Phone:   *user.Phone,
				Skills:  make([]model.Skills, 0),
			}
			skillsMap[user.ID] = &userWithSkill
		}
		var skills []model.Skills
		allSkills := Skills.SELECT(Skills.AllColumns)
		err = allSkills.Query(db, &skills)
		if err != nil {
			databaseError(stmt, w)
			return
		}
		for _, skill := range skills {
			(*skillsMap[int32(*skill.UserID)]).Skills = append(skillsMap[int32(*skill.UserID)].Skills, skill)
		}
		for _, value := range skillsMap {
			users_with_skills = append(users_with_skills, *value)
		}
		resp := response{
			Users: users_with_skills,
		}
		writeBody(r, w, resp)
	}
}

func handleGetOneUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id not integer"))
			return
		}
		var users []model.Users
		stmt := Users.SELECT(Users.AllColumns).WHERE(Users.ID.EQ(Int32(int32(id))))
		err = stmt.Query(db, &users)
		if err != nil {
			databaseError(stmt, w)
			return
		} else if len(users) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(stmt.DebugSql()))
			return
		}
		user := users[0]

		var skills []model.Skills
		stmt = Skills.SELECT(Skills.AllColumns).WHERE(Skills.UserID.EQ(Int32(int32(id))))
		stmt.Query(db, &skills)
		if err != nil {
			databaseError(stmt, w)
			return
		}
		combinedUserSkills := UserWithSkills{
			Id:      int(user.ID),
			Name:    *user.Name,
			Email:   *user.Email,
			Company: *user.Company,
			Phone:   *user.Phone,
			Skills:  skills,
		}
		user.HashedSecret = nil
		user.Salt = nil
		writeBody(r, w, combinedUserSkills)
	}
}

func handleUpdateOneUser(db *sql.DB) http.HandlerFunc {
	type request struct {
		Name    *string         `json:"name"`
		Email   *string         `json:"email"`
		Company *string         `json:"company"`
		Phone   *string         `json:"phone"`
		Skills  *[]model.Skills `json:"skills"`
	}
	type response = model.Users
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id not integer"))
			return
		}
		var input request
		body_bytes, _ := io.ReadAll(r.Body)
		json.Unmarshal(body_bytes, &input)

		/* bad code, consequence of go-jet */
		var exisingUsers []model.Users
		stmt := Users.SELECT(Users.AllColumns).WHERE(Users.ID.EQ(Int32(int32(id))))
		stmt.Query(db, &exisingUsers)
		existingUser := exisingUsers[0]

		if input.Name != nil {
			existingUser.Name = input.Name
		}
		if input.Company != nil {
			existingUser.Company = input.Company
		}
		if input.Phone != nil {
			existingUser.Phone = input.Phone
		}
		if input.Email != nil {
			existingUser.Email = input.Email
		}
		var updated []model.Users
		updt := Users.UPDATE(Users.MutableColumns).MODEL(existingUser).WHERE(Users.ID.EQ(Int32(int32(id)))).RETURNING(Users.AllColumns)
		updt.Query(db, &updated)

		var finalSkills []model.Skills
		var existingSkills []model.Skills
		stmt = Skills.SELECT(Skills.AllColumns).WHERE(Skills.UserID.EQ(Int32(int32(id))))
		stmt.Query(db, &existingSkills)

		if input.Skills != nil {
			input_skills := *input.Skills
			for idx := range input_skills {
				int32id := int32(id)
				input_skills[idx].UserID = &int32id
			}
			del := Skills.DELETE().WHERE(Skills.UserID.EQ(Int32(int32(id))))
			del.Exec(db)
			ins := Skills.INSERT(Skills.AllColumns).MODELS(*input.Skills).RETURNING(Skills.AllColumns)
			ins.Query(db, &finalSkills)
		} else {
			finalSkills = existingSkills
		}

		combinedUserSkills := UserWithSkills{
			Id:      int(updated[0].ID),
			Name:    *updated[0].Name,
			Email:   *updated[0].Email,
			Company: *updated[0].Company,
			Phone:   *updated[0].Phone,
			Skills:  finalSkills,
		}
		writeBody(r, w, combinedUserSkills)
	}
}

func handleGetAllSkills(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min_freq_str := r.URL.Query().Get("min_freq")
		max_freq_str := r.URL.Query().Get("max_freq")
		var min_freq, max_freq int
		if min_freq_str != "" {
			mf, err := strconv.Atoi(min_freq_str)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("malformed min_freq query parameters"))
				return
			}
			min_freq = mf
		}
		if max_freq_str != "" {
			mf, err := strconv.Atoi(max_freq_str)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("malformed max_freq query parameters"))
				return
			}
			max_freq = mf
		}
		grouped_skills, err := db.Query("SELECT skills.skill, COUNT(*) AS skill_count FROM public.skills GROUP BY skills.skill;")
		if err != nil {
			databaseError(nil, w)
			return
		}
		type groupedSkill struct {
			Skill       string `json:"skill"`
			Skill_count int    `json:"skill_count"`
		}
		var res []groupedSkill
		for grouped_skills.Next() {
			var skill groupedSkill
			_ = grouped_skills.Scan(&skill.Skill, &skill.Skill_count)
			if min_freq_str != "" && skill.Skill_count < min_freq {
				continue
			}
			if max_freq_str != "" && skill.Skill_count > max_freq {
				continue
			}
			res = append(res, skill)
		}
		writeBody(r, w, res)
	}
}
