package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
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

func main() {
	log.SetFlags(log.Lshortfile)
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
}
