package control

import (
	"../config"
	"../model"
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	// "github.com/bpowers/seshcookie"
	"encoding/json"
	"io/ioutil"
	"html/template"
	"net/http"
	"strings"
	// "net/url"
	"github.com/russross/blackfriday"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Dashboard")

	renderTemplate("home", w, r, func(data TemplateData) TemplateData {
		data.LanguageCompletion = model.GetLanguageCompletion()
		data.Issues, data.NumIssues = GetGithubIssues()
		data.WebsiteIssues, data.NumWebsiteIssues = GetWebsiteIssues()
		data.TranslatorIssues, data.NumTranslatorIssues = GetTranslatorIssues()
		return data
	})
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		name := r.FormValue("name")
		language := r.FormValue("language")
		user := &model.User{
			Email:    email,
			Name:     name,
			Language: language,
			Password: "",
			Secret:   "",
		}
		user.Save()

		http.Redirect(w, r, "/users", 303)
	} else {
		renderTemplate("users", w, r, func(data TemplateData) TemplateData {
			data.Users = model.GetUsers()
			data.UsersByLanguage = make(map[string][]*model.User, len(data.Languages))
			for _, user := range data.Users {
				data.UsersByLanguage[user.Language] = append(data.UsersByLanguage[user.Language], user)
			}
			return data
		})
	}
}

func UsersAddHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("users_add", w, r, nil)
}

func UsersMasqueradeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("users_masq", w, r, func(data TemplateData) TemplateData {
		// data.User = user
		return data
	})
}

func UsersDelHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := GetCurrentUser(r)
	if !currentUser.IsAdmin {
		http.Redirect(w, r, "/users", 303)
		return
	}

	email := r.FormValue("user")
	user := model.GetUserByEmail(email)
	if user == nil {
		http.Redirect(w, r, "/users", 303)
		return
	}

	gonow := r.FormValue("go")
	if r.Method == "POST" && gonow == "yes" {
		user.Delete()
		http.Redirect(w, r, "/users", 303)
		return
	} else {
		renderTemplate("users_del", w, r, func(data TemplateData) TemplateData {
			data.User = user
			return data
		})
	}
}

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)

	if r.Method == "POST" {
		user.Name = r.FormValue("name")
		language := r.FormValue("language")
		if language != "" {
			user.Language = language
		}
		user.Save()

		http.Redirect(w, r, "/home", 303)
	} else {
		renderTemplate("account", w, r, func(data TemplateData) TemplateData {
			return data
		})
	}
}

func AccountReclaimHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("account_reclaim", w, r, nil)
}

func SetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)

	if r.Method == "POST" {
		password := r.FormValue("password")
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err == nil {
			user.Password = string(hash)
			user.Save()
		}
		http.Redirect(w, r, "/account", 303)
	} else {
		renderTemplate("account_set_password", w, r, func(data TemplateData) TemplateData {
			return data
		})
	}
}

type Issue struct {
	Number          int    `json:"number"`
	Name            string `json:"title"`
	SummaryMarkdown string `json:"body"`
	SummaryHTML     template.HTML
	URL             string `json:"url"`
	CssClass        string
	Avatar          string
	User            struct{
		Avatar string `json:"avatar_url"`
	} `json:"user"`
	Labels          []struct {
		URL   string `json:"url"`
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"labels"`
}

// type issueLabel struct {
// 	URL   string `json:"url"`
// 	Name  string `json:"name"`
// 	Color string `json:"color"`
// }

func GetGithubIssues() ([]Issue, int) {
	issues, num := getGithubAPIIssues("charactersheets")
	return issues, num
}

func GetWebsiteIssues() ([]Issue, int) {
	issues, num := getGithubAPIIssues("charactersheets-website")
	return issues, num
}

func GetTranslatorIssues() ([]Issue, int) {
	issues, num := getGithubAPIIssues("charactersheets-translator")
	return issues, num
}

func getGithubAPIIssues(repo string) ([]Issue, int) {
	resp, err := http.Get("https://api.github.com/repos/marcusatbang/"+repo+"/issues?state=open&sort=updated&access_token="+config.Config.Github.AccessToken)
	if err != nil {
		fmt.Println("Error fetching issues from GitHub:", err)
		return []Issue{}, 0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading issues from GitHub:", err)
		return []Issue{}, 0
	}

	// fmt.Println(string(body))
	issues := make([]Issue, 0)
	err = json.Unmarshal(body, &issues)
	if err != nil {
		fmt.Println("Error decoding issues from GitHub:", err)
		return []Issue{}, 0
	}
	numIssues := len(issues)
	if len(issues) > 30 {
		issues = issues[0:30]
	}

	for i, issue := range issues {
		issues[i].URL = strings.Replace(issue.URL, "https://api.github.com/repos/", "https://www.github.com/", 1)

		if issue.SummaryMarkdown != "" {
			html := blackfriday.MarkdownCommon([]byte(issue.SummaryMarkdown))
			if html != nil {
				issues[i].SummaryHTML = template.HTML(html)
			}
		}
		// if issue.SummaryMarkdown != "" {
		// 	fmt.Println("Parsing Markdown:", issue.SummaryMarkdown)
		// 	resp, err = http.PostForm("https://api.github.com/markdown", url.Values{"text": {issue.SummaryMarkdown}})
		// 	if err != nil {
		// 		fmt.Println("Error parsing Markdown:", err)
		// 	} else {
		// 		html, err := ioutil.ReadAll(resp.Body)
		// 		resp.Body.Close()
		// 		if err != nil {
		// 			fmt.Println("Error parsing Markdown:", err)
		// 		} else {
		// 			issues[i].SummaryHTML = string(html)
		// 			fmt.Println("Parsed into HTML:", issues[i].SummaryHTML)
		// 		}
		// 	}
		// }

		for _, label := range issue.Labels {
			fmt.Println("Located label:", label.Name)
			if label.Name == "bug" {
				issues[i].CssClass = "danger"
			} else if label.Name == "enhancement" {
				issues[i].CssClass = "success"
			}
		}
	}

	fmt.Println("Loaded", len(issues), "issues from GitHub")
	return issues, numIssues
}
