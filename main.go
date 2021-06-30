package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instances.
var log = logrus.New()

func main() {
	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	log.Out = os.Stdout

	// You could set this to any `io.Writer` such as a file
	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	//  log.Out = file
	// } else {
	//  log.Info("Failed to log to file, using default stderr")
	// }

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Accept", "text/plain")
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)

	log.WithFields(logrus.Fields{
		"source": "https://icanhazdadjoke.com",
	}).Info(sb)
}
