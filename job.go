package go_jobpool

type Job interface {
	Before() error
	Run() error
	After() error
}


