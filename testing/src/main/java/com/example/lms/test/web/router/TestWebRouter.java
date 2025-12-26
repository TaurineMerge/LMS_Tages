package com.example.lms.test.web.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.test.web.controller.TestFormController;
import com.example.lms.tracing.SimpleTracer;

import static io.javalin.apibuilder.ApiBuilder.*;

public class TestWebRouter {

    private static final Logger logger = LoggerFactory.getLogger(TestWebRouter.class);

    public static void register(TestFormController controller) {

        path("/web", () -> {

            get("/", ctx -> ctx.redirect("/web/tests/new"));

            path("/tests", () -> {

                before(ctx -> {
                    logger.info("Web request started: {} {} (traceId: {})",
                            ctx.method(),
                            ctx.path(),
                            SimpleTracer.getCurrentTraceId()
                    );
                    ctx.attribute("redirect", ctx.queryParam("redirect"));
                });

                // ---------- CREATE ----------

                get("/new", controller::showNewTestForm);
                post("/save", controller::saveTestFromForm);
                post("/draft", controller::saveNewTestDraft);

                // ---------- ENTITY ----------

                path("/{id}", () -> {

                    get("/edit", controller::showEditTestForm);
                    get("/preview", controller::previewTest);

                    put("/", ctx -> {
                        String id = ctx.pathParam("id");
                        if (id.startsWith("draft-")) {
                            controller.updateDraftFromForm(ctx);
                        } else {
                            controller.updateTestFromForm(ctx);
                        }
                    });

                    delete("/", ctx -> {
                        String id = ctx.pathParam("id");
                        if (id.startsWith("draft-")) {
                            controller.deleteDraft(ctx);
                        } else {
                            controller.deleteTest(ctx);
                        }
                    });

                    post("/delete", ctx -> {
                        String id = ctx.pathParam("id");
                        if (id.startsWith("draft-")) {
                            controller.deleteDraft(ctx);
                        } else {
                            controller.deleteTest(ctx);
                        }
                    });

                    post("/create-draft", controller::createDraftFromTest);
                    post("/draft", controller::saveExistingTestDraft);
                    post("/publish", controller::publishDraft);
                });

                after(ctx -> {
                    logger.info("Web request completed: {} {} -> {} (traceId: {})",
                            ctx.method(),
                            ctx.path(),
                            ctx.status(),
                            SimpleTracer.getCurrentTraceId()
                    );
                });
            });
        });
    }
}