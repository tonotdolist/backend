cd ../

mockgen -source ./internal/repository/user.go -destination ./mocks/repository/user.go
mockgen -source ./internal/repository/activity.go -destination ./mocks/repository/activity.go
mockgen -source ./internal/repository/session.go -destination ./mocks/repository/session.go
mockgen -source ./pkg/clock/clock.go -destination ./mocks/clock/clock.go
mockgen -source ./internal/util/id.go -destination ./mocks/id/id.go

pause