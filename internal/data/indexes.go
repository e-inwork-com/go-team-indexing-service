// Copyright 2023, e-inwork.com. All rights reserved.

package data

type Indexes struct {
	Teams TeamIndex
}

func InitIndexes(solrURL string, solrTeam string) Indexes {
	return Indexes{
		Teams: TeamIndex{SolrURL: solrURL, SolrTeam: solrTeam},
	}
}
