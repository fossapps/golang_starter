//go:generate mockgen -destination=./mock_rate_limiter.go -package=mocks github.com/fossapps/starter/adapters RateLimiter
//go:generate mockgen -destination=./mock_db.go -package=mocks github.com/fossapps/starter/db DB
//go:generate mockgen -destination=./mock_devices.go -package=mocks github.com/fossapps/starter/db DeviceManager
//go:generate mockgen -destination=./mock_migration.go -package=mocks github.com/fossapps/starter/db MigrationManager
//go:generate mockgen -destination=./mock_permissions.go -package=mocks github.com/fossapps/starter/db PermissionManager
//go:generate mockgen -destination=./mock_refresh_tokens.go -package=mocks github.com/fossapps/starter/db RefreshTokenManager
//go:generate mockgen -destination=./mock_users.go -package=mocks github.com/fossapps/starter/db UserManager
//go:generate mockgen -destination=./mock_redis.go -package=mocks github.com/fossapps/starter/helpers RedisClient
//go:generate mockgen -destination=./mock_migration_implementation.go -package=mocks github.com/fossapps/starter/migrations Migration
//go:generate mockgen -destination=./mock_logger.go -package=mocks github.com/fossapps/starter Logger
//go:generate mockgen -destination=./mock_request_helper.go -package=mocks github.com/fossapps/starter RequestHelper

package mocks
