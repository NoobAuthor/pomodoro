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
	cycleWorkMinutes  int
	cycleBreakMinutes int
	cycles            int
	cycleCmd          = &cobra.Command{
		Use:   "cycle",
		Short: "Run a full Pomodoro cycle",
		Long:  `Run multiple Pomodoro cycles with work and break sessions.`,
		Run:   runPomodoroyCycle,
	}
)

func init() {
	rootCmd.AddCommand(cycleCmd)
	cycleCmd.Flags().IntVarP(&cycleWorkMinutes, "work", "w", 25, "Duration of work sessions in minutes")
	cycleCmd.Flags().IntVarP(&cycleBreakMinutes, "break", "b", 5, "Duration of break sessions in minutes")
	cycleCmd.Flags().IntVarP(&cycles, "cycles", "c", 4, "Number of work/break cycles to complete")
}

func runPomodoroyCycle(cmd *cobra.Command, args []string) {
	logger.Info("Starting Pomodoro cycle",
		zap.Int("work_minutes", cycleWorkMinutes),
		zap.Int("break_minutes", cycleBreakMinutes),
		zap.Int("cycles", cycles),
	)

	fmt.Printf("🍅 Starting %d Pomodoro cycles (%dm work, %dm break)\n",
		cycles, cycleWorkMinutes, cycleBreakMinutes)

	for i := 1; i <= cycles; i++ {
		fmt.Printf("\n--- Cycle %d/%d ---\n", i, cycles)

		// Work session
		runTimerSession("Work", cycleWorkMinutes, "🍅")

		if i < cycles {
			// Break session (skip break after last cycle)
			runTimerSession("Break", cycleBreakMinutes, "☕")
		}
	}

	logger.Info("Pomodoro cycles completed",
		zap.Int("completed_cycles", cycles),
	)

	fmt.Println("\n🎉 All Pomodoro cycles completed! Great work!")
}

func runTimerSession(sessionType string, minutes int, emoji string) {
	duration := time.Duration(minutes) * time.Minute
	durationSeconds := int(duration.Seconds())

	fmt.Printf("%s Starting %d-minute %s session...\n", emoji, minutes, sessionType)

	// Create progress bar
	p := mpb.New(mpb.WithWidth(60))
	bar := p.AddBar(int64(durationSeconds),
		mpb.PrependDecorators(
			decor.Name(fmt.Sprintf("%s: ", sessionType)),
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

	logger.Info("Session completed",
		zap.String("type", sessionType),
		zap.Int("minutes", minutes),
		zap.Duration("actual_duration", time.Since(start)),
	)
}
