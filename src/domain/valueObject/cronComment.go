package valueObject

import "errors"

type CronComment string

func NewCronComment(value string) (CronComment, error) {
	comment := CronComment(value)
	if !comment.isValid() {
		return "", errors.New("InvalidCronComment")
	}
	return comment, nil
}

func NewCronCommentPanic(value string) CronComment {
	comment, err := NewCronComment(value)
	if err != nil {
		panic(err)
	}
	return comment
}

func (comment CronComment) isValid() bool {
	isTooShort := len(string(comment)) < 2
	isTooLong := len(string(comment)) > 512
	return !isTooShort && !isTooLong
}

func (comment CronComment) String() string {
	return string(comment)
}
