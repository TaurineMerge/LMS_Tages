-- This script replaces the old dynamic data generation with a static,
-- deterministic set of data for consistent testing and development.
-- All UUIDs are hardcoded to ensure data integrity between runs.
-- Version 5: Content converted from JSON to HTML format.

BEGIN;

-- Truncate tables to ensure a clean slate, and restart identity for any sequence.
TRUNCATE TABLE knowledge_base.category_d, knowledge_base.course_b, knowledge_base.lesson_d RESTART IDENTITY CASCADE;

--
-- Category 1: Go Programming
--
INSERT INTO knowledge_base.category_d (id, title, created_at, updated_at) VALUES ('a1000000-0000-4000-8000-000000000001', 'Go Programming', '2024-01-10 10:00:00', '2024-01-10 10:00:00');

-- Course 1.1: Introduction to Go
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b1000001-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000001', 'Introduction to Go', 'Learn the fundamentals of the Go programming language from scratch. This course covers basic syntax, control structures, and the core philosophies behind Go.', 'easy', 'public', '2024-01-11 11:00:00', '2024-01-11 11:00:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES
('c1000001-0001-4000-8000-000000000001', 'b1000001-0000-4000-8000-000000000001', 'Module 1: Setting Up Your Environment', '<p>To start writing Go code, you first need to install the Go toolchain. Visit the official Go website at go.dev and download the installer for your operating system.</p><p>Once installed, verify the installation by opening a terminal and running `go version`. You should see the installed Go version printed to the console.</p><p>Next, configure your editor. VS Code with the Go extension is a popular choice, providing features like IntelliSense, debugging, and code navigation.</p>', '2024-01-11 12:00:00', '2024-01-11 12:00:00'),
('c1000001-0002-4000-8000-000000000001', 'b1000001-0000-4000-8000-000000000001', 'Module 2: Variables and Data Types', '<p>Go is a statically typed language. This means variable types are known at compile time. We will explore basic types like `int`, `string`, and `bool`.</p><p>You can declare a variable using `var name type = value` or use the short declaration `name := value` which is more common within function bodies.</p><p>Understanding the difference between value types and reference types is crucial for managing memory and performance in Go.</p><p>We will also look at composite types such as arrays, slices, and maps.</p>', '2024-01-12 13:00:00', '2024-01-12 13:00:00'),
('c1000001-0003-4000-8000-000000000001', 'b1000001-0000-4000-8000-000000000001', 'Module 3: Control Structures', '<p>This lesson covers `if/else` for conditional logic and the powerful `for` loop. Go''s `for` loop is its only looping construct, but it is very versatile.</p><p>We will also cover the `switch` statement, which provides a clean way to express complex conditionals.</p><p>Error handling is a key part of Go''s philosophy, often done with an explicit `if err != nil` check.</p>', '2024-01-13 14:00:00', '2024-01-13 14:00:00');

-- Course 1.2: Advanced Go Concurrency
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b1000002-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000001', 'Advanced Go Concurrency', 'Dive deep into Go''s powerful concurrency model, including goroutines, channels, and the select statement for complex coordination.', 'hard', 'draft', '2024-02-01 09:00:00', '2024-02-01 09:00:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES
('d1000002-0001-4000-8000-000000000001', 'b1000002-0000-4000-8000-000000000001', 'Unit 1: Goroutines Deep Dive', '<p>A goroutine is a lightweight thread managed by the Go runtime. We explore how the scheduler works and best practices for managing thousands of goroutines.</p><p>Learn about the `sync` package, including `WaitGroup` for synchronizing groups of goroutines.</p>', '2024-02-01 10:00:00', '2024-02-01 10:00:00'),
('d1000002-0002-4000-8000-000000000001', 'b1000002-0000-4000-8000-000000000001', 'Unit 2: Advanced Channel Patterns', '<p>Channels are the pipes that connect concurrent goroutines. This unit explores advanced patterns like fan-in, fan-out, and pipelines.</p><p>We will also discuss channel directionality and how to write safer, more readable concurrent code.</p><p>Closing channels correctly to signal completion is a critical skill we will master.</p>', '2024-02-02 11:00:00', '2024-02-02 11:00:00');

-- Course 1.3: Go Web Development
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b1000003-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000001', 'Go Web Development', 'Build scalable web applications using Go''s standard library and popular frameworks. Learn routing, middleware, and database integration.', 'medium', 'public', '2024-02-15 08:00:00', '2024-02-15 08:00:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES ('c1000003-0001-4000-8000-000000000001', 'b1000003-0000-4000-8000-000000000001', 'HTTP Servers and Routing', '<p>Learn how to build HTTP servers using Go''s net/http package and handle routing for different endpoints.</p>', '2024-02-16 09:00:00', '2024-02-16 09:00:00');

-- Course 1.4: Testing and Benchmarking in Go
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b1000004-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000001', 'Testing and Benchmarking in Go', 'Master writing unit tests, integration tests, and benchmarks to ensure code quality and performance in Go applications.', 'medium', 'public', '2024-02-20 09:30:00', '2024-02-20 09:30:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES ('c1000004-0001-4000-8000-000000000001', 'b1000004-0000-4000-8000-000000000001', 'Unit Testing Basics', '<p>Explore Go''s built-in testing package and learn to write effective unit tests using table-driven test patterns.</p>', '2024-02-21 10:00:00', '2024-02-21 10:00:00');

-- Course 1.5: Go REST APIs
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b1000005-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000001', 'Building REST APIs with Go', 'Design and implement RESTful APIs following best practices, including error handling, authentication, and JSON serialization.', 'medium', 'public', '2024-02-25 10:00:00', '2024-02-25 10:00:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES ('c1000005-0001-4000-8000-000000000001', 'b1000005-0000-4000-8000-000000000001', 'JSON and Request Handling', '<p>Learn how to parse incoming JSON requests and encode responses using Go''s encoding/json package.</p>', '2024-02-26 11:00:00', '2024-02-26 11:00:00');

-- Course 1.6: Go Microservices Architecture
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b1000006-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000001', 'Microservices Architecture with Go', 'Learn to design, build, and deploy microservices using Go, including service communication patterns and containerization strategies.', 'hard', 'public', '2024-03-05 08:30:00', '2024-03-05 08:30:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES ('c1000006-0001-4000-8000-000000000001', 'b1000006-0000-4000-8000-000000000001', 'Service Communication', '<p>Understand different patterns for inter-service communication including gRPC, message queues, and event-driven architectures.</p>', '2024-03-06 09:00:00', '2024-03-06 09:00:00');

-- Course 1.7: Go Performance Optimization
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b1000007-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000001', 'Performance Optimization in Go', 'Optimize Go applications for speed and memory efficiency using profiling tools, memory management techniques, and concurrency patterns.', 'hard', 'public', '2024-03-10 09:00:00', '2024-03-10 09:00:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES ('c1000007-0001-4000-8000-000000000001', 'b1000007-0000-4000-8000-000000000001', 'CPU and Memory Profiling', '<p>Use pprof and other profiling tools to identify and eliminate performance bottlenecks in your Go applications.</p>', '2024-03-11 10:00:00', '2024-03-11 10:00:00');

--
-- Category 2: Python
--
INSERT INTO knowledge_base.category_d (id, title, created_at, updated_at) VALUES ('a2000000-0000-4000-8000-000000000002', 'Python Programming', '2024-01-15 09:00:00', '2024-01-15 09:00:00');
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b2000001-0000-4000-8000-000000000002', 'a2000000-0000-4000-8000-000000000002', 'Python for Data Science', 'Learn how to use Python with Pandas, NumPy, and Matplotlib for data analysis and visualization.', 'medium', 'public', '2024-01-16 10:00:00', '2024-01-16 10:00:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES
('e2000001-0001-4000-8000-000000000002', 'b2000001-0000-4000-8000-000000000002', 'Week 1: NumPy Fundamentals', '<p>NumPy is the fundamental package for scientific computing with Python. We will explore its powerful N-dimensional array object.</p><p>Learn about array creation, indexing, slicing, and broadcasting.</p>', '2024-01-17 11:00:00', '2024-01-17 11:00:00'),
('e2000001-0002-4000-8000-000000000002', 'b2000001-0000-4000-8000-000000000002', 'Week 2: Introduction to Pandas', '<p>Pandas provides high-performance, easy-to-use data structures and data analysis tools. The primary data structures are Series and DataFrame.</p><p>This week focuses on reading data from various sources (CSV, Excel) and performing initial data exploration.</p>', '2024-01-24 11:00:00', '2024-01-24 11:00:00');

--
-- Category 3: JavaScript
--
INSERT INTO knowledge_base.category_d (id, title, created_at, updated_at) VALUES ('a3000000-0000-4000-8000-000000000003', 'JavaScript Ecosystem', '2024-02-10 10:00:00', '2024-02-10 10:00:00');
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b3000001-0000-4000-8000-000000000003', 'a3000000-0000-4000-8000-000000000003', 'Modern JavaScript (ES6+)', 'Master the features of modern JavaScript, including let/const, arrow functions, classes, modules, promises, and async/await.', 'medium', 'public', '2024-02-11 11:00:00', '2024-02-11 11:00:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES
('c3000001-0001-4000-8000-000000000003', 'b3000001-0000-4000-8000-000000000003', 'Day 1: New Variable Declarations', '<p>`let` and `const` provide block-scoping, which is a significant improvement over `var`''s function-scoping. We discuss the implications for code clarity and bug prevention.</p>', '2024-02-12 12:00:00', '2024-02-12 12:00:00'),
('c3000001-0002-4000-8000-000000000003', 'b3000001-0000-4000-8000-000000000003', 'Day 2: Arrow Functions', '<p>Arrow functions provide a more concise syntax for writing function expressions. They also lexically bind the `this` value, which simplifies many common scenarios.</p>', '2024-02-13 12:00:00', '2024-02-13 12:00:00'),
('c3000001-0003-4000-8000-000000000003', 'b3000001-0000-4000-8000-000000000003', 'Day 3: Promises and Async/Await', '<p>Asynchronous programming is fundamental to JavaScript. We will replace callback-hell with the elegance and readability of Promises and the `async/await` syntax.</p>', '2024-02-14 12:00:00', '2024-02-14 12:00:00');

-- Course 3.2: Building with React
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b3000002-0000-4000-8000-000000000003', 'a3000000-0000-4000-8000-000000000003', 'Building with React', 'An intermediate course on building interactive user interfaces with the React library. Requires basic knowledge of JavaScript.', 'medium', 'draft', '2024-03-01 10:00:00', '2024-03-01 10:00:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES ('d3000002-0001-4000-8000-000000000003', 'b3000002-0000-4000-8000-000000000003', 'React State Management', '<p>Exploring different state management solutions in React, from local `useState` to context API and libraries like Redux or Zustand.</p>', '2024-03-02 11:00:00', '2024-03-02 11:00:00');

--
-- Category 4: Databases
--
INSERT INTO knowledge_base.category_d (id, title, created_at, updated_at) VALUES ('a4000000-0000-4000-8000-000000000004', 'Databases', '2024-03-10 10:00:00', '2024-03-10 10:00:00');
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b4000001-0000-4000-8000-000000000004', 'a4000000-0000-4000-8000-000000000004', 'SQL Fundamentals', 'Learn the basics of Structured Query Language (SQL) to interact with relational databases.', 'easy', 'public', '2024-03-11 11:00:00', '2024-03-11 11:00:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES ('c4000001-0001-4000-8000-000000000004', 'b4000001-0000-4000-8000-000000000004', 'SELECT and WHERE clauses', '<p>The foundation of querying data is the SELECT statement. We will learn how to select specific columns and filter rows using the WHERE clause.</p>', '2024-03-12 12:00:00', '2024-03-12 12:00:00');

--
-- Category 5: DevOps
--
INSERT INTO knowledge_base.category_d (id, title, created_at, updated_at) VALUES ('a5000000-0000-4000-8000-000000000005', 'DevOps', '2024-04-01 10:00:00', '2024-04-01 10:00:00');
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b5000001-0000-4000-8000-000000000005', 'a5000000-0000-4000-8000-000000000005', 'Introduction to Docker', 'Learn how to containerize your applications with Docker for consistent development and deployment environments.', 'easy', 'public', '2024-04-02 11:00:00', '2024-04-02 11:00:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES ('c5000001-0001-4000-8000-000000000005', 'b5000001-0000-4000-8000-000000000005', 'Writing a Dockerfile', '<p>A Dockerfile is a script that contains instructions for building a Docker image. We will learn the basic commands like FROM, WORKDIR, COPY, RUN, and CMD.</p><p>We will build a simple image for a Go application.</p>', '2024-04-03 12:00:00', '2024-04-03 12:00:00');

--
-- Category 6: Design Principles
--
INSERT INTO knowledge_base.category_d (id, title, created_at, updated_at) VALUES ('a6000000-0000-4000-8000-000000000006', 'Design Principles', '2024-05-01 10:00:00', '2024-05-01 10:00:00');
INSERT INTO knowledge_base.course_b (id, category_id, title, description, level, visibility, created_at, updated_at) VALUES ('b6000001-0000-4000-8000-000000000006', 'a6000000-0000-4000-8000-000000000006', 'SOLID Principles in Practice', 'Understand and apply the five SOLID principles of object-oriented design to write more maintainable, flexible, and scalable code.', 'medium', 'draft', '2024-05-02 11:00:00', '2024-05-02 11:00:00');
INSERT INTO knowledge_base.lesson_d (id, course_id, title, content, created_at, updated_at) VALUES ('c6000001-0001-4000-8000-000000000006', 'b6000001-0000-4000-8000-000000000006', 'Single Responsibility Principle', '<p>A class should have only one reason to change. This principle helps to keep classes focused and small.</p>', '2024-05-03 12:00:00', '2024-05-03 12:00:00');

COMMIT;