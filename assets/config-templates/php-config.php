<?php
// PHP GLIDE Configuration Template
// Optimized for production web applications

use ValkeyGlide\ValkeyGlide;
use ValkeyGlide\ValkeyGlideCluster;

// Standalone Client Configuration
function createStandaloneClient(): ValkeyGlide {
    return new ValkeyGlide([
        'addresses' => [
            ['host' => 'localhost', 'port' => 6379]
        ],
        
        // Request timeout (500ms recommended for web apps)
        'request_timeout' => 500,
        
        // Connection retry strategy
        'reconnection_strategy' => [
            'number_of_retries' => 10,
            'factor' => 500,        // Base delay in ms
            'exponent_base' => 2,   // Exponential backoff
        ],
        
        // Client name for debugging
        'client_name' => 'my-app-client',
        
        // Lazy connect for serverless/PHP-FPM
        'lazy_connect' => true,  // Recommended for PHP-FPM
    ]);
}

// Cluster Client Configuration
function createClusterClient(): ValkeyGlideCluster {
    return new ValkeyGlideCluster([
        'addresses' => [
            ['host' => 'cluster.endpoint.cache.amazonaws.com', 'port' => 6379]
        ],
        
        // Request timeout
        'request_timeout' => 500,
        
        // AZ Affinity for cost optimization (read-heavy workloads)
        'read_strategy' => ValkeyGlideCluster::READ_FROM_AZ_AFFINITY,
        'client_az' => 'us-east-1a',  // Your application's AZ
        
        // Connection retry strategy
        'reconnection_strategy' => [
            'number_of_retries' => 10,
            'factor' => 500,
            'exponent_base' => 2,
        ],
        
        // Client name
        'client_name' => 'my-app-cluster-client',
        
        // Lazy connect
        'lazy_connect' => true,
    ]);
}

// Global client instances (reuse across requests in PHP-FPM)
global $glide_standalone, $glide_cluster;

if (!isset($glide_standalone)) {
    $glide_standalone = createStandaloneClient();
}

if (!isset($glide_cluster)) {
    $glide_cluster = createClusterClient();
}

// Usage in request handlers
function handleRequest() {
    global $glide_standalone;
    
    try {
        $value = $glide_standalone->get('key');
        return $value;
    } catch (ValkeyGlideException $e) {
        error_log("Valkey error: " . $e->getMessage());
        throw $e;
    }
}
