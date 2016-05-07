package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func main() {

	wh := New()
	wh.Handle("repo:push", handleRepoPushEvent)
	wh.Handle("repo:fork", handleRepoForkEvent)
	wh.Handle("repo:commit_comment_created", handleCommitCommentCreated)
	wh.Handle("issue:created", handleIssues)
	wh.Handle("issue:updated", handleIssues)
	wh.Handle("issue:comment_created", handleIssues)
	wh.Handle("pullrequest:created", handlePullRequests)
	wh.Handle("pullrequest:updated", handlePullRequests)
	wh.Handle("pullrequest:approved", handlePullRequests)
	wh.Handle("pullrequest:unapproved", handlePullRequests)
	wh.Handle("pullrequest:fulfilled", handlePullRequests)
	wh.Handle("pullrequest:rejected", handlePullRequests)
	wh.Handle("pullrequest:comment_created", handlePullRequestComments)
	wh.Handle("pullrequest:comment_updated", handlePullRequestComments)
	wh.Handle("pull_request:comment_deleted", handlePullRequestComments)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Printf("Dump error: %s", err)
			return
		}
		fmt.Printf("%s\n\n", dump)

		wh.ServeHTTP(w, r)
	})
	http.ListenAndServe(os.Getenv("LISTEN"), nil)
}

func handleRepoPushEvent(headers Headers, event interface{}) error {
	log.Println("This is a push. Continue ...")
	pushEvent := event.(*PushEvent)

	log.Println("Change count:", len(pushEvent.Push.Changes))

	if len(pushEvent.Push.Changes) > 0 {
		change := pushEvent.Push.Changes[0]
		log.Println("Account:", change.New.Target.Author.User.Username)
		log.Println("Repository:", change.New.Repository.Name)
		log.Println("Change:", change.New.Target.Hash)
		log.Println("Commit message", strings.TrimSpace(change.New.Target.Message))
		log.Println("Branch:", change.New.Name)
	}
	return nil
}

func handleRepoForkEvent(headers Headers, event interface{}) error {
	log.Println("This is a fork. Continue ...")

	forkEvent := event.(*ForkEvent)

	log.Println("Forked repo:", forkEvent.Repository.FullName)
	log.Println("New repo:", forkEvent.Fork.FullName)
	log.Println("Fork created by:", forkEvent.Actor.Username)

	return nil
}

func handleCommitCommentCreated(headers Headers, event interface{}) error {
	log.Println("This is a commit comment. Continue ...")

	cc := event.(*CommitCommentCreatedEvent)

	log.Println("Comment by:", cc.Actor.Username)
	log.Println("Commit:", cc.Commit.Hash)
	log.Println("Comment (raw):", cc.Comment.Content.Raw)
	log.Println("Comment (html):", cc.Comment.Content.HTML)
	log.Println("Comment (markup):", cc.Comment.Content.Markup)

	return nil
}

func handleIssues(headers Headers, event interface{}) error {
	log.Println("This is an issue. Continue ...")

	switch i := event.(type) {
	case *IssueCreatedEvent:
		printIssue(&i.IssueEvent)
		log.Println("Action: Created")

	case *IssueUpdatedEvent:
		printIssue(&i.IssueEvent)
		log.Println("Action: Updated")
		log.Println("Old status value:", i.Changes.Status.Old)
		log.Println("New status value:", i.Changes.Status.New)

	case *IssueCommentCreatedEvent:
		printIssue(&i.IssueEvent)
		log.Println("Action: Comment Created")
		log.Println("Comment:", i.Comment.Content.Raw)
	}

	return nil
}

func printIssue(i *IssueEvent) {

	log.Println("Issue ID:", i.Issue.ID)
	log.Println("Repository:", i.Repository.FullName)
	log.Println("Actor:", i.Actor.Username)
	log.Println("Title:", i.Issue.Title)
	log.Println("Type:", i.Issue.Type)
	log.Println("Priority:", i.Issue.Priority)
	log.Println("-----------------------------------")
}

func handlePullRequests(headers Headers, event interface{}) error {
	log.Println("This is a pull request. Continue ...")

	switch pr := event.(type) {
	case *PullRequestCreatedEvent:
		printPullRequest(&pr.PullRequestEvent)
		log.Println("Action: CREATED")
	case *PullRequestUpdatedEvent:
		printPullRequest(&pr.PullRequestEvent)
		log.Println("Action: UPDATED")
	case *PullRequestMergedEvent:
		printPullRequest(&pr.PullRequestEvent)
		log.Println("Action: MERGED")
	case *PullRequestDeclinedEvent:
		printPullRequest(&pr.PullRequestEvent)
		log.Println("Action: DECLINED")
	case *PullRequestApprovedEvent:
		printPullRequest(&pr.PullRequestEvent)
		log.Println("Action: APPROVED")
		log.Println("Approved by:", pr.Approval.User.Username)
	case *PullRequestApprovalRemovedEvent:
		printPullRequest(&pr.PullRequestEvent)
		log.Println("Action: APPROVAL REMOVED")
		log.Println("Removed by:", pr.Approval.User.Username)
	}
	return nil
}

func printPullRequest(pr *PullRequestEvent) {

	log.Println("PR ID:", pr.PullRequest.ID)
	log.Println("Repository:", pr.Repository.FullName)
	log.Println("Actor:", pr.Actor.Username)
	log.Println("Author:", pr.PullRequest.Author.Username)
	log.Println("Description:", pr.PullRequest.Description)
	log.Println("Merge commit:", pr.PullRequest.MergeCommit.Hash)
	log.Printf("Destination: %s#%s", pr.PullRequest.Destination.Repository.FullName, pr.PullRequest.Destination.Branch.Name)
	log.Println("Num participants:", len(pr.PullRequest.Participants))
	participants := []string{}
	for _, u := range pr.PullRequest.Participants {
		participants = append(participants, u.User.Username)
	}
	log.Println("Participants:", strings.Join(participants, ", "))
	users := []string{}
	for _, u := range pr.PullRequest.Reviewers {
		users = append(users, u.Username)
	}
	log.Println("Reviewers:", strings.Join(users, ", "))

	log.Println("-----------------------------------")
}

func handlePullRequestComments(headers Headers, event interface{}) error {
	log.Println("This is a pull request comment. Continue ...")

	switch c := event.(type) {
	case *PullRequestCommentCreatedEvent:
		printPullRequestComment(&c.PullRequestCommentEvent)
		log.Println("Action: CREATED")
	case *PullRequestCommentUpdatedEvent:
		printPullRequestComment(&c.PullRequestCommentEvent)
		log.Println("Action: UPDATED")
	case *PullRequestCommentDeletedEvent:
		printPullRequestComment(&c.PullRequestCommentEvent)
		log.Println("Action: DELETED")

	}

	return nil
}

func printPullRequestComment(c *PullRequestCommentEvent) {

	log.Println("PR ID:", c.PullRequest.ID)
	log.Println("Repository:", c.Repository.FullName)
	log.Println("Actor:", c.Actor.Username)
	log.Println("Author:", c.PullRequest.Author.Username)
	log.Println("Description:", c.PullRequest.Description)
	log.Println("Comment:", c.Comment.Content.Raw)

	participants := []string{}
	for _, u := range c.PullRequest.Participants {
		participants = append(participants, u.User.Username)
	}
	log.Println("Participants:", strings.Join(participants, ", "))
	log.Println("-----------------------------------")
}
