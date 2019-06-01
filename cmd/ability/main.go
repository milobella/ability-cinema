package main

import (
	"gitlab.milobella.com/milobella/ability-sdk-go/pkg/ability"
	"gitlab.milobella.com/milobella/cinema-ability/pkg/tools/allocine"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var allocineClient *allocine.Client
var allocineHost string
var allocinePort int

var additionalConfigPath string

//TODO: try to put some common stuff into a separate repository
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	additionalConfigPath = os.Getenv("ADDITIONAL_CONFIG_PATH")
	if len(additionalConfigPath) != 0 {
		viper.AddConfigPath(additionalConfigPath)
	}

	viper.AddConfigPath(".")
	viper.SetDefault("log-level", "info")

	logrus.SetFormatter(&logrus.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// TODO: read it in the config when move to viper
	logrus.SetLevel(logrus.DebugLevel)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		logrus.Errorf("Fatal error config file: %s \n", err)
	}

	if level, err := logrus.ParseLevel(viper.GetString("server.log-level")); err == nil {
		logrus.SetLevel(level)
	} else {
		logrus.Warn("Failed to parse the log level. Keeping the logrus default level.")
	}

	logrus.Debugf("Configuration -> %+v", viper.AllSettings())

	allocineHost = viper.GetString("tools.allocine.host")
	allocinePort = viper.GetInt("tools.allocine.port")
}

// fun main()
func main() {

	// Initialize client for allocine tool
	allocineClient = allocine.NewClient(allocineHost, allocinePort)

	// Initialize server
	server := ability.NewServer("Cinema Ability", viper.GetInt("server.port"))
	server.RegisterIntentRule("LAST_SHOWTIME", lastShowTimeHandler)
	server.Serve()
}

func lastShowTimeHandler(_ *ability.Request, resp *ability.Response) {
	result, err := allocineClient.GetLastShowTime("Mouans-Sartoux")
	if err != nil {
		resp.Nlg.Sentence = "Error"
		return
	}

	theater, _ := result.Path("feed.theaterShowtimes.place.theater.name").Children()
	location, _ := result.Path("feed.theaterShowtimes.place.theater.city").Children()

	resp.Nlg.Sentence = "Here are the movies in {{theater}} this evening, in the {{location}}'s theater"
	resp.Nlg.Params = []ability.NLGParam{
		{ Name: "theater", Value: theater[0].Data().(string), Type: "string" },
		{ Name: "location", Value: location[0].Data().(string), Type: "string" },
	}

	showTimesBug, _ := result.Path("feed.theaterShowtimes.movieShowtimes").Children()
	showTimes, _ := showTimesBug[0].Children()
	var visu []allocine.Show
	for _, show := range showTimes {
		visu = append(visu, allocine.Show{
			Title:   show.Path("onShow.movie.title").String(),
			Display: show.Path("display").String(),
		})
	}
	resp.Visu = visu
}
