# 🚀 TitanSQL

**TitanSQL** is a lightweight, distributed relational database designed for **performance**, **scalability**, and **learning**. It’s built on a **disk-based B+ Tree**, uses **Raft consensus** for reliability, and supports **standard SQL** operations. Perfect for experimenting, building small-scale distributed systems, or understanding how modern databases work.

---

## ✨ Features

- 🔍 **Efficient Storage**  
  Uses a **B+ Tree** for fast lookups, inserts, and range queries.

- 🧠 **SQL Engine**  
  Parses and executes standard SQL: `SELECT`, `INSERT`, `UPDATE`, `DELETE`  
  *(Powered by ANTLR-based SQL parser)*

- 🔄 **Concurrency Control**  
  Built-in **MVCC** (Multi-Version Concurrency Control) for safe, concurrent transactions.

- ⚙️ **Distributed Architecture**  
  **Raft consensus** ensures strong consistency and high availability across nodes.

- 🚀 **Optimized I/O**  
  Supports **memory-mapped files**, **page cache**, and **disk persistence**.

- 🛣️ **On the Roadmap**
  - ✅ Data sharding for horizontal scalability  
  - ✅ Query routing across nodes  
  - ⏳ Distributed query planning and optimization  
  - ⏳ Indexes and query cache

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
SELECT name, age FROM users WHERE height > 160 AND age < 30;
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

📚 License
TitanSQL is under the Apache 2.0 license. See the LICENSE file for details.

📫 Contact
Have questions? Reach out via GitHub Issues.

