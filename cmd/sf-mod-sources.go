package main

import (
	"context"
	"encoding/json"
	"fmt"
	"git.sr.ht/~emersion/gqlclient"
	"log"
	"net/http"
	"os"
	"sf-mod-sources/gql"
)

type ModSource struct {
	Name string
	URL  string
}

var c int

func main() {
	if len(os.Args) < 2 {
		log.Fatal("This program takes a path as an argument")
	}

	out := os.Args[1]

	client := gqlclient.New("https://api.ficsit.app/v2/query", http.DefaultClient)

	var r []ModSource

	var offset int

	fmt.Print("Getting page 1... ")
	mods := getMods(client, offset)
	count := appendModSources(&r, mods)
	fmt.Printf("Found %v mods of which %v have a source URL\n", len(mods), count)

	for len(mods) == 100 {
		offset += 100
		fmt.Printf("Getting page %v... ", offset+1)
		mods = getMods(client, offset)
		count := appendModSources(&r, mods)
		fmt.Printf("Found %v mods of which %v have a source URL\n", len(mods), count)
	}

	fmt.Println("Done searching mods!")

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal the data into JSON: %v", err)
	}

	err = os.WriteFile(out, data, 0644)
	if err != nil {
		log.Fatalf("Could not write the data to file: %v", err)
	}
}

func getMods(client *gqlclient.Client, offset int) []gql.Mod {
	data, err := gql.QGetMods(client, context.Background(), int32(offset))
	checkErr(err)
	return data.Mods
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Could not query mods: %v", err)
	}
}

func appendModSources(sources *[]ModSource, mods []gql.Mod) int {
	newSources := modsToModSourceSlice(mods)
	*sources = append(*sources, newSources...)
	return len(newSources)
}

func modsToModSourceSlice(mods []gql.Mod) []ModSource {
	r := make([]ModSource, 0, len(mods)/2)
	c += len(mods)

	for _, mod := range mods {
		source := modToModSource(mod)
		if source != nil {
			r = append(r, *source)
		}
	}
	return r
}

func modToModSource(mod gql.Mod) *ModSource {
	if mod.Source_url == nil || *mod.Source_url == "" {
		return nil
	}
	return &ModSource{
		Name: mod.Name,
		URL:  *mod.Source_url,
	}
}
