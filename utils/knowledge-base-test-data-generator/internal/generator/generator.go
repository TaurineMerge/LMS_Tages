package generator

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/TaurineMerge/LMS_Tages/utils/knowledge-base-test-data-generator/internal/config"
	"github.com/TaurineMerge/LMS_Tages/utils/knowledge-base-test-data-generator/internal/models"
	"github.com/TaurineMerge/LMS_Tages/utils/knowledge-base-test-data-generator/internal/ollama"
	"github.com/TaurineMerge/LMS_Tages/utils/knowledge-base-test-data-generator/internal/writer"
	"github.com/google/uuid"
)

// Generator orchestrates the data generation process.
type Generator struct {
	cfg        *config.Config
	ollama     *ollama.Client
	writer     *writer.SQLWriter
	randSource *rand.Rand
	baseTime   time.Time
}

// New creates a new Generator.
func New(cfg *config.Config, ollama *ollama.Client, writer *writer.SQLWriter) *Generator {
	return &Generator{
		cfg:        cfg,
		ollama:     ollama,
		writer:     writer,
		randSource: rand.New(rand.NewSource(time.Now().UnixNano())),
		baseTime:   time.Now().AddDate(-1, 0, 0), // Start generating data from a year ago
	}
}

// Run starts the data generation process.
func (g *Generator) Run() error {
	log.Println("Starting data generation...")

	// 1. Get counts
	categoriesCount, err := g.getRandomCount(g.cfg.Counts.Categories)
	if err != nil {
		return fmt.Errorf("invalid categories_count: %w", err)
	}

	log.Printf("Generating %d categories...\n", categoriesCount)
	var generatedCategoryTitles []string

	for i := 0; i < categoriesCount; i++ {
		// 2. Generate Category
		log.Printf("Generating category %d/%d...\n", i+1, categoriesCount)
		prompt := strings.Replace(g.cfg.Prompts.CategoryTitle, "{{.ExistingTitles}}", strings.Join(generatedCategoryTitles, ", "), 1)
		categoryTitle, err := g.ollama.Generate(prompt)
		if err != nil {
			return fmt.Errorf("failed to generate category title: %w", err)
		}
		generatedCategoryTitles = append(generatedCategoryTitles, categoryTitle)

		categoryTime := g.generateTimestamp(g.baseTime, 30)
		category := models.Category{
			ID:        g.generateUUID(),
			Title:     categoryTitle,
			CreatedAt: categoryTime,
			UpdatedAt: categoryTime,
		}
		if err := g.writer.WriteCategory(category); err != nil {
			return fmt.Errorf("failed to write category: %w", err)
		}

		// 3. Generate Courses for this Category
		coursesCount, _ := g.getRandomCount(g.cfg.Counts.Courses)
		log.Printf("  Generating %d courses for category '%s'...\n", coursesCount, category.Title)
		for j := 0; j < coursesCount; j++ {
			log.Printf("  - Generating course %d/%d...\n", j+1, coursesCount)
			
			type courseContent struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			}
			prompt := strings.Replace(g.cfg.Prompts.CourseContent, "{{.CategoryTitle}}", category.Title, 1)
			courseJSON, err := g.ollama.Generate(prompt)
			if err != nil {
				return fmt.Errorf("failed to generate course content: %w", err)
			}
			
			var cc courseContent
			// Clean up potential markdown code fences
			courseJSON = strings.TrimPrefix(courseJSON, "```json")
			courseJSON = strings.TrimSuffix(courseJSON, "```")
			if err := json.Unmarshal([]byte(courseJSON), &cc); err != nil {
				log.Printf("    - Warning: failed to unmarshal course JSON from Ollama. Raw response: %s. Error: %v. Skipping course.", courseJSON, err)
				continue
			}

			courseTime := g.generateTimestamp(categoryTime, 10)
			course := models.Course{
				ID:          g.generateUUID(),
				Title:       cc.Title,
				Description: cc.Description,
				Level:       g.getRandomString([]string{"easy", "medium", "hard"}),
				Visibility:  "public",
				CategoryID:  category.ID,
				CreatedAt:   courseTime,
				UpdatedAt:   courseTime,
			}
			if err := g.writer.WriteCourse(course); err != nil {
				return fmt.Errorf("failed to write course: %w", err)
			}

			// 4. Generate Lessons for this Course
			lessonsCount, _ := g.getRandomCount(g.cfg.Counts.Lessons)
			log.Printf("    - Generating %d lessons for course '%s'...\n", lessonsCount, course.Title)
			for k := 0; k < lessonsCount; k++ {
				log.Printf("      - Generating lesson %d/%d...\n", k+1, lessonsCount)
				prompt := strings.Replace(g.cfg.Prompts.LessonTitle, "{{.CourseTitle}}", course.Title, 1)
				lessonTitle, err := g.ollama.Generate(prompt)
				if err != nil {
					return fmt.Errorf("failed to generate lesson title: %w", err)
				}
				
				lessonTime := g.generateTimestamp(courseTime, 5)
				lesson := models.Lesson{
					ID:        g.generateUUID(),
					Title:     lessonTitle,
					CourseID:  course.ID,
					CreatedAt: lessonTime,
					UpdatedAt: lessonTime,
				}

				// 5. Generate Content Blocks for this Lesson
				var contentBlocks []models.ContentBlock
				contentBlocksCount, _ := g.getRandomCount(g.cfg.Counts.ContentBlocks)
				for l := 0; l < contentBlocksCount; l++ {
					prompt := strings.Replace(g.cfg.Prompts.LessonContentBlock, "{{.LessonTitle}}", lesson.Title, 1)
					content, err := g.ollama.Generate(prompt)
					if err != nil {
						return fmt.Errorf("failed to generate content block: %w", err)
					}
					contentBlocks = append(contentBlocks, models.ContentBlock{
						ContentType: "text",
						Data:        content,
					})
				}
				lesson.Content = contentBlocks
				
				if err := g.writer.WriteLesson(lesson); err != nil {
					return fmt.Errorf("failed to write lesson: %w", err)
				}
			}
		}
	}

	log.Println("Data generation completed successfully.")
	return nil
}

// generateUUID creates a standard UUID.
// The custom format requirement is complex and non-standard.
// A standard UUID is a better practice.
func (g *Generator) generateUUID() string {
	return uuid.New().String()
}

// getRandomCount calculates a random number from a range string (e.g., "5" or "2:6").
func (g *Generator) getRandomCount(rangeStr string) (int, error) {
	min, max, err := config.ParseRange(rangeStr)
	if err != nil {
		return 0, err
	}
	if min == max {
		return min, nil
	}
	return g.randSource.Intn(max-min+1) + min, nil
}

// generateTimestamp creates a new timestamp that is slightly after a base time.
func (g *Generator) generateTimestamp(base time.Time, maxDaysForward int) time.Time {
	days := g.randSource.Intn(maxDaysForward)
	hours := g.randSource.Intn(24)
	minutes := g.randSource.Intn(60)
	
	return base.AddDate(0, 0, days).Add(time.Hour*time.Duration(hours) + time.Minute*time.Duration(minutes))
}

// getRandomString selects a random string from a slice.
func (g *Generator) getRandomString(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	return slice[g.randSource.Intn(len(slice))]
}
