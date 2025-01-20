package notifier

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"k8s.io/client-go/transport"

	"github.com/openshift-kni/oran-o2ims/internal/controllers/utils"
)

// ClientFactory is a utility used to abstract building an HTTP client based on the type of callback
// URL supplied.
type ClientFactory struct {
	oauthConfig      *utils.OAuthClientConfig
	serviceTokenFile string
}

// ClientProvider defines the interface which any client factory must implement.  This exists for
// future unit test purposes so that the ClientFactory can be swapped out as needed.
type ClientProvider interface {
	NewClient(ctx context.Context, callbackURL string) (*http.Client, error)
}

// NewClientFactory creates a new factory
func NewClientFactory(oauthConfig *utils.OAuthClientConfig, serviceTokenFile string) ClientProvider {
	return &ClientFactory{
		oauthConfig:      oauthConfig,
		serviceTokenFile: serviceTokenFile,
	}
}

func (f *ClientFactory) newClusterClient(ctx context.Context) (*http.Client, error) {
	tlsConfig, _ := utils.GetDefaultTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12})
	baseClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 30 * time.Second,
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, baseClient)
	return oauth2.NewClient(ctx, transport.NewCachedFileTokenSource(f.serviceTokenFile)), nil
}

func (f *ClientFactory) newOAuthClient(ctx context.Context) (*http.Client, error) {
	client, err := utils.SetupOAuthClient(ctx, f.oauthConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to setup oauth client")
	}
	return client, nil
}

// NewClient creates a new Client based on the callback URL provided.  If the callback URL is a local
// service URL that contains "svc.cluster.local" then a Client will be created that uses the
// supplied service account token file; otherwise, it is assumed that the URL points to a public
// endpoint that requires the OAuth credentials.
func (f *ClientFactory) NewClient(ctx context.Context, callback string) (*http.Client, error) {
	if strings.Contains(callback, "svc.cluster.local") {
		return f.newClusterClient(ctx)
	}
	return f.newOAuthClient(ctx)
}