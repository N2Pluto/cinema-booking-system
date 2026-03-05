# Cinema Booking System

ระบบจองตั๋วหนังออนไลน์ รองรับ Real-time Seat Map, Distributed Lock, และป้องกัน Double Booking แบบ 100%

---

## Tech Stack

| Layer | Technology | หมายเหตุ |
|---|---|---|
| Language | Go 1.24 | Clean Architecture |
| HTTP Framework | Gin v1.12 | REST API |
| Frontend | Vue 3 + Vite | TypeScript, Composition API |
| Database | MongoDB 7 | Embedded seat document |
| Cache / Lock | Redis 7 | Distributed Lock + Pub-Sub |
| Real-time | WebSocket (gorilla/websocket) | Room-based Hub |
| Message Queue | Redis Pub-Sub | Async Audit Logging |
| Auth | Google OAuth 2.0 + JWT (HS256) | Role: USER / ADMIN |
| Container | Docker + Docker Compose | รันด้วยคำสั่งเดียว |

---

## โครงสร้างโปรเจกต์

```
cinema-booking-system/
├── docker-compose.yml
├── .env.example
├── README.md
├── backend/
│   ├── Dockerfile
│   ├── go.mod / go.sum
│   ├── cmd/
│   │   ├── api/main.go          # Entry point — wire ทุก dependency
│   │   └── seed/main.go         # Seed ข้อมูลตัวอย่าง
│   ├── internal/
│   │   ├── delivery/http/       # Handlers, Middleware, Router, WebSocket Hub
│   │   ├── domain/
│   │   │   ├── entity/          # Structs + Enums (User, Cinema, Booking, AuditLog)
│   │   │   ├── repository/      # Interfaces (CinemaRepo, SeatLocker, EventPublisher ...)
│   │   │   └── usecase/         # Interfaces (BookingUseCase, AuthUseCase ...)
│   │   ├── repository/
│   │   │   ├── mongodb/         # MongoDB implementations
│   │   │   └── redis/           # Redis implementations (SeatLocker, EventPublisher, AuditConsumer)
│   │   └── usecase/             # Concrete use case implementations
│   └── pkg/
│       ├── database/            # MongoDB connection helper
│       ├── jwt/                 # JWT generate + parse
│       ├── redis/               # Redis client helper
│       └── redislock/           # Distributed Lock (SET NX EX + Lua release script)
└── frontend/
    ├── Dockerfile
    ├── nginx.conf
    └── src/
        ├── views/               # LoginView, HomeView, CinemaDetailView, BookingsView
        │   └── admin/           # AdminCinemasView, AdminAuditLogsView, AdminCreateShowtimeView
        ├── components/          # CinemaCard, SeatGrid, SeatCell, AppNavbar ...
        ├── composables/         # useAuth, useBooking, useCinemas
        └── router/              # Vue Router + auth guards
```

---

## กลยุทธ์การ Lock — Distributed Lock Strategy

### ปัญหา: Double Booking

เมื่อผู้ใช้หลายคนเลือกที่นั่งเดียวกันพร้อมกัน (Race Condition) ระบบต้องรับประกันว่า **ที่นั่งหนึ่งจะถูกจองได้โดยคนเดียวเท่านั้น**

### วิธีแก้: Two-Layer Locking

ระบบใช้ Lock 2 ชั้นร่วมกัน ทำงานตามลำดับดังนี้:

#### ชั้นที่ 1 — Redis Distributed Lock (ป้องกัน Race Condition)

```
key:   lock:seat:{cinemaID}:{seatNo}
value: {userID}               ← token สำหรับ ownership verification
ttl:   5 นาที (300 วินาที)
cmd:   SET lock:seat:abc:A1 user123 NX EX 300
```

- `NX` (Not eXists) คือหัวใจสำคัญ — ถ้า key มีอยู่แล้ว คำสั่งจะ **ล้มเหลวทันที** โดยไม่ต้องรอ
- ผู้ใช้คนที่ 2 ที่พยายามล็อกที่นั่งเดียวกัน จะได้รับ `ErrNotAcquired` กลับทันที
- ถ้าล็อกที่นั่งหลายที่พร้อมกัน แต่ล็อกที่นั่งใดที่นั่งหนึ่งล้มเหลว ระบบจะ **rollback locks ที่ได้มาแล้วทั้งหมด** (all-or-nothing)
- Release ใช้ Lua Script เพื่อให้ atomic: `DEL key` เฉพาะเมื่อ value ตรงกับ token ของตัวเอง

#### ชั้นที่ 2 — MongoDB Atomic UpdateOne (ป้องกัน Inconsistency)

หลังจากได้ Redis Lock แล้ว ระบบจะ update status ที่นั่งใน MongoDB ด้วย:
- Filter: `$expr` ตรวจสอบว่าที่นั่งทุกตัวยัง `AVAILABLE` อยู่
- ถ้า `MatchedCount == 0` แสดงว่าสถานะใน DB ไม่ตรง → return error + release Redis locks

การใช้ 2 ชั้นนี้ทำให้มั่นใจได้ว่า:
1. Redis ป้องกัน concurrent requests ก่อนถึง DB
2. MongoDB ป้องกัน edge case ที่ Redis TTL หมดไปแล้วก่อน confirm

### Flow การจองที่นั่ง

```
User เลือกที่นั่ง
    │
    ▼
[1] Acquire Redis Lock สำหรับทุกที่นั่ง (SET NX EX 300)
    ├── ล้มเหลว → คืน error "seat already locked" ทันที
    └── สำเร็จ ↓
[2] MongoDB UpdateOne (atomic) → เปลี่ยน status เป็น LOCKED
    ├── ล้มเหลว → Release Redis locks ทั้งหมด + คืน error
    └── สำเร็จ ↓
[3] สร้าง Booking (status: PENDING)
[4] Publish AuditLog event → Redis Pub-Sub → Consumer เขียน MongoDB
[5] Broadcast seat_update → WebSocket → ทุก client เห็นทันที
    │
    ├── ภายใน 5 นาที: User กด Confirm
    │       ├── MongoDB status → BOOKED
    │       ├── Release Redis locks
    │       └── Publish BOOKING_SUCCESS audit event
    │
    └── เกิน 5 นาที: Background ticker (ทุก 30 วินาที)
            ├── Release Redis locks
            ├── MongoDB status → AVAILABLE
            ├── Booking status → TIMEOUT
            └── Publish TIMEOUT audit event
```

---

## Message Queue — Redis Pub-Sub

**Use Case:** Async Audit Logging

| Component | หน้าที่ |
|---|---|
| `EventPublisher` (redis) | Serialize AuditLog → `PUBLISH audit:events <json>` |
| `AuditConsumer` (redis) | `SUBSCRIBE audit:events` → deserialize → เขียน MongoDB |

การใช้ Pub-Sub แยก audit write ออกจาก booking path ทำให้:
- Booking response เร็วขึ้น (ไม่ต้องรอ DB write ของ audit)
- Audit log เขียนล้มเหลวไม่กระทบ user experience

---

## Domain Schema

### User
| Field | Type | หมายเหตุ |
|---|---|---|
| id | ObjectID | PK |
| google_id | string | unique |
| email | string | unique |
| display_name | string | |
| role | USER / ADMIN | กำหนดโดย ADMIN_EMAILS env |
| created_at | timestamp | |

### Cinema (รอบฉาย)
| Field | Type | หมายเหตุ |
|---|---|---|
| id | ObjectID | PK |
| movie_name | string | |
| theater_no | int | |
| start_time | timestamp | |
| end_time | timestamp | |
| price | int | ราคาพื้นฐาน |
| seats | []Seat | embedded — seat_no, status, price, user_id |
| created_at | timestamp | |

**Seat status:** `AVAILABLE` → `LOCKED` → `BOOKED` (หรือกลับไป `AVAILABLE` เมื่อ timeout)

### Booking
| Field | Type | หมายเหตุ |
|---|---|---|
| id | ObjectID | PK |
| user_id | ObjectID | ref → users |
| cinema_id | ObjectID | ref → cinemas |
| seat_numbers | []string | |
| total_amount | decimal | |
| status | PENDING / SUCCESS / TIMEOUT / CANCELLED | |
| created_at | timestamp | |

### AuditLog
| Field | Type | หมายเหตุ |
|---|---|---|
| id | ObjectID | PK |
| event_type | string | SEAT_LOCKED, BOOKING_SUCCESS, TIMEOUT, ... |
| user_id | ObjectID | ref → users |
| cinema_id | ObjectID | ref → cinemas |
| details | string | |
| payload | json | |
| timestamp | timestamp | |

---

## API Endpoints

### Public
| Method | Path | คำอธิบาย |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/auth/google` | Redirect ไป Google consent |
| GET | `/auth/google/callback` | รับ code จาก Google, คืน JWT token |
| GET | `/api/ws/seats/:cinemaId?token=<jwt>` | WebSocket — real-time seat updates |

### User (ต้องมี JWT)
| Method | Path | คำอธิบาย |
|---|---|---|
| GET | `/api/me` | ดูข้อมูลตัวเอง |
| GET | `/api/cinema` | รายการรอบฉาย (pagination, date filter) |
| GET | `/api/seats/:cinemaId` | สถานะที่นั่งทั้งหมดของรอบนั้น |
| POST | `/api/booking/lock` | ล็อกที่นั่ง (Redis Lock + PENDING booking) |
| POST | `/api/booking/confirm` | ยืนยันการจอง (BOOKED) |
| GET | `/api/booking/me` | รายการจองของตัวเอง |

### Admin (ต้องมี JWT + role ADMIN)
| Method | Path | คำอธิบาย |
|---|---|---|
| GET | `/api/admin/cinema` | รายการรอบฉายทั้งหมด |
| POST | `/api/admin/showtimes` | สร้างรอบฉายใหม่ |
| GET | `/api/admin/audit-logs` | ดู Audit Logs |

---

## Environment Variables

| Variable | ตัวอย่าง | คำอธิบาย |
|---|---|---|
| `MONGO_ROOT_USERNAME` | `root` | MongoDB root username |
| `MONGO_ROOT_PASSWORD` | `rootpassword` | MongoDB root password |
| `MONGO_DB` | `cinema_booking` | ชื่อ database |
| `MONGO_URI` | `mongodb://root:pass@mongodb:27017/...` | MongoDB connection URI |
| `REDIS_URL` | `redis://redis:6379` | Redis connection URL |
| `JWT_SECRET` | `<random 32+ chars>` | HS256 signing key |
| `GOOGLE_CLIENT_ID` | `xxx.apps.googleusercontent.com` | Google OAuth Client ID |
| `GOOGLE_CLIENT_SECRET` | `GOCSPX-xxx` | Google OAuth Client Secret |
| `GOOGLE_CALLBACK_URL` | `http://localhost:8080/auth/google/callback` | OAuth redirect URI |
| `ADMIN_EMAILS` | `admin@gmail.com,cto@company.com` | Comma-separated admin emails |
| `VITE_API_URL` | `http://localhost:8080` | Backend URL สำหรับ frontend |
| `FRONTEND_URL` | `http://localhost:3000` | Frontend URL สำหรับ OAuth redirect |
| `APP_PORT` | `8080` | Backend port (optional, default 8080) |

---

## วิธีรัน

### Prerequisites

- Docker Desktop
- ไฟล์ `.env` (copy จาก `.env.example` แล้วกรอก Google OAuth credentials)

```bash
cp .env.example .env
# แก้ไข GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, ADMIN_EMAILS, JWT_SECRET
```

### รันทั้งระบบ

```bash
docker compose up --build
```

| Service | URL |
|---|---|
| Frontend (Vue 3) | http://localhost:3000 |
| Backend API | http://localhost:8080 |
| MongoDB | localhost:27017 |
| Redis | localhost:6379 |

### Seed ข้อมูลตัวอย่าง

```bash
cd backend
MONGO_URI=mongodb://root:rootpassword@localhost:27017 \
MONGO_DB=cinema_booking \
go run cmd/seed/main.go
```

สร้าง: `users` (3 คน), `cinemas` (2 รอบ), `bookings` (1 รายการ)

---

## Security

- **JWT HS256** — ทุก request ต้องแนบ `Authorization: Bearer <token>`
- **Role-based access** — middleware `RequireRole()` แยก USER / ADMIN ชัดเจน
- **Rate limiting** — per-IP token bucket (20 req/s, burst 50)
- **Security headers** — X-Frame-Options, X-Content-Type-Options, Referrer-Policy
- **CORS** — จำกัดเฉพาะ FRONTEND_URL
- **No hardcode** — ทุก secret ใช้ environment variables

---

## MongoDB Compass

```
mongodb://root:rootpassword@localhost:27017/?authSource=admin
```
