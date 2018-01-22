package model

import "fmt"

type Result struct {
	Linter string
	Errors []Error
}

func (r Result) HasErrors() bool {
	return len(r.Errors) != 0
}

func (r Result) String() string {
	returnStr := r.Linter
	for _, e := range r.Errors {
		returnStr += fmt.Sprintf("\n\t%s\n\t%s\n", e.Message, e.Documentation)
	}
	return returnStr
}

type Error struct {
	Message       string
	Documentation string
}

func (r Error) String() string {
	return ""
}
