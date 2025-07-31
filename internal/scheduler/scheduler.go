package scheduler

import (
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	config *Config

	channel chan struct{}
	mutex   sync.Mutex
	running bool
	task    func() error
	ticker  *time.Ticker
	wg      sync.WaitGroup
}

func NewScheduler(config *Config, task func() error) *Scheduler {
	return &Scheduler{
		config: config,

		running: false,
		task:    task,
	}
}

func (s *Scheduler) Start() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running == true {
		return "Scheduler already running."
	}

	duration := time.Duration(s.config.Interval) * time.Second

	s.channel = make(chan struct{})
	s.ticker = time.NewTicker(duration)
	s.running = true

	log.Printf("[INFO] Scheduler is starting...")
	go s.run()

	return "Scheduler started successfully"
}

func (s *Scheduler) Stop() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running == false {
		return "Scheduler already stopped."
	}

	log.Println("[INFO] Scheduler is stopping...")
	close(s.channel)
	s.running = false
	s.wg.Wait()

	return "Scheduler stopped successfully"
}

func (s *Scheduler) run() {
	for {
		select {
		case <-s.ticker.C:
			s.wg.Add(1)
			func() {
				defer s.wg.Done()

				if err := s.task(); err != nil {
					log.Printf("[ERROR] Scheduler task failed: %v", err)
				}
			}()
		case <-s.channel:
			s.ticker.Stop()
			return
		}
	}
}
