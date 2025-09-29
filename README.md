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
│   │   └── redis/           # Redis caching
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
docker-compose up -d

# Start frontend development server
cd ../frontend
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

3. Run development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:5173`

## 🔧 API Documentation

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `GET /api/auth/me` - Get current user

### Premises & Cameras
- `GET /api/premises` - List all premises
- `GET /api/premises/:id/cameras` - Get cameras for premise
- `GET /api/cameras/:id/stream` - Get camera stream URL

### Alerts
- `GET /api/alerts` - List alerts
- `POST /api/alerts/:id/acknowledge` - Acknowledge alert
- `POST /api/alerts/:id/assign` - Assign alert to security guard

### Security Guards
- `GET /api/guards` - List security guards
- `POST /api/guards/:id/dispatch` - Dispatch guard to incident
- `GET /api/guards/me/incidents` - Get assigned incidents

### WebSocket Events
- `alert.created` - New alert notification
- `guard.dispatched` - Guard dispatch notification
- `incident.updated` - Incident status update

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./... -v
go test ./... -cover
```

### Frontend Tests
```bash
cd frontend
npm run check
npm run lint
```

## 🐳 Deployment

### Backend Docker Deployment
```bash
cd backend
docker-compose up -d --build
```

### Frontend Docker Deployment
```bash
cd frontend
docker build -t smart-city-frontend .
docker run -p 3000:3000 smart-city-frontend
```

### Full Stack Deployment
```bash
# From project root
docker-compose -f backend/docker-compose.yml up -d
cd frontend && docker build -t smart-city-frontend . && docker run -p 3000:3000 smart-city-frontend
```

### Kubernetes Deployment
```bash
kubectl apply -f k8s/
```

## 📊 Performance Optimizations

- **Caching**: Redis for session and data caching
- **Database**: Connection pooling and indexing
- **Concurrency**: Goroutines for async operations
- **Load Balancing**: Multiple backend instances
- **CDN**: Static asset delivery

## 🔒 Security Features

- JWT-based authentication
- Role-based access control
- Input validation and sanitization
- CORS configuration
- Rate limiting
- HTTPS enforcement

## 📱 Mobile Responsiveness

The frontend is fully responsive and optimized for:
- Desktop (1920x1080+)
- Tablet (768x1024)
- Mobile (375x667)

## 🎨 UI Components

The frontend uses shadcn-svelte components including:
- Buttons, Cards, Inputs, Labels
- Tables, Badges, Alerts
- Dialogs, Sheets, Tabs
- Avatars, Dropdown Menus
- And more...

## 🚨 Alert System

The system supports various alert types:
- Unauthorized access attempts
- Suspicious activities
- Equipment damage
- System failures

## 📞 Real-time Communication

- WebSocket for instant notifications
- Kafka for reliable message delivery
- Push notifications for mobile app
- Email notifications for critical alerts

## 🔄 Development Workflow

1. Feature branches from `main`
2. Code review required
3. Automated testing on PR
4. Staging deployment
5. Production deployment

## 📈 Monitoring & Logging

- Application metrics with Prometheus
- Log aggregation with ELK stack
- Health check endpoints
- Performance monitoring

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