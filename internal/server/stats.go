package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ab-testing-service/internal/storage"
)

type StatsResponse struct {
	TotalRequests int64                            `json:"total_requests"`
	TotalErrors   int64                            `json:"total_errors"`
	UniqueUsers   int64                            `json:"unique_users"`
	TargetStats   map[string][]storage.TargetStats `json:"target_stats,omitempty"`
	StartTime     time.Time                        `json:"start_time"`
	EndTime       time.Time                        `json:"end_time"`
}

func (s *Server) getStats(c *gin.Context) {
	// Parse time range from query parameters
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	var start, end time.Time
	var err error

	if startTime != "" {
		start, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format"})
			return
		}
	} else {
		start = time.Now().AddDate(0, 0, -7) // Default to last 7 days
	}

	if endTime != "" {
		end, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format"})
			return
		}
	} else {
		end = time.Now()
	}

	// Query overall stats
	totalRequests, totalErrors, err := s.storage.GetStats(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Query unique users count
	uniqueUsers, err := s.storage.GetUniqueUsersCount(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, StatsResponse{
		TotalRequests: totalRequests,
		TotalErrors:   totalErrors,
		UniqueUsers:   uniqueUsers,
		StartTime:     start,
		EndTime:       end,
	})
}

func (s *Server) getProxyStats(c *gin.Context) {
	proxyID := c.Param("proxy_id")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	var start, end time.Time
	var err error

	if startTime != "" {
		start, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"parse time error": "invalid start_time format"})
			return
		}
	} else {
		start = time.Now().AddDate(0, 0, -7) // Default to last 7 days
	}

	if endTime != "" {
		end, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"parse time error": "invalid end_time format"})
			return
		}
	} else {
		end = time.Now()
	}

	// Query target stats
	proxyStats, err := s.storage.GetTargetStats(c.Request.Context(), start, end, proxyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"query target stats error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, StatsResponse{
		TotalRequests: proxyStats.TotalRequests,
		TotalErrors:   proxyStats.TotalErrors,
		UniqueUsers:   proxyStats.TotalUniqueUsers,
		TargetStats:   proxyStats.TargetStats,
		StartTime:     start,
		EndTime:       end,
	})
}
