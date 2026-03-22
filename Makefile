.PHONY: run

run:
	cd web-app && go run ./cmd/web

test_all:
	cd web-app && gotestsum --format testname -- ./cmd/web/

test_routes:
	cd web-app && gotestsum --format testname -- ./cmd/web/ -run Test_application_routes

test_handlers:
	cd web-app && gotestsum --format testname -- ./cmd/web/ -run "Test_application_handlers|Test_application_home|Test_application_homev2|Test_application_render_bad_template|Test_app_login"

test_midd:
	cd web-app && gotestsum --format testname -- ./cmd/web/ -run "Test_application_addIPToContext|Test_application_ipFromContext"

test_cover:
	cd web-app && go test -coverprofile=coverage.out ./cmd/web/... && go tool cover -html=coverage.out