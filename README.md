# Proof of Reserves Verifier

## 背景 (Background)

本项目旨在提供一个开源的默克尔树验证工具，用于验证用户的资产是否被包含在交易所公布的储备金默克尔树中。通过该工具，用户可以独立审计储备金证明，确保资产的透明度和安全性。

## 简介 (Introduction)

### 源码构建 (Build from source)

请下载适用于您操作系统的最新版本。您也可以自行构建源代码。

**前置要求**：

- Go version >= 1.25.4

### 编译源代码 (Package and compile source code)

#### 1. 进入项目目录

```bash
cd ~/Projects/proof-of-reserves
```

#### 2. 下载依赖

```bash
go mod download
```

#### 3. 构建二进制文件

使用提供的构建脚本生成各平台的二进制文件：

```bash
sh scripts/build_release.sh
```

构建完成后，可执行文件将位于 `dist/` 目录下，例如 `dist/verifier-darwin-arm64`（macOS Apple Silicon）。

#### 4. 启动 (Start up)

```bash
./dist/verifier-darwin-arm64 --file test.json
```

或者直接使用 Go 运行：

```bash
go run ./cmd/verifier --file test.json
```

# 技术说明 (Technical Description)

## 什么是默克尔树？ (What is the Merkle Tree?)

默克尔树（Merkle Tree）是一种数据结构，也称为哈希树。默克尔树将数据存储在树结构的叶子节点中，通过一步步向上哈希数据直到顶部的根节点，叶子节点数据的任何变化都会传递到更高层级的节点，并最终显示为树根的变化。

### 1. 默克尔树的作用

- 零知识证明
- 确保数据不可篡改
- 确保数据隐私

### 2. 默克尔树定义 (Merkle Tree Definition)

本工具实现的默克尔树验证逻辑如下：

#### 2.1 节点信息

树中的每个节点包含：

1. 哈希值 (`merkelLeaf`)
2. 审计批次 ID (`auditId`)
3. 节点层级 (`level`)
4. 角色 (`role`)：1=左节点，2=右节点，3=根节点
5. 余额信息 (`balances`)：虽然包含在节点数据中，但在本工具的当前哈希计算实现中不参与父节点哈希的计算。

#### 2.2 哈希规则 (Hash Rules)

##### 父节点哈希计算

本工具采用标准的 SHA-256 哈希算法。

```
Parent node's hash = SHA256(LeftChildHash + RightChildHash)
```

- `LeftChildHash`: 左子节点的哈希值（字节形式）
- `RightChildHash`: 右子节点的哈希值（字节形式）
- 结果为 32 字节的哈希值，通常以 64 位十六进制字符串表示。

##### 填充节点规则

如果某一层节点数为奇数，且该节点没有兄弟节点，通常会生成一个填充节点（具体规则取决于生成树的逻辑，本验证工具主要负责验证给定的路径是否有效）。

### 验证原理 (Verification Principle)

#### 1. 验证逻辑

根据默克尔树路径数据，从用户自身的叶子节点开始，利用兄弟节点的哈希值，逐层向上计算父节点的哈希值，直到计算出根节点的哈希值。最后，将计算得到的根哈希与路径中提供的根节点哈希进行比对。如果两者一致，则验证通过；否则，验证失败。

#### 2. 示例

假设用户持有叶子节点 `h3`，路径中提供了兄弟节点 `h4`。

1. 计算父节点 `h6 = SHA256(h3 + h4)`
2. 利用 `h6` 和下一层的兄弟节点 `h5`，计算更高层父节点 `h7 = SHA256(h6 + h5)` (假设 h6 是左节点)
3. 重复此过程直到根节点。

#### 验证步骤 (Verification Steps)

1. **准备环境**：确保已安装 Go 环境或下载了预编译的二进制文件。
2. **准备数据**：获取包含默克尔树路径的 JSON 文件（例如 `test.json`）。该文件应包含 `path`（路径节点列表）和 `self`（用户自身节点信息）。
3. **运行验证**：
   ```bash
   # 使用源码运行
   go run ./cmd/verifier --file test.json
   ```
4. **查看结果**：

   - **验证成功**：控制台将输出 "Verification successful!" 以及根哈希信息。
   - **验证失败**：程序将报错并指出计算出的哈希与提供的哈希不匹配的层级。

   **成功示例**：

   ```text
   Verification successful!
     Audit ID : Au2024111208
     Leaf Hash: 30c92d6abf01e27426c41812de6fd708562a10d478e3355fb9a7a14ed8ef08d5
     Root Hash: e5f2729cdbc8c1e2989a4dfcce63ee14ef6c8891348cb14f218ae2659432c0ad
     Levels   : 4
   ```
