package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"time"
)

type GitHubEnv struct {
	Repository RepositoryTy
	Ref        string
	Commits    []CommitTy
}

type RepositoryTy struct {
	FullName string `json:"full_name"`
}

type CommitTy struct {
	Author  AuthorTy
	Id      string
	Message string
	Url     string
}

type AuthorTy struct {
	Name string
}

func parseEnv() *GitHubEnv {
	f, err := os.Open("env.json")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	e := &GitHubEnv{}
	err = json.Unmarshal(b, e)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", e)
	return e
}

var emailFmt = `Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=UTF-8
MIME-Version: 1.0
From: %s <%s>
To: %s
Subject: %s
Date: %s

%s`

func sendEmail() {
	f, err := os.Open("/tmp/diff.html")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	committer := "Sanjit Bhat"
	from := "mit.pdos.mailbot@gmail.com"
	to := "sanjit.bhat@gmail.com"
	subj := "test email 3"
	host := "smtp.gmail.com"
	pswd := os.Getenv("MAILBOT_SECRET")

	time := time.Now().Format(time.RFC1123Z)
	email := fmt.Sprintf(emailFmt, committer, from, to, subj, time, b)
	log.Print(email)

	auth := smtp.PlainAuth("", from, pswd, host)
	_ = auth
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(email))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	sendEmail()
}
