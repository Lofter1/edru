package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Lofter1/edru/simpleSftp"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/term"
)

var (
	romPath        string
	romsBaseFolder string
	emuSystem      string
	host           string
	port           string
	user           string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
	}
	err = loadConfig(path.Join(homeDir, ".edru.conf.json"))
	if err != nil {
		log.Println(err)
	}

	flag.StringVar(&romPath, "rom", "", "path to the rom (if the path is a directory, the entire directory will be copied)")
	flag.StringVar(&romsBaseFolder, "rom-folder", "/home/deck/Emulation/roms/", "EmuDeck ROM folder on the SteamDeck")
	flag.StringVar(&emuSystem, "system", "", "The Emulation system for the rom")
	flag.StringVar(&host, "sd-host", "steamdeck", "Hostname or IP of the SteamDeck")
	flag.StringVar(&port, "sd-ssh-port", "22", "SSH Port of the SteamDeck")
	flag.StringVar(&user, "sd-user", "deck", "Username of the SteamDeck user")

	flag.Parse()
}

func main() {
	remoteDestination := path.Join(romsBaseFolder, emuSystem, path.Base(romPath))

	passwd, err := readPassword()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connect to SteamDeck...")
	client, err := simpleSftp.ConnectWithPassword(user, passwd, host, port)
	if err != nil {
		log.Fatalf("Failed to connect to steamdeck: %v", err)
	}
	defer client.Close()

	fmt.Println("Upload file(s)...")

	var currentBar *progressbar.ProgressBar
	var previousFile string

	progressFunc := func(file string, bytesTransferred int, totalBytes int) {
		if file != previousFile {
			previousFile = file
			currentBar = progressbar.DefaultBytes(int64(totalBytes), file)
		}

		err := currentBar.Add(bytesTransferred)
		if err != nil {
			log.Print(err)
		}
	}

	err = client.PutProgress(romPath, remoteDestination, progressFunc)
	if err != nil {
		log.Fatalf("Failed to copy: %v", err)
	}
}

func readPassword() (string, error) {
	fmt.Print("Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return string(password), err
}

func loadConfig(location string) error {
	dat, err := os.ReadFile(location)
	if err != nil {
		return err
	}

	var config map[string]any
	if err := json.Unmarshal(dat, &config); err != nil {
		return err
	}

	romsBaseFolder = config["romFolder"].(string)

	return nil
}
