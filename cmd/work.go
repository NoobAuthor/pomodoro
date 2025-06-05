package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"go.uber.org/zap"
)

var (
	workMinutes int
	workCmd     = &cobra.Command{
		Use:   "work",
		Short: "Start a work session",
		Long:  `Start a Pomodoro work session with a countdown timer and progress bar.`,
		Run:   runWorkSession,
	}
)

func init() {
	rootCmd.AddCommand(workCmd)
	workCmd.Flags().IntVarP(&workMinutes, "minutes", "m", 25, "Duration of work session in minutes")
}

func runWorkSession(cmd *cobra.Command, args []string) {
	duration := time.Duration(workMinutes) * time.Minute
	durationSeconds := int(duration.Seconds())

	logger.Info("Starting work session",
		zap.Int("minutes", workMinutes),
		zap.Duration("duration", duration),
	)

	fmt.Printf("🍅 Starting %d-minute work session...\n", workMinutes)

	// Create progress bar
	p := mpb.New(mpb.WithWidth(60))
	bar := p.AddBar(int64(durationSeconds),
		mpb.PrependDecorators(
			decor.Name("Work: "),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60, decor.WCSyncSpace), "✅ Done!",
			),
		),
	)

	// Timer loop
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	start := time.Now()
	for elapsed := time.Duration(0); elapsed < duration; elapsed = time.Since(start) {
		select {
		case <-ticker.C:
			progress := int64(elapsed.Seconds())
			bar.SetCurrent(progress)
		}
	}

	bar.SetCurrent(int64(durationSeconds))
	p.Wait()

	logger.Info("Work session completed",
		zap.Int("minutes", workMinutes),
		zap.Duration("actual_duration", time.Since(start)),
	)

	fmt.Println("🎉 Work session complete! Time for a break.")
}
