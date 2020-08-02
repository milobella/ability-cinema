package main

import (
	"github.com/milobella/ability-cinema/internal/config"
	"github.com/milobella/ability-cinema/pkg/tools/allocine"
	"github.com/milobella/ability-sdk-go/pkg/ability"
	"os"

	"github.com/sirupsen/logrus"
)

var allocineClient *allocine.Client

//TODO: User location is for now hardly defined but we need to take from the request.
const userLocation = "Mouans-Sartoux"

//TODO: try to put some common stuff into a separate repository
func init() {

	logrus.SetFormatter(&logrus.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// TODO: read it in the config when move to viper
	logrus.SetLevel(logrus.DebugLevel)
}

// fun main()
func main() {
	// Read configuration
	conf, err := config.ReadConfiguration()
	if err != nil { // Handle errors reading the config file
		logrus.WithError(err).Fatalf("Error reading config.")
	} else {
		logrus.Infof("The configuration has been successfully ridden.")
		logrus.Debugf("-> %+v", conf)
	}

	// Initialize client for allocine tool
	allocineClient = allocine.NewClient(conf.Tools["allocine"].Host, conf.Tools["allocine"].Port)

	// Initialize server
	server := ability.NewServer("Cinema Ability", conf.Server.Port)
	server.RegisterIntentRule("LAST_SHOWTIME", lastShowTimeHandler)
	server.Serve()
}

func lastShowTimeHandler(_ *ability.Request, resp *ability.Response) {
	result, err := allocineClient.GetLastShowTime(userLocation)
	if err != nil {
		resp.Nlg.Sentence = "Error retrieving the last shows in theater of {{location}}."
		resp.Nlg.Params = []ability.NLGParam{
			{Name: "location", Value: userLocation, Type: "string"},
		}
		return
	}

	theater, _ := result.Path("feed.theaterShowtimes.place.theater.name").Children()
	location, _ := result.Path("feed.theaterShowtimes.place.theater.city").Children()

	resp.Nlg.Sentence = "Here are the movies in {{theater}} this evening, in the {{location}}'s theater"
	resp.Nlg.Params = []ability.NLGParam{
		{Name: "theater", Value: theater[0].Data().(string), Type: "string"},
		{Name: "location", Value: location[0].Data().(string), Type: "string"},
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
