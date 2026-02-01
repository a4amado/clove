package headers

import "fmt"

const (
	// Content Negotiation
	Accept         = "Accept"
	AcceptEncoding = "Accept-Encoding"
	AcceptLanguage = "Accept-Language"

	// Authentication
	Authorization      = "Authorization"
	ProxyAuthorization = "Proxy-Authorization"

	// Caching
	CacheControl      = "Cache-Control"
	IfMatch           = "If-Match"
	IfModifiedSince   = "If-Modified-Since"
	IfNoneMatch       = "If-None-Match"
	IfUnmodifiedSince = "If-Unmodified-Since"

	// Client Hints
	SecCHUA         = "Sec-CH-UA"
	SecCHUAMobile   = "Sec-CH-UA-Mobile"
	SecCHUAPlatform = "Sec-CH-UA-Platform"

	// Conditionals
	IfRange = "If-Range"

	// Connection Management
	Connection = "Connection"
	KeepAlive  = "Keep-Alive"

	// Content Type
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
	ContentEncoding = "Content-Encoding"
	ContentLanguage = "Content-Language"
	ContentLocation = "Content-Location"

	// Cookies
	Cookie = "Cookie"

	// CORS
	Origin                      = "Origin"
	AccessControlRequestMethod  = "Access-Control-Request-Method"
	AccessControlRequestHeaders = "Access-Control-Request-Headers"

	// Range Requests
	Range = "Range"

	// Referrer
	Referer        = "Referer" // Note: Misspelled in HTTP spec
	ReferrerPolicy = "Referrer-Policy"

	// Request Context
	From      = "From"
	Host      = "Host"
	UserAgent = "User-Agent"

	// Security
	SecFetchSite = "Sec-Fetch-Site"
	SecFetchMode = "Sec-Fetch-Mode"
	SecFetchUser = "Sec-Fetch-User"
	SecFetchDest = "Sec-Fetch-Dest"

	// Other
	Upgrade         = "Upgrade"
	Via             = "Via"
	XForwardedFor   = "X-Forwarded-For"
	XForwardedHost  = "X-Forwarded-Host"
	XForwardedProto = "X-Forwarded-Proto"
)

// Response headers - Standard HTTP response headers based on MDN Web Docs
const (
	// Authentication
	WWWAuthenticate   = "WWW-Authenticate"
	ProxyAuthenticate = "Proxy-Authenticate"

	// Caching
	Age           = "Age"
	ClearSiteData = "Clear-Site-Data"
	ETag          = "ETag"
	Expires       = "Expires"
	LastModified  = "Last-Modified"

	// Content Negotiation
	AcceptRanges = "Accept-Ranges"
	Vary         = "Vary"

	// Content
	ContentDisposition = "Content-Disposition"
	ContentRange       = "Content-Range"

	// Cookies
	SetCookie = "Set-Cookie"

	// CORS
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	AccessControlMaxAge           = "Access-Control-Max-Age"

	// Redirects
	Location = "Location"

	// Security
	ContentSecurityPolicy           = "Content-Security-Policy"
	ContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	CrossOriginEmbedderPolicy       = "Cross-Origin-Embedder-Policy"
	CrossOriginOpenerPolicy         = "Cross-Origin-Opener-Policy"
	CrossOriginResourcePolicy       = "Cross-Origin-Resource-Policy"
	PermissionsPolicy               = "Permissions-Policy"
	StrictTransportSecurity         = "Strict-Transport-Security"
	XContentTypeOptions             = "X-Content-Type-Options"
	XFrameOptions                   = "X-Frame-Options"
	XXSSProtection                  = "X-XSS-Protection"

	// Server Info
	Server = "Server"

	// Other
	Refresh    = "Refresh"
	RetryAfter = "Retry-After"
)

// ContentType values - Common MIME types
const (
	ContentTypeJSON           = "application/json"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	ContentTypeMultipart      = "multipart/form-data"
	ContentTypeTextPlain      = "text/plain"
	ContentTypeTextHTML       = "text/html"
	ContentTypeTextCSS        = "text/css"
	ContentTypeTextJS         = "text/javascript"
	ContentTypeXML            = "application/xml"
	ContentTypePDF            = "application/pdf"
	ContentTypeOctetStream    = "application/octet-stream"
)

// Cache-Control directives
const (
	CacheControlNoCache        = "no-cache"
	CacheControlNoStore        = "no-store"
	CacheControlMustRevalidate = "must-revalidate"
	CacheControlPublic         = "public"
	CacheControlPrivate        = "private"
)

// MaxAge returns a Cache-Control `max-age` directive using the provided seconds value.
// The seconds value is used as the directive's duration in seconds.
func MaxAge(seconds int) string {
	return fmt.Sprintf("max-age=%d", seconds)
}