# Training Project - Microservices Todo Application

A microservices-based todo application built with Go, featuring authentication, todo management, and a GraphQL BFF layer. This project demonstrates modern microservices architecture patterns including gRPC communication, JWT authentication, GraphQL API gateway, and containerized services.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Services](#services)

## Overview

This project is a training exercise demonstrating a complete microservices architecture with the following capabilities:

- **User Authentication**: JWT-based authentication with access and refresh tokens
- **Todo Management**: CRUD operations for todo items
- **GraphQL API**: Frontend-friendly GraphQL interface via BFF pattern
- **Microservices Communication**: gRPC for inter-service communication
- **Database Migrations**: Automated database schema management
- **Testing**: Unit and integration tests
- **Containerization**: Docker-based development environment

## Architecture

### System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                        │
│                    (Web/Mobile Frontend)                    │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  │ GraphQL (HTTP)
                  ▼
┌────────────────────────────────────────────────────────────────────┐
│                            Todo-BFF Service                        │
│                           (GraphQL Gateway)                        │
│                            Port: 5006                              │
└──────────────┬─────────────────────────────┬───────────────────────┘
               │                             │
               │ gRPC                        │ gRPC
               ▼                             ▼
┌──────────────────────────┐      ┌─────────────────────────────┐
│    Auth Service          │      │     Todo Service            │
│  (Authentication & JWT)  │      │  (Todo Management)          │
│     Port: 5007           │      │     Port: 5005              │
└──────────┬───────────────┘      └───────────┬─────────────────┘
           │                                  │
           │ MySQL                            │ MySQL
           ▼                                  ▼
┌──────────────────────────┐      ┌─────────────────────────────┐
│   Auth Database          │      │    Todo Database            │
│   Port: 33061            │      │    Port: 33062              │
└──────────────────────────┘      └─────────────────────────────┘
```

## Services

### 1. Auth Service
**Location**: `go/services/auth`  
**Port**: 5007  
**Database Port**: 33061  
**Protocol**: gRPC

Authentication and authorization service using JWT tokens.

**Key Features**:
- User registration with email validation
- Login with JWT token generation
- Token verification and refresh
- Bcrypt password hashing

[Auth Service README](go/services/auth/README.md)

### 2. Todo Service
**Location**: `go/services/todo`  
**Port**: 5005  
**Database Port**: 33062  
**Protocol**: gRPC

Todo item management service.

**Key Features**:
- CRUD operations for todos
- User-scoped todo management
- Status tracking (Pending, In Progress, Completed)
- Pagination support

[Todo Service README](go/services/todo/README.md)

### 3. Todo-BFF Service
**Location**: `go/services/todo-bff`  
**Port**: 5006  
**Protocol**: GraphQL over HTTP

Backend-For-Frontend service providing a GraphQL API.

**Key Features**:
- GraphQL API for frontend consumption
- Authentication via JWT
- Real-time schema introspection
- GraphQL Playground for development

[Todo-BFF Service README](go/services/todo-bff/README.md)
