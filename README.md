# tinySvc 

A lightweight, self-hosted utility service collection providing IP detection, pastebin, URL shortener, and JSON storage. Built with Go and Clean Architecture principles, designed to run efficiently on minimal hardware.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Architecture](https://img.shields.io/badge/architecture-Clean-green.svg)](docs/ARCHITECTURE.md)

## Features

### **What's My IP**
- Instant public IP detection
- Support for Cloudflare headers (CF-Connecting-IP)
- X-Forwarded-For and X-Real-IP support
- IPv4 and IPv6 detection
- Clean, responsive interface

### **Pastebin**
- Text and code snippet sharing
- **Markdown rendering** with syntax highlighting
- Configurable expiration (1 day to never)
- **10MB size limit** per paste
- Copy, download, and delete functionality
- Raw and rendered view toggle
- Direct paste URLs
- Automatic cleanup of expired pastes

### **URL Shortener**
- Create short, memorable links
- **Custom short codes** (3-20 characters)
- **View tracking** (click statistics)
- QR code generation for sharing
- Flexible expiration options
- Same expiration system as pastebin
- Direct redirect endpoint (`/s/{code}`)

### **JSON Bin**
- Store and share JSON data
- **Real-time JSON validation**
- Syntax highlighting
- Format/Minify tools
- **Direct API access** (`/api/v1/json/{id}/raw`)
- **Update support** (PUT requests)
- Code examples (cURL, Fetch, Axios)
- Perfect for API testing and webhooks

## Frontend

- **Frameworkless** vanilla JavaScript (fast & lightweight)
- **Responsive design** (mobile, tablet, desktop)
- **Dark theme** with smooth animations
- **Dropdown navigation** with mobile hamburger menu
- **Modal dialogs** for QR codes and updates
- **Toast notifications** for user feedback
- **Syntax highlighting** (highlight.js)
- **Markdown rendering** (marked.js)
- **QR code generation** (qrcode.js)

## Architecture

TinySvc follows **Clean Architecture** principles with **SOLID** design:
