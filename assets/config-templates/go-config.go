// Go GLIDE Configuration Template
// Optimized for production web applications

package main

import (
	"context"
	"time"

	"github.com/valkey-io/valkey-glide/go/glide/api"
)

// Standalone Client Configuration
func standaloneConfig() *api.GlideClientConfiguration {
	return &api.GlideClientConfiguration{
		Addresses: []api.NodeAddress{
			{Host: "localhost", Port: 6379},
		},
		
		// Request timeout (500ms recommended for web apps)
		RequestTimeout: 500 * time.Millisecond,
		
		// Connection retry strategy
		ConnectionBackoff: &api.BackoffStrategy{
			NumberOfRetries: 10,
			Factor:          500,  // Base delay in ms
			ExponentBase:    2,    // Exponential backoff
		},
		
		// Client name for debugging
		ClientName: "my-app-client",
		
		// Lazy connect for serverless/Lambda
		LazyConnect: false,  // Set to true for Lambda
		
		// High-throughput configuration
		InflightRequestsLimit: 2000,  // Default: 1000
	}
}

// Cluster Client Configuration
func clusterConfig() *api.GlideClusterClientConfiguration {
	return &api.GlideClusterClientConfiguration{
		Addresses: []api.NodeAddress{
			{Host: "cluster.endpoint.cache.amazonaws.com", Port: 6379},
		},
		
		// Request timeout
		RequestTimeout: 500 * time.Millisecond,
		
		// AZ Affinity for cost optimization (read-heavy workloads)
		ReadFrom: api.AZAffinity,
		ClientAZ: "us-east-1a",  // Your application's AZ
		
		// Connection retry strategy
		ConnectionBackoff: &api.BackoffStrategy{
			NumberOfRetries: 10,
			Factor:          500,
			ExponentBase:    2,
		},
		
		// Client name
		ClientName: "my-app-cluster-client",
		
		// High-throughput configuration
		InflightRequestsLimit: 2000,
	}
}

// Clients holds both standalone and cluster clients
type Clients struct {
	Standalone api.GlideClient
	Cluster    api.GlideClusterClient
}

// CreateClients creates both clients (do this once at application startup)
func CreateClients(ctx context.Context) (*Clients, error) {
	standalone, err := api.NewGlideClient(ctx, standaloneConfig())
	if err != nil {
		return nil, err
	}
	
	cluster, err := api.NewGlideClusterClient(ctx, clusterConfig())
	if err != nil {
		standalone.Close()
		return nil, err
	}
	
	return &Clients{
		Standalone: standalone,
		Cluster:    cluster,
	}, nil
}

// Close gracefully shuts down both clients
func (c *Clients) Close() {
	if c.Standalone != nil {
		c.Standalone.Close()
	}
	if c.Cluster != nil {
		c.Cluster.Close()
	}
}
