.PHONY: tailwind-watch
tailwind-watch:
	./tailwindcss -i ./public/css/index.css -o ./public/css/tw.css --watch

.PHONY: tailwind-build
tailwind-build:
	./tailwindcss -i ./public/css/index.css -o ./public/css/tw.min.css --minify

.PHONY: dev-no-reload
dev-no-reload:
	air & \
	./tailwindcss \
	  -i 'public/css/index.css' \
	  -o 'public/css/tw.css' \
	  --watch


.PHONY: dev
dev:
	air & \
	./tailwindcss \
	  -i 'public/css/index.css' \
	  -o 'public/css/tw.css' \
	  --watch & \
	browser-sync start \
	  --files 'views/**/*.html, public/**/*.css' \
	  --port 3000 \
	  --proxy 'localhost:5050' \
	  --reload-delay 1000 \
	  --middleware 'function(req, res, next) { \
		res.setHeader("Cache-Control", "no-cache, no-store, must-revalidate"); \
		return next(); \
	  }'

.PHONY: build
build:
	make tailwind-build
	go build -ldflags "-X main.Environment=production" -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/*.go
