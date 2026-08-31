package httptransport

import (
	"encoding/json"
	"net/http"

	api "github.com/6ivkin/test.git/internal/transport/http/api"
)

func openAPIHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	spec, err := api.GetSwagger()
	if err != nil {
		http.Error(
			w,
			"failed to load OpenAPI specification",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(spec); err != nil {
		return
	}
}

func swaggerHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Library API</title>

    <link
        rel="stylesheet"
        href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"
    >
</head>

<body>
    <div id="swagger-ui"></div>

    <script
        src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js">
    </script>

    <script>
        SwaggerUIBundle({
            url: "/openapi.json",
            dom_id: "#swagger-ui"
        });
    </script>
</body>
</html>`))
}
