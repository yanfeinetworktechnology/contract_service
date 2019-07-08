package common

// application version
const VERSION = "1.0.0-alpha"

// Cross-sites resource sharing settings
var CORS_ALLOW_ORIGINS = []string{
	"http://*",
	"https://*",
}

var CORS_ALLOW_DEBUG_ORIGINS = []string{
	"http://*",
	"https://*",
}

var CORS_ALLOW_HEADERS = []string{
	"Origin",
	"Content-Length",
	"Content-Type",
	"token",
	"X-CSRF-TOKEN",
	"withCredentials",
	"access-token",
}

var CORS_ALLOW_METHODS = []string{
	"GET",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
	"HEAD",
}

var CORS_EXPOSE_HEADERS = []string{
	"X-CSRF-TOKEN",
	"token",
	"access-token",
}
