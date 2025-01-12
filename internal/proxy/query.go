package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type RedirectInfo struct {
	RID   string // Redirect ID (same for all users within proxy)
	RRID  string // Redirect Request ID (unique per click)
	RUID  string // Redirect User ID (unique per user)
	Query url.Values
}

func (p *Proxy) appendRedirectParams(targetURL string, info *RedirectInfo) string {
	u, err := url.Parse(targetURL)
	if err != nil {
		return targetURL
	}

	// Get existing query parameters
	query := u.Query()

	// Add redirect info parameters
	query.Set("rid", info.RID)
	query.Set("rrid", info.RRID)
	query.Set("ruid", info.RUID)

	// Add all original query parameters
	for key, values := range info.Query {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	// Set the updated query string
	u.RawQuery = query.Encode()

	return u.String()
}

func (p *Proxy) getOrCreateRedirectInfo(r *http.Request) (*RedirectInfo, error) {
	// Get or generate RUID from cookie
	ruidCookie, err := r.Cookie("ruid")
	var ruid string
	if errors.Is(err, http.ErrNoCookie) || ruidCookie == nil {
		ruid = uuid.New().String()
	} else {
		ruid = ruidCookie.Value
	}

	// Get RID from proxy ID
	rid := fmt.Sprintf("rid_%s", p.ID)

	// Generate new RRID for this request
	rrid := uuid.New().String()

	// Get original query parameters
	query := r.URL.Query()

	return &RedirectInfo{
		RID:   rid,
		RRID:  rrid,
		RUID:  ruid,
		Query: query,
	}, nil
}
