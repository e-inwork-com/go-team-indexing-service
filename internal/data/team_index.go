package data

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/golang-module/carbon"
)

type TeamIndex struct {
	SolrURL  string
	SolrTeam string
}

// Indexing a Team to the Solr Collection
func (i TeamIndex) Update(team *Team) (*http.Response, error) {
	// Delete or Update
	if !team.IsDeleted {
		createAt := carbon.Parse(team.CreatedAt.String()).ToRfc3339String("UTC")
		record := fmt.Sprintf(
			`{"id": "%v", "created_at":  "%v", "team_user": "%v", "team_name": "%v", "team_picture": "%v", "version": "%v"}`,
			team.ID, createAt, team.TeamUser, team.TeamName, team.TeamPicture, team.Version)

		res, err := http.Post(i.SolrURL+"/api/collections/"+i.SolrTeam+"/update?commit=true", "application/json", bytes.NewReader([]byte(record)))
		if err != nil {
			return res, err
		}
		if res.StatusCode != http.StatusOK {
			return res, res.Request.Context().Err()
		}

		return res, nil
	} else {
		res, err := http.Post(i.SolrURL+"/solr/"+i.SolrTeam+"/update?commit=true", "application/json", bytes.NewReader([]byte("{'delete': {'query': 'id:"+team.ID.String()+"'}}")))
		if err != nil {
			return res, err
		}
		if res.StatusCode != http.StatusOK {
			return res, res.Request.Context().Err()
		}

		return res, nil
	}
}
