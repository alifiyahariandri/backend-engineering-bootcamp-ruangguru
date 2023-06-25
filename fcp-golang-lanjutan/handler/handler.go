package handler

import (
	"a21hc3NpZ25tZW50/client"
	"a21hc3NpZ25tZW50/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GetProgramByCode(user model.User) string {
	file, err := os.Open("data/list-study.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	content, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		panic(err)
	}

	listprog := strings.Split(string(content), "\n")

	for _, prog := range listprog {
		splitprog := strings.Split(prog, "_")
		if user.StudyCode == splitprog[0] {
			return splitprog[0]
		}
	}

	return ""
}

func GetByUserModel(user model.User) []string {
	file, err := os.Open("data/users.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	content, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		panic(err)
	}

	if string(content) == "" {
		baru := []string{}

		return baru
	}

	listuser := strings.Split(string(content), "\n")

	for _, perUser := range listuser {
		splitprog := strings.Split(perUser, "_")
		if user.ID == splitprog[0] {
			return splitprog
		}
	}

	baru := []string{}

	return baru

}

func GetByIdAndName(user model.UserLogin) []string {
	file, err := os.Open("data/users.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	content, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		panic(err)
	}

	if string(content) == "" {
		baru := []string{}

		return baru
	}

	listuser := strings.Split(string(content), "\n")

	for _, perUser := range listuser {
		splitprog := strings.Split(perUser, "_")
		if user.ID == splitprog[0] && user.Name == splitprog[1] {
			return splitprog
		}
	}

	baru := []string{}

	return baru

}

func GetById(id string) []string {
	file, err := os.Open("data/users.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	content, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		panic(err)
	}

	if string(content) == "" {
		baru := []string{}
		return baru
	}

	listprog := strings.Split(string(content), "\n")

	for _, prog := range listprog {
		splitprog := strings.Split(prog, "_")
		if id == splitprog[0] {
			return splitprog
		}
	}

	baru := []string{}
	return baru

}

var UserLogin = make(map[string]model.User)

// DESC: func Auth is a middleware to check user login id, only user that already login can pass this middleware
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("user_login_id")
		if err != nil {
			fmt.Println("AAAAAAAAAAAAAAAAAAAAa")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
			return
		}

		if _, ok := UserLogin[c.Value]; !ok || c.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user login id not found"})
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", c.Value)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DESC: func AuthAdmin is a middleware to check user login role, only admin can pass this middleware
func AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // your code here }) // TODO: replace this
		c, _ := r.Cookie("user_login_role")

		if c.Value != "admin" {
			w.WriteHeader(401)
			msg := model.ErrorResponse{"user login role not Admin"}
			msgresp, err := json.Marshal(msg)

			if err != nil {
				panic(err)
			}
			w.Write(msgresp)
		}
		next.ServeHTTP(w, r)
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	body, err := io.ReadAll(r.Body)
	fmt.Println(string(body))
	if err != nil {
		panic(err)
	}

	var user model.UserLogin

	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		panic(err)
	}

	if user.ID == "" || user.Name == "" {
		w.WriteHeader(400)
		msg := model.ErrorResponse{"ID or name is empty"}
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)
	} else if len(GetByIdAndName(user)) == 0 {
		w.WriteHeader(400)
		msg := model.ErrorResponse{"user not found"}
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)
	} else {
		var userAsli model.User
		getuser := GetById(user.ID)
		userAsli.ID = getuser[0]
		userAsli.Name = getuser[1]
		userAsli.StudyCode = getuser[2]
		userAsli.Role = getuser[3]

		UserLogin[user.ID] = userAsli

		cookie := &http.Cookie{
			Name:  "user_login_id",
			Value: user.ID,
		}
		http.SetCookie(w, cookie)

		cookie2 := &http.Cookie{
			Name:  "user_login_role",
			Value: userAsli.Role,
		}
		http.SetCookie(w, cookie2)

		w.WriteHeader(200)
		msg := model.SuccessResponse{}
		msg.Username = user.Name
		msg.Message = "login success"
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	fmt.Println(string(body))
	if err != nil {
		panic(err)
	}

	var user model.User

	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		panic(err)
	}

	fmt.Println(GetByUserModel(user))

	if user.ID == "" || user.Name == "" || user.StudyCode == "" || user.Role == "" {

		w.WriteHeader(400)
		msg := model.ErrorResponse{"ID, name, study code or role is empty"}
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)

	} else if user.Role != "admin" && user.Role != "user" {

		w.WriteHeader(400)
		msg := model.ErrorResponse{"role must be admin or user"}
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)

	} else if len(GetByUserModel(user)) != 0 {
		w.WriteHeader(400)
		msg := model.ErrorResponse{"user id already exist"}
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)
	} else {
		w.WriteHeader(200)
		msg := model.SuccessResponse{}
		msg.Username = user.Name
		msg.Message = "register success"
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)

		data := user.ID + "_" + user.Name + "_" + user.StudyCode + "_" + user.Role + "\n"

		file, err := os.OpenFile("data/users.txt", os.O_APPEND, 0644)
		if _, err := file.Write([]byte(data)); err != nil {
			file.Close()
		}

	}
}

func Logout(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(string)

	if userID == "" {
		w.WriteHeader(401)
		msg := model.ErrorResponse{"user login id not found"}
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)
	} else {
		w.WriteHeader(200)
		msg := model.SuccessResponse{}
		msg.Username = userID
		msg.Message = "logout success"
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)

		http.SetCookie(w, &http.Cookie{
			Name:  "user_login_id",
			Value: "",
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "user_login_role",
			Value: "",
		})

		for k := range UserLogin {
			delete(UserLogin, k)
		}
	}
	// TODO: answer here
}

func GetStudyProgram(w http.ResponseWriter, r *http.Request) {
	// list study program
	// TODO: answer here
	file, err := os.Open("data/list-study.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	content, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		panic(err)
	}

	listprog := strings.Split(string(content), "\n")
	var programs []model.StudyData

	for _, prog := range listprog {
		splitprog := strings.Split(prog, "_")
		program := model.StudyData{}
		program.Code = splitprog[0]
		program.Name = strings.Trim(splitprog[1], "\r")
		programs = append(programs, program)
	}

	var jsonData, err3 = json.Marshal(programs)
	if err3 != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

func AddUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	fmt.Println(string(body))
	if err != nil {
		panic(err)
	}

	var user model.User

	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		panic(err)
	}

	if GetProgramByCode(user) == "" {

		w.WriteHeader(400)
		msg := model.ErrorResponse{"study code not found"}
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)
	} else {

		w.WriteHeader(200)
		msg := model.SuccessResponse{}
		msg.Username = user.Name
		msg.Message = "add user success"
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)

		data := user.ID + "_" + user.Name + "_" + user.StudyCode + "_" + user.Role + "\n"

		file, err := os.OpenFile("data/users.txt", os.O_APPEND, 0644)
		if _, err := file.Write([]byte(data)); err != nil {
			file.Close()
		}
	}
	// TODO: answer here
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	if r.URL.Query().Get("id") == "" {
		w.WriteHeader(400)
		msg := model.ErrorResponse{"user id is empty"}
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)
	} else {
		w.WriteHeader(200)
		msg := model.SuccessResponse{}
		msg.Username = r.URL.Query().Get("id")
		msg.Message = "delete success"
		msgresp, err := json.Marshal(msg)

		if err != nil {
			panic(err)
		}
		w.Write(msgresp)
	}
}

func GetWeatherByRegion(region string, ch chan model.MainWeather) (model.MainWeather, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://rg-weather-api.fly.dev/weather?region="+region, nil)
	if err != nil {
		return model.MainWeather{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return model.MainWeather{}, err
	}

	defer resp.Body.Close()

	var weather model.MainWeather

	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return model.MainWeather{}, err
	}
	ch <- weather
	return weather, nil
}

// DESC: Gunakan variable ini sebagai goroutine di handler GetWeather
var GetWetherByRegionAPI = client.GetWeatherByRegion
var GetWetherByRegionAPI2 = GetWeatherByRegion

func GetWeather(w http.ResponseWriter, r *http.Request) {

	var listRegion = []string{"jakarta", "bandung", "surabaya", "yogyakarta", "medan", "makassar", "manado", "palembang", "semarang", "bali"}
	ch := make(chan model.MainWeather, len(listRegion))

	for _, region := range listRegion {
		go GetWetherByRegionAPI2(region, ch)
	}

	baru := []model.MainWeather{}

	for i := 0; i < len(listRegion); i++ {
		result := <-ch
		baru = append(baru, result)
	}

	jsonData, err := json.Marshal(baru)
	if err != nil {

	}
	w.Write(jsonData)
	// DESC: dapatkan data weather dari 10 data di atas menggunakan goroutine
	// TODO: answer here
}
