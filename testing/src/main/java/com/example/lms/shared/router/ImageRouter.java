package com.example.lms.shared.router;

import com.example.lms.shared.controller.ImageController;

import static io.javalin.apibuilder.ApiBuilder.*;

public class ImageRouter {

    public static void register(ImageController controller) {
        path("/api", () -> {
            // Изображения вопросов
            path("/questions/{questionId}/images", () -> {
                post(controller::uploadQuestionImage);
                path("/{imageId}", () -> {
                    get(controller::getQuestionImage);
                    delete(controller::deleteQuestionImage);
                });
            });

            // Изображения ответов
            path("/answers/{answerId}/images", () -> {
                post(controller::uploadAnswerImage);
                path("/{imageId}", () -> {
                    get(controller::getAnswerImage);
                    delete(controller::deleteAnswerImage);
                });
            });
        });
    }
}