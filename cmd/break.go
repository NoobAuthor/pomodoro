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
	breakMinutes int
	breakCmd     = &cobra.Command{
		Use:   "break",
		Short: "Start a break session",
		Long:  `Start a Pomodoro break session with a countdown timer and progress bar.`,
		Run:   runBreakSession,
	}
)

func init() {
	rootCmd.AddCommand(breakCmd)
	breakCmd.Flags().IntVarP(&breakMinutes, "minutes", "m", 5, "Duration of break session in minutes")
}

func runBreakSession(cmd *cobra.Command, args []string) {
	duration := time.Duration(breakMinutes) * time.Minute
	durationSeconds := int(duration.Seconds())

	logger.Info("Starting break session",
		zap.Int("minutes", breakMinutes),
		zap.Duration("duration", duration),
	)

	fmt.Printf("☕ Starting %d-minute break session...\n", breakMinutes)

	// Create progress bar
	p := mpb.New(mpb.WithWidth(60))
	bar := p.AddBar(int64(durationSeconds),
		mpb.PrependDecorators(
			decor.Name("Break: "),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60, decor.WCSyncSpace), "⏰ Time's up!",
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

	logger.Info("Break session completed",
		zap.Int("minutes", breakMinutes),
		zap.Duration("actual_duration", time.Since(start)),
	)

	fmt.Println("🚀 Break time over! Ready to get back to work?")
}
