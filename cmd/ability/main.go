package main

import (
	"milobella/abilities/ability-sdk-go/pkg/ability"
	"milobella/abilities/cinema-ability/pkg/tools/allocine"
)

var allocineClient = allocine.NewClient("http://0.0.0.0", 8000)

// fun main()
func main() {
	server := ability.NewServer(10200)
	server.RegisterIntent("LAST_SHOWTIME", lastShowTimeHandler)
	server.Serve()
}

func lastShowTimeHandler(req ability.Request, resp *ability.Response) {
	result, err := allocineClient.GetLastShowTime("Mouans-Sartoux")
	if err != nil {
		resp.Nlg.Sentence = "Error"
		return
	}

	theater, _ := result.Path("feed.theaterShowtimes.place.theater.name").Children()
	location, _ := result.Path("feed.theaterShowtimes.place.theater.city").Children()

	resp.Nlg.Sentence = "Here are the movies in {{theater}} this evening, in the {{location}}'s theater"
	resp.Nlg.Params = map[string]string{
		"theater":  theater[0].Data().(string),
		"location": location[0].Data().(string),
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
