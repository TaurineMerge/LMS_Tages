package com.example.lms.test.web.dto;

import java.util.List;
import java.util.UUID;

/**
 * DTO для данных формы создания/редактирования теста
 */
public class TestFormData {
    private String testId;
    private String title;
    private String description;
    private Integer minPoint;
    private List<QuestionFormData> questions;
    
    // Геттеры и сеттеры для TestFormData
    public String getTestId() {
        return testId;
    }
    
    public void setTestId(String testId) {
        this.testId = testId;
    }
    
    public String getTitle() {
        return title;
    }
    
    public void setTitle(String title) {
        this.title = title;
    }
    
    public String getDescription() {
        return description;
    }
    
    public void setDescription(String description) {
        this.description = description;
    }
    
    public Integer getMinPoint() {
        return minPoint;
    }
    
    public void setMinPoint(Integer minPoint) {
        this.minPoint = minPoint;
    }
    
    public List<QuestionFormData> getQuestions() {
        return questions;
    }
    
    public void setQuestions(List<QuestionFormData> questions) {
        this.questions = questions;
    }
    
    public static class QuestionFormData {
        private UUID id;
        private String textOfQuestion;
        private Integer order;
        private List<AnswerFormData> answers;
        
        // Геттеры и сеттеры для QuestionFormData
        public UUID getId() {
            return id;
        }
        
        public void setId(UUID id) {
            this.id = id;
        }
        
        public String getTextOfQuestion() {
            return textOfQuestion;
        }
        
        public void setTextOfQuestion(String textOfQuestion) {
            this.textOfQuestion = textOfQuestion;
        }
        
        public Integer getOrder() {
            return order;
        }
        
        public void setOrder(Integer order) {
            this.order = order;
        }
        
        public List<AnswerFormData> getAnswers() {
            return answers;
        }
        
        public void setAnswers(List<AnswerFormData> answers) {
            this.answers = answers;
        }
    }
    
    public static class AnswerFormData {
        private UUID id;
        private String text;
        private Integer score;
        private Boolean isCorrect;
        // УБИРАЕМ поле order из AnswerFormData
        // private Integer order;
        
        // Геттеры и сеттеры для AnswerFormData
        public UUID getId() {
            return id;
        }
        
        public void setId(UUID id) {
            this.id = id;
        }
        
        public String getText() {
            return text;
        }
        
        public void setText(String text) {
            this.text = text;
        }
        
        public Integer getScore() {
            return score;
        }
        
        public void setScore(Integer score) {
            this.score = score;
        }
        
        public Boolean getIsCorrect() {
            return isCorrect;
        }
        
        public void setIsCorrect(Boolean isCorrect) {
            this.isCorrect = isCorrect;
        }
        
        // УБИРАЕМ методы getOrder и setOrder
        // public Integer getOrder() {
        //     return order;
        // }
        // 
        // public void setOrder(Integer order) {
        //     this.order = order;
        // }
    }
}