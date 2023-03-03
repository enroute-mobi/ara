package siri_tests

import (
	"math/rand"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri"
)

var benchmarkResult string // Store benchmark result to avoid complilator optimisation

// Fill a single SIRINotifyProductionTimetable with a variable number of RecordedCalls
func benchmarkPTTNotifyBuildXML(pc int, b *testing.B) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	ptt := &siri.SIRINotifyProductionTimetable{
		ProducerRef:            "ProducerRef",
		SubscriptionIdentifier: "SubscriptionIdentifier",
		ResponseTimestamp:      time.Now(),
		Status:                 true,
		DatedTimetableVersionFrames: []*siri.SIRIDatedTimetableVersionFrame{
			&siri.SIRIDatedTimetableVersionFrame{
				LineRef:        "LineRef",
				RecordedAtTime: time.Now(),
				Attributes:     make(map[string]string),
				DatedVehicleJourneys: []*siri.SIRIDatedVehicleJourney{
					&siri.SIRIDatedVehicleJourney{
						DataFrameRef:           "DataFrameRef",
						DatedVehicleJourneyRef: "DatedVehicleJourneyRef",
						PublishedLineName:      "PublishedLineName",
						Attributes:             make(map[string]string),
						References:             make(map[string]string),
					},
				},
			},
		},
	}

	var dcs []*siri.SIRIDatedCall

	for i := 0; i != pc; i++ {
		dc := &siri.SIRIDatedCall{
			StopPointRef:       randSeq(10),
			StopPointName:      "StopPointName",
			DestinationDisplay: "DestinationDisplay",
			Order:              i,
			AimedArrivalTime:   time.Now(),
			AimedDepartureTime: time.Now(),
		}
		dcs = append(dcs, dc)
	}

	ptt.DatedTimetableVersionFrames[0].DatedVehicleJourneys[0].DatedCalls = dcs

	for n := 0; n < b.N; n++ {
		benchmarkResult, _ = ptt.BuildXML()
	}
}

func BenchmarkPTTNotifyBuildXML10(b *testing.B)     { benchmarkPTTNotifyBuildXML(10, b) }
func BenchmarkPTTNotifyBuildXML50(b *testing.B)     { benchmarkPTTNotifyBuildXML(50, b) }
func BenchmarkPTTNotifyBuildXML100(b *testing.B)    { benchmarkPTTNotifyBuildXML(100, b) }
func BenchmarkPTTNotifyBuildXML1000(b *testing.B)   { benchmarkPTTNotifyBuildXML(1000, b) }
func BenchmarkPTTNotifyBuildXML5000(b *testing.B)   { benchmarkPTTNotifyBuildXML(5000, b) }
func BenchmarkPTTNotifyBuildXML10000(b *testing.B)  { benchmarkPTTNotifyBuildXML(10000, b) }
func BenchmarkPTTNotifyBuildXML100000(b *testing.B) { benchmarkPTTNotifyBuildXML(100000, b) }

// To generate a random string
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
