package proxy

import (
	"net/http"
)

func (p *Proxy) getTargetFromCookie(r *http.Request) *Target {
	cookie, err := r.Cookie(p.cookieName)
	if err != nil {
		return nil
	}

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, target := range p.Targets {
		if target.IsActive && target.URL == cookie.Value {
			return &target
		}
	}
	return nil
}

func (p *Proxy) setCookies(w http.ResponseWriter, info *RedirectInfo) {
	// Set ab_target cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "ab_target_",
		Value:    p.Targets[0].URL,
		Path:     "/",
		MaxAge:   3600 * 24, // 24 hours
		HttpOnly: true,
	})

	// Set RID cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "rid",
		Value:    info.RID,
		Path:     "/",
		MaxAge:   3600 * 24 * 30, // 30 days
		HttpOnly: true,
	})

	// Set RRID cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "rrid",
		Value:    info.RRID,
		Path:     "/",
		MaxAge:   3600 * 24, // 24 hours
		HttpOnly: true,
	})

	// Set RUID cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "ruid",
		Value:    info.RUID,
		Path:     "/",
		MaxAge:   3600 * 24 * 365, // 1 year
		HttpOnly: true,
	})
}
