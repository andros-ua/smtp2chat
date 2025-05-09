// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"smtp2chat/smtp"
	"smtp2chat/teams"
	"smtp2chat/telegram"
)

var (
	service string
	token   string
	chatID  string
	webhook string
	bind    string
	verbose bool
)

func init() {
	flag.StringVar(&service, "s", "telegram", "Message service (telegram or teams)")
	flag.StringVar(&service, "service", "telegram", "Message service (telegram or teams)")

	flag.StringVar(&token, "t", "", "Telegram bot token")
	flag.StringVar(&token, "token", "", "Telegram bot token")

	flag.StringVar(&chatID, "c", "", "Telegram chat ID")
	flag.StringVar(&chatID, "chatid", "", "Telegram chat ID")

	flag.StringVar(&webhook, "w", "", "Teams webhook URL")
	flag.StringVar(&webhook, "webhook", "", "Teams webhook URL")

	flag.StringVar(&bind, "b", ":2525", "Address to bind SMTP server")
	flag.StringVar(&bind, "bind", ":2525", "Address to bind SMTP server")

	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "smtp2chat is a simple SMTP to Telegram/Teams relay server.")
		fmt.Fprintln(flag.CommandLine.Output(), "Available Options:")
		flag.PrintDefaults()
	}

	flag.Parse()
}

func main() {
	if service == "telegram" && (token == "" || chatID == "") {
		log.Fatal("--token and --chatid must be set for telegram")
	}
	if service == "teams" && webhook == "" {
		log.Fatal("--webhook must be set for teams")
	}

	listener, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	log.Printf("SMTP server listening on %s", bind)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if verbose {
				log.Printf("accept error: %v", err)
			}
			continue
		}
		if verbose {
			log.Printf("new connection from: %s", conn.RemoteAddr())
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	email := smtp.HandleConnection(conn, verbose)
	if email == nil {
		return
	}

	if verbose {
		fmt.Printf("Received email: Subject=%q From=%q To=%q\n", email.Subject, email.From, email.To)
	}

	switch service {
	case "telegram":
		telegram.SendMessage(token, chatID, email)
	case "teams":
		teams.SendMessage(webhook, email, verbose)
	default:
		log.Printf("unknown service: %s", service)
	}
}
