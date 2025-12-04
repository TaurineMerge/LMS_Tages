package com.example;

import io.javalin.Javalin;

public class Main {
	public static void main(String[] args) {
		Javalin app = Javalin.create()
				.get("/", ctx -> ctx.result("Hello Javalin!"))
				.get("/api/hello", ctx -> ctx.json(new Message("Hello from Javalin API!")))
				.start(8085);

		System.out.println("Server started at http://localhost:8085");
	}

	static class Message {
		private String text;

		public Message(String text) {
			this.text = text;
		}

		public String getText() {
			return text;
		}
	}
}