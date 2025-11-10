# Quest Auth - Domain Events

## 📡 Event System Overview

Quest Auth uses domain events to track authentication activities and enable system integration.

---

## 🎯 Event Catalog

### UserRegistered

Emitted when a new user registers.

**Fields:**
- `user_id` - User UUID
- `email` - User email
- `phone` - User phone
- `name` - User name
- `created_at` - Timestamp

**Example:**
```json
{
  "event_type": "user.registered",
  "aggregate_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "phone": "+1234567890",
    "name": "John Doe",
    "created_at": "2025-11-10T20:00:00Z"
  }
}
```

---

### UserLoggedIn

Emitted when a user successfully logs in.

**Fields:**
- `user_id` - User UUID
- `login_at` - Login timestamp

---

### UserPhoneChanged

Emitted when a user changes their phone number.

**Fields:**
- `user_id` - User UUID
- `old_phone` - Previous phone
- `new_phone` - New phone
- `changed_at` - Timestamp

---

### UserNameChanged

Emitted when a user changes their name.

**Fields:**
- `user_id` - User UUID
- `old_name` - Previous name
- `new_name` - New name
- `changed_at` - Timestamp

---

### UserPasswordChanged

Emitted when a user changes their password.

**Fields:**
- `user_id` - User UUID
- `changed_at` - Timestamp

---

## 🔄 Event Flow

```
1. Domain Operation
   ↓
2. Aggregate adds event to internal list
   ↓
3. Use case handler publishes events
   ↓
4. Events persisted in PostgreSQL
   ↓
5. Transaction commits
   ↓
6. Events cleared from aggregate
```

---

## 💾 Event Storage

Events are stored in PostgreSQL `events` table:

```sql
CREATE TABLE events (
    id UUID PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    aggregate_id UUID NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL
);
```

---

## 🔮 Future: Message Queue

Phase 2 will integrate message broker (RabbitMQ/Kafka) for:
- Real-time event streaming
- External system integration
- Event-driven microservices

---

**Last Updated:** November 10, 2025
