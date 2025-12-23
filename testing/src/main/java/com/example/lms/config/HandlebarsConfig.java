package com.example.lms.config;

import com.github.jknack.handlebars.Handlebars;
import com.github.jknack.handlebars.helper.ConditionalHelpers;
import com.github.jknack.handlebars.io.ClassPathTemplateLoader;
import com.github.jknack.handlebars.Helper;
import com.github.jknack.handlebars.Options;

import java.io.IOException;

public class HandlebarsConfig {
    public static Handlebars configureHandlebars() {
        try {
            // Основной загрузчик шаблонов
            ClassPathTemplateLoader loader = new ClassPathTemplateLoader();
            loader.setPrefix("/templates");
            loader.setSuffix(".hbs");

            Handlebars handlebars = new Handlebars(loader);
            handlebars.setInfiniteLoops(true);

            // Регистрация стандартных хелперов
            handlebars.registerHelpers(ConditionalHelpers.class);
            
            // Регистрация кастомных хелперов
            registerCustomHelpers(handlebars);

            // Регистрация helpers для partials
            handlebars.registerHelpers(new ClassPathTemplateLoader("/templates/partials", ".hbs"));

            return handlebars;
        } catch (Exception e) {
            System.err.println("Ошибка при конфигурации Handlebars: " + e.getMessage());
            e.printStackTrace();
            return new Handlebars();
        }
    }
    
    private static void registerCustomHelpers(Handlebars handlebars) {
        // Хелпер для сложения двух чисел
        handlebars.registerHelper("add", new Helper<Object>() {
            @Override
            public Object apply(Object context, Options options) throws IOException {
                if (context instanceof Number && options.params.length > 0 && options.params[0] instanceof Number) {
                    Number a = (Number) context;
                    Number b = (Number) options.params[0];
                    return a.intValue() + b.intValue();
                }
                return context;
            }
        });
        
        // Хелпер для сравнения двух значений
        handlebars.registerHelper("eq", new Helper<Object>() {
            @Override
            public Object apply(Object context, Options options) throws IOException {
                if (options.params.length > 0) {
                    Object param = options.params[0];
                    if (context == null && param == null) return true;
                    if (context != null && context.equals(param)) return true;
                }
                return false;
            }
        });
        
        // Хелпер для сравнения чисел
        handlebars.registerHelper("eqNumber", new Helper<Object>() {
            @Override
            public Object apply(Object context, Options options) throws IOException {
                if (context instanceof Number && options.params.length > 0 && options.params[0] instanceof Number) {
                    Number a = (Number) context;
                    Number b = (Number) options.params[0];
                    return a.intValue() == b.intValue();
                }
                return false;
            }
        });
        
        // Хелпер для проверки, больше ли первое число
        handlebars.registerHelper("gt", new Helper<Object>() {
            @Override
            public Object apply(Object context, Options options) throws IOException {
                if (context instanceof Number && options.params.length > 0 && options.params[0] instanceof Number) {
                    Number a = (Number) context;
                    Number b = (Number) options.params[0];
                    return a.intValue() > b.intValue();
                }
                return false;
            }
        });
        
        // Хелпер для проверки, меньше ли первое число
        handlebars.registerHelper("lt", new Helper<Object>() {
            @Override
            public Object apply(Object context, Options options) throws IOException {
                if (context instanceof Number && options.params.length > 0 && options.params[0] instanceof Number) {
                    Number a = (Number) context;
                    Number b = (Number) options.params[0];
                    return a.intValue() < b.intValue();
                }
                return false;
            }
        });
        
        // Хелпер для проверки, не равно ли
        handlebars.registerHelper("neq", new Helper<Object>() {
            @Override
            public Object apply(Object context, Options options) throws IOException {
                if (options.params.length > 0) {
                    Object param = options.params[0];
                    if (context == null && param == null) return false;
                    if (context == null || param == null) return true;
                    return !context.equals(param);
                }
                return true;
            }
        });
        
        // Хелпер для деления (возвращает целую часть)
        handlebars.registerHelper("divide", new Helper<Object>() {
            @Override
            public Object apply(Object context, Options options) throws IOException {
                if (context instanceof Number && options.params.length > 0 && options.params[0] instanceof Number) {
                    Number a = (Number) context;
                    Number b = (Number) options.params[0];
                    if (b.intValue() == 0) return 0;
                    return a.intValue() / b.intValue();
                }
                return context;
            }
        });
        
        // Хелпер для умножения
        handlebars.registerHelper("multiply", new Helper<Object>() {
            @Override
            public Object apply(Object context, Options options) throws IOException {
                if (context instanceof Number && options.params.length > 0 && options.params[0] instanceof Number) {
                    Number a = (Number) context;
                    Number b = (Number) options.params[0];
                    return a.intValue() * b.intValue();
                }
                return context;
            }
        });
        
        // Хелпер для проверки, содержит ли строка подстроку
        handlebars.registerHelper("contains", new Helper<Object>() {
            @Override
            public Object apply(Object context, Options options) throws IOException {
                if (context != null && options.params.length > 0 && options.params[0] != null) {
                    String str = context.toString();
                    String substring = options.params[0].toString();
                    return str.contains(substring);
                }
                return false;
            }
        });
        
        // Хелпер для проверки длины массива/списка
        handlebars.registerHelper("length", new Helper<Object>() {
            @Override
            public Object apply(Object context, Options options) throws IOException {
                if (context instanceof java.util.List) {
                    return ((java.util.List<?>) context).size();
                } else if (context instanceof Object[]) {
                    return ((Object[]) context).length;
                } else if (context instanceof String) {
                    return ((String) context).length();
                }
                return 0;
            }
        });
    }
}