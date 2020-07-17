package go_jobpool

type Job interface {
	Before()
	Run()
	After()
}


