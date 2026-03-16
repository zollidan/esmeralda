package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/queue"
)

func (b *Bot) ListTasks() ([]models.Task, error) {
	var taskList []models.Task

	req, err := http.NewRequest("GET", b.cfg.TelegramBot.APIBaseURL+"/api/tasks", nil)
	if err != nil {
		return nil, err
	}

	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&taskList)
	if err != nil {
		return nil, err
	}

	return taskList, nil
}

func (b *Bot) CreateTask(taskDate string) (*queue.ParseTask, error) {
	body := map[string]string{
		"date": taskDate,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		b.cfg.TelegramBot.APIBaseURL+"/api/tasks",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("api error: %s", string(body))
	}

	var task queue.ParseTask
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (b *Bot) GetTask(taskID string) (*models.Task, error) {
	req, err := http.NewRequest("GET", b.cfg.TelegramBot.APIBaseURL+"/api/tasks/"+taskID, nil)
	if err != nil {
		return nil, err
	}

	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("api error: %s", string(body))
	}

	var task models.Task
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (b *Bot) WaitForTask(taskID string, pollInterval, timeout time.Duration) (*models.Task, error) {
	deadline := time.Now().Add(timeout)

	for {
		task, err := b.GetTask(taskID)
		if err != nil {
			return nil, err
		}

		switch task.Status {
		case string(queue.StatusDone), string(queue.StatusError), string(queue.StatusFailed), string(queue.StatusCancelled):
			return task, nil
		}

		if timeout > 0 && time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for task %s", taskID)
		}

		time.Sleep(pollInterval)
	}
}

func (b *Bot) DownloadExport(date string) ([]byte, string, error) {
	url := fmt.Sprintf("%s/api/export?date_start=%s&date_end=%s", b.cfg.TelegramBot.APIBaseURL, date, date)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	res, err := b.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, "", fmt.Errorf("api error: %s", string(body))
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("esmeralda_%s.xlsx", date)
	if header := res.Header.Get("Content-Disposition"); header != "" {
		_, _ = fmt.Sscanf(header, "attachment; filename=%q", &filename)
	}

	return content, filename, nil
}