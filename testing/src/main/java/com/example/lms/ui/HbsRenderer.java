package com.example.lms.ui;

import com.github.jknack.handlebars.Handlebars;
import com.github.jknack.handlebars.Template;
import com.github.jknack.handlebars.io.ClassPathTemplateLoader;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public final class HbsRenderer {

    private static final Handlebars handlebars;
    private static final Map<String, Template> cache = new ConcurrentHashMap<>();

    static {
        // Ищем шаблоны в classpath: /templates/*.hbs
        var loader = new ClassPathTemplateLoader("/templates", ".hbs");
        handlebars = new Handlebars(loader);
    }

    private HbsRenderer() { }

    public static String render(String templateName, Object model) {
        try {
            Template t = cache.computeIfAbsent(templateName, name -> {
                try {
                    return handlebars.compile(name);
                } catch (Exception e) {
                    throw new RuntimeException("Cannot compile template: " + name, e);
                }
            });
            return t.apply(model);
        } catch (Exception e) {
            throw new RuntimeException("Cannot render template: " + templateName, e);
        }
    }
}