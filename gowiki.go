package main

import (
	"html/template"
	"net/http"
	"fmt"
	"time"
	"strconv"
	"strings"
	"regexp"
	"io/ioutil"
	"encoding/json"
)

type Page struct{
	Title string
	Ctitle string
	Uid int
	Static int
	Nuerons []neuron
	Synapse []int
	Timestamp time.Time
	TimestampDisplay string
}

func loadPage(title string) (*Page, error) {
	moaddr := MOAddr{"localhost:27017", "gowiki", "pantheon"}
	noaddr := MOAddr{"localhost:27017", "gowiki", "hades"}
	a := mongo_find(moaddr, title)
	c := []neuron{}
	for i := range a.Synapse{
		c = append(c, mongo_export(noaddr, a.Synapse[i]))
	}
	return &Page{title, a.Ctitle, a.Uid, 0, c, a.Synapse, a.Timestamp, a.TimestampDisplay}, nil
}

func loadResult(title string) (*Page, error) {
	noaddr := MOAddr{"localhost:27017", "gowiki", "olympus"}
	moaddr := MOAddr{"localhost:27017", "gowiki", "hades"}
	a := mongo_findRelate(noaddr, title)
	c := []neuron{}
	t := false
	for i := range a.Uids{
		if a.Uids[i] == 0 {}else{
			fmt.Println(c)
			c = append(c, mongo_export(moaddr, a.Uids[i]))
			t = true
		}
	}
	if t {
		return &Page{a.Title, a.Title, 1, 0, c, nil, time.Now(), ""}, nil
	} else {
		return &Page{a.Title, a.Title, 1, 0, nil, nil, time.Now(), ""}, nil
	}
	
}

func loadParadox(uid string) (*Page, error) {
	noaddr := MOAddr{"localhost:27017", "gowiki", "hades"}
	u, _ := strconv.Atoi(uid)
	a := mongo_locate(noaddr, u)
	fmt.Println(a)
	return &Page{"", a[0].Title, a[0].Uid, 0, a, a[0].Synapse, a[0].Timestamp, a[0].TimestampDisplay}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	cookie, _ := r.Cookie("gowiki")
	cookiea, _ := r.Cookie("gowiki-a")
	if CheckLogin(cookie, cookiea) {
		p, err := loadPage(title)
		if err != nil {
			http.Redirect(w, r, "/edit/" + title, http.StatusFound)
			return
		}
		renderTemplate(w, "view", p)
	}else{
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	cookie, _ := r.Cookie("gowiki")
	cookiea, _ := r.Cookie("gowiki-a")
	if CheckLogin(cookie, cookiea) {
		p, _ := loadPage(title)
		renderTemplate(w, "edit", p)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func editsmallHandler(w http.ResponseWriter, r *http.Request, uid string) {
	cookie, _ := r.Cookie("gowiki")
	cookiea, _ := r.Cookie("gowiki-a")
	if CheckLogin(cookie, cookiea) {
		p, _ := loadParadox(uid)
		renderTemplate(w, "editsmall", p)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func resultsHandler(w http.ResponseWriter, r *http.Request, title string){
	cookie, _ := r.Cookie("gowiki")
	cookiea, _ := r.Cookie("gowiki-a")
	if CheckLogin(cookie, cookiea) {
		fmt.Println(title)
		p, _ := loadResult(title)
		renderTemplate(w, "results", p)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("gowiki")
	renderTemplate(w, "login", p)
}

func loginAttempt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow- Origin", "*") 
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	var jsonData = []byte(body)
	var t Login
	json.Unmarshal(jsonData, &t)
	moaddr := MOAddr{"localhost:27017", "gowiki", "hermes"}
	f := mongo_login(moaddr, t.Username)
	if decrypt(f.Password) == t.Password{
		v := CreateSessionID()
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "gowiki", Value: v, Path: "/",Expires: expiration}
    http.SetCookie(w, &cookie)
		cookiea := http.Cookie{Name: "gowiki-a", Value: encrypt(t.Username), Path: "/",Expires: expiration}
    http.SetCookie(w, &cookiea)
		mongo_loginSuccess(moaddr, Login{t.Username, f.Password, v})
		http.Redirect(w, r, "/view/gowiki", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login/", http.StatusFound)
	}
}


var templates = template.Must(template.ParseFiles("edit.html", "view.html", "editsmall.html","results.html", "login.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|editsmall|view|results)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func CheckLogin(check *http.Cookie, checka *http.Cookie) bool{
	if check == nil || checka == nil {
		return false
	} else {
		moaddr := MOAddr{"localhost:27017", "gowiki", "hermes"}
		f := mongo_login(moaddr, decrypt(checka.Value))
		if f.Auth == check.Value{
			return true
		} else {
			return false
		}
		return true
	}
}

func Save(w http.ResponseWriter, r *http.Request) {
	//Allows the call to the RESTFUL API to come from across domains
	w.Header().Set("Access-Control-Allow- Origin", "*") 
	//Configures addr information 
	moaddr := MOAddr{"localhost:27017", "gowiki", "pantheon"}
	noaddr := MOAddr{"localhost:27017", "gowiki", "hades"}
	body, _ := ioutil.ReadAll(r.Body)
	var jsonData = []byte(body)
	var t neuron
	json.Unmarshal(jsonData, &t)
	f := mongo_find(moaddr, "LINDEX")
	mongo_init(moaddr, axion{"LINDEX", "", (f.Uid+1), 0, nil, time.Now(), ""})
	g := mongo_find(moaddr, t.Tags[0])
	g.Synapse = append(g.Synapse, (f.Uid+1))
	mongo_init(moaddr, axion{t.Tags[0], t.Ctitle,  g.Uid, 0, g.Synapse, time.Now(), ""})
	t.Uid = (f.Uid + 1)
	t.Synapse = append(t.Synapse, t.Uid)
	mongo_insert(noaddr, t)
}

func SmallSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow- Origin", "*") 
	noaddr := MOAddr{"localhost:27017", "gowiki", "hades"}
	body, _ := ioutil.ReadAll(r.Body)
	var jsonData = []byte(body)
	var t neuron
	json.Unmarshal(jsonData, &t)
	t.Timestamp = time.Now()
	t.TimestampDisplay = t.Timestamp.Format("Mon Jan _2 15:04:05 2006")
	mongo_update(noaddr, t)
}

func TagSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow- Origin", "*") 
	noaddr := MOAddr{"localhost:27017", "gowiki", "hades"}
	body, _ := ioutil.ReadAll(r.Body)
	var jsonData = []byte(body)
	var t neur
	json.Unmarshal(jsonData, &t)
	a := mongo_locateone(noaddr, t.Uid)
	a.Tags = append(a.Tags, t.Tags)
	mongo_update(noaddr, a)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow- Origin", "*") 
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	var jsonData = []byte(body)
	var t search
	json.Unmarshal(jsonData, &t)
	moaddr := MOAddr{"localhost:27017", "gowiki", "pantheon"}
	hoaddr := MOAddr{"localhost:27017", "gowiki", "hades"}
	noaddr := MOAddr{"localhost:27017", "gowiki", "olympus"}
	t.Searchterms = strings.ToLower(t.Searchterms)
	t.Searchables = strings.Split(t.Searchterms, " ")
	g := []axion{}
	h := []axion{}
	for i := range t.Searchables{
		h = mongo_seekfind(hoaddr, string(t.Searchables[i]))
		for j := range h {
			g = append(g, h[j])
		}
	}
	if len(g) == 0 {
		f := mongo_find(moaddr, "INDEX")
		j := CreateSessionID()
		mongo_init(moaddr, axion{"INDEX", "", (f.Uid+1), 0, nil, time.Now(), ""})
		k := axion{j, strings.ToLower(t.Searchterms), (f.Uid+1), 0, nil, time.Now(), ""}
		mongo_insertAxion(moaddr, k)
		js, _ := json.Marshal(k)
		w.Write(js)
	}else{
		j := CreateSessionID()
		c := make([]int, 1)
		for i := range g{
			for j := range g[i].Synapse{
				c = append(c, g[i].Synapse[j])	
			}
		}
		d := related{c, nil, t.Searchterms, j}
		mongo_insertRelate(noaddr, d)
		js, _ := json.Marshal(d)
		w.Write(js)
	}
}

func forwardHandler(w http.ResponseWriter, r *http.Request){
	cookie, _ := r.Cookie("gowiki")
	if cookie == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else{
		if cookie.Name == "true" {
			a := strings.Split(r.URL.String(), "/")
			if a[1] == "" {
				http.Redirect(w, r, "/view/" + "gowiki", http.StatusFound)
			}else{
				http.Redirect(w, r, "/view/" + a[1], http.StatusFound)	
			}
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}

}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/editsmall/", makeHandler(editsmallHandler))
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/process/", Save)
	http.HandleFunc("/tagprocess/", TagSave)
	http.HandleFunc("/subprocess/", SmallSave)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/loginAttempt/", loginAttempt)
	http.HandleFunc("/results/", makeHandler(resultsHandler))
	http.HandleFunc("/", forwardHandler)
	http.Handle("/css/",http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/fonts/",http.StripPrefix("/fonts/", http.FileServer(http.Dir("./fonts"))))
	http.Handle("/js/",http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/img/",http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	http.ListenAndServeTLS(":8085", "/etc/letsencrypt/live/wiki.rebirtharmitage.com/cert.pem", "/etc/letsencrypt/live/wiki.rebirtharmitage.com/privkey.pem", nil)
}