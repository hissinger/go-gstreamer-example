package main

import (
	"log"
	"os"

	gst "github.com/spreadspace/go-gstreamer"
	"github.com/ziutek/glib"
)

func onMessage(bus *gst.Bus, message *gst.Message) {
	if message.GetType() == gst.MESSAGE_EOS {
		log.Println("End of stream")
		os.Exit(0)
	} else if message.GetType() == gst.MESSAGE_ERROR {
		_, debug := message.ParseError()
		log.Printf("ERROR: %s\n", debug)
	} else if message.GetType() == gst.MESSAGE_WARNING {
		_, debug := message.ParseWarning()
		log.Printf("WARNING: %s\n", debug)
	}
}

func main() {
	gst.Init(nil)

	pipeline, err := gst.PipelineNew("videotest-pipeline")
	if err != nil {
		log.Fatal(err)
	}

	source, err := gst.ElementFactoryMake("videotestsrc", "video-src")
	if err != nil {
		log.Fatal(err)
	}
	source.SetProperty("num-buffers", 50)
	sink, err := gst.ElementFactoryMake("autovideosink", "video-sink")
	if err != nil {
		log.Fatal(err)
	}

	pipeline.Bin.Add(source)
	pipeline.Bin.Add(sink)
	source.Link(sink)

	bus, err := pipeline.GetBus()
	if err != nil {
		log.Fatal(err)
	}

	bus.AddSignalWatch()
	_, err = bus.Connect("message", onMessage, nil)
	if err != nil {
		log.Fatal(err)
	}

	var state = pipeline.SetState(gst.STATE_PLAYING)
	log.Printf("StateChangeReturn: %d\n", state)

	glib.NewMainLoop(nil).Run()
}
