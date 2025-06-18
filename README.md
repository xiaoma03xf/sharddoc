# 🚀 TitanSQL

**TitanSQL** is a lightweight, distributed relational database designed for **performance**, **scalability**, and **learning**. Built with modern database principles, it combines a **disk-based B+ Tree** storage engine, **MVCC** transaction management, and **Raft consensus** for distributed reliability.

TitanSQL aims to be both a practical database for small to medium workloads and an educational platform for understanding database internals. Its modular architecture makes it ideal for studying how modern databases implement core features like indexing, transaction isolation, and distributed consensus.

---

## ✨ Features

- 🔍 **Efficient Storage**  
  Uses a **B+ Tree** for fast lookups, inserts, and range queries with optimized page management.
  
- 📊 **Advanced Indexing**  
  Supports **primary keys** and **secondary indexes** for fast data retrieval with composite keys.

- 🧠 **SQL Engine**  
  Parses and executes standard SQL: `SELECT`, `INSERT`, `UPDATE`, `DELETE`  
  *(Powered by ANTLR4-based SQL parser with custom AST visitor)*

- 💼 **ACID Transactions**  
  Full support for **atomic**, **consistent**, **isolated**, and **durable** transactions with proper rollback.

- 🔄 **Concurrency Control**  
  Built-in **MVCC** (Multi-Version Concurrency Control) for safe, concurrent transactions with conflict detection.

- 🗄️ **Memory Management**  
  Efficient **free list** implementation for page reuse and memory optimization.

- 📸 **Snapshot & Recovery**  
  Point-in-time snapshots for backup and fast recovery with WAL (Write-Ahead Logging).

- ⚙️ **Distributed Architecture**  
  **Raft consensus** ensures strong consistency and high availability across nodes with automatic leader election.

- 🔌 **Client-Server Protocol**  
  Custom TCP-based protocol for efficient client-server communication with binary encoding.

- 🚀 **Optimized I/O**  
  Supports **memory-mapped files**, **page cache**, and **disk persistence** for high-performance operations.

- 🛣️ **On the Roadmap**
  - ✅ Data sharding for horizontal scalability  
  - ✅ Query routing across nodes  
  - ⏳ Distributed query planning and optimization  
  - ⏳ Advanced query cache mechanisms

---

## 🔧 Getting Started

### 📥 Clone & Build
```bash
git clone https://github.com/xiaoma03xf/titansql.git
cd titansql
go build -o titansql ./cmd


# Start three nodes with different configs
./titansql -config node1.yaml
./titansql -config node2.yaml
./titansql -config node3.yaml
```

### 🧪SQL Example
```SQL 
CREATE TABLE users (
    id INT64,
    name BYTES,
    age INT64,
    height INT64,
    PRIMARY KEY (id),
    INDEX (age, height)
);

INSERT INTO users (id, name, age, height) VALUES (1, 'Alice', 25, 170);
SELECT name, age FROM users WHERE age = 30 AND height > 175;
UPDATE users SET age = 30 WHERE name = 'Alice';
DELETE FROM users WHERE id = 1;
```

### 🧭 Architecture Overview
```
    +-------------+       +-------------------+
    |   Client    | <---> | Query Router Node |
    +-------------+       +-------------------+
                                    |
        +------------+------------+------------+
        |            |                         |
    +---------+  +---------+               +---------+
    | Shard 1 |  | Shard 2 |   ...         | Shard N |
    +---------+  +---------+               +---------+
        | Raft |     | Raft |                  | Raft |
```

## 🔬 Technical Details

### 💾 Storage Engine

TitanSQL uses a custom B+ Tree implementation for efficient data storage and retrieval:

- **Page-Based Structure**: Fixed-size pages (4KB) optimize disk I/O
- **Key-Value Storage**: Supports variable-length keys and values
- **Efficient Operations**: O(log n) lookups, inserts, and deletes
- **Range Queries**: Fast sequential scans through leaf node links
- **Node Splitting**: Automatic node splitting and merging to maintain balance
- **Memory Management**: Efficient page allocation with free list tracking

### 🔄 Transaction Management

Transactions in TitanSQL use Multi-Version Concurrency Control (MVCC):

- **Snapshot Isolation**: Each transaction works with a consistent snapshot
- **Conflict Detection**: Prevents write-after-read conflicts with range tracking
- **No Blocking**: Readers never block writers and vice versa
- **Abort Handling**: Clean rollback of failed transactions
- **Versioning**: Maintains multiple versions of data for concurrent access
- **Key Range Tracking**: Records read ranges for precise conflict detection

### 🧠 SQL Processing

SQL statements are processed through several stages:

1. **Parsing**: ANTLR4-based parser converts SQL text to AST
2. **Planning**: Simple query planning for optimal execution
3. **Execution**: Executes the query against the storage engine
4. **Result Processing**: Formats and returns query results

### 🌐 Distributed Consensus

TitanSQL uses Hashicorp's Raft implementation for distributed consensus:

- **Leader Election**: Automatic leader election for write operations
- **Log Replication**: Ensures all nodes have the same data
- **Membership Changes**: Dynamic cluster membership
- **Snapshots**: Periodic state snapshots for efficient recovery
- **Failure Handling**: Automatic recovery from node failures

## 📁 Project Structure

```
titansql/
├── cmd/                  # Command-line interface
├── lib/                  # Utility libraries
│   ├── logger/           # Logging utilities
│   ├── targz/            # Compression utilities
│   └── utils/            # General utilities
├── parser/               # SQL parsing
│   ├── antlr/            # ANTLR grammar files
│   └── ast/              # Abstract Syntax Tree
├── storage/              # Storage engine
│   ├── bplustree.go      # B+ Tree implementation
│   ├── bplustree_iter.go # Iter for range
│   ├── free_list.go      # Control node page cycle
│   ├── tx.go             # Transaction management
│   ├── kv.go             # Key-value store
│   ├── table.go          # Table operations
│   └── parser.go         # SQL parser integration
└── tcp/                  # Network layer
    ├── server.go         # TCP server
    ├── client.go         # TCP client
    ├── raft.go           # Raft consensus
    └── fsm.go            # Finite State Machine
```

## 🛠️ Development Guide

### Prerequisites

- Go 1.24 or higher
- ANTLR4 (for SQL grammar modifications)

### Configuration

TitanSQL uses YAML configuration files for node setup:

```yaml
# Example node configuration (node1.yaml)
node_id: node1
raft_addr: 127.0.0.1:28001
http_addr: 127.0.0.1:29001
raft_dir: ./clusterdb/raft/node1
db_path: ./clusterdb/data/node1.db
bootstrap: true  # First node in cluster

# Common settings
raft_timeout: 10s
snapshot_interval: 30s
snapshot_threshold: 100000
```

### Adding a New Node to Cluster

1. Create a configuration file (e.g., `node4.yaml`)
2. Set `bootstrap: false` and `join_addr` to an existing node's HTTP address
3. Start the new node with `./titansql -config node4.yaml`

## 🧪 Testing

TitanSQL includes comprehensive tests:

- **Unit Tests**: Test individual components
- **Integration Tests**: Test component interactions
- **Distributed Tests**: Test multi-node scenarios

Run tests with:

```bash
# Run all tests
go test ./...

# Run specific tests
go test ./storage/...
go test ./tcp/...
```

## 📊 Performance Characteristics

- **Storage**: Efficient B+ Tree with O(log n) operations
- **Concurrency**: MVCC allows high concurrent throughput
- **Distributed**: Raft consensus with reasonable latency
- **Memory Usage**: Configurable cache sizes for different workloads
- **Disk I/O**: Optimized page management reduces disk operations
- **Query Performance**: Indexes accelerate common query patterns

## 🧩 Design Principles

TitanSQL follows these core design principles:

1. **Simplicity**: Clean, understandable code over complex optimizations
2. **Reliability**: Data integrity and consistency are paramount
3. **Modularity**: Well-defined interfaces between components
4. **Learnability**: Code structure and documentation that facilitates learning
5. **Performance**: Efficient algorithms and data structures for core operations

## 🤝 Contributing

Contributions are welcome! Here's how to contribute:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code follows the project's coding standards and includes appropriate tests.

## 📚 License
TitanSQL is under the Apache 2.0 license. See the LICENSE file for details.

## 📫 Contact
Have questions? Reach out via GitHub Issues.
