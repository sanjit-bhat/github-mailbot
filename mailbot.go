package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
	"time"
)

type GitHubEvent struct {
	Repository RepositoryTy
	Ref        string
	Commits    []*CommitTy
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

type EmailConfig struct {
	Host     string
	Port     string
	From     string
	To       string
	Password string
	Repo     string
	Branch   string
}

func main() {
	log.SetFlags(log.Lshortfile)

	host := os.Getenv("MAILBOT_HOST")
	port := os.Getenv("MAILBOT_PORT")
	from := os.Getenv("MAILBOT_FROM")
	to := os.Getenv("MAILBOT_TO")
	pswd := os.Getenv("MAILBOT_PASSWORD")
	eventB := []byte(os.Getenv("MAILBOT_GH_EVENT"))
	if len(pswd) == 0 {
		log.Fatal("empty password. is the password configured correctly?")
	}

	event := unmarshalEvent(eventB)
	config := &EmailConfig{
		Host:     host,
		Port:     port,
		From:     from,
		To:       to,
		Password: pswd,
		Repo:     event.Repository.FullName,
		Branch:   event.Ref,
	}
	for _, c := range event.Commits {
		f := genDiff(c.Id)
		config.sendEmail(f, c)
	}
}

func unmarshalEvent(b []byte) *GitHubEvent {
	e := &GitHubEvent{}
	err := json.Unmarshal([]byte(b), e)
	if err != nil {
		log.Fatal(err)
	}
	// example ref: "refs/heads/main".
	sp := strings.Split(e.Ref, "/")
	e.Ref = sp[len(sp)-1]
	return e
}

const diffFmt = `git show --compact-summary --patch --pretty=format:"%%h|%%B" %s | \
	delta --no-gitconfig --light | \
	aha > %s`

func genDiff(commitId string) (fName string) {
	f, err := os.CreateTemp("", "*.html")
	if err != nil {
		log.Fatal(err)
	}
	fName = f.Name()

	diff := fmt.Sprintf(diffFmt, commitId, fName)
	_, err = exec.Command("bash", "-c", diff).Output()
	if err != nil {
		log.Fatal(err)
	}
	return
}

const emailFmt = `Content-Type: text/html; charset=UTF-8
From: %s <%s>
To: %s
Subject: %s
Date: %s

%s`

func (cfg *EmailConfig) sendEmail(diffFile string, commit *CommitTy) {
	f, err := os.Open(diffFile)
	if err != nil {
		log.Fatal(err)
	}
	html, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	msg, _, _ := strings.Cut(commit.Message, "\n")
	subj := fmt.Sprintf("%s %s: %s", cfg.Repo, cfg.Branch, msg)
	time := time.Now().Format(time.RFC1123Z)
	email := fmt.Sprintf(emailFmt, commit.Author.Name, cfg.From, cfg.To, subj, time, html)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.From, cfg.Password, cfg.Host)
	toSplit := strings.Split(cfg.To, ",")
	err = smtp.SendMail(addr, auth, cfg.From, toSplit, []byte(email))
	if err != nil {
		log.Fatal(err)
	}

	os.Remove(diffFile)
}
