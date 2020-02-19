package command

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/google/go-github/v29/github"
)

const commentRegex = `@sre-bot \/((.|\n)*)`

var (
	commentPattern = regexp.MustCompile(commentRegex)
)

func (c *CommandRedeliver) ProcessIssueCommentEvent(event *github.IssueCommentEvent) {
	// only PR author can trigger this comment
	if event.GetComment().GetUser().GetLogin() == "sre-bot" {
		return
	}
	if event.GetIssue().GetUser().GetLogin() != event.GetComment().GetUser().GetLogin() {
		return
	}
	// match comment
	m := commentPattern.FindStringSubmatch(event.GetComment().GetBody())
	if len(m) != 3 {
		return
	}
	comment := fmt.Sprintf("/%s", m[1])
	comment = strings.ReplaceAll(comment, "/merge", "")
	comment = strings.ReplaceAll(comment, "/run-auto-merge", "")
	comment = strings.TrimSpace(comment)
	if !strings.Contains(comment, "/run") &&
		!strings.Contains(comment, "/test") &&
		!strings.Contains(comment, "/bench") {
		comment = ""
	}
	if comment == "" {
		comment = fmt.Sprintf("@%s No command or invalid command", event.GetComment().GetUser().GetLogin())
	}
	githubComment := &github.IssueComment{
		Body: &comment,
	}
	IssueInfo := fmt.Sprintf("%s/%s #%d", c.repo.Owner, c.repo.Repo, event.GetIssue().GetNumber())
	if _, _, err := c.opr.Github.Issues.CreateComment(context.Background(),
		c.repo.Owner, c.repo.Repo, event.GetIssue().GetNumber(), githubComment); err != nil {
		log.Println("error occured when redeliver command in %s, %v", IssueInfo, err)
	} else {
		log.Println("redeliver command success, pull %s", IssueInfo)
	}
}
