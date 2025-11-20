# Code Review Guidelines

> **Purpose:** This document provides guidelines for conducting effective code reviews in the Zarinpal Platform project.

> **Important:** Code formatting, import organization, and basic linting are **automated** by git hooks and CI/CD.
> Reviewers should focus on architecture, logic, security, and maintainability rather than formatting issues.

## Table of Contents
- [Review Process](#review-process)
- [Review Checklist](#review-checklist)
- [Code Standards](#code-standards)
- [Common Issues](#common-issues)

---

## Review Process

### For Authors
1. **Before Submitting:**
   - Run `make lint` and fix all issues
   - Run `make test` and ensure all tests pass
   - Run `make format` to format code
   - Ensure commit messages follow the standard format
   - Self-review your code first

2. **Creating Merge Request:**
   - Write a clear description of what changed and why
   - Reference related issues/tickets
   - Add screenshots/examples if applicable
   - Ensure CI/CD pipeline passes

3. **During Review:**
   - Respond to all comments
   - Don't take feedback personally - it's about the code, not you
   - Ask questions if feedback is unclear
   - Mark conversations as resolved when addressed

### For Reviewers
1. **Be Constructive:**
   - Praise good code
   - Explain "why" when requesting changes
   - Suggest alternatives
   - Use "we" instead of "you"

2. **Be Timely:**
   - Review within 24 hours if possible
   - If busy, let the author know when you can review

3. **Focus on (Manual Review):**
   - Correctness & Logic
   - Architecture/Design
   - Security
   - Test Coverage
   - Maintainability
   - Performance

4. **Don't Focus on (Automated):**
   - Code formatting (gofmt/goimports handles this)
   - Import order (automated by git hooks)
   - Basic lint issues (golangci-lint catches these)
   - If you see these, the git hooks weren't run - ask author to run `make format` and `make lint`

---

## Review Checklist

### 🏗️ Architecture & Design

- [ ] **Follows Clean Architecture**
  - Domain layer is independent
  - Business logic is in domain, not handlers
  - Proper separation of concerns

- [ ] **Dependency Injection**
  - Dependencies are injected, not created
  - Interfaces are used for external dependencies
  - No global variables for state

- [ ] **Service Layer**
  - Business logic in domain services
  - Handlers are thin (just mapping)
  - Repository pattern followed

- [ ] **Code Organization**
  - Files in correct directories
  - Package names are clear
  - No circular dependencies

---

### 🐛 Error Handling

- [ ] **Errors Are Checked**
  - All errors are handled or explicitly ignored
  - No silent failures
  - Errors include context

- [ ] **Error Messages**
  - User-facing errors use locale (i18n)
  - Error messages are helpful
  - Internal errors logged with details

- [ ] **gRPC Status Codes**
  - Appropriate status codes used
  - InvalidArgument for validation errors
  - NotFound for missing resources
  - Internal for server errors
  - Unavailable for external service failures

- [ ] **Error Wrapping**
  - Errors wrapped with context
  - Stack traces preserved
  - No information loss

**Example:**
```go
// ❌ Bad
if err != nil {
    return nil, err
}

// ✅ Good
if err != nil {
    return nil, errors.NewAppError(ctx).DatabaseError()
}
```

---

### 🔐 Security

- [ ] **Input Validation**
  - All user input validated
  - Length limits enforced
  - Format validation (email, phone, etc.)
  - Type checking

- [ ] **SQL Injection Prevention**
  - Parameterized queries used
  - No string concatenation for SQL
  - ORM/query builder used correctly

- [ ] **Authentication & Authorization**
  - Endpoints require auth if needed
  - Permissions checked
  - User context propagated

- [ ] **Secrets Management**
  - No hardcoded secrets
  - Environment variables or config
  - Secrets not logged

- [ ] **Data Exposure**
  - Sensitive data not in logs
  - PII handled carefully
  - Error messages don't leak info

**Example:**
```go
// ❌ Bad - SQL Injection risk
query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID)

// ✅ Good - Parameterized query
query := "SELECT * FROM users WHERE id = $1"
row := db.QueryRow(ctx, query, userID)
```

---

### 📊 Observability

- [ ] **Tracing**
  - Operations traced with OpenTelemetry
  - Span names are descriptive
  - Important attributes added
  - Errors recorded in spans

- [ ] **Logging**
  - Structured logging used
  - Appropriate log levels
  - No sensitive data in logs
  - Contextual information included

- [ ] **Metrics**
  - Important operations measured
  - Counters for events
  - Histograms for durations
  - Metrics properly labeled

**Example:**
```go
// ✅ Good tracing
ctx, span := trace.GetTracer().Start(ctx, "UserService.CreateUser")
defer span.End()

if err != nil {
    span.RecordError(err)
    return nil, err
}
```

---

### 🧪 Testing

- [ ] **Test Coverage**
  - New code has tests
  - Coverage > 70%
  - Important paths tested
  - Edge cases covered

- [ ] **Test Quality**
  - Tests are clear and focused
  - One assertion per test (generally)
  - Tests are independent
  - No flaky tests

- [ ] **Test Types**
  - Unit tests for business logic
  - Integration tests for repositories
  - Mocks for external dependencies

- [ ] **Test Naming**
  - Clear test names
  - Follows pattern: `Test<Function>_<Scenario>_<Expected>`

**Example:**
```go
// ✅ Good test naming
func TestUserService_CreateUser_WithValidData_ReturnsSuccess(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)

    // Act
    user, err := service.CreateUser(ctx, validUserData)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

---

### 📝 Code Style & Readability

> **Note:** Formatting (gofmt, imports) is checked automatically by git hooks and linters.
> Focus on readability and maintainability during manual review.

- [ ] **Naming**
  - Variables named clearly
  - Functions describe what they do
  - No abbreviations (unless common)
  - Consistent naming

- [ ] **Function Size**
  - Functions are focused (single responsibility)
  - Functions < 50 lines (generally)
  - Complex logic extracted
  - No deep nesting (< 4 levels)

- [ ] **Comments**
  - Complex logic explained
  - "Why" not "what"
  - No commented-out code
  - Public functions documented

**Naming Examples:**
```go
// ❌ Bad naming
func GetU(id int) (*User, error)
var d time.Duration
func process(x string) error

// ✅ Good naming
func GetUserByID(userID int) (*User, error)
var sessionTimeout time.Duration
func ValidateEmail(email string) error
```

---

### 🗄️ Database & Performance

- [ ] **Database Queries**
  - Efficient queries (no N+1)
  - Appropriate indexes
  - Pagination for lists
  - Limits on results

- [ ] **Transactions**
  - Used when needed
  - Properly committed/rolled back
  - Not held too long

- [ ] **Caching**
  - Appropriate use of cache
  - Cache invalidation considered
  - TTL set correctly

- [ ] **Resource Management**
  - Connections closed
  - Contexts used with timeouts
  - No goroutine leaks
  - defer used correctly

**Example:**
```go
// ✅ Good resource management
func (r *userRepository) GetUser(ctx context.Context, id int) (*User, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    conn, err := r.db.Acquire(ctx)
    if err != nil {
        return nil, err
    }
    defer conn.Release()

    // ... query
}
```

---

### 🔄 Concurrency

- [ ] **Race Conditions**
  - No data races
  - Proper synchronization
  - Channels used correctly

- [ ] **Goroutines**
  - Not leaked
  - Proper error handling
  - Context cancellation handled

- [ ] **Mutexes**
  - Used when needed
  - Not held too long
  - Deadlock prevention

---

### 🌐 API Design

- [ ] **Proto Definitions**
  - Messages well-named
  - Fields properly numbered
  - Backwards compatible
  - Comments for fields

- [ ] **Request/Response**
  - Clear structure
  - Validation rules defined
  - Optional vs required clear
  - Consistent patterns

- [ ] **Versioning**
  - API version considered
  - Breaking changes avoided
  - Migration path if needed

---

### 🔧 Configuration

- [ ] **Config Management**
  - Config in config.yaml
  - Environment-specific values
  - Sensible defaults
  - Validation on startup

- [ ] **Feature Flags**
  - Used for gradual rollouts
  - GrowthBook integration
  - Proper naming

---

## Code Standards

### File Organization
```
services/<service-name>/
├── api/grpc/              # API layer (handlers)
├── domain/                # Business logic
├── data/
│   ├── model/            # Database models
│   └── repository/       # Data access
├── config/               # Configuration
├── errors/               # Error definitions
├── initializer/          # Service setup
└── locales/              # Translations
```

### Import Ordering
```go
import (
    // Standard library
    "context"
    "fmt"

    // External dependencies
    "github.com/jackc/pgx/v5/pgxpool"
    "google.golang.org/grpc"

    // Internal packages
    "zarinpal-platform/core/trace"
    "zarinpal-platform/services/user/domain"
)
```

### Error Handling Pattern
```go
// In domain layer
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    ctx, span := trace.GetTracer().Start(ctx, "UserService.CreateUser")
    defer span.End()

    // Validation
    if err := s.validateUser(req); err != nil {
        span.RecordError(err)
        return nil, errors.NewAppError(ctx).InvalidArgumentError()
    }

    // Business logic
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        span.RecordError(err)
        return nil, errors.NewAppError(ctx).DatabaseError()
    }

    return user, nil
}
```

---

## Common Issues

### ❌ Missing Error Checks
```go
// Bad
result, _ := someFunction()

// Good
result, err := someFunction()
if err != nil {
    return nil, err
}
```

### ❌ Not Using Context
```go
// Bad
func GetUser(id int) (*User, error)

// Good
func GetUser(ctx context.Context, id int) (*User, error)
```

### ❌ Business Logic in Handlers
```go
// Bad - in handler
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
    // validation logic
    // database queries
    // business rules
}

// Good - handler delegates to service
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
    user, err := h.service.CreateUser(ctx, req)
    if err != nil {
        return nil, err
    }
    return toProto(user), nil
}
```

### ❌ SQL Injection
```go
// Bad
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)

// Good
query := "SELECT * FROM users WHERE email = $1"
row := db.QueryRow(ctx, query, email)
```

### ❌ Not Closing Resources
```go
// Bad
file, _ := os.Open("file.txt")
// ... use file

// Good
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close()
```

---

## Review Comments Examples

### Constructive Feedback
```
❌ "This is wrong"
✅ "This might cause a race condition. Consider using a mutex here."

❌ "Why did you do it this way?"
✅ "Could we use dependency injection here instead? It would make testing easier."

❌ "This function is too long"
✅ "This function has multiple responsibilities. Could we extract the validation logic into a separate function?"
```

### Praise Good Code
```
✅ "Nice error handling here!"
✅ "Good test coverage for edge cases"
✅ "I like how this is organized"
✅ "Great use of dependency injection"
```

---

## Quick Reference

**Before Requesting Review:**
```bash
make format     # Format code (automated)
make lint       # Run linter (automated)
make test       # Run tests
make coverage   # Check coverage
```
> Git hooks will automatically check formatting and linting before commit/push.
> The author should run these commands before requesting review.

**Review Priority (Manual Review Focus):**
1. Does it work correctly?
2. Is it secure?
3. Is it tested?
4. Is it maintainable?
5. Is it performant?

**Don't Focus On (Automated):**
- ❌ Code formatting (gofmt handles this)
- ❌ Import organization (goimports handles this)
- ❌ Basic linting issues (golangci-lint catches these)

---

## Questions?

If you're unsure about any guideline, ask the team. These guidelines are meant to help, not hinder. Use your judgment, and when in doubt, discuss with the team.

**Happy reviewing! 🚀**
