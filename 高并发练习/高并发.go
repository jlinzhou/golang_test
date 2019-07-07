package main
import "fmt"
import "time"

/*
代码说明
高并发任务
很多任务从一个入口进入，放到对外任务的channel，然后放到内部任务的channel，
多个work同时去内部任务的channel，进行工作

*/


//----------------有关Task任务角色的功能-----------
type Task struct{
	f func() error
}
//创建一个Task任务
func NewTask(args_f func()error)*Task{
	t:=Task{
		f:args_f,
	}
	return &t
}
//----------------Task也需要一个执行业务的方法----------------
func (t *Task)Execute(){
	t.f()//调用任务中已经绑定好的业务方法
}

//----------------有关携程池Pool角色的功能----------------
type Pool struct{
	//对外的Task入口 EntryChannel
	EntryChannel chan *Task

	//内部的Task队列 JobsChannel
	JobsChannel chan *Task

	//协程池中最大的worker的数量
	worker_num int
}
//创建Pool的函数
func NewPool(cap int) *Pool{
	//创建一个Pool
	p:=Pool{
		EntryChannel:make(chan *Task),
		JobsChannel:make(chan *Task),
		worker_num:cap,
	}
	//返回一个Pool
	return &p
}

//协程池创建一个worker，并且让这个worker去工作
func(p *Pool)worker(worker_ID int){
	//一个worker具体的工作
	//1永久的从JobsChannel去取任务
	for task:=range p.JobsChannel{
		//task就是当前worker从JobsChannel中拿到的任务
		//2一旦取到任务，执行这个任务
		task.Execute()
		fmt.Println("worker ID",worker_ID,"执行完了一个任务")
	}
}

//让协程池，开始真正的工作，协程池一个启动方法
func (p*Pool)run(){
	//1根据worker_num来创建worker去工作
	for i:=0;i<p.worker_num;i++{
		//每个worker都应该是一个goroutine
		go p.worker(i)

	}
	//2从EntryChannel中去取任务,将取到的任务，发送给JobsChannel
	for task := range p.EntryChannel{
		p.JobsChannel<-task
	}
}
//主函数 来测试协程池的工作
func main(){

	//1创建一些任务
	t:=NewTask(func() error{
		fmt.Println(time.Now())
		return nil;
	})
	tasknum:=0
	//2创建一个Pool协程池，这个协程池最大的worker数量是4
	p:=NewPool(4)
	//3将这些任务，交给协程池Pool
	go func(){
		for{
			tasknum=tasknum+1
			fmt.Println("任务数为:",tasknum)
			//不断的向p中写入任务t,每个任务就是打印当前的时间
			p.EntryChannel<-t
		}
	}()
	//4启动Pool
	p.run()
}