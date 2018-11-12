package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func mongo_insert(moaddress MOAddr, t neuron) bool{
  session, err := mgo.Dial(moaddress.session)
  if err != nil {
    return false
  }
  defer session.Close()
  c := session.DB(moaddress.table).C(moaddress.doc)
  c.Insert(t)
	if err != nil{
		return false
	}
	return true
}

func mongo_insertAxion(moaddress MOAddr, t axion) bool{
  session, err := mgo.Dial(moaddress.session)
  if err != nil {
    return false
  }
  defer session.Close()
  c := session.DB(moaddress.table).C(moaddress.doc)
  c.Insert(t)
	if err != nil{
		return false
	}
	return true
}

func mongo_insertRelate(moaddress MOAddr, t related) bool{
  session, err := mgo.Dial(moaddress.session)
  if err != nil {
    return false
  }
  defer session.Close()
  c := session.DB(moaddress.table).C(moaddress.doc)
  c.Insert(t)
	if err != nil{
		return false
	}
	return true
}

func mongo_init(moaddress MOAddr, t axion) bool{
  session, err := mgo.Dial(moaddress.session)
  if err != nil {
    return false
  }
  defer session.Close()
  c := session.DB(moaddress.table).C(moaddress.doc)
  c.Update(bson.M{"title": t.Title}, t)
	if err != nil{
		return false
	}
	return true
}

func mongo_export(moaddress MOAddr, t int) neuron {
	session, err := mgo.Dial(moaddress.session)
	if err != nil {
			panic(err)
	}
	defer session.Close()
	c := session.DB(moaddress.table).C(moaddress.doc)
	result := neuron{}
	erro := c.Find(bson.M{"uid": t}).One(&result)
	if erro != nil {
		log.Fatal(erro)
	}
	return result
}

func mongo_multiexport(moaddress MOAddr, t int) []neuron {
	session, err := mgo.Dial(moaddress.session)
	if err != nil {
			panic(err)
	}
	defer session.Close()
	c := session.DB(moaddress.table).C(moaddress.doc)
	result := []neuron{}
	iter := c.Find(bson.M{"uid": t}).Limit(1000).Iter()
	err = iter.All(&result)
	if err != nil {
			log.Fatal(err)
	}
	return result
}

func mongo_find(moaddress MOAddr, f string) axion {
	session, err := mgo.Dial(moaddress.session)
	if err != nil {
			panic(err)
	}
	defer session.Close()
	c := session.DB(moaddress.table).C(moaddress.doc)
	result := axion{}
	erro := c.Find(bson.M{"title": f}).One(&result)
	if erro != nil {
			return result
	}
	return result
}

func mongo_truefind(moaddress MOAddr, f string) axion {
	session, err := mgo.Dial(moaddress.session)
	if err != nil {
			panic(err)
	}
	defer session.Close()
	c := session.DB(moaddress.table).C(moaddress.doc)
	result := axion{}
	erro := c.Find(bson.M{"ctitle": f}).One(&result)
	if erro != nil {
			return result
	}
	return result
}

func mongo_seekfind(moaddress MOAddr, f string) []axion {
	session, err := mgo.Dial(moaddress.session)
	if err != nil {
			panic(err)
	}
	defer session.Close()
	c := session.DB(moaddress.table).C(moaddress.doc)
	result := []axion{}
	erro := c.Find(bson.M{"ctitle": bson.RegEx{f, ""}}).All(&result)
	if erro != nil {
			return result
	}
	return result
}

func mongo_locate(moaddress MOAddr, f int) []neuron {
	session, err := mgo.Dial(moaddress.session)
	if err != nil {
			panic(err)
	}
	defer session.Close()
	c := session.DB(moaddress.table).C(moaddress.doc)
	result := []neuron{}
	iter := c.Find(bson.M{"uid": f}).Limit(100).Iter()
	err = iter.All(&result)
	if err != nil {
			log.Fatal(err)
	}
	return result
}

func mongo_locateone(moaddress MOAddr, f int) neuron {
	session, err := mgo.Dial(moaddress.session)
	if err != nil {
			panic(err)
	}
	defer session.Close()
	c := session.DB(moaddress.table).C(moaddress.doc)
	result := neuron{}
	iter := c.Find(bson.M{"uid": f}).Limit(1).Iter()
	err = iter.All(&result)
	if err != nil {
			log.Fatal(err)
	}
	return result
}

func mongo_update(moaddress MOAddr, t neuron) bool{
  session, err := mgo.Dial(moaddress.session)
  if err != nil {
    return false
  }
  defer session.Close()
  c := session.DB(moaddress.table).C(moaddress.doc)
  c.Update(bson.M{"uid": t.Uid}, t)
  return true
}

func mongo_login(moaddress MOAddr, f string) Login {
	session, err := mgo.Dial(moaddress.session)
	if err != nil {
			panic(err)
	}
	defer session.Close()
	c := session.DB(moaddress.table).C(moaddress.doc)
	result := Login{}
	erro := c.Find(bson.M{"username": f}).One(&result)
	if erro != nil {
			return result
	}
	return result
}

func mongo_loginSuccess(moaddress MOAddr, login Login) bool {
  session, err := mgo.Dial(moaddress.session)
  if err != nil {
    return false
  }
  defer session.Close()
  c := session.DB(moaddress.table).C(moaddress.doc)
  c.Update(bson.M{"username": login.Username}, login)
  return true
}

func mongo_findRelate(moaddress MOAddr, f string) related {
	session, err := mgo.Dial(moaddress.session)
	if err != nil {
			panic(err)
	}
	defer session.Close()
	c := session.DB(moaddress.table).C(moaddress.doc)
	result := related{}
	erro := c.Find(bson.M{"uid": f}).One(&result)
	if erro != nil {
			return result
	}
	return result
}