package routes

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/duncanleo/plex-dvr-hls/config"
	"github.com/gin-gonic/gin"
)

type ChannelSimplified struct {
	ID   int
	Name string
}

type Programme struct {
	HourStr       string
	DateTimeStart string
	DateTimeEnd   string
}

func XMLTV(c *gin.Context) {
	var channels []ChannelSimplified

	for index, channel := range config.Channels {
		channels = append(
			channels,
			ChannelSimplified{
				ID:   index + 1,
				Name: channel.Name,
		},
		)
	}

	var programmes []Programme
	var now = time.Now()

	for i := 0; i < 4; i++ {
                var nextSunday = now.AddDate(0, 0, 7-int(now.Weekday()))

		if now.Weekday() == time.Sunday {
			var start = time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
			var end = time.Date(now.Year(), now.Month(), now.Day(), 8, 90, 0, 0, now.Location())

			var dateTimeStart = start.Format("20060102150405 -0700")
			var dateTimeEnd = end.Format("20060102150405 -0700")
			var hourStr = start.Format("Jan-2-2006 3PM")

                programmes = append(
                        programmes,
                        Programme{
                                HourStr:       hourStr,
                                DateTimeStart: dateTimeStart,
                                DateTimeEnd:   dateTimeEnd,
                        },
                )
		} else {
			var start = time.Date(now.Year(), nextSunday.Month(), nextSunday.Day(), 8, 0, 0, 0, now.Location())
			var end = time.Date(now.Year(), nextSunday.Month(), nextSunday.Day(), 8, 90, 0, 0, now.Location())

			var dateTimeStart = start.Format("20060102150405 -0700")
			var dateTimeEnd = end.Format("20060102150405 -0700")

			var hourStr = start.Format("Jan-2-2006 3PM")

		programmes = append(
			programmes,
			Programme{
				HourStr:       hourStr,
				DateTimeStart: dateTimeStart,
				DateTimeEnd:   dateTimeEnd,
			},
		)
		}
               now = nextSunday
	}

	t := template.Must(template.New("xmltv.tmpl").ParseFiles("templates/xmltv.tmpl"))

	var b bytes.Buffer
	var err = t.Execute(
		&b,
		gin.H{
			"channels":   channels,
			"programmes": programmes,
		},
	)

	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "application/xml", b.Bytes())
}
