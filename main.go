package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
)

// Structure for settings - these can be modified by the user
type Settings struct {
	// Interval (seconds) defining how often a mock log entry is created
	LogInterval int

	// File to log into
	LogFile string
}

// Create a new instance of the logger. You can have any number of instances.
var stdout_log = logrus.New()
var file_log = logrus.New()

func log_generic(log *logrus.Logger) {
	log.SetLevel(logrus.TraceLevel)
	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Accept", "text/plain")
	resp, _ := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)

	log.WithFields(logrus.Fields{
		"source": "https://icanhazdadjoke.com",
	}).Info(sb)
}

// Start logging every _interval_ seconds
func log_with_interval(settings *Settings) {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	f, err := os.OpenFile(settings.LogFile, os.O_WRONLY|os.O_CREATE, 0755)

	if err != nil {
		log.Printf("Could not open file: %s - E: %s, exiting...", settings.LogFile, err.Error())
		os.Exit(1)
	}
	duration := time.Duration(settings.LogInterval * int(time.Second))

	// Initialize loggers
	file_formatter := &logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true, DisableColors: true}
	stdout_formatter := &logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true}
	file_log.SetOutput(f)
	stdout_log.SetOutput(os.Stdout)
	file_log.SetFormatter(file_formatter)
	stdout_log.SetFormatter(stdout_formatter)

	for range time.Tick(duration) {
		log_generic(stdout_log)
		log_generic(file_log)
	}
}

// Wrapper for the actual mock logger functionality
func start_logging(settings *Settings) {
	go log_with_interval(settings)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until signal (os.Interrupt)
	<-c
}

func get_default_settings() Settings {
	return Settings{LogInterval: 3, LogFile: "./mocklog.log"}
}

// Entry
func main() {
	settings := get_default_settings()
	flag.IntVar(&settings.LogInterval, "interval", 5, "Interval defining how often a mock log entry is created")
	flag.StringVar(&settings.LogFile, "logfile", "./mock_log.log", "Logfile")
	flag.Parse()
	//fmt.Printf("Using settings: \n%+v\n", settings)
	start_logging(&settings)
}
