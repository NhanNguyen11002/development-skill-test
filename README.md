# Smart City Surveillance System

A comprehensive surveillance system for ST Engineering Smart City, providing real-time monitoring, alert management, and security guard dispatch capabilities.

## 🏗️ Architecture

### Backend (Go)
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Message Queue**: Apache Kafka
- **Real-time**: WebSocket
- **Authentication**: JWT
- **Containerization**: Docker & Docker Compose

### Frontend (SvelteKit 5)
- **Framework**: SvelteKit 5 with TypeScript
- **UI**: shadcn-svelte components with Tailwind CSS
- **Real-time**: WebSocket client
- **State Management**: Svelte stores
- **Mobile Responsive**: Yes
- **Authentication**: JWT-based

## 🚀 Features

### SCS Operator Dashboard
- 2x2 grid camera view with live video feeds
- Real-time alert notifications
- Security guard dispatch system
- Role-based access control (view all cameras)

### Security Guard Mobile App
- Real-time incident notifications
- Camera access (assigned cameras only)
- Incident reporting with text/video updates
- Location tracking

### System Features
- Real-time video streaming (simulated)
- Alert management system
- Multi-threading for high performance
- Caching for improved performance
- Containerized deployment

## 📁 Project Structure

```
smart-city-surveillance/
├── backend/                 # Go backend application
│   ├── cmd/server/          # Main server entry point
│   ├── internal/            # Internal application code
│   │   ├── config/          # Configuration management
│   │   ├── database/        # Database connection and migrations
│   │   ├── handlers/        # HTTP request handlers
│   │   ├── middleware/      # Authentication and authorization
│   │   └── models/          # Data models and types
│   ├── pkg/                 # Reusable packages
│   │   ├── websocket/       # WebSocket implementation
│   │   ├── kafka/           # Kafka integration
│   │   ├── response/        # Response share lib
│   │   └── redis/           # Redis caching
│   ├── uml/                 # Diagram folder
│   ├── Dockerfile           # Backend containerization
│   └── docker-compose.yml   # Backend services orchestration
├── frontend/                # SvelteKit 5 frontend application
│   ├── src/
│   │   ├── lib/             # Shared libraries
│   │   │   ├── components/  # shadcn-svelte UI components
│   │   │   ├── stores/      # Svelte stores for state management
│   │   │   ├── api.ts       # API client
│   │   │   └── types.ts     # TypeScript type definitions
│   │   ├── routes/          # SvelteKit routes
│   │   │   ├── login/       # Authentication pages
│   │   │   └── dashboard/   # Main dashboard
│   │   └── app.css          # Global styles with Tailwind
│   ├── static/              # Static assets
│   ├── Dockerfile           # Frontend containerization
│   └── package.json         # Frontend dependencies
└── README.md               # Project documentation
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Node.js 18+

### Quick Setup
```bash
# Clone the repository
git clone <repository-url>
cd smart-city-surveillance

# Start backend services
cd backend
make up 
OR 
docker-compose up -d

# Start frontend development server
cd frontend
npm install
npm run dev
```

### Demo Credentials
- **SCS Operator**: `operator1` / `password`
- **Security Guard**: `guard1` / `password`

## 🛠️ Detailed Setup Instructions

### Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
```

4. Run with Docker Compose:
```bash
make up
OR
docker-compose up -d
```

### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Set up environment variables:
```bash
cp .env.example .env
```

4. Run development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:5173`

## 🔧 API Documentation

- Go to http://localhost:8081/swagger/index.html after start backend container for more info


## 🎨 UI Components

The frontend uses shadcn-svelte components including:
- Buttons, Cards, Inputs, Labels
- Tables, Badges, Alerts
- Dialogs, Sheets, Tabs
- Avatars, Dropdown Menus
- And more...

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## 📄 License

This project is proprietary to ST Engineering Smart City.

## 📞 Support

For technical support, contact the Group IT Team. 