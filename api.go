package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

type child struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	ChildName string        `bson:"name"`
	School    string        `bson:"school"`
}
type Message struct {
	text string
	to   []string
	from string
}

type Board struct {
	Txt    string
	Posted string
	Header string
	Date   string
}

type data struct {
	Dat   []*child
	board []*Board
}
type Parent struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Phone    string        `bson:"phone"`
	Name     string        `bson:"name"`
	Children []child       `bson:"children"`
	Pin      string        `bson:"pin"`
}

func RegParent(w http.ResponseWriter, r *http.Request) {
	//result := []string{}
	tmp := r.URL.Query().Get("no")
	//result = append(result, tmp)
	req, err := http.Get("http://api.clickatell.com/http/sendmsg?api_id=3593341&user=smilecs&password=NDEPPCLYOLXOJT&to=" + tmp + "&msg_type=SMS_FLASH&text=0999099")
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		defer req.Body.Close()
		contents, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("%s", err)
		}
		fmt.Printf("%s\n", string(contents))
		w.Write(contents)
	}

	/*m := new(Message)
	m.text = "09859584"
	m.to = result
	//Message{"09859584", result}
	/*	b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}*/
	/*h := "<?xml version 1.0?><request><data><text>Test Message</text><to>" + tmp + "</to></data></request>"
	var jsonStr = []byte(h)
	req, err := http.NewRequest("POST", "http://api.clickatell.com/rest/message", bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Version", "1")
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Authorization", "Bearer ["+Token.AuthToken+"]")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response:", string(body))
	w.Write(body)*/
}
