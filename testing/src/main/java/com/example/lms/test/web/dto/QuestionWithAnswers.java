package com.example.lms.test.web.dto;

import java.util.List;

import com.example.lms.answer.api.dto.Answer;
import com.example.lms.question.api.dto.Question;

/**
 * Обертка для вопроса с его ответами
 */
public class QuestionWithAnswers {
    private Question question;
    private List<Answer> answers;
    
    public QuestionWithAnswers(Question question, List<Answer> answers) {
        this.question = question;
        this.answers = answers;
    }
    
    public Question getQuestion() { return question; }
    public void setQuestion(Question question) { this.question = question; }
    
    public List<Answer> getAnswers() { return answers; }
    public void setAnswers(List<Answer> answers) { this.answers = answers; }
}