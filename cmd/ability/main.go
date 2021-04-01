package main

import (
	"github.com/milobella/ability-cinema/pkg/tools/allocine"
	"github.com/milobella/ability-sdk-go/pkg/ability"
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
