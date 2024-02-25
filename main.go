package main

import (
	"encoding/csv"
	"fmt"
	"github.com/twpayne/go-kml/v3"
	"image/color"
	"log"
	"os"
	"strconv"
)

func write(ggs map[string]kml.Coordinate, vr map[string]kml.Coordinate) {
	f, err := os.OpenFile("new.kml", os.O_CREATE|os.O_TRUNC, 0777)
	catchErr(err)
	defer f.Close()
	k := kml.Folder()
	style := kml.SharedStyle(
		"green",
		kml.LineStyle(
			kml.Color(color.RGBA{R: 170, G: 255, B: 0, A: 200}),
			kml.Width(1),
		))
	k.Append(style)
	for key, v := range ggs {
		el := kml.Placemark(
			kml.Name(key),
			kml.Point(
				kml.Coordinates(v),
			),
		)
		k.Append(el)
	}
	for key, v := range vr {
		el := kml.Placemark(
			kml.Name(key),
			kml.Point(
				kml.Coordinates(v),
			),
		)
		k.Append(el)
		for _, val := range ggs {
			el = kml.Placemark(kml.Name("Absolute Extruded"),
				kml.StyleURL("#green"),
				kml.LineString(
					kml.Coordinates([]kml.Coordinate{
						{Lon: val.Lon, Lat: val.Lat, Alt: val.Alt},
						{Lon: v.Lon, Lat: v.Lat, Alt: v.Alt}}...),
				))
			k.Append(el)
		}
	}
	result := kml.KML(kml.Document(k))
	err = result.WriteIndent(f, "", " ")
	if err != nil {
		log.Fatal("Couldn't write xml: ", err)
	}
}

func catchErr(err error) {
	if err != nil {
		log.Fatal("Something went wrong", err)
	}
}

func createCoordinates(points [][]string) map[string]kml.Coordinate {
	c := make(map[string]kml.Coordinate, 0)
	for _, val := range points {
		lat, err := strconv.ParseFloat(val[2], 32)
		catchErr(err)
		lon, err := strconv.ParseFloat(val[3], 32)
		catchErr(err)
		c[val[1]] = kml.Coordinate{Lat: lat, Lon: lon, Alt: 0}
	}
	return c
}

func splitCoords(coords [][]string) ([][]string, [][]string) {
	ggs := make([][]string, 0)
	vr := make([][]string, 0)
	for _, row := range coords {
		if row[0] == "ggs" {
			ggs = append(ggs, row)
		} else {
			vr = append(vr, row)
		}
	}
	return ggs, vr
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("error! execution is \"./kmlCreator <path_to_file>\"")
		os.Exit(-1)
	}
	input := os.Args[1]
	f, err := os.Open(input)
	if err != nil {
		fmt.Println("error! no input file:", err)
		os.Exit(-1)
	}
	defer f.Close()
	r := csv.NewReader(f)
	coords, err := r.ReadAll()
	if err != nil {
		fmt.Println("error! something went wrong while reading file:", err)
		os.Exit(-1)
	}
	ggs, vr := splitCoords(coords)
	cGgs := createCoordinates(ggs)
	cVr := createCoordinates(vr)
	write(cGgs, cVr)
}
