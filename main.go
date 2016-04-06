// $ go get -u all
package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

type Msg struct {
	hostname string
	stdout   string
	stderr   string
	err      error
}

func executeCmd(index int, hostname string) *Msg {
	resp, err := http.Get(fmt.Sprint("http://www.", hostname))
	//	if err != nil {
	//		fmt.Println(err)
	//		return ""
	//	}
	defer resp.Body.Close()

	return &Msg{err: err, hostname: fmt.Sprint(index, "] ", hostname, ": ", resp.Status)}
	//	return fmt.Sprint(index, "] ", hostname, ": ", resp.Status)

	//	body, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		fmt.Println(err)
	//		return ""
	//	}
	//	return string(body)
}

func start() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("NumCPU", runtime.NumCPU())
	hosts := []string{
		"vk.ru",
		"apriori-vk.ru",
		"delioro.ru",
		"shokoladki.ru",
		"chocolatevk.ru",
		"shokoeshka.com",
	} //os.Args[2:]
	// записываем результаты в буфферизированный список
	responseChannel := make(chan *Msg, 10)
	//	results := make(chan string, 10)
	// через 5 сек в канал timeout придет сообщение
	timeout := time.After(time.Second * 5)
	timeStart := time.Now()
	// запускаем по одной goroutine на сервер
	for index, hostname := range hosts {
		//		fmt.Println(index)
		go func(index int, hostname string) {
			responseChannel <- executeCmd(index, hostname)
		}(index, hostname)
		//		fmt.Println(executeCmd(index, hostname))
	}
	// соберем результаты со всех серверов или напишем "Время вышло"
	for i := 0; i < len(hosts); i++ {
		select {
		case msg := <-responseChannel:
			fmt.Println(msg.hostname)
		case <-timeout:
			fmt.Println("Time out")
			return
		}
	}
	timeTotal := time.Since(timeStart)
	fmt.Println(timeTotal) // 564.789389ms / 3.861487534s
}

func main() {
	start()
}
