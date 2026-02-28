package personalization

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// InterestClassifier analyzes user behavior and explicit preferences
// to build a weighted interest profile for personalized analogies.
type InterestClassifier struct {
	repo *Repository
	log  *zap.Logger
}

func NewInterestClassifier(repo *Repository, log *zap.Logger) *InterestClassifier {
	return &InterestClassifier{repo: repo, log: log}
}

// InterestCategory represents a high-level grouping of interests.
type InterestCategory string

const (
	CategoryAnime      InterestCategory = "anime_manga"
	CategoryGaming     InterestCategory = "gaming"
	CategorySports     InterestCategory = "sports"
	CategoryMusic      InterestCategory = "music"
	CategoryMovies     InterestCategory = "movies_tv"
	CategoryScience    InterestCategory = "science"
	CategoryTechnology InterestCategory = "technology"
	CategoryArt        InterestCategory = "art_design"
	CategoryBooks      InterestCategory = "books_literature"
	CategoryFood       InterestCategory = "food_cooking"
	CategoryTravel     InterestCategory = "travel"
	CategoryFitness    InterestCategory = "fitness"
	CategoryBusiness   InterestCategory = "business"
	CategoryNature     InterestCategory = "nature"
)

// interestKeywords maps keywords to categories for classification.
var interestKeywords = map[string]InterestCategory{
	// Anime & Manga
	"anime": CategoryAnime, "manga": CategoryAnime, "naruto": CategoryAnime,
	"one piece": CategoryAnime, "attack on titan": CategoryAnime, "dragon ball": CategoryAnime,
	"studio ghibli": CategoryAnime, "demon slayer": CategoryAnime, "jujutsu kaisen": CategoryAnime,
	"my hero academia": CategoryAnime, "cosplay": CategoryAnime,

	// Gaming
	"gaming": CategoryGaming, "video games": CategoryGaming, "esports": CategoryGaming,
	"rpg": CategoryGaming, "fps": CategoryGaming, "mmorpg": CategoryGaming,
	"playstation": CategoryGaming, "xbox": CategoryGaming, "nintendo": CategoryGaming,
	"minecraft": CategoryGaming, "fortnite": CategoryGaming, "league of legends": CategoryGaming,
	"valorant": CategoryGaming, "zelda": CategoryGaming, "pokemon": CategoryGaming,

	// Sports
	"basketball": CategorySports, "nba": CategorySports, "football": CategorySports,
	"soccer": CategorySports, "tennis": CategorySports, "baseball": CategorySports,
	"hockey": CategorySports, "golf": CategorySports, "swimming": CategorySports,
	"running": CategorySports, "cycling": CategorySports, "martial arts": CategorySports,
	"boxing": CategorySports, "mma": CategorySports, "olympics": CategorySports,

	// Music
	"music": CategoryMusic, "hip hop": CategoryMusic, "rock": CategoryMusic,
	"pop": CategoryMusic, "jazz": CategoryMusic, "classical": CategoryMusic,
	"electronic": CategoryMusic, "k-pop": CategoryMusic, "guitar": CategoryMusic,
	"piano": CategoryMusic, "drums": CategoryMusic, "concerts": CategoryMusic,

	// Movies & TV
	"movies": CategoryMovies, "film": CategoryMovies, "cinema": CategoryMovies,
	"tv shows": CategoryMovies, "netflix": CategoryMovies, "marvel": CategoryMovies,
	"star wars": CategoryMovies, "dc comics": CategoryMovies, "horror": CategoryMovies,
	"comedy": CategoryMovies, "drama": CategoryMovies, "documentaries": CategoryMovies,

	// Science
	"physics": CategoryScience, "chemistry": CategoryScience, "biology": CategoryScience,
	"astronomy": CategoryScience, "space": CategoryScience, "nasa": CategoryScience,
	"experiments": CategoryScience, "research": CategoryScience,

	// Technology
	"programming": CategoryTechnology, "coding": CategoryTechnology, "ai": CategoryTechnology,
	"machine learning": CategoryTechnology, "robotics": CategoryTechnology,
	"gadgets": CategoryTechnology, "smartphones": CategoryTechnology, "computers": CategoryTechnology,
	"startups": CategoryTechnology, "crypto": CategoryTechnology,

	// Art & Design
	"art": CategoryArt, "painting": CategoryArt, "drawing": CategoryArt,
	"photography": CategoryArt, "graphic design": CategoryArt, "architecture": CategoryArt,
	"sculpture": CategoryArt, "digital art": CategoryArt,

	// Books & Literature
	"books": CategoryBooks, "reading": CategoryBooks, "novels": CategoryBooks,
	"fantasy": CategoryBooks, "sci-fi": CategoryBooks, "history": CategoryBooks,
	"philosophy": CategoryBooks, "poetry": CategoryBooks,

	// Food & Cooking
	"cooking": CategoryFood, "baking": CategoryFood, "cuisine": CategoryFood,
	"restaurants": CategoryFood, "recipes": CategoryFood, "chef": CategoryFood,

	// Nature
	"nature": CategoryNature, "animals": CategoryNature, "wildlife": CategoryNature,
	"hiking": CategoryNature, "camping": CategoryNature, "ocean": CategoryNature,
	"mountains": CategoryNature, "gardening": CategoryNature,
}

// Classify builds an interest profile from user data and behavior.
func (c *InterestClassifier) Classify(userID uuid.UUID) (*InterestProfile, error) {
	// Get explicit interests from database
	explicitInterests, err := c.repo.GetUserInterests(userID)
	if err != nil {
		c.log.Warn("failed to get explicit interests", zap.Error(err))
		explicitInterests = []InterestWeight{}
	}

	// Get behavioral signals indicating interest
	since := time.Now().AddDate(0, 0, -60)
	interestSignals, err := c.repo.GetSignalsByType(userID, SignalInterestIndication, since)
	if err != nil {
		c.log.Warn("failed to get interest signals", zap.Error(err))
	}

	satisfactionSignals, err := c.repo.GetSignalsByType(userID, SignalAnalogySatisfaction, since)
	if err != nil {
		c.log.Warn("failed to get satisfaction signals", zap.Error(err))
	}

	// Aggregate weights by category
	categoryScores := make(map[InterestCategory]float64)
	tagScores := make(map[string]float64)

	// Process explicit interests (highest weight)
	for _, interest := range explicitInterests {
		tag := strings.ToLower(interest.Tag)
		category := c.categorizeInterest(tag)

		tagScores[tag] += interest.Weight * 2.0 // Explicit interests get 2x weight
		categoryScores[category] += interest.Weight * 2.0
	}

	// Process behavioral interest signals
	for _, signal := range interestSignals {
		if topic, ok := signal.Context["interest"].(string); ok {
			tag := strings.ToLower(topic)
			category := c.categorizeInterest(tag)

			// Apply recency decay
			daysSince := time.Since(signal.CreatedAt).Hours() / 24
			recencyFactor := math.Exp(-daysSince / 30)

			tagScores[tag] += signal.Value * recencyFactor
			categoryScores[category] += signal.Value * recencyFactor
		}
	}

	// Process analogy satisfaction signals (indicates which interests resonate)
	for _, signal := range satisfactionSignals {
		if domain, ok := signal.Context["analogy_domain"].(string); ok {
			tag := strings.ToLower(domain)
			category := c.categorizeInterest(tag)

			daysSince := time.Since(signal.CreatedAt).Hours() / 24
			recencyFactor := math.Exp(-daysSince / 30)

			// Satisfaction signals are strong indicators
			tagScores[tag] += signal.Value * recencyFactor * 1.5
			categoryScores[category] += signal.Value * recencyFactor * 1.5
		}
	}

	// Build weighted interest list
	var interests []InterestWeight
	for tag, score := range tagScores {
		if score > 0.1 { // Filter out very low scores
			category := c.categorizeInterest(tag)
			interests = append(interests, InterestWeight{
				Tag:      tag,
				Category: string(category),
				Weight:   score,
			})
		}
	}

	// Sort by weight
	sort.Slice(interests, func(i, j int) bool {
		return interests[i].Weight > interests[j].Weight
	})

	// Limit to top 20 interests
	if len(interests) > 20 {
		interests = interests[:20]
	}

	// Determine top categories
	var topCategories []string
	type catScore struct {
		cat   InterestCategory
		score float64
	}
	var catScores []catScore
	for cat, score := range categoryScores {
		catScores = append(catScores, catScore{cat, score})
	}
	sort.Slice(catScores, func(i, j int) bool {
		return catScores[i].score > catScores[j].score
	})
	for i := 0; i < len(catScores) && i < 5; i++ {
		topCategories = append(topCategories, string(catScores[i].cat))
	}

	// Generate analogy source domains
	analogySources := c.generateAnalogySources(interests, topCategories)

	profile := &InterestProfile{
		UserID:         userID,
		Interests:      interests,
		TopCategories:  topCategories,
		AnalogySources: analogySources,
		LastUpdatedAt:  time.Now(),
	}

	c.log.Info("classified interests",
		zap.String("user_id", userID.String()),
		zap.Int("interest_count", len(interests)),
		zap.Strings("top_categories", topCategories),
	)

	return profile, nil
}

// categorizeInterest maps an interest tag to a category.
func (c *InterestClassifier) categorizeInterest(tag string) InterestCategory {
	tagLower := strings.ToLower(tag)

	// Direct match
	if cat, exists := interestKeywords[tagLower]; exists {
		return cat
	}

	// Substring match
	for keyword, cat := range interestKeywords {
		if strings.Contains(tagLower, keyword) || strings.Contains(keyword, tagLower) {
			return cat
		}
	}

	// Default to general technology (as this is a learning platform)
	return CategoryTechnology
}

// generateAnalogySources creates a list of domains to draw analogies from.
func (c *InterestClassifier) generateAnalogySources(interests []InterestWeight, topCategories []string) []string {
	sources := make(map[string]bool)

	// Add top interests as direct sources
	for i := 0; i < len(interests) && i < 5; i++ {
		sources[interests[i].Tag] = true
	}

	// Add category-based analogy domains
	categoryDomains := map[InterestCategory][]string{
		CategoryAnime:      {"anime storylines", "manga plot devices", "shonen training arcs", "anime power systems"},
		CategoryGaming:     {"game mechanics", "RPG skill trees", "game strategy", "quest systems", "boss battles"},
		CategorySports:     {"sports strategy", "team dynamics", "athletic training", "game statistics"},
		CategoryMusic:      {"musical composition", "rhythm and tempo", "band dynamics", "concert performances"},
		CategoryMovies:     {"movie plots", "character arcs", "cinematic techniques", "storytelling"},
		CategoryScience:    {"scientific experiments", "research methods", "natural phenomena"},
		CategoryTechnology: {"software systems", "engineering principles", "tech products"},
		CategoryArt:        {"artistic techniques", "creative processes", "design principles"},
		CategoryBooks:      {"literary themes", "story structures", "historical narratives"},
		CategoryFood:       {"cooking techniques", "recipe construction", "flavor combinations"},
		CategoryNature:     {"ecosystems", "animal behavior", "natural cycles"},
	}

	for _, catStr := range topCategories {
		cat := InterestCategory(catStr)
		if domains, exists := categoryDomains[cat]; exists {
			for _, d := range domains {
				sources[d] = true
			}
		}
	}

	var result []string
	for source := range sources {
		result = append(result, source)
	}

	return result
}

// RecordInterestSignal records that a user showed interest in a topic.
func (c *InterestClassifier) RecordInterestSignal(userID uuid.UUID, interest string, intensity float64) error {
	signal := &BehaviorSignal{
		UserID:     userID,
		SignalType: SignalInterestIndication,
		Value:      intensity,
		Context:    map[string]interface{}{"interest": interest},
	}
	return c.repo.RecordSignal(signal)
}

// RecordAnalogySatisfaction records that a user responded well to an analogy.
func (c *InterestClassifier) RecordAnalogySatisfaction(userID uuid.UUID, analogyDomain string, satisfaction float64) error {
	signal := &BehaviorSignal{
		UserID:     userID,
		SignalType: SignalAnalogySatisfaction,
		Value:      satisfaction,
		Context:    map[string]interface{}{"analogy_domain": analogyDomain},
	}
	return c.repo.RecordSignal(signal)
}
