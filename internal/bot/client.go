package bot

type APIClient struct {}

func NewRequestClient(apiBaseURL, apiToken string) *APIClient {
	return &APIClient{}
}

func (c *APIClient) ListTasks() ([]string, error) 
func (c *APIClient) CreateTask(taskName string) error