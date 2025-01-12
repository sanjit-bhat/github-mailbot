package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/smtp"
	"os"
	"os/exec"
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

var emailFmt = `Content-Type: text/html; charset=UTF-8
From: %s <%s>
To: %s
Subject: %s
Date: %s

%s`

func genTestSubj() string {
	return fmt.Sprintf("test email %v", rand.IntN(10000))
}

func sendEmail(diffFile string) {
	f, err := os.Open(diffFile)
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
	subj := genTestSubj()
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

	os.Remove(diffFile)
}

var diffFmt = `git show --compact-summary --patch --format="" 560b72c | \
	delta --no-gitconfig --light | \
	aha > %s`

func genDiff() string {
	f, err := os.CreateTemp("", "*.html")
	if err != nil {
		log.Fatal(err)
	}
	fName := f.Name()

	diff := fmt.Sprintf(diffFmt, fName)
	_, err = exec.Command("bash", "-c", diff).Output()
	if err != nil {
		log.Fatal(err)
	}
	return fName
}

func main() {
	log.SetFlags(log.Lshortfile)
	// f := genDiff()
	// sendEmail(f)
	host := os.Getenv("MAILBOT_HOST")
	port := os.Getenv("MAILBOT_PORT")
	from := os.Getenv("MAILBOT_FROM")
	to := os.Getenv("MAILBOT_TO")
	pswd := os.Getenv("MAILBOT_PASSWORD")

	log.Print(host)
	log.Print(port)
	log.Print(from)
	log.Print(to)
	if pswd == "" {
		log.Print("empty pswd")
	} else {
		log.Print("non-empty pswd")
	}
}
