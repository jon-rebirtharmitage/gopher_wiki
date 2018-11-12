package main

import (
	"html/template"
	"math/rand"
	"time"
)

type MOAddr struct {
  session, table, doc string
}

type MOValue struct {
  Title, Body string
}

type MOGValue struct {
	Title string
}

type neuron struct{
  Uid int								`json:"uid"` 
	ContentType int				`json:"contenttype"` 
  Title string					`json:"title"` 
	Ctitle string					`json:"ctitle"` 
	Parent string					`json:"parent"` 
	Content template.HTML `json:"content"` 
  Tags []string					`json:"tags"` 
  Synapse []int					`json:"synapse"` 
	Timestamp time.Time			`json:"timestamp"`
	TimestampDisplay string	`json:"timestampdisplay"`
}

type neur struct{
	Uid int								`json:"uid"` 
	Title string					`json:"title"` 
	Tags string						`json:"tag"`
	Synapse int						`json:"synapse"`
}

type axion struct{
  Title string			`json:"title"` 
	Ctitle string			`json:"ctitle"`
	Uid int						`json:"uid"`
	Static int				`json:"static"`
  Synapse []int			`json:"synapse"`
	Timestamp time.Time 	`json:"timestamp"`
	TimestampDisplay string	`json:"timestampdisplay"`
}	

type relate struct{
	Uids int		`json:"uids"`
}

type related struct{
	Uids []int		`json:"uids"`
	Links []string	`json:"links"`
	Title string	`json:"title"`
	Uid string		`json:"uid"`
}

type search struct{
	Searchterms string		`json:"searchterms"`
	Searchables []string	`json:"searchables"`
}

type view struct{
	Viewit string
}

func CreateSessionID() (string){
	source := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 24; i++{
		s = s + string(source[rand.Intn(50)])
	}
	return s
}

type Login struct {
	Username string			`json:"username"` 
	Password string			`json:"password"` 
	Auth string					`json:"auth"` 
}

type Cookie struct {
		Name       string
		Value      string
		Path       string
		Domain     string
		Expires    time.Time
		RawExpires string

// MaxAge=0 means no 'Max-Age' attribute specified.
// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// MaxAge>0 means Max-Age attribute present and given in seconds
		MaxAge   int
		Secure   bool
		HttpOnly bool
		Raw      string
		Unparsed []string // Raw text of unparsed attribute-value pairs
}
