// MongoDB init script — runs automatically on fresh volume (docker compose down -v && up --build)
// Seeded collections: users, cinemas, bookings

db = db.getSiblingDB('cinema_booking');

// ── Users ──────────────────────────────────────────────────────────────────
db.users.drop();
db.users.insertMany([
  { _id: new ObjectId(), google_id: 'google-001', email: 'alice@example.com',  display_name: 'Alice', role: 'USER',  created_at: new Date() },
  { _id: new ObjectId(), google_id: 'google-002', email: 'bob@example.com',    display_name: 'Bob',   role: 'USER',  created_at: new Date() },
  { _id: new ObjectId(), google_id: 'google-admin', email: 'admin@cinema.com', display_name: 'Admin', role: 'ADMIN', created_at: new Date() },
]);
print('✓ users seeded:', db.users.countDocuments());

// ── Cinemas ────────────────────────────────────────────────────────────────
db.cinemas.drop();

const now = new Date();
const h = (n) => new Date(now.getTime() + n * 3600000);

function makeSeats(rows, cols, basePrice, premiumRows, booked, locked) {
  const seats = [];
  for (let r = 0; r < rows; r++) {
    const rowLetter = String.fromCharCode(65 + r);
    const price = r >= rows - premiumRows ? basePrice + 50 : basePrice;
    for (let c = 1; c <= cols; c++) {
      const seatNo = rowLetter + c;
      let status = 'AVAILABLE';
      if (booked.includes(seatNo)) status = 'BOOKED';
      else if (locked.includes(seatNo)) status = 'LOCKED';
      seats.push({ seat_no: seatNo, status: status, price: price });
    }
  }
  return seats;
}

db.cinemas.insertMany([
  {
    _id: new ObjectId(),
    movie_name: 'Avengers: Endgame',
    theater_no: 1,
    start_time: h(1), end_time: h(4),
    price: 200,
    seats: makeSeats(5, 10, 200, 1, ['A1','A2','A5','B3','B4','C1'], ['D2','D3']),
    created_at: now,
  },
  {
    _id: new ObjectId(),
    movie_name: 'Avengers: Endgame',
    theater_no: 1,
    start_time: h(6), end_time: h(9),
    price: 200,
    seats: makeSeats(5, 10, 200, 1, ['B1','B2'], []),
    created_at: now,
  },
  {
    _id: new ObjectId(),
    movie_name: 'Interstellar',
    theater_no: 2,
    start_time: h(2), end_time: h(5),
    price: 220,
    seats: makeSeats(6, 8, 220, 2, ['A1','A2','A3','B5','B6','C2','C3','C4'], ['E1','E2']),
    created_at: now,
  },
  {
    _id: new ObjectId(),
    movie_name: 'Dune: Part Two',
    theater_no: 2,
    start_time: h(7), end_time: h(10),
    price: 250,
    seats: makeSeats(6, 8, 250, 2, ['A1','A2','B1','B2'], ['C5','C6']),
    created_at: now,
  },
  {
    _id: new ObjectId(),
    movie_name: 'Oppenheimer',
    theater_no: 3,
    start_time: h(3), end_time: h(6),
    price: 180,
    seats: makeSeats(4, 12, 180, 1, ['A1','A2','A3','A4','B6','B7'], []),
    created_at: now,
  },
  {
    _id: new ObjectId(),
    movie_name: 'Spider-Man: No Way Home',
    theater_no: 3,
    start_time: h(8), end_time: h(11),
    price: 180,
    seats: makeSeats(4, 12, 180, 1, [], []),
    created_at: now,
  },
]);
print('✓ cinemas seeded:', db.cinemas.countDocuments());

// ── Bookings ───────────────────────────────────────────────────────────────
db.bookings.drop();
db.bookings.insertMany([
  {
    _id: new ObjectId(),
    user_id: new ObjectId(),
    cinema_id: new ObjectId(),
    seat_numbers: ['A3'],
    total_amount: 200,
    status: 'SUCCESS',
    created_at: new Date(),
  },
]);
print('✓ bookings seeded:', db.bookings.countDocuments());
print('✓ Seed completed!');
