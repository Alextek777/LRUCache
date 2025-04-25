# Thread Safe LRUCache Implementation in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/yourrepo.svg)](https://pkg.go.dev/github.com/yourusername/yourrepo)
[![Tests](https://github.com/yourusername/yourrepo/actions/workflows/tests.yml/badge.svg)](https://github.com/yourusername/yourrepo/actions/workflows/tests.yml)

A thread-safe LRU (Least Recently Used) cache implementation in Go with TTL support.

## Features

- 🚀 Thread-safe operations using sync.RWMutex
- ⏳ Time-based eviction (TTL)
- 📏 Size-based eviction (LRU)
- 🔄 Automatic cleanup of expired items
- 📊 Metrics and statistics (optional)

## Installation


```sh
go get github.com/Alextek777/LRUCache