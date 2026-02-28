package personalization

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewRepository(db *sql.DB, log *zap.Logger) *Repository {
	return &Repository{db: db, log: log}
}

// ─── Behavior Signals ───────────────────────────────────────────────────────

func (r *Repository) RecordSignal(signal *BehaviorSignal) error {
	contextJSON, err := json.Marshal(signal.Context)
	if err != nil {
		contextJSON = []byte("{}")
	}

	_, err = r.db.Exec(`
		INSERT INTO behavior_signals (user_id, session_id, signal_type, value, context, topic, subject)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, signal.UserID, signal.SessionID, signal.SignalType, signal.Value, contextJSON, signal.Topic, signal.Subject)
	return err
}

func (r *Repository) GetUserSignals(userID uuid.UUID, limit int) ([]BehaviorSignal, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, session_id, signal_type, value, context, topic, subject, created_at
		FROM behavior_signals
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signals []BehaviorSignal
	for rows.Next() {
		var s BehaviorSignal
		var contextJSON []byte
		var sessionID sql.NullString
		var topic, subject sql.NullString

		if err := rows.Scan(&s.ID, &s.UserID, &sessionID, &s.SignalType, &s.Value, &contextJSON, &topic, &subject, &s.CreatedAt); err != nil {
			return nil, err
		}

		if sessionID.Valid {
			sid, _ := uuid.Parse(sessionID.String)
			s.SessionID = &sid
		}
		s.Topic = topic.String
		s.Subject = subject.String
		json.Unmarshal(contextJSON, &s.Context)
		signals = append(signals, s)
	}
	return signals, rows.Err()
}

func (r *Repository) GetSignalsByType(userID uuid.UUID, signalType SignalType, since time.Time) ([]BehaviorSignal, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, session_id, signal_type, value, context, topic, subject, created_at
		FROM behavior_signals
		WHERE user_id = $1 AND signal_type = $2 AND created_at >= $3
		ORDER BY created_at DESC
	`, userID, signalType, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signals []BehaviorSignal
	for rows.Next() {
		var s BehaviorSignal
		var contextJSON []byte
		var sessionID sql.NullString
		var topic, subject sql.NullString

		if err := rows.Scan(&s.ID, &s.UserID, &sessionID, &s.SignalType, &s.Value, &contextJSON, &topic, &subject, &s.CreatedAt); err != nil {
			return nil, err
		}

		if sessionID.Valid {
			sid, _ := uuid.Parse(sessionID.String)
			s.SessionID = &sid
		}
		s.Topic = topic.String
		s.Subject = subject.String
		json.Unmarshal(contextJSON, &s.Context)
		signals = append(signals, s)
	}
	return signals, rows.Err()
}

func (r *Repository) GetSignalCounts(userID uuid.UUID, since time.Time) (map[SignalType]int, error) {
	rows, err := r.db.Query(`
		SELECT signal_type, COUNT(*) as count
		FROM behavior_signals
		WHERE user_id = $1 AND created_at >= $2
		GROUP BY signal_type
	`, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[SignalType]int)
	for rows.Next() {
		var signalType SignalType
		var count int
		if err := rows.Scan(&signalType, &count); err != nil {
			return nil, err
		}
		counts[signalType] = count
	}
	return counts, rows.Err()
}

// ─── Learning Style Profile ─────────────────────────────────────────────────

func (r *Repository) GetLearningStyleProfile(userID uuid.UUID) (*LearningStyleProfile, error) {
	var p LearningStyleProfile
	err := r.db.QueryRow(`
		SELECT id, user_id, primary_style, visual_score, auditory_score, reading_score, 
		       kinesthetic_score, confidence, sample_size, last_calculated_at, created_at, updated_at
		FROM learning_style_profiles
		WHERE user_id = $1
	`, userID).Scan(
		&p.ID, &p.UserID, &p.PrimaryStyle, &p.VisualScore, &p.AuditoryScore, &p.ReadingScore,
		&p.KinestheticScore, &p.Confidence, &p.SampleSize, &p.LastCalculatedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) UpsertLearningStyleProfile(p *LearningStyleProfile) error {
	_, err := r.db.Exec(`
		INSERT INTO learning_style_profiles 
			(user_id, primary_style, visual_score, auditory_score, reading_score, 
			 kinesthetic_score, confidence, sample_size, last_calculated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id) DO UPDATE SET
			primary_style = $2,
			visual_score = $3,
			auditory_score = $4,
			reading_score = $5,
			kinesthetic_score = $6,
			confidence = $7,
			sample_size = $8,
			last_calculated_at = $9,
			updated_at = NOW()
	`, p.UserID, p.PrimaryStyle, p.VisualScore, p.AuditoryScore, p.ReadingScore,
		p.KinestheticScore, p.Confidence, p.SampleSize, p.LastCalculatedAt)
	return err
}

// ─── Topic Engagement Tracking ──────────────────────────────────────────────

func (r *Repository) GetTopicEngagement(userID uuid.UUID, limit int) ([]struct {
	Topic      string
	Subject    string
	Count      int
	AvgValue   float64
	LastAccess time.Time
}, error) {
	rows, err := r.db.Query(`
		SELECT topic, subject, COUNT(*) as count, AVG(value) as avg_value, MAX(created_at) as last_access
		FROM behavior_signals
		WHERE user_id = $1 AND topic != ''
		GROUP BY topic, subject
		ORDER BY count DESC, last_access DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Topic      string
		Subject    string
		Count      int
		AvgValue   float64
		LastAccess time.Time
	}
	for rows.Next() {
		var item struct {
			Topic      string
			Subject    string
			Count      int
			AvgValue   float64
			LastAccess time.Time
		}
		if err := rows.Scan(&item.Topic, &item.Subject, &item.Count, &item.AvgValue, &item.LastAccess); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, rows.Err()
}

// ─── Interest Tags ──────────────────────────────────────────────────────────

func (r *Repository) GetUserInterests(userID uuid.UUID) ([]InterestWeight, error) {
	rows, err := r.db.Query(`
		SELECT tag, category, weight
		FROM interest_tags
		WHERE user_id = $1
		ORDER BY weight DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interests []InterestWeight
	for rows.Next() {
		var i InterestWeight
		if err := rows.Scan(&i.Tag, &i.Category, &i.Weight); err != nil {
			return nil, err
		}
		interests = append(interests, i)
	}
	return interests, rows.Err()
}

func (r *Repository) UpdateInterestWeight(userID uuid.UUID, tag, category string, weight float64) error {
	_, err := r.db.Exec(`
		INSERT INTO interest_tags (user_id, tag, category, weight)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, tag) DO UPDATE SET
			weight = interest_tags.weight * 0.9 + $4 * 0.1,
			category = $3
	`, userID, tag, category, weight)
	return err
}

// ─── Engagement Metrics ─────────────────────────────────────────────────────

func (r *Repository) GetEngagementMetrics(userID uuid.UUID) (*EngagementMetrics, error) {
	var m EngagementMetrics
	m.UserID = userID

	err := r.db.QueryRow(`
		SELECT 
			COUNT(*) as total_sessions,
			COALESCE(AVG(duration_sec), 0) as avg_duration,
			COALESCE(AVG(
				(SELECT COUNT(*) FROM ai_interaction_history WHERE session_id = ls.id)
			), 0) as avg_questions
		FROM learning_sessions ls
		WHERE user_id = $1
	`, userID).Scan(&m.TotalSessions, &m.AvgSessionDuration, &m.AvgQuestionsPerSess)

	if err != nil {
		return nil, err
	}

	// Calculate completion rate
	r.db.QueryRow(`
		SELECT COALESCE(
			CAST(COUNT(*) FILTER (WHERE status = 'completed') AS FLOAT) / NULLIF(COUNT(*), 0),
			0
		)
		FROM learning_sessions WHERE user_id = $1
	`, userID).Scan(&m.CompletionRate)

	// Determine engagement level
	if m.TotalSessions >= 20 && m.AvgSessionDuration > 600 {
		m.EngagementLevel = "high"
	} else if m.TotalSessions >= 5 && m.AvgSessionDuration > 300 {
		m.EngagementLevel = "medium"
	} else {
		m.EngagementLevel = "low"
	}

	return &m, nil
}
