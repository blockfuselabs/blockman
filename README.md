# BlockMan: Blockchain ABI Interaction Toolkit

BlockMan is a developer tool for interacting with Ethereum smart contracts via their Application Binary Interface (ABI). Inspired by Postman, it allows developers to upload ABIs, list contract functions, and call them directly through a simple API.

![BlockMan Logo](https://via.placeholder.com/150?text=BlockMan)

## 🚀 Features

- **Upload ABIs**: Easily upload and parse Ethereum contract ABIs
- **List Contract Functions**: View all functions and their details, including inputs, outputs, and state mutability
- **Call Functions**: Interact with smart contract functions, passing parameters and receiving responses
- **Persistent Storage**: Store your ABIs for future use with automatic cleanup options
- **Secure Access**: API key authentication to protect your contracts
- **Simple Integration**: RESTful API for easy integration with your development workflow

## 📋 Requirements

- Go 1.20 or later
- Access to an Ethereum JSON-RPC node (e.g., Infura, Alchemy, or a self-hosted node)

## 🛠️ Installation Options

### Self-Hosted

1. Clone the repository:
   ```bash
   git clone https://github.com/blockfuselabs/blockman.git
   cd blockman
   ```

2. Install dependencies:
   ```bash
   make install
   ```

3. Create a `.env` file in the project root:
   ```env
   ETH_NODE_URL=<Your Ethereum Node URL>
   PORT=8080  # Optional, defaults to 8080
   CLEANUP_ENABLED=true  # Optional, defaults to true
   CLEANUP_HOURS=24  # Optional, defaults to 24
   ```

4. Build and run the application:
   ```bash
   make build
   ./blockman
   ```

### Docker

1. Pull the Docker image:
   ```bash
   docker pull blockfuselabs/blockman:latest
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 -e ETH_NODE_URL=<Your Ethereum Node URL> blockfuselabs/blockman:latest
   ```

### Cloud Deployment

We offer detailed deployment guides for major cloud providers:
- [AWS Deployment Guide](https://github.com/blockfuselabs/blockman/docs/aws-deployment.md)
- [Google Cloud Deployment Guide](https://github.com/blockfuselabs/blockman/docs/gcp-deployment.md)
- [Azure Deployment Guide](https://github.com/blockfuselabs/blockman/docs/azure-deployment.md)

## 🔒 Security Considerations

When deploying BlockMan for public consumption, consider the following:

1. **API Authentication**: Implement API keys or OAuth2 for secure access
2. **HTTPS**: Always use HTTPS in production environments
3. **Rate Limiting**: Set up rate limiting to prevent abuse
4. **Access Controls**: Restrict access to sensitive operations
5. **Node Security**: Secure your Ethereum node connection with credentials

## 📚 API Documentation

### Upload ABI

Upload and parse an ABI for future interaction.

**Endpoint**: `POST /upload-abi`
**Authentication**: Required

**Request Body**:
```json
{
  "abi": "[{...ABI JSON...}]"
}
```

**Response**: 
```json
{
  "message": "ABI uploaded successfully",
  "abi_id": "unique-abi-id"
}
```

### List Contract Functions

Retrieve all functions from an uploaded ABI.

**Endpoint**: `POST /list-functions`
**Authentication**: Required

**Request Body**: 
```json
{
  "abi_id": "unique-abi-id"
}
```

**Response**: 
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

### Call Function

Call a specific function from the ABI.

**Endpoint**: `POST /call-function`
**Authentication**: Required

**Request Body**: 
```json
{
  "abi_id": "unique-abi-id",
  "contract_address": "0xYourContractAddress",
  "function_name": "functionName",
  "function_input": ["param1", "param2"]
}
```

**Response**: 
```json
{
  "result": "0x..."
}
```

### List Stored ABIs

Get a list of all ABIs in your account.

**Endpoint**: `GET /abis`
**Authentication**: Required

**Response**: 
```json
{
  "abis": [
    {
      "abi_id": "unique-abi-id",
      "created_at": "2023-04-01T12:00:00Z",
      "last_used": "2023-04-02T14:30:00Z"
    }
  ]
}
```

### Remove ABI

Delete a stored ABI.

**Endpoint**: `DELETE /abis/:id`
**Authentication**: Required

**Response**: 
```json
{
  "message": "ABI removed successfully"
}
```

## 🖥️ Client Libraries

For easier integration, we provide client libraries in multiple languages:

- [JavaScript/Node.js](https://github.com/blockfuselabs/blockman-js)
- [Python](https://github.com/blockfuselabs/blockman-python)
- [Go](https://github.com/blockfuselabs/blockman-go)

Example JavaScript usage:
```javascript
const BlockmanClient = require('blockman-client');

const client = new BlockmanClient({
  apiUrl: 'https://your-blockman-instance.com',
  apiKey: 'your-api-key'
});

// Upload an ABI
const abiId = await client.uploadABI(contractABI);

// List functions
const functions = await client.listFunctions(abiId);

// Call a function
const result = await client.callFunction({
  abiId,
  contractAddress: '0xContractAddress',
  functionName: 'balanceOf',
  functionInput: ['0xUserAddress']
});

console.log(result);
```

## 📊 Monitoring & Logging

BlockMan provides built-in monitoring endpoints:

- `GET /health` - Service health check
- `GET /metrics` - Prometheus-compatible metrics
- `GET /status` - Service status and statistics

## 🔄 Production Best Practices

For running BlockMan in production, we recommend:

1. Deploy multiple instances behind a load balancer
2. Set up persistent storage for ABI data (database configuration)
3. Configure monitoring and alerts
4. Implement regular backups
5. Use a CDN for high-traffic deployments

## 🤝 Contributing

We welcome contributions to improve BlockMan! To get started:

1. Fork the repository
2. Create a new branch for your feature/fix
3. Submit a pull request with a clear description of your changes

Please follow our [Code of Conduct](CODE_OF_CONDUCT.md) and review our [Contributing Guidelines](CONTRIBUTING.md).

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋 Support & Community

- [Documentation](https://docs.blockman.io)
- [Discord Community](https://discord.gg/blockman)
- [Issue Tracker](https://github.com/blockfuselabs/blockman/issues)
- [Email Support](mailto:support@blockman.io)