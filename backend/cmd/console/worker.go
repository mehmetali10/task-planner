package main

import (
	"context"
	"sync"
	"time"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/internal/pkg/repository"
	"github.com/mehmetali10/task-planner/pkg/log"
)

// WorkerPool management
type WorkerPool struct {
	taskQueue chan payload.CreateTaskRequest
	workerNum int
	repo      repository.Repository
	wg        sync.WaitGroup // WaitGroup to track workers
	logger    log.Logger
}

// Create a new Worker Pool
func NewWorkerPool(workerNum int, repo repository.Repository) *WorkerPool {
	return &WorkerPool{
		taskQueue: make(chan payload.CreateTaskRequest, 100), // 100 buffer size
		workerNum: workerNum,
		repo:      repo,
		logger:    log.NewLogger("worker-pool", "trace"),
	}
}

// Start workers
func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workerNum; i++ {
		wp.wg.Add(1) // Add a worker
		go wp.worker(ctx, i)
	}
}

// Worker function
func (wp *WorkerPool) worker(ctx context.Context, workerID int) {
	defer wp.wg.Done() // Remove from WaitGroup when worker finishes

	wp.logger.Trace("Worker %d started...", workerID)
	for {
		select {
		case task, ok := <-wp.taskQueue:
			if !ok {
				wp.logger.Trace("Worker %d: Task queue closed, exiting...", workerID)
				return
			}

			resp, err := wp.repo.CreateTask(ctx, task)
			if err != nil {
				wp.logger.Error("Worker %d: Error creating task: %v", workerID, err)
			} else {
				wp.logger.Debug("Worker %d: Task created successfully: %+v", workerID, resp)
			}

		case <-ctx.Done():
			wp.logger.Info("Worker %d stopping...", workerID)
			time.Sleep(time.Millisecond * 100)
			return
		}
	}
}

// Function to submit a task
func (wp *WorkerPool) SubmitTask(task payload.CreateTaskRequest) {
	wp.taskQueue <- task
}

// Stop the WorkerPool
func (wp *WorkerPool) Stop() {
	close(wp.taskQueue) // Close the queue, workers will start shutting down
	wp.wg.Wait()        // Wait for all workers to finish
}
