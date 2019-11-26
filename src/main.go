package main


import (
	"net"
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"context"
	"time"
	"github.com/spf13/pflag"
)

var timeOut time.Duration

func init() {
	pflag.DurationVarP(&timeOut, "timeout", "t", 5 * time.Second, "timeout connection")
}


func connection(ctx context.Context, ip, port string, chanReq chan string, chanRes chan string, timeOut time.Duration) {
	ctx2, _ := context.WithTimeout(ctx, timeOut)
	var d net.Dialer
	fmt.Fprint(os.Stdout, "Trying ", pflag.Arg(0)+":"+pflag.Arg(1), "...\n")
	conn, err := d.DialContext(ctx2, "tcp", ip+":"+port)
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
	pflag.Parse()
	if pflag.NArg() < 2 {
		fmt.Println("USAGE: go-telnet [OPTIONS] <HOST> <PORT>")
		os.Exit(1)
	}

	var requestCh = make(chan string)
	var responseCh = make(chan string)
	var signalCh = make(chan os.Signal, 1)

	signal.Notify(signalCh, os.Interrupt)
	fmt.Println("Press CTRL + C to interrupt the program")

	ctx := context.Background()
	
	go connection(ctx, pflag.Arg(0), pflag.Arg(1), requestCh, responseCh, timeOut)

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