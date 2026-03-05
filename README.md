# Cinema Booking System

ระบบจองตั๋วหนังออนไลน์ รองรับ Real-time Seat Map, Distributed Lock, และป้องกัน Double Booking

---

## System Architecture Diagram

```
[ Browser / Frontend (Vue 3) ]
              │
              ▼
      [ Backend API (Go + Gin) ]
          │        │         │
          │        │         ├────────► [ Redis Pub-Sub ]
          │        │                      │
          │        │                      ▼
          │        └────────► [ Redis (Cache + Distributed Lock) ]
          │
          └────────────────► [ MongoDB (Primary Data Store) ]

Real-time:
Browser ── WebSocket ──► Backend ── seat_update ──► Broadcast to all clients
```

---

## Tech Stack Overview

| Layer | Technology | หมายเหตุ |
|---|---|---|
| Language | Go 1.24 | Clean Architecture |
| HTTP Framework | Gin v1.12 | REST API |
| Frontend | Vue 3 + Vite | TypeScript, Composition API |
| Database | MongoDB 7 | Embedded seat document |
| Cache / Lock | Redis 7 | ACL Auth + Distributed Lock + Pub-Sub |
| Real-time | WebSocket (gorilla/websocket) | Room-based Hub |
| Message Queue | Redis Pub-Sub | Async Audit Logging |
| Auth | Google OAuth 2.0 + JWT (HS256) | Role: USER / ADMIN |
| Container | Docker + Docker Compose | รันด้วยคำสั่งเดียว |

---

## Booking Flow — Step by Step

1. **Login**
   - ผู้ใช้กดปุ่ม "Sign in with Google" ที่ frontend → redirect ไป Google OAuth
   - Backend รับ callback, สร้าง/อัปเดต `User` แล้วออก JWT token ส่งกลับให้ frontend
2. **Browse Cinemas**
   - Frontend เรียก `GET /api/cinema` เพื่อดึกรอบฉาย (filter ตามวันที่ได้)
   - เมื่อเลือกเรื่อง/รอบ ระบบจะโหลดสถานะที่นั่งด้วย `GET /api/seats/:cinemaId`
3. **Subscribe Real-time Seat Map**
   - Frontend เปิด WebSocket `GET /api/ws/seats/:cinemaId?token=<jwt>`
   - เมื่อมีใครล็อก/จองที่นั่ง ระบบจะ broadcast `seat_update` ให้ทุก client ในห้องเดียวกัน
4. **Lock Seats (Hold)**
   - ผู้ใช้เลือกที่นั่ง → Frontend เรียก `POST /api/booking/lock`
   - Backend พยายาม acquire Redis Lock สำหรับทุกที่นั่ง + อัปเดต MongoDB เป็น `LOCKED`
   - ถ้าสำเร็จ: สร้าง booking สถานะ `PENDING` และ broadcast seat_update
5. **Confirm Booking**
   - ภายในเวลาที่กำหนด (เช่น 5 นาที) ผู้ใช้กดปุ่มยืนยัน → `POST /api/booking/confirm`
   - Backend เปลี่ยนสถานะที่นั่งเป็น `BOOKED`, ปลด Redis Lock, อัปเดต booking เป็น `SUCCESS`
6. **Timeout / Cancel**
   - Background job ตรวจ booking ที่ `PENDING` เกิน TTL:
     - ปลด Redis Lock, เปลี่ยนที่นั่งกลับเป็น `AVAILABLE`, ตั้ง booking เป็น `TIMEOUT`
   - ผู้ใช้สามารถยกเลิกเองผ่าน `POST /api/booking/cancel` (เปลี่ยนเป็น `CANCELLED`)
7. **Audit & Admin View**
   - ทุก action สำคัญ (lock, confirm, timeout, cancel) จะถูก publish เป็น event ไป Redis Pub-Sub
   - Consumer อ่าน event แล้วเขียน `AuditLog` ลง MongoDB → Admin ดูย้อนหลังได้ผ่านหน้า Admin

---

## โครงสร้างโปรเจกต์

```
cinema-booking-system/
├── docker-compose.yml
├── .env.example
├── README.md
├── redis/
│   └── entrypoint.sh        # Generate ACL file จาก env vars ตอน startup
├── mongo-init/
│   └── 01-seed.js           # Seed ข้อมูลตัวอย่าง (รันอัตโนมัติตอน fresh volume)
├── backend/
│   ├── Dockerfile
│   ├── go.mod / go.sum
│   ├── cmd/
│   │   ├── api/main.go      # Entry point — wire ทุก dependency + run migrations
│   │   └── seed/main.go     # Seed ข้อมูลตัวอย่าง
│   ├── internal/
│   │   ├── delivery/http/   # Handlers, Middleware, Router, WebSocket Hub
│   │   ├── domain/
│   │   │   ├── entity/      # Structs + Enums (User, Cinema, Booking, AuditLog)
│   │   │   ├── repository/  # Interfaces (CinemaRepo, SeatLocker, EventPublisher ...)
│   │   │   └── usecase/     # Interfaces (BookingUseCase, AuthUseCase ...)
│   │   ├── repository/
│   │   │   ├── mongodb/     # MongoDB implementations
│   │   │   └── redis/       # Redis implementations (SeatLocker, EventPublisher, AuditConsumer)
│   │   └── usecase/         # Concrete use case implementations
│   └── pkg/
│       ├── database/        # MongoDB connection + Migrate() (index setup)
│       ├── jwt/             # JWT generate + parse
│       ├── redis/           # Redis client helper
│       └── redislock/       # Distributed Lock (SET NX EX + Lua release script)
└── frontend/
    ├── Dockerfile
    ├── nginx.conf
    └── src/
        ├── views/           # LoginView, HomeView, CinemaDetailView, BookingsView
        │   └── admin/       # AdminBookingsView, AdminCinemasView, AdminAuditLogsView, AdminCreateShowtimeView
        ├── components/      # CinemaCard, SeatGrid, SeatCell, AppNavbar ...
        ├── composables/     # useAuth, useBooking, useCinemas
        └── router/          # Vue Router + auth guards
```

---

## Redis Lock Strategy — Distributed Lock

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
| seats | []Seat | embedded — seat_no (string), status, price, user_id |
| created_at | timestamp | |

**Seat status:** `AVAILABLE` → `LOCKED` → `BOOKED` (หรือกลับไป `AVAILABLE` เมื่อ timeout)

### Booking
| Field | Type | หมายเหตุ |
|---|---|---|
| id | ObjectID | PK |
| user_id | ObjectID | ref → users |
| cinema_id | ObjectID | ref → cinemas |
| seat_numbers | []string | |
| total_amount | float64 | |
| status | PENDING / SUCCESS / TIMEOUT / CANCELLED | |
| created_at | timestamp | |

### AuditLog
| Field | Type | หมายเหตุ |
|---|---|---|
| id | ObjectID | PK |
| event_type | string | BOOKING_SUCCESS, TIMEOUT, BOOKING_CANCELLED, SEAT_LOCKED, SEAT_RELEASED |
| user_id | ObjectID | ref → users |
| cinema_id | ObjectID | ref → cinemas |
| details | string | |
| payload | json | |
| timestamp | timestamp | |

---

## MongoDB Indexes (Auto-created on startup)

`pkg/database/Migrate()` สร้าง indexes ทุกครั้งที่ server start (idempotent — ปลอดภัยถ้ามีอยู่แล้ว)

| Collection | Index | หมายเหตุ |
|---|---|---|
| `users` | `google_id` | unique + sparse |
| `users` | `email` | unique |
| `cinemas` | `start_time` | date-range filter |
| `bookings` | `user_id` | query |
| `bookings` | `cinema_id` | query |
| `bookings` | `status + created_at` | compound — timeout ticker |
| `audit_logs` | `timestamp` (desc) | sort |
| `audit_logs` | `user_id` | query |
| `audit_logs` | `cinema_id` | query |

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
| POST | `/api/booking/cancel` | ยกเลิกการจอง |
| GET | `/api/booking/me` | รายการจองของตัวเอง |

### Admin (ต้องมี JWT + role ADMIN)
| Method | Path | คำอธิบาย |
|---|---|---|
| GET | `/api/admin/cinema` | รายการรอบฉายทั้งหมด |
| POST | `/api/admin/showtimes` | สร้างรอบฉายใหม่ |
| GET | `/api/admin/bookings` | รายการจองทั้งหมด (filter by movie, date, status) |
| GET | `/api/admin/audit-logs` | ดู Audit Logs |

#### Query Parameters — `GET /api/admin/bookings`
| Parameter | ตัวอย่าง | คำอธิบาย |
|---|---|---|
| `movie` | `Avengers` | ค้นหาชื่อหนัง (case-insensitive, partial match) |
| `date` | `2026-03-05` | กรองตามวันที่จอง (YYYY-MM-DD) |
| `status` | `SUCCESS` | กรองตามสถานะ (PENDING / SUCCESS / TIMEOUT / CANCELLED) |
| `page` | `1` | หน้าที่ต้องการ (default: 1) |
| `limit` | `20` | จำนวนรายการต่อหน้า (default: 20) |

---

## Environment Variables

| Variable | ตัวอย่าง | คำอธิบาย |
|---|---|---|
| `MONGO_ROOT_USERNAME` | `root` | MongoDB root username |
| `MONGO_ROOT_PASSWORD` | `rootpassword` | MongoDB root password |
| `MONGO_DB` | `cinema_booking` | ชื่อ database |
| `MONGO_URI` | `mongodb://root:pass@mongodb:27017/...` | MongoDB connection URI |
| `REDIS_USERNAME` | `cinema_user` | Redis ACL username |
| `REDIS_PASSWORD` | `<strong-password>` | Redis ACL password |
| `REDIS_URL` | `redis://cinema_user:pass@redis:6379` | Redis connection URL (รวม credentials) |
| `JWT_SECRET` | `<random 32+ chars>` | HS256 signing key |
| `GOOGLE_CLIENT_ID` | `xxx.apps.googleusercontent.com` | Google OAuth Client ID |
| `GOOGLE_CLIENT_SECRET` | `GOCSPX-xxx` | Google OAuth Client Secret |
| `GOOGLE_CALLBACK_URL` | `http://localhost:8080/auth/google/callback` | OAuth redirect URI |
| `ADMIN_EMAILS` | `admin@gmail.com,cto@company.com` | Comma-separated admin emails |
| `VITE_API_URL` | `http://localhost:8080` | Backend URL สำหรับ frontend |
| `FRONTEND_URL` | `http://localhost:3000` | Frontend URL สำหรับ OAuth redirect |
| `APP_PORT` | `8080` | Backend port (optional, default 8080) |

---

## วิธีรันระบบ

### Prerequisites

- Docker Desktop
- ไฟล์ `.env` (copy จาก `.env.example` แล้วกรอก credentials)

```bash
cp .env.example .env
# แก้ไข: GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, ADMIN_EMAILS, JWT_SECRET
# แก้ไข: REDIS_USERNAME, REDIS_PASSWORD (และ REDIS_URL ให้ตรงกัน)
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

สร้าง: `users` (3 คน), `cinemas` (6 รอบ), `bookings` (1 รายการ)

---

## Assumptions & Trade-offs

- **MongoDB เป็น single primary database**
  - เลือก MongoDB เพราะเหมาะกับ document ที่ฝังรายการที่นั่ง (`seats`) ใน `Cinema` ได้
  - Trade-off: ถ้าขนาด seats โตมากขึ้นมาก ๆ อาจต้องแยก collection หรือ shard เพิ่ม
- **Embedded seat document แทนการแยกเป็น collection ใหม่**
  - ข้อดี: อ่านสถานะที่นั่งทั้งรอบได้ใน query เดียว, เขียนแบบ atomic ง่ายกว่า
  - ข้อเสีย: ขนาด document โตตามจำนวนที่นั่ง, ต้องออกแบบ index ดี ๆ
- **Redis ใช้ทั้งเป็น cache, lock และ message queue**
  - ลดจำนวน dependency และทำงานได้ดีในระบบขนาดเล็ก-กลาง
  - Trade-off: ถ้า traffic ใหญ่ขึ้นมาก อาจต้องแยก message broker จริงจัง (เช่น Kafka / RabbitMQ)
- **Redis Lock TTL คงที่ (เช่น 300 วินาที)**
  - สมมติว่า user ส่วนใหญ่ยืนยันการจองภายในเวลานี้
  - ถ้า business ต้องการเวลานานขึ้น อาจต้องเพิ่ม mechanism สำหรับ extend lock แทนการ fixed TTL
- **Audit Logging แบบ async**
  - เลือกให้ booking path เร็วที่สุด แม้ audit จะล้มเหลว log บางส่วนอาจหายได้
  - ถ้าต้องการ audit แบบ strong guarantee อาจต้องใช้ transactional outbox แทน Pub-Sub ตรง ๆ

---


## Optional Features (Extra Credit)

- **Postman Collection**
  - มีไฟล์ collection พร้อมใช้งานที่ `postman/cinema-booking.postman_collection.json`
  - ตั้งค่า `base_url`, `jwt_token`, `admin_jwt_token` เป็น Postman variables แล้วสามารถยิง endpoint หลักได้ครบ (`/auth/*`, `/api/cinema`, `/api/seats/:cinemaId`, `/api/booking/*`, `/api/admin/*`)
- **Simple Test Case**
  - มี unit test สำหรับ `BookingUseCase` (happy-path lock + confirm) ที่ไฟล์ `backend/internal/usecase/booking/usecase_test.go`
  - รันเทสต์ทั้งหมดของ backend ได้ด้วยคำสั่ง:
    - `cd backend && go test ./...`
- **Notification (Email / Line / Mock)**
  - มี abstraction `NotificationService` ที่ถูก inject เข้า `BookingUseCase` และถูกเรียกหลัง booking `SUCCESS` และ `TIMEOUT`
  - ปัจจุบันมี implementation แบบ mock คือ `LoggerNotificationService` ที่ log ข้อความแจ้งเตือน (ต่อยอดไปเป็น Email / Line Notify ได้ในอนาคต)

---

## Security

- **JWT HS256** — ทุก request ต้องแนบ `Authorization: Bearer <token>`
- **Role-based access** — middleware `RequireRole()` แยก USER / ADMIN ชัดเจน
- **Rate limiting** — per-IP token bucket (20 req/s, burst 50)
- **Security headers** — X-Frame-Options, X-Content-Type-Options, Referrer-Policy
- **CORS** — จำกัดเฉพาะ FRONTEND_URL
- **Redis ACL** — ปิด `default` user, ใช้ named user พร้อม password แทน
- **No hardcode** — ทุก secret ใช้ environment variables

---

## เชื่อมต่อ Database Tools

### MongoDB Compass
```
mongodb://root:<MONGO_ROOT_PASSWORD>@localhost:27017/?authSource=admin
```

### Redis Insight
| Field | ค่า |
|---|---|
| Host | `localhost` |
| Port | `6379` |
| Username | ค่าจาก `REDIS_USERNAME` ใน `.env` |
| Password | ค่าจาก `REDIS_PASSWORD` ใน `.env` |

### Redis CLI
```bash
docker exec -it redis_container redis-cli -u redis://<REDIS_USERNAME>:<REDIS_PASSWORD>@localhost:6379
```
