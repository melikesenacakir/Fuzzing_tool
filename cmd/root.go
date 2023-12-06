package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Site struct {
	URL string
}
type Result struct {
	URL    string
	Status int
}

func worker(jobs <-chan Site, results chan<- Result) {
	for site := range jobs {
		resp, err := http.Get(site.URL)
		if err != nil {
			color.Red("\nWebsite not found!")
			os.Exit(0)
		}
		results <- Result{URL: site.URL, Status: resp.StatusCode}
	}
}
func fuzz(url string, st int, file string, speed int, output string) {
	list, err := os.OpenFile(file, os.O_APPEND, 0644)
	file2, err2 := os.OpenFile(output, os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		color.Red("\nFile not found!")
	} else {
		defer list.Close()
		scanner := bufio.NewScanner(list)
		jobs := make(chan Site)
		results := make(chan Result)
		for i := 1; i <= 3; i++ {
			go worker(jobs, results)
		}

		for scanner.Scan() {
			line := scanner.Text()
			full_url := url + "/" + line
			jobs <- Site{URL: full_url}
			res := <-results
			if res.Status == st {
				fmt.Println(res.URL)
				file2.WriteString("URL: " + res.URL + "\n" + "STATUS CODE: " + fmt.Sprintf("%d", res.Status) + "\n\n")
			} else if st == 0 {
				fmt.Println(res.URL + " - " + fmt.Sprintf("%d", res.Status))
				file2.WriteString("URL: " + res.URL + "\n" + "STATUS CODE: " + fmt.Sprintf("%d", res.Status) + "\n\n")
			}

			time.Sleep(time.Second / time.Duration(speed))
		}
		close(jobs)
	}
	if err2 == nil {
		color.Yellow(">> Output saved in " + output + " <<")
	}
	color.Cyan("\nFuzzing completed")
}

var rootCmd = &cobra.Command{
	Use:   "fuzzer_go",
	Short: "Fuzz program",
	Long:  `Scans all paths in the given url`,

	Run: func(cmd *cobra.Command, args []string) {
		txtfile, _ := cmd.Flags().GetString("txt")
		geturl, _ := cmd.Flags().GetString("u")
		status, _ := cmd.Flags().GetInt("st")
		speed, _ := cmd.Flags().GetInt("s")
		output, _ := cmd.Flags().GetString("save")
		fuzz(geturl, status, txtfile, speed, output)

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(0)
	}
}

func init() {
	rootCmd.Flags().String("txt", "wordlist.txt", "Path of txt which will be used")
	rootCmd.Flags().String("u", "http://localhost", "Define url to scan")
	rootCmd.Flags().Int("s", 72, "Scanning speed")
	rootCmd.Flags().Int("st", 0, "Filters given status code")
	rootCmd.Flags().String("save", "", "Save output in a file")
}
