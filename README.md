# TinySQL

**TinySQL** is a distributed relational database built on a disk-based B+ tree. It aims to provide high-performance, scalable, and reliable data storage and query capabilities. At its core lies a custom-designed key-value storage engine that powers the entire database. TinySQL combines modern database techniques—including MVCC, copy-on-write, and memory-mapped files—with an ANTLR-based SQL parser and Raft consensus algorithm to support fault tolerance and horizontal scalability in distributed environments.

Whether you're interested in learning how databases work under the hood or you're building a lightweight distributed system, TinySQL is an ideal choice for research, experimentation, or educational use.

---

## 🚀 Features

### 🧱 Key-Value Storage Engine

- **Disk-Based B+ Tree**  
  Supports efficient point lookups and range queries with predictable performance.

- **Memory-Mapped I/O (mmap)**  
  Leverages OS-level memory mapping to optimize disk I/O performance.

- **Multi-Version Concurrency Control (MVCC)**  
  Enables snapshot isolation and concurrent read/write transactions.

- **Copy-on-Write (CoW)**  
  Ensures data consistency during updates while reducing write amplification.

- **Free List Management**  
  Handles disk page allocation and recycling, minimizing fragmentation.

---

### 🧾 SQL Engine

- **ANTLR-Based SQL Parser**  
  Parses standard SQL syntax, including `SELECT`, `INSERT`, `UPDATE`, and `DELETE`.

- **Schema Definition**  
  Supports table creation with primary keys, indexes, and multi-column definitions.

---

### 🌐 Distributed Capabilities

- **Raft Consensus Protocol**  
  Provides leader election, log replication, and strong consistency guarantees.

- **Dynamic Cluster Management**  
  Supports node membership changes and automatic recovery in case of failures.

---

### ⚙️ Performance Optimizations

- **B+ Tree Iterators**  
  Enables efficient sequential scans and ordered range queries.

- **Fine-Grained Transaction Control**  
  Reduces lock contention and improves concurrent throughput.

- **Page Caching & Prefetching**  
  Boosts I/O efficiency by intelligently caching and preloading disk pages.

---

## 📌 Roadmap

Coming soon:

- 🔄 **Data Sharding**  
  Horizontal partitioning of tables across nodes for better scalability.

- 🧭 **Query Routing**  
  Intelligent routing of SQL queries to the appropriate data partitions.

- 📊 **Distributed Query Planner & Optimizer**  
  Execution planning across multiple nodes.

---

## 📚 License

This project is licensed under the MIT License.

---

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or pull requests to improve TinySQL.

---

## 📫 Contact

For questions or discussions, feel free to reach out via [issues](https://github.com/your-username/tinysql/issues) or submit a PR.

---

> TinySQL — A lightweight yet powerful engine for modern distributed SQL experimentation.

### The descriptions above were all written by AI, and I’m not sure what to write myself.