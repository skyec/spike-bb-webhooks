package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	wh := New()
	assert.NotNil(t, wh)
	assert.IsType(t, &WebHooks{}, wh)
}

func TestEvents(t *testing.T) {
	wh := New()

	type fixture struct {
		event       string
		handler     WebHookHandler
		payloadFile string
	}
	for _, fix := range []fixture{

		{"repo:push", func(h Headers, e interface{}) error {
			pushEvent := e.(*PushEvent)
			assert.Equal(t, "repo:push", h["X-Event-Key"])
			assert.Equal(t, "spike-bb-hook-source-repo", pushEvent.Repository.Name)
			return nil
		}, "push_event.json"},

		{"repo:fork", func(h Headers, e interface{}) error {
			forkEvent := e.(*ForkEvent)
			assert.NotNil(t, forkEvent, "forkEvent")
			assert.Equal(t, "spike-bb-hook-source-repo-forked", forkEvent.Fork.Name)
			return nil
		}, "fork_event.json"},

		{"repo:commit_comment_created", func(h Headers, e interface{}) error {
			commitCommentCreated := e.(*CommitCommentCreatedEvent)
			assert.NotNil(t, commitCommentCreated)
			assert.Equal(t, "This is a comment", commitCommentCreated.Comment.Content.Raw)
			assert.Equal(t, "markdown", commitCommentCreated.Comment.Content.Markup)
			return nil
		}, "commit_comment_created_event.json"},

		// TODO: repo:commit_status_created
		// TODO: repo:commit_status_updated

		{"issue:created", func(h Headers, e interface{}) error {
			ic := e.(*IssueCreatedEvent)
			assert.NotNil(t, ic)
			assert.Equal(t, "This is the first issue", ic.Issue.Title)
			assert.Equal(t, "This is the issue description.", ic.Issue.Content.Raw)
			assert.Equal(t, "skyec", ic.Actor.Username)
			assert.Equal(t, "spike-bb-hook-source-repo", ic.Repository.Name)
			return nil
		}, "issue_created_event.json"},

		{"issue:updated", func(h Headers, e interface{}) error {
			ic := e.(*IssueUpdatedEvent)
			assert.NotNil(t, ic)
			assert.Equal(t, "This is the first issue", ic.Issue.Title)
			assert.Equal(t, "This is the issue description.", ic.Issue.Content.Raw)
			assert.Equal(t, "skyec", ic.Actor.Username)
			assert.Equal(t, "spike-bb-hook-source-repo", ic.Repository.Name)
			return nil
		}, "issue_updated_event.json"},

		// TODO: issue:comment_created
		// TODO: pullrequest:created
		// TODO: pullrequest:updated
		// TODO: pullrequest:approved
		// TODO: pullrequest:unapproved
		// TODO: pullrequest:fulfilled
		// TODO: pullrequest:rejected
		// TODO: pullrequest:comment_created
		// TODO: pullrequest:comment_updated
		// TODO: pull_request:comment_deleted

	} {
		called := false
		wh.Handle(fix.event, func(h Headers, e interface{}) error {
			called = true
			return fix.handler(h, e)
		})

		jsn, err := ioutil.ReadFile("fixtures/" + fix.payloadFile)
		require.Nil(t, err)

		rec := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "http://example.com", bytes.NewReader(jsn))
		require.Nil(t, err)
		req.Header.Add("X-Event-Key", fix.event)

		log.Println("Test event:", fix.event)
		wh.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.True(t, called, "Event is called: "+fix.event)
	}

}
