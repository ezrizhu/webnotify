package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gregdel/pushover"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/button", buttonHandler)
	log.Fatal(http.ListenAndServe(":8887", nil))
}

func buttonHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	app := pushover.New(os.Getenv("PO_APP"))
	recipient := pushover.NewRecipient(os.Getenv("PO_RECIPIENT"))
	title := r.FormValue("title")
	message := r.FormValue("message")
	secret := r.FormValue("secret")
	ip := r.Header.Get("X-Forwarded-For")
	if r.Form.Get("secret") != os.Getenv("SECRET") {
		println("UNAUTHORIZED Button pressed from " + ip + " with title " + title + " message " + message + " and secret " + secret)
		time.Sleep(10 * time.Second)
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "Unauthorized")
		return
	}
	if (title == "") || (message == "") {
		io.WriteString(w, "Title and message are required")
		return
	}
	pushMessage := &pushover.Message{
		Title:    title,
		Message:  message,
		Priority: pushover.PriorityHigh,
	}
	response, err := app.SendMessage(pushMessage, recipient)
	if err != nil {
		log.Panic(err)
	}
	if response.Status != 1 {
		log.Panic(response.Status)
		io.WriteString(w, "Error")
	} else {
		io.WriteString(w, "Success")
	}
	println("Button pressed from " + ip + " with title " + title + " and message " + message)
}
