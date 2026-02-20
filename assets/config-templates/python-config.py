# Python GLIDE Configuration Template
# Optimized for production web applications

from glide import (
    GlideClient,
    GlideClusterClient,
    GlideClientConfiguration,
    GlideClusterClientConfiguration,
    NodeAddress,
    BackoffStrategy,
)

# Standalone Client Configuration (Async)
async_standalone_config = GlideClientConfiguration(
    addresses=[NodeAddress("localhost", 6379)],
    
    # Request timeout (500ms recommended for web apps)
    request_timeout=500,
    
    # Connection retry strategy
    connection_backoff=BackoffStrategy(
        number_of_retries=10,
        factor=500,        # Base delay in ms
        exponent_base=2,   # Exponential backoff
    ),
    
    # Client name for debugging
    client_name="my-app-client",
    
    # Lazy connect for serverless/Lambda
    lazy_connect=False,  # Set to True for Lambda
    
    # High-throughput configuration
    inflight_requests_limit=2000,  # Default: 1000
)

# Cluster Client Configuration (Async)
async_cluster_config = GlideClusterClientConfiguration(
    addresses=[NodeAddress("cluster.endpoint.cache.amazonaws.com", 6379)],
    
    # Request timeout
    request_timeout=500,
    
    # AZ Affinity for cost optimization (read-heavy workloads)
    read_from="AZAffinity",
    client_az="us-east-1a",  # Your application's AZ
    
    # Connection retry strategy
    connection_backoff=BackoffStrategy(
        number_of_retries=10,
        factor=500,
        exponent_base=2,
    ),
    
    # Client name
    client_name="my-app-cluster-client",
    
    # High-throughput configuration
    inflight_requests_limit=2000,
)

# Create clients (do this once at application startup)
async def create_async_clients():
    standalone = await GlideClient.create(async_standalone_config)
    cluster = await GlideClusterClient.create(async_cluster_config)
    return standalone, cluster

# Graceful shutdown
async def close_async_clients(standalone, cluster):
    await standalone.close()
    await cluster.close()

# Sync Client Configuration
from glide import GlideClientSync, GlideClusterClientSync

sync_standalone_config = GlideClientConfiguration(
    addresses=[NodeAddress("localhost", 6379)],
    request_timeout=500,
    connection_backoff=BackoffStrategy(
        number_of_retries=10,
        factor=500,
        exponent_base=2,
    ),
    client_name="my-app-sync-client",
)

# Create sync clients
def create_sync_clients():
    standalone = GlideClientSync.create(sync_standalone_config)
    return standalone

# Graceful shutdown
def close_sync_clients(standalone):
    standalone.close()
