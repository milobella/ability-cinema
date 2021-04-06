package main

import (
	"github.com/milobella/ability-cinema/pkg/tools/allocine"
	"github.com/milobella/ability-sdk-go/pkg/ability"
	"github.com/sirupsen/logrus"
)

var allocineClient *allocine.Client

//TODO: User location is for now hardly defined but we need to take from the request.
const userLocation = "Mouans-Sartoux"

// fun main()
func main() {
	// Read configuration
	conf := ability.ReadConfiguration()

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
		lastShowTimeError(err, resp)
		return
	}
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.Debug(result.String())
	}

	theater, err := result.Path("feed.theaterShowtimes.place.theater.name").Children()
	if err != nil {
		lastShowTimeError(err, resp)
		return
	}
	location, err := result.Path("feed.theaterShowtimes.place.theater.city").Children()
	if err != nil {
		lastShowTimeError(err, resp)
		return
	}

	showTimesBug, err := result.Path("feed.theaterShowtimes.movieShowtimes").Children()
	if err != nil {
		lastShowTimeError(err, resp)
		return
	}
	if len(showTimesBug) <= 0 {
		resp.Nlg.Sentence = "There is no movie at {{theater}}, the {{location}}'s theater this evening"
		resp.Nlg.Params = []ability.NLGParam{
			{Name: "theater", Value: theater[0].Data().(string), Type: "string"},
			{Name: "location", Value: location[0].Data().(string), Type: "string"},
		}
	}
	showTimes, err := showTimesBug[0].Children()
	if err != nil {
		lastShowTimeError(err, resp)
		return
	}

	resp.Nlg.Sentence = "Here are the movies in {{theater}} this evening, in the {{location}}'s theater"
	resp.Nlg.Params = []ability.NLGParam{
		{Name: "theater", Value: theater[0].Data().(string), Type: "string"},
		{Name: "location", Value: location[0].Data().(string), Type: "string"},
	}


	var visu []allocine.Show
	for _, show := range showTimes {
		visu = append(visu, allocine.Show{
			Title:   show.Path("onShow.movie.title").String(),
			Display: show.Path("display").String(),
		})
	}
	resp.Visu = visu
}

func lastShowTimeError(err error, resp *ability.Response) {
	logrus.WithError(err).WithField("location", userLocation).Error("An error occurred while handling LAST_SHOWTIME intent.")
	resp.Nlg.Sentence = "Error retrieving the last shows in theater of {{location}}."
	resp.Nlg.Params = []ability.NLGParam{
		{Name: "location", Value: userLocation, Type: "string"},
	}
	return
}
