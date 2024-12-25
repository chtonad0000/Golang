package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/slon/shad-go/gitfame/internal/blameparser"
	"gitlab.com/slon/shad-go/gitfame/internal/git"
	"gitlab.com/slon/shad-go/gitfame/internal/output"
	"gitlab.com/slon/shad-go/gitfame/internal/stats"
)

func main() {
	// Флаги
	repoPath := flag.String("repository", ".", "Path to the Git repository")
	revision := flag.String("revision", "HEAD", "Git commit revision")
	orderBy := flag.String("order-by", "lines", "Key to sort results (lines, commits, files)")
	useCommitter := flag.Bool("use-committer", false, "Use committer instead of author")
	outputFormat := flag.String("format", "tabular", "Output format (tabular, csv, json, json-lines)")
	extensions := flag.String("extensions", "", "Comma-separated list of file extensions to include")
	languages := flag.String("languages", "", "Comma-separated list of languages to include")
	excludePatterns := flag.String("exclude", "", "Comma-separated list of glob patterns to exclude files")
	restrictPatterns := flag.String("restrict-to", "", "Comma-separated list of glob patterns to include files")
	showProgress := flag.Bool("progress", false, "Show progress")

	flag.Parse()

	// Проверка флагов
	if *orderBy != "lines" && *orderBy != "commits" && *orderBy != "files" {
		fmt.Fprintln(os.Stderr, "Invalid value for --order-by. Must be one of: lines, commits, files.")
		os.Exit(1)
	}

	if *outputFormat != "tabular" && *outputFormat != "csv" && *outputFormat != "json" && *outputFormat != "json-lines" {
		fmt.Fprintln(os.Stderr, "Invalid value for --format. Must be one of: tabular, csv, json, json-lines.")
		os.Exit(1)
	}

	// Получение списка файлов
	files, err := git.GetFiles(*repoPath, *extensions, *languages, *excludePatterns, *restrictPatterns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting files:", err)
		os.Exit(1)
	}

	// Получение данных blame
	blameData, err := blameparser.GetBlameData(*repoPath, *revision, files, *useCommitter, *showProgress)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting blame data:", err)
		os.Exit(1)
	}

	// Подсчёт статистик
	statistics := stats.CalculateStatistics(blameData)

	// Вывод результатов
	if err := output.Render(statistics, *outputFormat, *orderBy); err != nil {
		fmt.Fprintln(os.Stderr, "Error rendering output:", err)
		os.Exit(1)
	}
}
