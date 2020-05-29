package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"
	"log"
	"encoding/csv"
	"io/ioutil"

	"github.com/Rhymen/go-whatsapp/binary/proto"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

func main() {
	//create new WhatsApp connection
	res := parseCsv("input.csv")

	wac, err := whatsapp.NewConn(10 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	err = login(wac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		return
	}

	<-time.After(3 * time.Second)

	for _, r := range res {
		fmt.Printf("Sending Message to Phone number %s \n", r[0])
		sendDocumentMessage(wac, r[0], "document.pdf")
		sendTextMessage(wac, r[0])
		sendImageMessage(wac, r[0], "image.jpg")
		time.Sleep(3 * time.Second)
	}
}


func parseCsv(fileName string) [][]string {
	csvfile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	rows, err := r.ReadAll()

	return rows
}

func sendTextMessage(wac *whatsapp.Conn, phoneNumber string) {

	content, err := ioutil.ReadFile("welcome.txt")
	text := string(content)

	previousMessage := "😘"
	quotedMessage := proto.Message{
		Conversation: &previousMessage,
	}

	ContextInfo := whatsapp.ContextInfo{
		QuotedMessage:   &quotedMessage,
		QuotedMessageID: "",
		Participant:     "", //Whot sent the original message
	}

	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: phoneNumber + "@s.whatsapp.net",
		},
		ContextInfo: ContextInfo,
		// Text:        "As-salamu alaykum \n\nThanks for Contacting Hayder Ali & Company \n\n*Following are the Project Details for Defence Skyline*",
		Text: text,
	}

	msgId, err := wac.Send(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending image message: %v", err)
		os.Exit(1)
	} else {
		fmt.Println("Message Sent -> ID : " + msgId)
	}
}

func sendImageMessage(wac *whatsapp.Conn, phoneNumber string, imageName string) {

	content, err := ioutil.ReadFile("image.txt")
	text := string(content)

	img, err := os.Open(imageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	msg2 := whatsapp.ImageMessage{
		Info: whatsapp.MessageInfo {
			RemoteJid: phoneNumber + "@s.whatsapp.net",
		},
		Type:    "image/jpeg",
		Caption: text,
		Content: img,
	}

	msgId, err := wac.Send(msg2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending image message: %v", err)
		os.Exit(1)
	} else {
		fmt.Println("Message Sent -> ID : " + msgId)
	}
}

func sendDocumentMessage(wac *whatsapp.Conn, phoneNumber string, documentName string) {

	content, err := ioutil.ReadFile("image.txt")
	text := string(content)

	img, err := os.Open(documentName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	msg2 := whatsapp.DocumentMessage{
		Info: whatsapp.MessageInfo {
			RemoteJid: phoneNumber + "@s.whatsapp.net",
		},
		Type: "application/pdf",
		Title: text + "-Document",
		Content: img,
	}

	msgId, err := wac.Send(msg2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending image message: %v", err)
		os.Exit(1)
	} else {
		fmt.Println("Message Sent -> ID : " + msgId)
	}
}

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}
