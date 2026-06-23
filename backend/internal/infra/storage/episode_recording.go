package storage

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	s3client "github.com/airoa-org/yubi-app/backend/internal/s3"
)

type episodeRecording struct {
	s3 *s3client.Client
}

func NewEpisodeRecording(s3 *s3client.Client) *episodeRecording {
	return &episodeRecording{s3: s3}
}

var (
	pathValuePattern  = regexp.MustCompile(`[^\w.-]+`)
	multiDashPattern  = regexp.MustCompile(`-{2,}`)
	whitespacePattern = regexp.MustCompile(`\s+`)
)

// normalizePathValue normalizes a string for use in S3 path segments.
func normalizePathValue(value string) string {
	s := strings.TrimSpace(norm.NFKC.String(value))
	s = strings.ToLower(s)
	if s == "" {
		return "unknown"
	}
	s = strings.ReplaceAll(s, "/", "-")
	s = whitespacePattern.ReplaceAllString(s, "-")
	s = pathValuePattern.ReplaceAllString(s, "-")
	s = multiDashPattern.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-.")
	if s == "" {
		return "unknown"
	}
	return s
}

// canonicalPrefix builds the S3 prefix that matches the workflow's canonical path format:
// preview/org={org}/site={site}/location={loc}/date={date}/robot_type={type}/robot_id={id}/ts={ts}/uuid={uuid}
func canonicalPrefix(p model.EpisodePreviewPath) string {
	ts := p.StartedAt.UTC()
	return fmt.Sprintf("preview/org=%s/site=%s/location=%s/date=%s/robot_type=%s/robot_id=%s/ts=%s/uuid=%s",
		normalizePathValue(p.Organization),
		normalizePathValue(p.Site),
		normalizePathValue(p.Location),
		ts.Format("2006-01-02"),
		normalizePathValue(p.RobotType),
		normalizePathValue(p.RobotID),
		ts.Format("2006-01-02T15:04:05.000Z"),
		strings.ToLower(p.UUID),
	)
}

func (g *episodeRecording) GetRecordingURLs(ctx context.Context, path model.EpisodePreviewPath) (map[string]string, error) {
	prefix := canonicalPrefix(path) + "/videos/chunk-000/"

	keys, err := g.s3.ListObjectKeys(ctx, prefix)
	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeInternal, "failed to list recording objects"))
	}

	result := make(map[string]string)
	for _, key := range keys {
		if !strings.HasSuffix(key, "/episode_000000.mp4") {
			continue
		}
		feature, ok := extractFeatureFromVideoKey(key, prefix)
		if !ok {
			continue
		}

		url, err := g.s3.GetPresignedURL(ctx, key)
		if err != nil {
			return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeInternal, "failed to generate presigned URL"))
		}
		result[feature] = url
	}

	return result, nil
}

func extractFeatureFromVideoKey(key, chunkPrefix string) (string, bool) {
	rel := strings.TrimPrefix(key, chunkPrefix)

	featureEnd := strings.Index(rel, "/")
	if featureEnd < 0 {
		return "", false
	}
	return rel[:featureEnd], true
}

type statsLine struct {
	EpisodeIndex int                        `json:"episode_index"`
	Stats        map[string]rawFeatureStats `json:"stats"`
}

type rawFeatureStats struct {
	Min   json.RawMessage `json:"min"`
	Max   json.RawMessage `json:"max"`
	Mean  json.RawMessage `json:"mean"`
	Std   json.RawMessage `json:"std"`
	Count json.RawMessage `json:"count"`
}

func (g *episodeRecording) GetStats(ctx context.Context, path model.EpisodePreviewPath) (model.EpisodeRecordingStats, error) {
	key := canonicalPrefix(path) + "/meta/episodes_stats.jsonl"

	data, err := g.s3.GetObjectBody(ctx, key)
	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeInternal, "failed to get stats file"))
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 10*1024*1024), 10*1024*1024)
	if !scanner.Scan() {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeInternal, "stats file is empty"))
	}

	var line statsLine
	if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeInternal, "failed to parse stats file"))
	}

	result := make(model.EpisodeRecordingStats, len(line.Stats))
	for feature, raw := range line.Stats {
		result[feature] = model.EpisodeFeatureStats{
			Min:   flattenFloats(raw.Min),
			Max:   flattenFloats(raw.Max),
			Mean:  flattenFloats(raw.Mean),
			Std:   flattenFloats(raw.Std),
			Count: extractCount(raw.Count),
		}
	}

	return result, nil
}

func flattenFloats(raw json.RawMessage) []float64 {
	if len(raw) == 0 {
		return nil
	}
	var scalar float64
	if err := json.Unmarshal(raw, &scalar); err == nil {
		return []float64{scalar}
	}
	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err != nil {
		return nil
	}
	var result []float64
	for _, elem := range arr {
		result = append(result, flattenFloats(elem)...)
	}
	return result
}

func extractCount(raw json.RawMessage) int {
	vals := flattenFloats(raw)
	if len(vals) == 0 {
		return 0
	}
	return int(vals[0])
}
