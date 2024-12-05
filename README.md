# BlockMan: Blockchain ABI Interaction Toolkit

Blockman is a developer tool for interacting with Ethereum smart contracts via their Application Binary Interface (ABI). Inspired by Postman, it allows developers to upload ABIs, list contract functions, and call them directly through a simple API interface.

### Features
- Upload ABIs: Easily upload and parse Ethereum contract ABIs.
- List Contract Functions: View all functions and their details, including inputs, outputs, and state mutability.
- Call Functions: Interact with smart contract functions, passing parameters and receiving responses.

### Getting Started

#### Prerequisites
1. Go 1.20 or later
2. Access to an Ethereum JSON-RPC node (e.g., Infura, Alchemy, or a self-hosted node)

#### Installation
1. Clone the repository:
```bash
git clone https://github.com/blockfuselabs/blockman.git
cd blockman
```

2. Clone the repository:
```bash
go mod tidy
```

3. Create a .env file in the project root:
```env
ETH_NODE_URL=<Your Ethereum Node URL>
```

4. Build and run the application:
```bash
go run cmd/web/main.go
```

5. The server will be available at: http://localhost:8080


#### API Endpoints

1. Upload ABI
Upload and parse an ABI for future interaction.

- Endpoint: POST /upload-abi
- Request Body:
```json
{
  "abi": "[{...ABI JSON...}]"
}
```
- Response: 
```json
{
  "message": "ABI uploaded successfully",
  "abi_id": "unique-abi-id"
}
```

2. List Contract Functions
Retrieve all functions from an uploaded ABI.
- Endpoint: POST /list-functions
- Request Body: 
```json
{
  "abi_id": "unique-abi-id"
}
```
- Response: 
```json
{
  "functions": [
    {
      "name": "functionName",
      "inputs": [{"name": "param1", "type": "uint256"}],
      "outputs": [{"name": "", "type": "string"}],
      "constant": true,
      "payable": false,
      "stateful": false
    }
  ],
  "total_functions": 1
}
```

3. Call Function
Call a specific function from the ABI.
- Endpoint: POST /call-function
- Request Body: 
```json
{
  "abi_id": "unique-abi-id",
  "contract_address": "0xYourContractAddress",
  "function_name": "functionName",
  "function_input": ["param1", "param2"]
}
```
- Response: 
```json
{
  "result": "0x..."
}
```

### Contributing
We welcome contributions to improve Blockman! To get started:

1. Fork the repository.
2. Create a new branch for your feature/fix.
3. Submit a pull request with a clear description of your changes.

### Roadmap
1. Database Integration: Persistent ABI storage using SQLite or Redis.
2. Authentication: Secure API endpoints with API keys or OAuth.
3. Web UI: Build an interactive interface for managing and testing ABIs.
4. Enhanced Input Parsing: Support complex input types like structs and arrays.


### License
This project is licensed under the MIT License. 




