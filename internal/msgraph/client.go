package msgraph

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"roombooker/internal/config"
)

type Client struct {
	graphClient *msgraphsdk.GraphServiceClient
}

func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.Graph.ClientID == "" || cfg.Graph.ClientSecret == "" || cfg.Graph.TenantID == "" {
		return nil, nil // No Graph integration
	}

	cred, err := azidentity.NewClientSecretCredential(cfg.Graph.TenantID, cfg.Graph.ClientID, cfg.Graph.ClientSecret, nil)
	if err != nil {
		return nil, err
	}

	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return nil, err
	}

	return &Client{graphClient: graphClient}, nil
}

func (c *Client) CreateEvent(resourceID, subject, startTime, endTime string) (*models.Eventable, error) {
	if c == nil || c.graphClient == nil {
		return nil, nil
	}

	event := models.NewEvent()
	event.SetSubject(&subject)

	start := models.NewDateTimeTimeZone()
	start.SetDateTime(&startTime)
	start.SetTimeZone(&[]string{"UTC"}[0])
	event.SetStart(start)

	end := models.NewDateTimeTimeZone()
	end.SetDateTime(&endTime)
	end.SetTimeZone(&[]string{"UTC"}[0])
	event.SetEnd(end)

	result, err := c.graphClient.Users().ByUserId(resourceID).Events().Post(context.Background(), event, nil)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Add more methods for updating, deleting events
