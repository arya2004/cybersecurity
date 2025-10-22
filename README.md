# Cybersecurity Lab Codes in Go

Welcome to the Cybersecurity Lab Codes repository! This repository contains a comprehensive collection of Go (Golang) programs designed to explore various aspects of cybersecurity. These programs are intended for educational purposes and can be used as part of cybersecurity labs, exercises, or research.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Lab Modules](#lab-modules)
- [Installation](#installation)
- [Usage](#usage)
- [Security Tools](#security-tools)
- [Educational Value](#educational-value)
- [Contributing](#contributing)
- [License](#license)

## Introduction

This repository aims to provide hands-on experience with cybersecurity concepts using Go. Each program focuses on a specific area of cybersecurity, including cryptography, network security, vulnerability analysis, web application security, and penetration testing. By working through these programs, users can gain a deeper understanding of how cybersecurity mechanisms are implemented and how to defend against various types of attacks.

## Features

- **Comprehensive Coverage**: 10+ cybersecurity domains with practical implementations
- **Educational Focus**: Designed for learning and understanding security concepts
- **Real-World Applications**: Tools that demonstrate actual security vulnerabilities
- **Interactive Interfaces**: User-friendly menu-driven applications
- **Detailed Documentation**: Extensive explanations and usage examples
- **Security Analysis**: Vulnerability detection and remediation guidance
- **Modern Go Implementation**: Clean, efficient, and well-structured code

## Lab Modules

### 🔐 Cryptography Modules
1. **[01_des](./01_des/)** - Data Encryption Standard implementation and analysis
2. **[02_aes](./02_aes/)** - Advanced Encryption Standard with multiple modes
3. **[03_rsa](./03_rsa/)** - RSA public-key cryptography implementation
4. **[05_message-digest-5](./05_message-digest-5/)** - MD5 hash function and vulnerabilities
5. **[10_advanced_crypto](./10_advanced_crypto/)** - Comprehensive cryptography tool with multiple algorithms

### 🔑 Key Exchange & Authentication
6. **[04.1_ecc_diiffie_hellman](./04.1_ecc_diiffie_hellman/)** - Elliptic Curve Diffie-Hellman
7. **[04.2_diffie_hellman](./04.2_diffie_hellman/)** - Classic Diffie-Hellman key exchange
8. **[08_jwt_security](./08_jwt_security/)** - JWT security analysis and testing tool

### 🛡️ Network Security
9. **[07_network_scanner](./07_network_scanner/)** - Advanced network security scanner
10. **[06_password_security](./06_password_security/)** - Password security analysis

### 🌐 Web Application Security
11. **[09_sql_injection_tester](./09_sql_injection_tester/)** - SQL injection vulnerability testing tool

## Installation

To get started with these programs, you need to have Go installed on your machine. Follow the instructions below to set up your environment:

1. **Install Go**:
   - Download and install Go from the official [Go website](https://golang.org/dl/).
   - Follow the installation instructions for your operating system.

2. **Clone the repository**:
   ```sh
   git clone https://github.com/arya2004/cybersecurity.git
   cd cybersecurity
   ```

3. **Verify installation**:
   ```sh
   go version
   ```

## Usage

Each program is located in its own directory with comprehensive documentation. To run any program:

```sh
cd [module_directory]
go run main.go
```

### Quick Start Examples

```sh
# Network Security Scanner
cd 07_network_scanner
go run main.go

# JWT Security Tool
cd 08_jwt_security
go run main.go

# SQL Injection Tester
cd 09_sql_injection_tester
go run main.go

# Advanced Cryptography Tool
cd 10_advanced_crypto
go run main.go
```

## Security Tools

### 🔍 Network Security Scanner
**Location**: `07_network_scanner/`
- TCP port scanning with service detection
- Network host discovery
- OS fingerprinting
- Vulnerability assessment
- Interactive scanning interface

### 🔐 JWT Security Analysis Tool
**Location**: `08_jwt_security/`
- JWT token parsing and validation
- Algorithm confusion testing
- Weak secret detection
- Signature bypass testing
- Comprehensive security reporting

### 💉 SQL Injection Testing Tool
**Location**: `09_sql_injection_tester/`
- Error-based injection testing
- Boolean-based blind injection
- Time-based injection detection
- Multi-database support (MySQL, MSSQL, PostgreSQL)
- 50+ specialized payloads

### 🛡️ Advanced Cryptography Tool
**Location**: `10_advanced_crypto/`
- Multiple encryption algorithms (AES, DES, RSA)
- Hash function analysis
- Key strength assessment
- Timing attack demonstration
- Cryptographic weakness detection

## Educational Value

This repository provides hands-on learning opportunities for:

### Core Security Concepts
- **Symmetric Cryptography**: AES, DES implementations and analysis
- **Asymmetric Cryptography**: RSA, ECC key exchange protocols
- **Hash Functions**: MD5, SHA family with vulnerability analysis
- **Network Security**: Port scanning, service enumeration
- **Web Security**: SQL injection, JWT vulnerabilities

### Practical Skills Development
- **Vulnerability Assessment**: Automated security testing tools
- **Penetration Testing**: Real-world attack simulation
- **Security Analysis**: Weakness identification and remediation
- **Secure Coding**: Defensive programming practices
- **Tool Development**: Building custom security tools

### Professional Applications
- **Security Auditing**: Comprehensive security assessment
- **Incident Response**: Understanding attack vectors
- **Security Architecture**: Implementing proper defenses
- **Compliance**: Meeting security standards and requirements


## Contributing

We welcome contributions from the community! If you have a program or improvement you'd like to share, please follow these steps:

1. **Fork the repository**
2. **Create a new branch** with a descriptive name:
   ```sh
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes** and test thoroughly
4. **Follow coding standards**:
   - Include comprehensive documentation
   - Add educational comments explaining security concepts
   - Implement proper error handling
   - Include usage examples in README files
5. **Submit a pull request** with a detailed description of your changes

### Contribution Guidelines

- **Security Focus**: All contributions should have educational security value
- **Code Quality**: Follow Go best practices and include tests where appropriate
- **Documentation**: Provide clear README files and inline documentation
- **Ethical Use**: Ensure tools are designed for authorized testing only
- **Educational Value**: Include explanations of security concepts and vulnerabilities

### Ideas for Contributions

- Additional cryptographic algorithms
- Web application security testing tools
- Blockchain security analysis
- IoT security assessment tools
- Machine learning security applications
- Cloud security configuration tools

## Security Notice

⚠️ **Important**: These tools are designed for educational purposes and authorized security testing only. Users are responsible for:

- Obtaining proper authorization before testing any systems
- Complying with applicable laws and regulations
- Using tools responsibly and ethically
- Not using tools for malicious purposes

## License

This repository is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

## Acknowledgments

- Go community for excellent cryptographic libraries
- Security researchers for vulnerability disclosure
- Educational institutions promoting cybersecurity learning
- Open source contributors advancing security knowledge

---

**Disclaimer**: This repository is for educational and authorized testing purposes only. The authors are not responsible for any misuse of the tools or information provided.

## License

This repository is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

