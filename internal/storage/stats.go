package storage

import (
	"context"
	"time"
)

func (s *Storage) GetStats(ctx context.Context, start time.Time, end time.Time) (totalRequests, totalErrors int64, err error) {
	result, err := s.q.GetStats(ctx, &GetStatsParams{
		FromTime: start.Format("2006-01-02 15:04:05.000"),
		ToTime:   end.Format("2006-01-02 15:04:05.000"),
	})
	if err != nil {
		return 0, 0, err
	}

	return int64(result.Requests), int64(result.Errors), nil
}

func (s *Storage) GetUniqueUsersCount(ctx context.Context, start time.Time, end time.Time) (uniqueUsers int64, err error) {
	count, err := s.q.GetUniqueUsersCount(ctx, &GetUniqueUsersCountParams{
		FromTime: start.Format("2006-01-02 15:04:05.000"),
		ToTime:   end.Format("2006-01-02 15:04:05.000"),
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

type ProxyStats struct {
	TargetStats      map[string][]TargetStats
	TotalRequests    int64
	TotalErrors      int64
	TotalUniqueUsers int64
}

type TargetStats struct {
	Requests   int32  `json:"requests"`
	Errors     int32  `json:"errors"`
	UsersCount int32  `json:"users_count"`
	Timestamp  string `json:"timestamp"`
}

func (s *Storage) GetTargetStats(ctx context.Context, start time.Time, end time.Time, proxyID string) (*ProxyStats, error) {
	stats, err := s.q.GetTargetStats(ctx, &GetTargetStatsParams{
		ProxyID:  proxyID,
		FromTime: start.Format("2006-01-02 15:04:05.000"),
		ToTime:   end.Format("2006-01-02 15:04:05.000"),
	})
	if err != nil {
		return nil, err
	}

	targetStats := make(map[string][]TargetStats)

	var totalRequests, totalErrors, totalUniqueUsers int32

	for _, t := range stats {

		targetStats[t.TargetID] = append(targetStats[t.TargetID], TargetStats{
			Requests:   t.Requests,
			Errors:     t.Errors,
			UsersCount: t.UsersCount,
			Timestamp:  t.Timestamp.Time.Format("2006-01-02 15:04:05.000"),
		})

		totalRequests += t.Requests
		totalErrors += t.Errors
		totalUniqueUsers += t.UsersCount
	}

	return &ProxyStats{
		TargetStats:      targetStats,
		TotalRequests:    int64(totalRequests),
		TotalErrors:      int64(totalErrors),
		TotalUniqueUsers: int64(totalUniqueUsers),
	}, nil
}
