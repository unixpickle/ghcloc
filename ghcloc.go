package main

import (
	"fmt"
	"github.com/unixpickle/ghcloc/lib"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "Usage: ghcloc user/repo [path]")
		os.Exit(1)
	}
	
	repo := getRepo()
	
	directory := "/"
	if len(os.Args) == 3 {
		directory = os.Args[2]
	}
	
	fmt.Println("Counting...")
	if counts, err := ghcloc.CountInDir(repo, directory); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to count lines in "+directory+": "+
			err.Error())
		os.Exit(1)
	} else {
		fmt.Printf("\nTotal line counts (%d files):\n\n", counts.FileCount)
		printTable(counts.TotalLines)
		fmt.Println("")
	}
}

func getRepo() *ghcloc.Repository {
	var username string
	var password string
	fmt.Print("Github username: ")
	fmt.Scanln(&username)
	fmt.Print("Github password: ")
	setTTYEcho(false)
	fmt.Scanln(&password)
	setTTYEcho(true)
	fmt.Println("")

	repo, err := ghcloc.ParseRepository(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid repository: "+err.Error())
		os.Exit(1)
	}
	repo.Authenticate(username, password)
	return repo
}

func printTable(table map[string]int) {
	maxLen := 0
	for key := range table {
		if len(key) > maxLen {
			maxLen = len(key)
		}
	}
	for key, value := range table {
		for i := 0; i < maxLen-len(key); i++ {
			fmt.Print(" ")
		}
		fmt.Printf("  %s %d", key, value)
		fmt.Println("")
	}
}

func setTTYEcho(enabled bool) {
	stty, err := exec.LookPath("stty")
	if err != nil {
		fmt.Println("popy")
		return
	}
	arg := "echo"
	if !enabled {
		arg = "-echo"
	}
	cmd := exec.Command(stty, arg)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Run()
}
