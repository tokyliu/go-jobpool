package go_jobpool


type ErrPoolClosed struct {

}

func (ErrPoolClosed)Error() string {
	return "job pool is closed"
}
