package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"multi-tenant-messaging-app/internal/repository"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type TenantService struct {
	conn       *amqp.Connection
	consumers  map[string]context.CancelFunc
	workerPool map[string]int
	mu         sync.Mutex
	repo       *repository.MessageRepository
}

func NewTenantService(repo *repository.MessageRepository, conn *amqp.Connection) *TenantService {
	return &TenantService{
		conn:       conn,
		repo:       repo,
		consumers:  make(map[string]context.CancelFunc),
		workerPool: make(map[string]int),
	}
}

func (tm *TenantService) StartConsumer(tenantID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// // Buat partisi untuk tenant
	// if err := tm.repo.CreateTenantPartition(tenantID); err != nil {
	// 	return fmt.Errorf("failed to create tenant partition: %w", err)
	// }

	if _, exists := tm.consumers[tenantID]; exists {
		log.Printf("consumer already exists for tenant %s", tenantID)
		return nil
	}

	ch, err := tm.conn.Channel()
	if err != nil {
		return err
	}

	queueName := fmt.Sprintf("tenant_%s_queue", tenantID)
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	tm.consumers[tenantID] = cancel

	go tm.consumeTenantQueue(ctx, tenantID, queueName)

	return nil
}

func (tm *TenantService) consumeTenantQueue(ctx context.Context, tenantID, queueName string) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopped consumer for tenant %s", tenantID)
			return
		default:
		}

		ch, err := tm.conn.Channel()
		if err != nil {
			log.Println("channel error:", err)
			time.Sleep(time.Second)
			continue
		}

		msgs, err := ch.Consume(queueName, fmt.Sprintf("consumer_%s", tenantID), true, false, false, false, nil)
		if err != nil {
			log.Println("consume error:", err)
			ch.Close()
			time.Sleep(time.Second)
			continue
		}

		for msg := range msgs {
			var payload map[string]interface{}
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				log.Printf("Failed to unmarshal message for tenant %s: %v", tenantID, err)
				continue
			}

			if err := tm.repo.SaveMessage(tenantID, payload); err != nil {
				log.Printf("Failed to save message for tenant %s: %v", tenantID, err)
			}
		}

		ch.Close()
	}
}

func (tm *TenantService) GetAllTenantIDs() []string {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	ids := make([]string, 0, len(tm.consumers))
	for id := range tm.consumers {
		ids = append(ids, id) // id sudah string
	}
	return ids
}

func (tm *TenantService) PublishToTenantQueue(tenantID string, payload map[string]interface{}) error {
	ch, err := tm.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	queueName := fmt.Sprintf("tenant_%s_queue", tenantID)
	return ch.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// UpdateWorkerCount memperbarui jumlah worker untuk tenant tertentu
func (tm *TenantService) UpdateWorkerCount(tenantID string, workerCount int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Periksa apakah tenant memiliki konsumen aktif
	cancelFunc, exists := tm.consumers[tenantID]
	if !exists {
		return fmt.Errorf("no active consumer found for tenant %s", tenantID)
	}

	// Hentikan worker lama
	cancelFunc()

	// Buat context baru untuk worker
	ctx, cancel := context.WithCancel(context.Background())
	tm.consumers[tenantID] = cancel

	// Mulai worker baru sesuai jumlah yang diatur
	for i := 0; i < workerCount; i++ {
		go tm.worker(ctx, tenantID)
	}

	// Perbarui jumlah worker di workerPool
	tm.workerPool[tenantID] = workerCount

	log.Printf("Updated worker count for tenant %s to %d", tenantID, workerCount)
	return nil
}

// worker adalah fungsi yang dijalankan oleh setiap worker untuk memproses pesan
func (tm *TenantService) worker(ctx context.Context, tenantID string) {
	queueName := fmt.Sprintf("tenant_%s_queue", tenantID)

	ch, err := tm.conn.Channel()
	if err != nil {
		log.Printf("Failed to open channel for tenant %s: %v", tenantID, err)
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to consume messages for tenant %s: %v", tenantID, err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker stopped for tenant %s", tenantID)
			return
		case msg := <-msgs:
			// Proses pesan
			log.Printf("Processing message for tenant %s: %s", tenantID, string(msg.Body))
			// Tambahkan logika pemrosesan pesan di sini
		}
	}
}

func (tm *TenantService) StopConsumer(tenantID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	cancel, exists := tm.consumers[tenantID]
	if !exists {
		return fmt.Errorf("consumer does not exist for tenant %s", tenantID)
	}

	cancel()
	delete(tm.consumers, tenantID)
	log.Printf("Stopped consumer for tenant %s", tenantID)

	ch, err := tm.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queueName := fmt.Sprintf("tenant_%s_queue", tenantID)
	if _, err := ch.QueueDelete(queueName, false, false, false); err != nil {
		return err
	}

	return nil
}
