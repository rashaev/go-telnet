package main


import (
	"net"
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"context"
	"time"
)


func connection(ctx context.Context, ip string, chanReq chan string, chanRes chan string) {
	ctx, _ = context.WithTimeout(ctx, 3 * time.Second)
	var d net.Dialer
	fmt.Fprint(os.Stdout, "Trying ", ip, "...\n")
	conn, err := d.DialContext(ctx, "tcp", ip)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		fmt.Fprint(os.Stdout, "Connected to ", ip, "\n")
		message, _ := bufio.NewReader(conn).ReadString('\n')
		chanRes<- message
	}
	defer conn.Close()

	for {
		x := <-chanReq
		conn.Write([]byte(x+"\n"))
		message, _ := bufio.NewReader(conn).ReadString('\n')
		chanRes<- message
		}
	}



func main() {
	var requestCh = make(chan string)
	var responseCh = make(chan string)
	var signalCh = make(chan os.Signal, 1)

	signal.Notify(signalCh, os.Interrupt)
	fmt.Println("Press CTRL + C to interrupt the program")

	ctx := context.Background()
	
	go connection(ctx, "smtp.yandex.ru:25", requestCh, responseCh)

	fmt.Printf("%s", <-responseCh)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan() 
		requestCh<- scanner.Text()
		
		select {
		case servRepl := <-responseCh:
			fmt.Printf("%s", servRepl)
		case <-signalCh:
			os.Exit(0)
		}	
	}
	
}