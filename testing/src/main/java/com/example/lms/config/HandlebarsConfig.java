package com.example.lms.config;

import com.github.jknack.handlebars.Handlebars;
import com.github.jknack.handlebars.io.ClassPathTemplateLoader;

public class HandlebarsConfig {
    public static Handlebars configureHandlebars() {
        try {
            // Основной загрузчик шаблонов
            ClassPathTemplateLoader loader = new ClassPathTemplateLoader();
            loader.setPrefix("/templates");
            loader.setSuffix(".hbs");

            Handlebars handlebars = new Handlebars(loader);
            handlebars.setInfiniteLoops(true);

            // Регистрация helpers для partials
            handlebars.registerHelpers(new ClassPathTemplateLoader("/templates/partials", ".hbs"));

            return handlebars;
        } catch (Exception e) {
            System.err.println("Ошибка при конфигурации Handlebars: " + e.getMessage());
            e.printStackTrace();
            return new Handlebars();
        }
    }
}
