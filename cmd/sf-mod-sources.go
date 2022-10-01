package main

import (
	"context"
	"fmt"
	"git.sr.ht/~emersion/gqlclient"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sf-mod-sources/gql"
)

type ModVersions struct {
	Name     string
	Versions []ModVersion
}

type ModVersion struct {
	Version string
	URL     string
}

var c int

func main() {
	//if len(os.Args) < 2 {
	//	log.Fatal("This program takes a path as an argument")
	//}

	//out := os.Args[1]

	client := gqlclient.New("https://api.ficsit.app/v2/query", http.DefaultClient)

	var r []ModVersions

	var offset int

	fmt.Print("Getting page 1... ")
	mods := getMods(client, offset)
	count := appendModVersions(&r, mods)
	fmt.Printf("Found %v mods of which %v have a source URL\n", len(mods), count)

	for len(mods) == 100 {
		offset += 100
		fmt.Printf("Getting page %v... ", offset/100+1)
		mods = getMods(client, offset)
		appendModVersions(&r, mods)
		fmt.Printf("Found %v mods\n", len(mods))
	}

	fmt.Println("Done searching mods!")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get the working directory: %v", err)
	}

	//errChan := make(chan error)
	//doneChan := make(chan struct{})
	//pool := workerpool.New(runtime.NumCPU() * 4)
	for _, mod := range r {
		for _, version := range mod.Versions {
			versionDir := filepath.Join(wd, mod.Name, version.Version)

			err := os.MkdirAll(versionDir, 0777)
			if err != nil && !os.IsExist(err) {
				log.Fatalf("Error making a mod's folder: %v", err)
			}
			err = downloadLinkTo("https://api.ficsit.app"+version.URL, filepath.Join(versionDir, version.Version+".zip"))
			if err != nil {
				log.Fatalf("Error downloading a mod: %v", err)
			}
		}
	}

	//go func() {
	//	pool.StopWait()
	//	doneChan <- struct{}{}
	//}()

	fmt.Println("Done!")
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

func appendModVersions(sources *[]ModVersions, mods []gql.Mod) int {
	newSources := modsToModVersionsSlice(mods)
	*sources = append(*sources, newSources...)
	return len(newSources)
}

func modsToModVersionsSlice(mods []gql.Mod) []ModVersions {
	r := make([]ModVersions, 0, len(mods)/2)

	for _, mod := range mods {
		source := modToModVersions(mod)
		if source != nil {
			r = append(r, *source)
		}
	}
	return r
}

func modToModVersions(mod gql.Mod) *ModVersions {
	if len(mod.Versions) == 0 {
		return nil
	}
	return &ModVersions{
		Name:     mod.Name,
		Versions: versionsToModVersionSlice(mod.Versions),
	}
}

func versionsToModVersionSlice(versions []gql.Version) []ModVersion {
	//r := make([]ModVersion, len(versions))
	//for i, version := range versions {
	//	r[i] = versionToModVersion(version)
	//}
	//return r
	return []ModVersion{versionToModVersion(versions[0])}
}

func versionToModVersion(version gql.Version) ModVersion {
	return ModVersion{
		Version: version.Version,
		URL:     version.Link,
	}
}

func downloadLinkTo(link, to string) error {
	resp, err := http.Get(link)
	if err != nil {
		return errors.Wrap(err, "could not downloads the file")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "could not downloads the file")
	}

	err = os.WriteFile(to, b, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write the file")
	}

	return nil
}
