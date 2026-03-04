# Cinema Booking System

ระบบจองตั๋วหนัง — Backend เขียนด้วย Go แบบ Clean Architecture

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP Framework | Gin v1.12 |
| Database | MongoDB 7 |
| Container | Docker + Docker Compose |

---

## โครงสร้างโปรเจกต์

```
cinema-booking-system/
├── docker-compose.yml
└── backend/
    ├── main.go                        # Entry point + Gin server
    ├── Dockerfile                     # Multi-stage production build
    ├── go.mod / go.sum
    ├── cmd/
    │   └── seed/
    │       └── main.go                # Seed ข้อมูลตัวอย่างลง MongoDB
    ├── internal/
    │   └── domain/
    │       └── entity/                # Domain layer — structs + enums
    │           ├── user.go
    │           ├── cinema.go
    │           ├── booking.go
    │           └── audit_log.go
    └── pkg/
        └── database/
            └── mongodb.go             # MongoDB connection
```

---

## Domain Schema

### User
| Field | Type | หมายเหตุ |
|---|---|---|
| id | ObjectID | PK |
| google_id | string | unique |
| email | string | unique |
| display_name | string | |
| role | USER / ADMIN | default USER |
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

**Seat status:** `AVAILABLE` / `LOCKED` / `BOOKED`

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
| event_type | string | BOOKING_SUCCESS, TIMEOUT, ... |
| user_id | ObjectID | ref → users |
| cinema_id | ObjectID | ref → cinemas |
| details | string | |
| payload | json | |
| timestamp | timestamp | |

---

## วิธีรัน

### รันทั้งระบบ

```bash
docker compose up -d
```

- App: http://localhost:8080
- MongoDB: localhost:27017

### Seed ข้อมูลตัวอย่าง

```bash
cd backend
MONGO_URI=mongodb://root:rootpassword@localhost:27017 \
MONGO_DB=cinema_booking \
go run cmd/seed/main.go
```

จะสร้าง collection: `users` (3), `cinemas` (2), `bookings` (1)

---

## MongoDB Compass

เชื่อมต่อด้วย URI:

```
mongodb://root:rootpassword@localhost:27017/?authSource=admin
```

---

## API Endpoints

| Method | Path | คำอธิบาย |
|---|---|---|
| GET | `/health` | Health check |

---

## Environment Variables

| Variable | Default | คำอธิบาย |
|---|---|---|
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection URI |
| `MONGO_DB` | `cinema_booking` | ชื่อ database |
| `APP_PORT` | `8080` | Port ของ server |
