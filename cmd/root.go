package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "see",
	Short: "gosee is a fuzzy finder that will allow you to preview files in terminal.",
	Long:  `gosee is a fuzzy finder that will allow you to preview files in terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		var filelist []string

		filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err == nil {
				if !d.IsDir() {
					filelist = append(filelist, path)
				}

			}

			return nil
		})

		idx, err := fuzzyfinder.FindMulti(
			filelist,
			func(i int) string {
				return filelist[i]
			},
			fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
				if i == -1 {
					return ""
				}
				contents, readErr := os.ReadFile(filelist[i])

				if readErr != nil {
					log.Fatal(readErr)
				}
				return fmt.Sprintf("%s %s",
					filelist[i], string(markdown.Render("```\n"+string(contents)+"\n```", 200, 2)))
			}))

		if err != nil {
			log.Fatal(err)
		}

		selected := filelist[idx[len(idx)-1]]

		contents, readErr := os.ReadFile(selected)

		if readErr != nil {
			log.Fatal(readErr)
		}

		out := markdown.Render("```\n"+string(contents)+"\n```", 200, 2)
		fmt.Print(string(out))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
