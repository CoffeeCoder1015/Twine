package main

import (
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type DebounceRefresh struct{ }
type DebounceLoading struct{ }

// Debouncer manages scheduling and executing a debounced function.
type Debouncer struct {
    updateRequests bool
    mu             sync.Mutex    // Protects internal state, specifically the activeWorker flag
    activeWorker   bool          // True if a worker goroutine is currently running
    debounceDelay  time.Duration // The time to wait before executing the update
    task func()
    d *tea.Program
}

func NewDebouncer(delay time.Duration) *Debouncer {
    return &Debouncer{
        // A buffered channel helps absorb rapid calls without blocking the caller.
        // A size of 1 is often sufficient for simple debouncing.
        debounceDelay:  delay,
    }
}

// ScheduleUpdate is called by the "update stream" to request an update.
func (d *Debouncer) ScheduleUpdate(task func()) {
    d.updateRequests = true
    d.startWorker()
    d.task = task
}

// startWorker ensures that only one worker goroutine is running.
func (d *Debouncer) startWorker() {
    d.mu.Lock()
    defer d.mu.Unlock()

    if d.activeWorker {
        return // Worker already active
    }
    d.activeWorker = true
    go d.worker() // Start the worker goroutine
}

func (d *Debouncer) worker() bool{
    defer func() {
            d.mu.Lock()
            d.activeWorker = false // Mark worker as inactive when it exits
            d.mu.Unlock()
    }()
    for {
        <-time.After(d.debounceDelay)
        if d.updateRequests {
            d.updateRequests = false
            d.d.Send(DebounceLoading{})
            continue
        }
        d.task()
        d.d.Send(DebounceRefresh{})
        break
    }
    return true
}
