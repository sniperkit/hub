package progress

import (
	"os"
	"os/signal"
	"syscall"

	pb "github.com/gosuri/uiprogress"

	"github.com/sniperkit/hub/plugin/assitant/core"
)

// ProgressBar Implements core.Progress
type ProgressBar struct {
	Impl        *pb.Bar
	isCancelled bool
}

func (p *ProgressBar) Init(total int) {
	if !core.IsQuietModeEnabled() {
		pb.Start()
	}
	p.Impl = pb.AddBar(total)
	p.Impl.AppendCompleted()
	p.Impl.PrependElapsed()
	p.isCancelled = false

	// TODO: move that code into a dedicated function
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-c
		p.isCancelled = true
		p.Done()
		core.Info("Cancelling...")
	}()
}

func (p *ProgressBar) Increment() {
	p.Impl.Incr()
}

func (p *ProgressBar) Done() {
	if !core.IsQuietModeEnabled() {
		pb.Stop()
	}
}

func (p *ProgressBar) IsCancelled() bool {
	return p.isCancelled
}
