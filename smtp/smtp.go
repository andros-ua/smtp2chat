// smtp/smtp.go
package smtp

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"strings"

	"smtp2chat/types"
)

func HandleConnection(conn net.Conn, verbose bool) *types.Email {
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	send(w, "220 smtp2chat Service Ready\r\n")

	hasMail := false
	hasRcpt := false
	email := &types.Email{}

	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF && verbose {
				log.Printf("read error: %v", err)
			}
			return nil
		}
		cmd := strings.TrimSpace(string(line))
		if verbose {
			log.Printf("client command: %s", cmd)
		}

		switch {
		case strings.HasPrefix(cmd, "EHLO") || strings.HasPrefix(cmd, "HELO"):
			send(w, "250-smtp2chat Hello\r\n250 HELP\r\n")
		case strings.HasPrefix(cmd, "MAIL FROM:"):
			hasMail = true
			email.From = parseEmailField(cmd)
			send(w, "250 OK\r\n")
		case strings.HasPrefix(cmd, "RCPT TO:"):
			if !hasMail {
				send(w, "503 Need MAIL first\r\n")
				continue
			}
			hasRcpt = true
			email.To = parseEmailField(cmd)
			send(w, "250 OK\r\n")
		case strings.HasPrefix(cmd, "DATA"):
			if !hasMail || !hasRcpt {
				send(w, "503 Need MAIL and RCPT first\r\n")
				continue
			}
			send(w, "354 End data with .<CR><LF>\r\n")
			email.Body = readEmailData(r, verbose)
			email.Subject = extractHeader(email.Body, "Subject:")
			send(w, "250 OK\r\n")
			return email
		case strings.HasPrefix(cmd, "QUIT"):
			send(w, "221 Bye\r\n")
			return nil
		default:
			send(w, "502 Command not implemented\r\n")
		}
	}
}

func parseEmailField(line string) string {
	start := strings.Index(line, "<")
	end := strings.Index(line, ">")
	if start == -1 || end == -1 || end <= start {
		return ""
	}
	return line[start+1 : end]
}

func extractHeader(data, header string) string {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToLower(line), strings.ToLower(header)) {
			return strings.TrimSpace(line[len(header):])
		}
	}
	return ""
}

func readEmailData(r *bufio.Reader, verbose bool) string {
	var buf bytes.Buffer
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if verbose {
				log.Printf("error reading DATA: %v", err)
			}
			return ""
		}
		if bytes.Equal(line, []byte(".\r\n")) {
			break
		}
		buf.Write(line)
	}
	if verbose {
		log.Println("finished reading email data")
	}
	return buf.String()
}

func send(w *bufio.Writer, msg string) {
	_, _ = w.WriteString(msg)
	_ = w.Flush()
}
