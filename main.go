package main

import (
	"crypto/tls"
	"flag"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	// Structure for settings - these can be modified by the user
	Settings struct {
		// Interval (seconds) defining how often a mock log entry is created
		LogInterval int
		// File to log into
		LogFile string
	}
	LogOutput struct {
		logger    *logrus.Logger
		fields    *logrus.Fields
		formatter *logrus.TextFormatter
	}
)

// Create a new instance of the logger. You can have any number of instances.
//var stdout_log = logrus.Logger{}
//var file_log = logrus.Logger{}

func getNewRandomJoke() string {
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com", nil)
	if err != nil {
		log.Fatalln("Error: Unable to prepare HTTP request:", err)
	}
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}
	req.Header.Set("Accept", "text/plain")
	resp, reqErr := client.Do(req)
	if reqErr != nil {
		log.Fatalln("Error: Unable to send HTTP request:", reqErr)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatalln("Error: Unable to read HTTP response:", readErr)
	}
	return string(body)
}

func log_generic(output *LogOutput) {
	output.logger.SetLevel(logrus.TraceLevel)
	output.logger.WithFields(*output.fields).Info(getNewRandomJoke())
}

// Start logging every _interval_ seconds
func log_with_interval(settings *Settings) {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	duration := time.Duration(settings.LogInterval * int(time.Second))

	hostname, ok := os.LookupEnv("HOSTNAME")
	if !ok {
		hostname = "unset"
	}

	// Initialize loggers
	// audit log
	auditOutput := LogOutput{
		logger: logrus.New(),
		fields: &logrus.Fields{
			"Source":     hostname,
			"SourceType": "audit",
			"EventType":  "privilege",
		},
		formatter: &logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true, DisableColors: true},
	}
	auditOutput.logger.SetFormatter(auditOutput.formatter)
	// standard log
	standardOutput := LogOutput{
		logger: logrus.New(),
		fields: &logrus.Fields{
			"Source":     "https://icanhazdadjoke.com",
			"SourceType": "api",
			"EventType":  "getNewRandomJoke",
		},
		formatter: &logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true, DisableColors: true},
	}
	standardOutput.logger.SetOutput(os.Stdout)
	standardOutput.logger.SetFormatter(standardOutput.formatter)

	// Log file was defined, output audit log to file
	if settings.LogFile != "" {
		auditOutput.logger.SetOutput(&lumberjack.Logger{
			Filename:   settings.LogFile,
			MaxSize:    3, // megabytes
			MaxBackups: 1,
			MaxAge:     7,    //days
			Compress:   true, // disabled by default
		})
	} else { // Log file was not defined, output audit log to console
		auditOutput.logger.SetOutput(os.Stdout)
	}
	var outputs []LogOutput
	outputs = append(outputs, auditOutput)
	outputs = append(outputs, standardOutput)

	for range time.Tick(duration) {
		for output := range outputs {
			log_generic(&outputs[output])
		}
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
	return Settings{LogInterval: 5, LogFile: ""}
}

// Entry
func main() {
	settings := get_default_settings()
	flag.IntVar(&settings.LogInterval, "interval", settings.LogInterval, "Interval defining how often a mock log entry is created")
	flag.StringVar(&settings.LogFile, "logfile", settings.LogFile, "Logfile")
	flag.Parse()
	//fmt.Printf("Using settings: \n%+v\n", settings)
	start_logging(&settings)
}
