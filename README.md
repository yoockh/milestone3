# Donation & Auction Platform

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Echo Framework](https://img.shields.io/badge/Echo-v4-00ADD8?style=flat&logo=go)](https://echo.labstack.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=flat&logo=redis&logoColor=white)](https://redis.io/)
[![Google Cloud](https://img.shields.io/badge/Google_Cloud-4285F4?style=flat&logo=google-cloud&logoColor=white)](https://cloud.google.com/)
[![Midtrans](https://img.shields.io/badge/Midtrans-Payment-FF6B00?style=flat)](https://midtrans.com/)
[![Swagger](https://img.shields.io/badge/Swagger-API_Docs-85EA2D?style=flat&logo=swagger)](https://yourdonaterise-278016640112.asia-southeast2.run.app/swagger/index.html)
[![Coverage](https://img.shields.io/badge/Coverage-71.9%25-green?style=flat&logo=go)](be/coverage.html)

A comprehensive donation and auction management system that transforms donated goods into meaningful impact through transparent auctions and direct donations to institutions in need.

---

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [System Flow](#system-flow)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Running Migrations](#running-migrations)
- [Future Features](#future-features)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

This platform enables efficient management of donated goods through a dual-channel system: high-value items are auctioned to generate funds, while items unsuitable for auction are directly donated to institutions. The system ensures transparency through weekly reporting and automated workflow management.

### Core Workflow

1. **Donation Submission**: Donors submit items with photos and descriptions
2. **Verification**: Verifiers assess item condition and determine auction eligibility
3. **Dual Processing**:
   - **Auction Track**: Eligible items enter bidding sessions
   - **Direct Donation Track**: Non-auction items go directly to institutions
4. **Transaction Management**: Payment processing and distribution tracking
5. **Transparency**: Weekly articles reporting all activities and fund allocation

---

## Key Features

### For User
- Submit donations with detailed descriptions and photos
- Track donation status in real-time
- View complete donation history
- Receive notifications on item outcomes
- Browse active auction items by category
- Place real-time bids on items
- View bid history and current highest bid
- Secure payment processing after winning
- Track auction participation history

### For Admins
- Review pending donation submissions
- Assess physical condition and categorization
- Determine auction eligibility
- Upload verification photos and notes
- Comprehensive dashboard with analytics
- Manage auction sessions and scheduling
- User and institution management
- Approve donation distributions
- Generate transparency reports
- Publish weekly articles

### System Features
- Automated auction winner determination
- Real-time bid tracking with Redis caching
- Secure payment integration via Midtrans
- Email notifications for all key events
- Photo storage via Google Cloud Storage
- Rate limiting and spam prevention

---

## System Flow

![flowchart](assets/flow.jpg)

### 1. Donation Flow
```
Donor Submits → Pending Status → Verifier Assigned → Physical Inspection
    ↓
Verification Decision
    ├─→ Auction Eligible → Auction Processing
    └─→ Direct Donation → Institution Distribution
```

### 2. Auction Flow
```
Verified Item → Admin Creates Session → Scheduled Auction → Bidding Opens
    ↓
Real-time Bidding → Session Ends → Winner Selected → Payment → Delivery
```

### 3. Direct Donation Flow
```
Admin Approves (verified_for_donation) → Auto-create Final Donation Entry
    ↓
User Adds Notes → Final Donation Complete → Transparency Report
```

---

## Tech Stack

### Backend
- **Framework**: Echo (Golang) - High-performance, minimalist web framework
- **Language**: Go 1.21+ - Efficient concurrency for real-time bidding

### Database & Cache
- **Primary Database**: PostgreSQL 14+ - ACID compliance for financial transactions
- **Cache Layer**: Redis 7+ - Real-time bid caching and rate limiting
- **ORM**: GORM or sqlx - Type-safe database operations

### Infrastructure
- **Cloud Platform**: Google Cloud Platform
  - **Compute**: Cloud Run - Serverless, auto-scaling containers
  - **Storage**: Cloud Storage - Secure photo storage
  - **Database**: Cloud SQL - Managed PostgreSQL
- **Containerization**: Docker

### Third-Party Services
- **Payment Gateway**: Midtrans - Indonesian payment processing
- **Email Service**: Mailjet/Resend - Transactional email delivery
- **File Upload**: Google Cloud Storage - Scalable object storage

### DevOps & Tools
- **Migration Tool**: golang-migrate or custom SQL migrations
- **API Documentation**: Swagger/OpenAPI (optional)
- **Monitoring**: Google Cloud Monitoring

---

## Project Structure

```
be/
├── app/
│   └── main.go                          # Application entry point
│
├── config/
│   ├── connectionDb.go                  # PostgreSQL connection setup
│   └── redis.go                         # Redis client configuration
│
├── internal/
│   ├── controller/                      # HTTP request handlers
│   │   ├── admin_controller.go
│   │   ├── article_controller.go
│   │   ├── auction_item_controller.go
│   │   ├── auction_session_controller.go
│   │   ├── bid_controller.go
│   │   ├── donation_controller.go
│   │   ├── final_donation_controller.go
│   │   ├── payment_controller.go
│   │   └── user_controller.go
│   │
│   ├── service/                         # Business logic layer
│   │   ├── admin_service.go
│   │   ├── article_service.go
│   │   ├── auction_item_service.go
│   │   ├── auction_session_service.go
│   │   ├── bid_service.go
│   │   ├── donation_service.go
│   │   ├── final_donation_service.go
│   │   ├── payment_service.go
│   │   ├── user_service.go
│   │   └── errors.go
│   │
│   ├── repository/                      # Data access layer
│   │   ├── admin_repo.go
│   │   ├── ai_repo.go
│   │   ├── article_repo.go
│   │   ├── auction_item_repo.go
│   │   ├── auction_session_repo.go
│   │   ├── auction_session_redis.go
│   │   ├── bid_repo.go
│   │   ├── bid_redis_repo.go
│   │   ├── donation_repo.go
│   │   ├── final_donation.go
│   │   ├── gcp_storage_repo.go
│   │   ├── payment_repo.go
│   │   └── user_repo.go
│   │
│   ├── entity/                          # Domain models
│   │   ├── article.go
│   │   ├── auction.go
│   │   ├── bid.go
│   │   ├── donation.go
│   │   ├── final_donation.go
│   │   ├── payment.go
│   │   └── user.go
│   │
│   ├── dto/                             # Data transfer objects
│   │   ├── admin_dto.go
│   │   ├── article_dto.go
│   │   ├── auction_dto.go
│   │   ├── bid_dto.go
│   │   ├── donation_dto.go
│   │   ├── final_donation_dto.go
│   │   ├── payment_dto.go
│   │   └── user_dto.go
│   │
│   ├── mocks/                           # Generated mock repositories
│   │   └── mock_*.go
│   │
│   └── utils/                           # Utility functions
│       ├── aplouder.go                  # File upload helper
│       ├── auth.go                      # Auth utilities
│       ├── jwt.go                       # JWT token utilities
│       ├── response.go                  # Standardized API responses
│       └── validator.go                 # Validation utilities
│
├── api/
│   ├── middleware/
│   │   ├── admin.go                     # Admin role check
│   │   ├── auth.go                      # JWT authentication
│   │   ├── logging.go                   # Request logging
│   │   └── role.go                      # Role-based access control
│   │
│   └── routes/                          # API route definitions
│       ├── admin_routes.go
│       ├── article_routes.go
│       ├── auction_routes.go
│       ├── bid_routes.go
│       ├── donation_routes.go
│       ├── final_donation_routes.go
│       ├── payment_routes.go
│       ├── routes.go
│       └── user_routes.go
│
├── cron/
│   └── main.go                          # Scheduled jobs
│
├── docs/                                # Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── migrations/
│   ├── 001_init.sql                     # Schema and enum definitions
│   ├── 002_triggers.sql                 # Database triggers and functions
│   └── 003_seed.sql                     # Seed data for testing
│
├── .env.example                         # Environment variables template
├── go.mod                               # Go module dependencies
├── go.sum                               # Dependency checksums
├── Makefile                             # Build and test commands
└── generate_mocks.sh                    # Mock generation script
```

---

## Database Schema

![erd](assets/erd.jpg)

### User Roles
```sql
CREATE TYPE user_role AS ENUM ('donor', 'verifikator', 'admin', 'bidder');
```

### Donation Status
```sql
CREATE TYPE donation_status AS ENUM (
    'pending',
    'verified_for_auction',
    'verified_for_donation'
);
```

### Verification Decision
```sql
CREATE TYPE verification_decision AS ENUM ('auction', 'donation');
```

### Auction Item Status
```sql
CREATE TYPE auction_item_status AS ENUM ('scheduled', 'ongoing', 'finished');
```

### Payment Status
```sql
CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed');
```

### Core Tables

#### users
- Manages all system users (donors, verifiers, bidders, admins)
- Stores authentication credentials and role assignments

#### donations
- Records all submitted donation items
- Tracks verification status and item details

#### donation_photos
- Stores multiple photos per donation item
- Links to Google Cloud Storage URLs

#### verifications
- Records verification assessments
- Links verifiers to donations with decisions

#### auction_sessions
- Manages scheduled auction periods
- Defines start and end times for bidding

#### auction_items
- Lists items approved for auction
- Includes starting price and session assignment

#### bids
- Records all bid attempts
- Maintains complete bid history per item

#### payments
- Tracks payment transactions
- Links winners to their payment obligations

#### final_donations
- Records items distributed directly to institutions
- Maintains distribution notes and tracking

#### articles
- Stores weekly transparency reports
- Documents auction results and fund allocation

---

## API Endpoints

### Authentication (2 endpoints)
```
POST   /register               Register new user
POST   /login                  User authentication
```

### Donations (6 endpoints)
```
POST   /donations              Create donation submission
GET    /donations              List donations (admin: all, user: own)
GET    /donations/{id}         Get donation details
PUT    /donations/{id}         Update donation
PATCH  /donations/{id}         Update donation status (admin only)
DELETE /donations/{id}         Delete donation
```

### Auction Items (5 endpoints)
```
GET    /auction/items          List auction items
GET    /auction/items/{id}     Get item details
POST   /auction/items          Create auction item (admin only)
PUT    /auction/items/{id}     Update auction item (admin only)
DELETE /auction/items/{id}     Remove auction item (admin only)
```

### Auction Sessions (5 endpoints)
```
POST   /auction/sessions       Create auction session (admin only)
GET    /auction/sessions       List all sessions
GET    /auction/sessions/{id}  Get session details
PUT    /auction/sessions/{id}  Update session (admin only)
DELETE /auction/sessions/{id}  Delete session (admin only)
```

### Bidding (3 endpoints)
```
POST   /auction/sessions/{sessionID}/items/{itemID}/bid         Place bid on item
GET    /auction/sessions/{sessionID}/items/{itemID}/highest-bid Get highest bid
POST   /auction/sessions/{sessionID}/items/{itemID}/sync        Sync highest bid from Redis
```

### Final Donations (4 endpoints)
```
GET    /donations/final              List final donations (admin: all, user: own)
GET    /donations/final/me           Get my final donations
GET    /donations/final/user/{id}    Get final donations by user (admin only)
POST   /donations/final/notes        Add notes to final donation
```

### Articles (2 endpoints)
```
POST   /articles               Publish article (admin only)
GET    /articles               List all articles
GET    /articles/{id}          Get article details
```

### Payments (4 endpoints)
```
POST   /payments/{auctionId}   Create payment for auction
GET    /payments               Get all payments
GET    /payments/{id}          Get payment details
GET    /payments/status/{id}   Check payment status via Midtrans
```

### Admin (1 endpoint)
```
GET    /admin/dashboard        Get dashboard analytics
```

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- Redis 7+
- Docker & Docker Compose (optional)
- Google Cloud account (for Cloud Storage)
- Midtrans account (for payment processing)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd donation-auction-platform
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start infrastructure services**
   ```bash
   docker-compose up -d postgres redis
   ```

5. **Run database migrations**
   ```bash
   go run migrations/migrate.go up
   ```

6. **Start the application**
   ```bash
   go run app/main.go
   ```

The API will be available at `http://localhost:8080`

### Docker Deployment

```bash
# Build and run all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

---

## Environment Variables

Create a `.env` file in the root directory:

```env
# Database
POSTGRE_URL=postgres://user:password@host:5432/dbname?sslmode=disable

# Application
PORT=8000

# JWT
SECRET_KEY=your_jwt_secret_key
EXPIRED_JWT=24

# Google Cloud Storage
BUCKET_PUBLIC=your-gcp-public-bucket-name
BUCKET_PRIVATE=your-gcp-private-bucket-name
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json

# AI Service
GEMINI_API_KEY=your_gemini_api_key

# Redis
REDIS_URL=rediss://default:password@host:6379

# Midtrans
MIDTRANS_SERVER_KEY=your_server_key
MIDTRANS_CLIENT_KEY=your_client_key
```

---

## Running Migrations

### Up Migrations
```bash
# Run all pending migrations
go run migrations/migrate.go up

# Run specific migration
go run migrations/migrate.go up 001
```

### Down Migrations
```bash
# Rollback last migration
go run migrations/migrate.go down

# Rollback to specific version
go run migrations/migrate.go down 001
```

### Create New Migration
```bash
# Generate new migration file
go run migrations/migrate.go create add_new_feature
```

---

## Future Features

The following features are planned for post-MVP implementation to enhance system capabilities and user experience.

### 1. Real-Time Bidding (WebSocket + Redis)

**Implementation**: WebSocket connections with Redis pub/sub

**Key Benefits**:
- Live bid updates without page refresh
- Reduced server load through event-driven architecture
- Enhanced user experience with instant feedback
- Redis-based current highest bid synchronization

**Technical Highlights**:
- Demonstrates understanding of concurrent systems
- Event-driven architecture implementation
- Real-time data synchronization patterns
- Scalable WebSocket connection management

**Use Cases**:
- Live auction dashboards
- Real-time bid notifications
- Instant winner announcements
- Active bidder count tracking

---

### 2. AI Image Classification (Google Vision API / OpenAI Vision)

**Implementation**: Automated item assessment using computer vision

**Key Benefits**:
- Automated item quality assessment
- Reduced manual verification workload
- Faster processing time for donations
- Consistent evaluation criteria

**Technical Highlights**:
- Integration with modern AI services
- Practical application of machine learning
- Image processing pipeline implementation
- Confidence scoring and fallback logic

**Use Cases**:
- Automatic item condition grading
- Category suggestion based on image content
- Duplicate item detection
- Quality threshold enforcement

**Assessment Criteria**:
- Item condition (new, good, fair, poor)
- Category classification accuracy
- Auction eligibility recommendation
- Estimated value range suggestion

---

### 3. Fraud & Abuse Detection

**Implementation**: Redis-based rate limiting with rules engine

**Key Benefits**:
- Prevention of bid manipulation
- Spam and bot protection
- Fair auction environment
- System integrity maintenance

**Technical Highlights**:
- Advanced rate limiting strategies
- Pattern recognition for suspicious behavior
- Real-time abuse detection
- Automated response mechanisms

**Detection Mechanisms**:

**a. Bid Spam Prevention**
- Rate limit: Maximum 5 bids per minute per user
- Progressive cooldown for repeated violations
- Temporary account suspension for severe cases

**b. Suspicious Bidding Patterns**
- Detection of coordinated bid manipulation
- Identification of shill bidding (fake bids to inflate price)
- Monitoring of last-second bid sniping patterns
- Analysis of bid timing and amount patterns

**c. Account Abuse Prevention**
- Multiple account detection using device fingerprinting
- IP-based rate limiting for registration
- Email/phone verification requirements
- Behavioral analysis across sessions

**Rules Engine Examples**:
```
Rule 1: If user bids > 10 times on same item → Flag for review
Rule 2: If 90% of bids occur in last 30 seconds → Suspicious pattern
Rule 3: If bid amount exactly matches competitor + 1 repeatedly → Possible automation
Rule 4: If account created < 24h ago places high-value bid → Require verification
```

**Response Actions**:
- Soft limit: Warning notification to user
- Medium: Temporary bid cooldown (5-15 minutes)
- Hard: Account temporary suspension
- Permanent: Ban for repeated severe violations

---

### Implementation Priority

1. **Phase 1 (Post-MVP)**: Real-Time Bidding
   - Most impactful for user experience
   - Establishes WebSocket infrastructure for future features

2. **Phase 2**: AI Image Classification
   - Reduces operational burden significantly
   - Improves verification consistency

3. **Phase 3**: Fraud & Abuse Detection
   - Critical for system integrity at scale
   - Protects revenue and reputation

---

## Development Roadmap

### MVP (Current)
- Core donation and auction workflows
- User management and authentication
- Payment processing integration
- Basic reporting and transparency

### Version 1.1
- Real-time bidding implementation
- Enhanced notification system
- Mobile app support

### Version 1.2
- AI-powered image classification
- Advanced analytics dashboard
- Multi-language support

### Version 2.0
- Fraud detection system
- Advanced reporting tools
- Institution portal
- API for third-party integrations

---

## Contributing

Contributions are welcome. Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -m 'Add some feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Open a Pull Request

### Code Style
- Follow standard Go formatting (`gofmt`)
- Write meaningful commit messages
- Add tests for new features
- Update documentation as needed

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## Contact & Support

For questions, issues, or contributions, please open an issue on the repository or contact the development team.

**Project Status**: Active Development

**Documentation**: Full API documentation available at `/api/docs` when running in development mode.

**Testing**: Run tests with `go test ./...`

--- 
## Contributor
- Rafly Ade Kusuma
- Deden Ruslan
- Aisiya Qutwatunnada