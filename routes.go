package main

type Route int

const (
	QuestionRoute = iota
	LoadingRoute
	ResultsRoute
)

var routeName = map[Route]string{
	QuestionRoute: "question",
	LoadingRoute:  "loading",
	ResultsRoute:  "results",
}

func (s Route) String() string {
	return routeName[s]
}
