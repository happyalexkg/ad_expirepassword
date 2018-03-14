package main

import (
	"log"
	"gopkg.in/ldap.v2"
	"fmt"
	"strings"
	"strconv"
	"time"
	"net/smtp"
)

type Users struct {
	Username string
	Email string
	pwdChangedTime string

}


func Conn_Search( )  []Users {

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "dc01.example.com", 389))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = l.Bind("CN=testuser,CN=Users,DC=example,DC=com", "bublik")
	if err != nil {
		log.Fatal(err)
	}


	searchRequest := ldap.NewSearchRequest(
		"DC=example,DC=kg", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=*))", // The filter to apply
		[]string{ "mail","cn","pwdLastSet"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	var users []Users

	for _, entry := range sr.Entries {
if (strings.Contains(entry.GetAttributeValue("mail"), "@") && len (entry.GetAttributeValue("pwdLastSet")) > 0 && entry.GetAttributeValue("pwdLastSet") != "0") {
	users = append(users, Users{
		entry.GetAttributeValue("cn"),
		entry.GetAttributeValue("mail"),
		entry.GetAttributeValue("pwdLastSet"),
	})

	}

	}
	return users
}
//waring 14 days 
const warAge int64 = 1209600
// 30 days password age
const pwdAge int64 = 2592000

func convertit ( value string ) int64  {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func main () {
	users:= Conn_Search()

	for _,v := range users {
		i:=convertit (v.pwdChangedTime)
		i=i/ 10000000
		k:=((1970 - 1601) * 365 - 3 + ((1970 - 1601) / 4)) * 86400
		timestamp:=i-int64(k)
		currtime:=time.Now().Unix()
		difftime:=(currtime-timestamp)+warAge
		//fmt.Println(difftime)
		if (difftime > pwdAge ) {
				sendmail (v.Email,v.Username,(currtime-timestamp))
			}
		}
	}


func sendmail(email string, username string, expired int64) {

	expired1:=30-(expired/86400)

	if (expired1 > 0 ) {

		toAddresses := []string{email}

		toHeader := strings.Join(toAddresses, ", ")

		header := make(map[string]string)
		header["From"] = "it@ipc.com"
		header["To"] = toHeader

		header["Subject"] = "The validity of your password for the user " + username + " ends!"
		header["Content-Type"] = `text/html; charset="UTF-8"`

		msg := ""
		for k, v := range header {
			msg += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		msg += "\r\n"

		body := "The validity of your password for the user " + username + " ends in " + strconv.FormatInt(expired1, 10) + " days. \r\n"
		msg += "\r\n" + body
		fmt.Println(msg)

		bMsg := []byte(msg)

	err := smtp.SendMail("smtp.example.com:25", nil, "it@ipc.com", toAddresses, bMsg)
	if err != nil {
		log.Fatal("OPS", err)
	}

	}
}
