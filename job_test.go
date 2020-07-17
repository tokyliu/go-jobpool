package go_jobpool


type JobTest struct {
	jobId 		string
	done  		bool
	error 		error
	params 		[]interface{}
	result 		interface{}
}


func (job *JobTest)Before() {

}


func (job *JobTest)Run() {
	if job.error != nil {
		return
	}
	job.done = true
	//do you job
}


func (job *JobTest)After() {
	if job.error == nil {

	}

	//上传任务执行结果
	job.pushJobResult()
}


func (job JobTest)pushJobResult() {

}