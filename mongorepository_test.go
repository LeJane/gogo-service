package main

import (
	"testing"

	"github.com/cloudnativego/cfmgo"
	"github.com/cloudnativego/gogo-engine"
	"github.com/cloudnativego/gogo-service/fakes"
)

var (
	fakeDBURI = "mongodb://fake.uri@addr:port/guid"
)

func TestAddMatchShowsUpInMongoRepository(t *testing.T) {
	var fakeMatches = []matchRecord{}
	var matchesCollection = cfmgo.Connect(
		fakes.FakeNewCollectionDialer(fakeMatches),
		fakeDBURI,
		MatchesCollectionName)

	repo := NewMongoMatchRepository(matchesCollection)
	match := gogo.NewMatch(19, "bob", "alfred")
	err := repo.addMatch(match)
	if err != nil {
		t.Error("Got an error adding a match to mongo, should not have.")
	}

	matches, err := repo.getMatches()
	if err != nil {
		t.Errorf("Got an error retrieving matches: %v", err)
	}
	if len(matches) != 1 {
		t.Errorf("Expected matches length to be 1; received %d", len(matches))
	}
}

func TestGetMatchRetrievesProperMatchFromMongo(t *testing.T) {
	fakes.TargetCount = 1
	var fakeMatches = []matchRecord{}
	var matchesCollection = cfmgo.Connect(
		fakes.FakeNewCollectionDialer(fakeMatches),
		fakeDBURI,
		MatchesCollectionName)

	repo := NewMongoMatchRepository(matchesCollection)
	match := gogo.NewMatch(19, "bob", "alfred")
	err := repo.addMatch(match)
	if err != nil {
		t.Errorf("Got an error adding a match to mongo: %v", err)
	}

	targetID := match.ID
	foundMatch, err := repo.getMatch(targetID)
	if err != nil {
		t.Errorf("Unable to find match with ID: %v... %s", targetID, err)
	}

	if foundMatch.GridSize != 19 || foundMatch.PlayerBlack != "bob" {
		t.Errorf("Unexpected match results: %v", foundMatch)
	}
}

func TestGetNonExistentMatchReturnsError(t *testing.T) {
	fakes.TargetCount = 0
	var fakeMatches = []matchRecord{}
	var matchesCollection = cfmgo.Connect(
		fakes.FakeNewCollectionDialer(fakeMatches),
		fakeDBURI,
		MatchesCollectionName)

	repo := NewMongoMatchRepository(matchesCollection)

	_, err := repo.getMatch("buckshank")
	if err == nil {
		t.Errorf("Expected getMatch to error with incorrect match details")
	}

	if err.Error() != "Match not found" {
		t.Errorf("Expected 'Match not found' error; received: '%v'", err)
	}

}
