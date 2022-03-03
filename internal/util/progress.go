package util

type Runnable interface {
	Run()
}

type Progresser interface {
	Runnable
	Progress() int
	Completed() int
	Total() int
	Failed() int
}
