package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/ab-testing-service/internal/models"
)

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	redirectInfo, err := p.getOrCreateRedirectInfo(r)
	if err != nil {
		http.Error(w, "Failed to process redirect info", http.StatusInternalServerError)
		p.stats.IncrementErrors(p.ID, "")
		return
	}

	// Check if this is a redirect request from another proxy
	if r.Header.Get("X-Internal-Redirect") == "true" {
		// Remove the header to prevent redirect loops
		r.Header.Del("X-Internal-Redirect")
		target, err := p.selectTarget(r)
		if err != nil {
			http.Error(w, "Error selecting target", http.StatusInternalServerError)
			p.stats.IncrementErrors(p.ID, redirectInfo.RUID)
			return
		}
		targetURL, err := url.Parse(target.URL)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			p.stats.IncrementErrors(target.ID, redirectInfo.RUID)
			http.Error(w, "Error forwarding request", http.StatusBadGateway)
		}
		proxy.ServeHTTP(w, r)
		return
	}

	target, err := p.selectTarget(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to select target: %s", err), http.StatusInternalServerError)
		p.stats.IncrementErrors(p.ID, redirectInfo.RUID)
		return
	}
	log.Printf("Selected target: %s", target.URL)

	p.setCookies(w, redirectInfo)

	// Get user identifier (prefer X-User-ID header, fallback to IP)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = redirectInfo.RUID
	}

	// Track request with user ID
	p.stats.IncrementRequestsWithUser(target.ID, userID)

	if p.Mode == models.ProxyModeRedirect || p.Mode == models.ProxyModePath {
		// Check if the target URL has a different host
		//targetURL := p.appendRedirectParams(target.URL, redirectInfo)
		parsedTarget, err := url.Parse(target.URL)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}

		if parsedTarget.Scheme == "" {
			parsedTarget.Scheme = "https"
		}

		// If target host is different from current host, do external redirect
		if parsedTarget.Host != r.Host {
			log.Printf("Redirecting to %s", parsedTarget.String())
			http.Redirect(w, r, parsedTarget.String(), http.StatusMovedPermanently)
			return
		}

		// For same host, do internal redirect
		r.URL.Path = parsedTarget.Path
		r.URL.RawQuery = parsedTarget.RawQuery
		r.Header.Set("X-Internal-Redirect", "true")
		http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
		return
	}

	// For reverse proxy mode
	targetURL, err := url.Parse(target.URL)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		p.stats.IncrementErrors(target.ID, redirectInfo.RUID)
		return
	}

	// Add redirect info to request headers
	r.Header.Set("X-Redirect-ID", redirectInfo.RID)
	r.Header.Set("X-Redirect-Request-ID", redirectInfo.RRID)
	r.Header.Set("X-Redirect-User-ID", redirectInfo.RUID)
	r.Header.Set("X-Redirect-Query-Params", targetURL.RawQuery)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		p.stats.IncrementErrors(target.ID, redirectInfo.RUID)
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
	}

	proxy.ServeHTTP(w, r)

	duration := time.Since(start).Seconds()
	p.metrics.LatencyHistogram.WithLabelValues(target.URL).Observe(duration)
	p.metrics.RequestsTotal.WithLabelValues(target.URL).Inc()
}

func (p *Proxy) getTargetByCondition(r *http.Request) *Target {
	var value string
	switch p.Config.Condition.Type {
	case models.ConditionTypeHeader:
		value = r.Header.Get(p.Config.Condition.ParamName)
	case models.ConditionTypeQuery:
		value = r.URL.Query().Get(p.Config.Condition.ParamName) // fixme
	case models.ConditionTypeCookie:
		cookie, err := r.Cookie(p.Config.Condition.ParamName)
		if err == nil {
			value = cookie.Value
		}
	case models.ConditionTypeUserAgent:
		ua := r.Header.Get("User-Agent")
		switch p.Config.Condition.ParamName {
		case "platform":
			value = detectPlatform(ua)
		case "browser":
			value = detectBrowser(ua)
		}
	case models.ConditionTypeLanguage:
		value = parseAcceptLanguage(r.Header.Get("Accept-Language"))
	default:
		return p.getTargetById(p.Config.Condition.Default)
	}

	// Check if the value matches any of the specified values
	if targetID, ok := p.Config.Condition.Values[value]; ok {
		if target := p.getTargetById(targetID); target != nil {
			return target
		}
	}

	// Fall back to default target
	return p.getTargetById(p.Config.Condition.Default)
}

func (p *Proxy) getTargetById(id string) *Target {
	for _, target := range p.Targets {
		if target.ID == id && target.IsActive {
			return &target
		}
	}
	return nil
}

// detectPlatform detects the platform (mobile/desktop) from User-Agent
func detectPlatform(ua string) string {
	ua = strings.ToLower(ua)
	mobileKeywords := []string{
		"mobile", "android", "iphone", "ipad", "ipod",
		"windows phone", "blackberry", "opera mini",
	}

	for _, keyword := range mobileKeywords {
		if strings.Contains(ua, keyword) {
			return "mobile"
		}
	}
	return "desktop"
}

// detectBrowser detects the browser from User-Agent
func detectBrowser(ua string) string {
	ua = strings.ToLower(ua)
	browsers := map[string]string{
		"firefox":   "firefox",
		"chrome":    "chrome",
		"safari":    "safari",
		"edge":      "edge",
		"opera":     "opera",
		"msie":      "ie",
		"trident/7": "ie",
	}

	for keyword, browser := range browsers {
		if strings.Contains(ua, keyword) {
			return browser
		}
	}
	return "other"
}

// parseAcceptLanguage gets the preferred language from Accept-Language header
func parseAcceptLanguage(acceptLang string) string {
	if acceptLang == "" {
		return ""
	}

	// Split the Accept-Language header into parts
	parts := strings.Split(acceptLang, ",")
	if len(parts) == 0 {
		return ""
	}

	// Get the first (most preferred) language
	lang := strings.Split(parts[0], ";")[0]
	return strings.ToLower(strings.TrimSpace(lang))
}
