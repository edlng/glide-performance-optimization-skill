---
name: GLIDE Performance Optimization
description: Expert guidance for optimizing Valkey GLIDE clients across Node.js, Python, Java, and Go with progressive disclosure
version: 1.0.0
author: Valkey Maintainers
tags:
  - performance
  - optimization
  - valkey
  - glide
  - redis
  - caching
languages:
  - javascript
  - typescript
  - python
  - java
  - go
---

# GLIDE Performance Optimization Skill

Expert guidance for optimizing Valkey GLIDE clients across Node.js, Python, Java, and Go. This skill uses progressive disclosure—loading only the content relevant to your current task to minimize context usage.

## How This Skill Works

This skill automatically detects the programming language you're working with and loads language-specific optimization patterns on-demand:

- **Core Principles** (this file): Always loaded, provides universal anti-patterns and optimization strategies
- **Language-Specific Patterns**: Loaded only when reviewing code in that language
  - `reference/nodejs-patterns.md` - Node.js/TypeScript specific guidance
  - `reference/python-patterns.md` - Python async/sync specific guidance  
  - `reference/java-patterns.md` - Java specific guidance
  - `reference/go-patterns.md` - Go specific guidance

**Context Efficiency**: Reviewing Node.js code loads only Node.js patterns. Python/Java/Go patterns remain unloaded, reducing context usage.

## When to Use This Skill

- P99 latencies >100ms or timeout errors
- Migrating from legacy clients (node-redis, redis-py, Jedis, go-redis)
- High cross-AZ data transfer costs in AWS
- High-throughput applications (>100K ops/sec)
- Production performance incidents
- Real-time features requiring sub-10ms latency
- ElastiCache/Valkey cost optimization
- Serverless/Lambda deployments

## Language Detection & Loading

When you review code, this skill detects the language from:

1. **File extensions**: `.js`, `.ts`, `.py`, `.java`, `.go`
2. **Import statements**: 
   - Node.js: `import { GlideClient }`, `require('@valkey/valkey-glide')`
   - Python: `from glide import`, `import glide`
   - Java: `import glide.api.*`
   - Go: `import "github.com/valkey-io/valkey-glide/go"`
3. **Syntax patterns**: Language-specific keywords and structures

**Action Required**: When language is detected, load the corresponding reference file:
- Node.js/TypeScript → Load `.kiro/skills/glide-performance-optimization/reference/nodejs-patterns.md`
- Python → Load `.kiro/skills/glide-performance-optimization/reference/python-patterns.md`
- Java → Load `.kiro/skills/glide-performance-optimization/reference/java-patterns.md`
- Go → Load `.kiro/skills/glide-performance-optimization/reference/go-patterns.md`

## Universal Anti-Patterns (All Languages)

These anti-patterns cause performance issues regardless of language:

### 1. Per-Request Client Creation (CRITICAL)

**Problem**: Creating a new client for each request adds massive connection overhead.

**Impact**: 
- Latency increases
- Connection pool exhaustion
- Memory leaks from unclosed connections

**Detection**: Look for client creation inside request handlers, loops, or frequently-called functions.

**Fix**: Create client once at application startup, reuse everywhere.

### 2. Missing Request Timeouts

**Problem**: Operations can hang indefinitely without timeouts, causing cascading failures.

**Impact**:
- Requests hang forever on network issues
- Thread/connection pool exhaustion
- Cascading failures across services

**Detection**: Client configuration missing `requestTimeout` or `timeout` parameter.

**Fix**: Always configure timeouts based on use case:
- Sub-10ms apps: 20-50ms
- Web apps: 200-500ms
- Batch processing: 1000-5000ms

### 3. Sequential Operations (High Latency)

**Problem**: Executing commands one-by-one multiplies network latency (N commands = N × roundtrip time).

**Impact**:
- Latency scales linearly with operation count
- Poor throughput
- Wasted network capacity

**Detection**: Multiple await/get calls in sequence without batching or concurrency.

**Fix**: Use batching (pipeline/transaction) or concurrent execution patterns.

### 4. Blocking Commands on Shared Client

**Problem**: BLPOP, BRPOP, BLMOVE, BZPOPMIN, BZPOPMAX block the connection, preventing other operations.

**Impact**:
- All operations on that client are blocked
- Timeouts on concurrent requests
- Poor resource utilization

**Detection**: Blocking commands (BLPOP, BRPOP, etc.) used on same client as regular commands.

**Fix**: Use dedicated client instance for blocking operations with longer timeout.

### 5. Large Batch Sizes (>1000 operations)

**Problem**: Batches with >1000 operations or >10MB payload cause memory issues and timeouts.

**Impact**:
- Memory pressure
- Request timeouts
- Poor error recovery

**Detection**: Batch/pipeline with >1000 commands or very large payloads.

**Fix**: Keep batches between 10-100 commands for optimal balance.

### 6. Missing Error Handling & Retries

**Problem**: Network failures without retry logic cause immediate failures.

**Impact**:
- Poor reliability
- Unnecessary error propagation
- User-facing failures for transient issues

**Detection**: No try-catch blocks or retry configuration around client operations.

**Fix**: Configure connection backoff and implement retry strategies for transient failures.

## Core Optimization Strategies

### 1. Batching (Pipeline & Transactions)

**Concept**: Execute multiple commands in a single network roundtrip.

**When to Use**:
- Multiple independent commands → Pipeline (non-atomic)
- Atomic operations → Transaction (atomic)
- Bulk data operations (MGET, MSET)

**Impact**: Reduces latency from N × roundtrip to 1 × roundtrip.

**Batch Size Guidelines**:
- Optimal: 10-100 commands
- Avoid: >1000 commands or >10MB payload

### 2. Cluster-Aware Operations

**Hash Tags for Co-location**: Use `{tag}` syntax to ensure related keys map to same slot.

**Example**: `{user:123}:name`, `{user:123}:email`, `{user:123}:age` all map to same slot.

**Benefit**: Single roundtrip for multi-key operations on related data.

### 3. AZ Affinity (Cost Optimization)

**Concept**: Route read operations to replicas in same availability zone.

**Requirements**:
- Valkey 8.0+ or AWS ElastiCache for Valkey 7.2+
- Cluster mode with replicas
- Read-heavy workload (>80% reads)

**Impact**:
- Lower read latency
- Reduced cross-AZ data transfer costs
- Improved read throughput

**When NOT to use**:
- Write-heavy workloads
- Strong read consistency required
- Single-AZ deployments

### 4. Async/Concurrent Patterns

**Concept**: Execute independent operations concurrently instead of sequentially.

**Benefit**: Reduces wall-clock time from sum(latencies) to max(latencies).

**Note**: Batching is usually more efficient than concurrent individual operations.

### 5. Data Size Optimization

**Guidelines**:
- Keep values <100KB
- Use compression for values >10KB
- Split large objects across keys
- Use Hash data structures for structured data instead of JSON strings

## Configuration Recommendations

### Request Timeouts

Configure based on your use case:
- Sub-10ms real-time apps: 20-50ms
- Web applications: 200-500ms
- Batch processing: 1000-5000ms

### Connection Retry Strategy

Configure exponential backoff for resilience:
- Number of retries: 5-10
- Base delay: 500ms
- Exponential base: 2

### Connection Pool Sizing

For high-throughput scenarios (>100K ops/sec):
- Increase `inflightRequestsLimit` from default 1000 to 2000+

### Serverless/Lambda

Use `lazyConnect` to defer connection until first command, reducing cold start time.

## Performance Checklist

- [ ] Client reuse (not per-request creation)
- [ ] Request timeout configured (500ms recommended)
- [ ] Batching for bulk operations (10-100 commands)
- [ ] AZ Affinity for read-heavy workloads (>80% reads)
- [ ] Async/concurrent patterns where appropriate
- [ ] Error handling with retry strategy
- [ ] Hash data structures for structured data
- [ ] Values <100KB
- [ ] Hash tags for related keys in cluster
- [ ] Dedicated client for blocking commands
- [ ] lazyConnect for serverless/Lambda

## Troubleshooting Guide

| Symptom | Likely Cause | Solution |
|---------|-------------|----------|
| P99 >100ms | Per-request client creation | Reuse client instances |
| Timeout errors | No timeout configured | Set appropriate timeout |
| Low throughput | Sequential operations | Use batching or concurrency |
| High data transfer costs | Cross-AZ traffic | Enable AZ affinity |
| Connection errors | Network issues | Configure connectionBackoff |
| Blocked operations | Blocking commands on shared client | Use dedicated client |
| Memory issues | Large batches (>1000 ops) | Reduce batch size to 10-100 |

## Additional Resources

- [GLIDE Wiki - Batching](https://github.com/valkey-io/valkey-glide/wiki/General-Concepts#batching-pipeline-and-transaction)
- [GLIDE Wiki - Connection Management](https://github.com/valkey-io/valkey-glide/wiki/General-Concepts#connection-management)
- [AZ Affinity Blog](https://valkey.io/blog/az-affinity-strategy/)
- [Benchmarks](https://github.com/valkey-io/valkey-glide/tree/main/benchmarks)
- [Examples](https://github.com/valkey-io/valkey-glide/tree/main/examples)
