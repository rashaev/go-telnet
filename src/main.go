package main


import (
	"net"
	"bufio"
	//"io"
	"fmt"
	"os"
	"log"
	"context"
	"time"
)


func connection(ctx context.Context, ip string, chanReq chan string, chanRes chan string) {
	var d net.Dialer
	fmt.Fprint(os.Stdout, "Trying ", ip, "...\n")
	conn, err := d.DialContext(ctx, "tcp", ip)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Fprint(os.Stdout, "Connected to ", ip, "\n")
		message, _ := bufio.NewReader(conn).ReadString('\n')
		chanRes<- message
	}
	defer conn.Close()

	for {
		select {
		case x := <-chanReq:
			conn.Write([]byte(x+"\n"))
			message, _ := bufio.NewReader(conn).ReadString('\n')
			chanRes<- message
		//case <-ctx.Done():
			//return
		}
	}

}

func main() {
	var requestCh = make(chan string)
	var responseCh = make(chan string)

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 10000	 * time.Millisecond)
	go connection(ctx, "smtp.yandex.ru:25", requestCh, responseCh)

	fmt.Println(<-responseCh)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan() 
		requestCh<- scanner.Text()

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		fmt.Println(<-responseCh)	
	}
	
}