// $ go get -u all
package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

func executeCmd(index int, hostname string) string {
	resp, err := http.Get(fmt.Sprint("http://www.", hostname))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	return fmt.Sprint(index, "] ", hostname, ": ", resp.Status)
	//	body, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		fmt.Println(err)
	//		return ""
	//	}
	//	return string(body)
}

func main() {
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
	results := make(chan string, 10)
	// через 5 сек в канал timeout придет сообщение
	timeout := time.After(time.Second * 5)
	timeStart := time.Now()
	// запускаем по одной goroutine на сервер
	for index, hostname := range hosts {
		//		fmt.Println(index)
		go func(index int, hostname string) {
			results <- executeCmd(index, hostname)
		}(index, hostname)
		//		fmt.Println(executeCmd(index, hostname))
	}
	// соберем результаты со всех серверов или напишем "Время вышло"
	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-results:
			fmt.Println(res)
		case <-timeout:
			fmt.Println("Time out")
			return
		}
	}
	timeTotal := time.Since(timeStart)
	fmt.Println(timeTotal) // 564.789389ms / 3.861487534s
}
