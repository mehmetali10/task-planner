package main

import (
	"context"
	"log"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	postgresRepo "github.com/mehmetali10/task-planner/internal/pkg/repository/postgres"
)

// WorkerPool yönetimi
type WorkerPool struct {
	taskQueue chan payload.CreateTaskRequest
	workerNum int
	repo      *postgresRepo.PostgresRepo
}

// Yeni bir Worker Pool oluştur
func NewWorkerPool(workerNum int, repo *postgresRepo.PostgresRepo) *WorkerPool {
	return &WorkerPool{
		taskQueue: make(chan payload.CreateTaskRequest, 100), // 100 buffer size
		workerNum: workerNum,
		repo:      repo,
	}
}

// Worker'ları başlat
func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workerNum; i++ {
		go wp.worker(ctx, i)
	}
}

// Worker fonksiyonu
func (wp *WorkerPool) worker(ctx context.Context, workerID int) {
	log.Printf("Worker %d started...\n", workerID)

	for {
		select {
		case task := <-wp.taskQueue:
			_, err := wp.repo.CreateTask(ctx, task)
			if err != nil {
				log.Printf("Worker %d: Error creating task: %v\n", workerID, err)
			} else {
				log.Printf("Worker %d: Task created successfully: %+v\n", workerID, task)
			}

		case <-ctx.Done():
			log.Printf("Worker %d stopping...\n", workerID)
			return
		}
	}
}

// Task gönderme fonksiyonu
func (wp *WorkerPool) SubmitTask(task payload.CreateTaskRequest) {
	wp.taskQueue <- task
}

// WorkerPool'u durdur
func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
}
