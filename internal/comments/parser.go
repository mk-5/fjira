package comments

import (
	"bytes"
	"fmt"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
)

// TODO - could be optimized a bit
func ParseCommentsFromIssue(issue *jira.Issue, limitX, limitY int) []Comment {
	cs := make([]Comment, 0, 100)
	var commentsBuffer bytes.Buffer
	if len(issue.Fields.Comment.Comments) > 0 {
		for _, comment := range issue.Fields.Comment.Comments {
			title := fmt.Sprintf("%s, %s", comment.Created, comment.Author.DisplayName)
			body := fmt.Sprintf("\n%s", comment.Body)
			lines := app.DrawTextLimited(nil, 0, 0, limitX, limitY, app.DefaultStyle, comment.Body) + 2
			cs = append(cs, Comment{
				Title: title,
				Body:  body,
				Lines: lines,
			})
			commentsBuffer.Reset()
		}
	}
	return cs
}
