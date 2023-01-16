// Copyright 2023, e-inwork.com. All rights reserved.

package grpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/e-inwork-com/go-team-indexing-service/internal/data"
	"github.com/e-inwork-com/go-team-indexing-service/internal/jsonlog"
	"github.com/stretchr/testify/assert"
)

func TestE2E(t *testing.T) {
	// Set Configuration
	var cfg Config

	cfg.Db.Dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	cfg.Db.MaxOpenConn = 25
	cfg.Db.MaxIdleConn = 25
	cfg.Db.MaxIdleTime = "15m"
	cfg.GRPCPort = "5001"
	cfg.SolrURL = "http://localhost:8983"
	cfg.SolrTeam = "teams"

	// Set logger
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// Set Database
	db, err := OpenDB(cfg)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()

	// Read a SQL file for the deleting all records
	script, err := os.ReadFile("./test/sql/delete_all.sql")
	if err != nil {
		t.Fatal(err)
	}

	// Delete all records
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Delete all indexing in the Solr Team Collection
	res, err := http.Post("http://localhost:8983/solr/"+cfg.SolrTeam+"/update?commit=true", "application/json", bytes.NewReader([]byte("{'delete': {'query': '*:*'}}")))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal(res)
	}

	// Set the application
	app := Application{
		Config:  cfg,
		Logger:  logger,
		Models:  data.InitModels(db),
		Indexes: data.InitIndexes(cfg.SolrURL, cfg.SolrTeam),
	}

	// Run gRPC Server
	go app.GRPCListen()

	// Register
	email := "jon@doe.com"
	password := "pa55word"
	user := fmt.Sprintf(`{"email": "%v", "password":  "%v", "first_name": "Jon", "last_name": "Doe"}`, email, password)
	res, err = http.Post("http://localhost:8000/service/users", "application/json", bytes.NewReader([]byte(user)))
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusAccepted)

	// Sign in
	login := fmt.Sprintf(`{"email": "%v", "password":  "%v"}`, email, password)
	res, err = http.Post("http://localhost:8000/service/users/authentication", "application/json", bytes.NewReader([]byte(login)))
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)

	// Read a token
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	assert.Nil(t, err)

	type authType struct {
		Token string `json:"token"`
	}
	var authResult authType
	err = json.Unmarshal(body, &authResult)
	assert.Nil(t, err)
	assert.NotNil(t, authResult.Token)

	// Create Team
	// Create body buffer
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// Add team name
	teamName := "Doe's Team"
	bodyWriter.WriteField("team_name", teamName)

	// Add team picture
	filename := "./test/images/team.jpg"
	fileWriter, err := bodyWriter.CreateFormFile("team_picture", filename)
	if err != nil {
		t.Fatal(err)
	}

	// Open file
	fileHandler, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	// Copy file
	_, err = io.Copy(fileWriter, fileHandler)
	if err != nil {
		t.Fatal(err)
	}

	// Put on body
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// Post a new team
	req, _ := http.NewRequest("POST", "http://localhost:8000/service/teams", bodyBuf)
	req.Header.Add("Content-Type", contentType)

	bearer := fmt.Sprintf("Bearer %v", authResult.Token)
	req.Header.Set("Authorization", bearer)

	client := &http.Client{}
	res, err = client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)

	// Get all indexing in the Solr Team Collection
	res, err = http.Get("http://localhost:8983/api/collections/" + cfg.SolrTeam + "/select?q=*:*")
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	// Read a response
	body, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	assert.Nil(t, err)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(body), &result)
	assert.Nil(t, err)

	// Found 1 record
	response := result["response"].(map[string]interface{})
	assert.NotNil(t, response)
	assert.Equal(t, response["numFound"], float64(1))

	// Response docs is not empty
	docs := response["docs"].([]interface{})
	assert.NotNil(t, docs)

	// Read first doc
	doc := docs[0].(map[string]interface{})
	assert.NotNil(t, doc)
	assert.Equal(t, doc["team_name"], teamName)
}
