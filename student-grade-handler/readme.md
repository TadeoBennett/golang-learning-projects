# Project Overview

Welcome to our project repository!

This project is a collaborative effort between **Tadeo Bennett** and **Cahlil Tillet**. Through this course, we have gained invaluable insights into key programming concepts and best practices. Our learnings have been instrumental in shaping this project, particularly in the realms of coding practices, security, and more.

## Key Features

Our project demonstrates the practical application of several important concepts and technologies:

- **PostgreSQL Database:** Utilizes PostgreSQL for robust and scalable database management.
- **Middleware:** Implements middleware to enhance functionality and manage request processing.
- **Centralized Error Handling:** Ensures consistent and efficient error management across the application.
- **Dependency Injection:** Employs dependency injection to promote modularity and testability.
- **Self-Signed Certification:** Uses self-signed certificates to handle secure communications.
- **Custom Error and Feedback Logging:** Provides custom logging for errors and user feedback to improve troubleshooting and user experience.

## Collaboration

This project is the result of a collaborative effort between:

- **Tadeo Bennett**
- **Cahlil Tillet**

Working together, we have applied the concepts learned throughout the course to build a comprehensive solution that showcases our understanding and implementation of key programming practices.

## Learnings

The teachings of this course have been fundamental in guiding us through:

- **Best Coding Practices:** Emphasizing clean, maintainable, and efficient code.
- **Security Measures:** Implementing robust security practices to safeguard the application.
- **Effective Error Handling:** Developing strategies for handling and logging errors effectively.
- **Advanced Programming Concepts:** Applying advanced techniques such as dependency injection and middleware.

Our project not only demonstrates these concepts in action but also reflects our commitment to applying them to real-world scenarios.

## Getting Started

To get started with this project, please follow the instructions below:

1. **Clone the Repository:**

   ```bash
    git clone https://github.com/yourusername/yourproject.git

   ```

2. **Setup PostgreSQL Database:**

   **Step 1: Install PostgreSQL**

   Install PostgreSQL by following this video guide: [PostgreSQL Installation Guide](https://www.youtube.com/watch?v=0Il040ExA_Q)

   **Step 2: Configure the Database**

   - **Login to PostgreSQL:**
     Open your terminal and switch to the `postgres` user:

     ```bash
     sudo -i -u postgres psql
     ```

   - **Run the SQL Script:**
     Execute the commands in the `students.sql` file located in the `db` directory. You can run the SQL script with the following command:

     ```bash
     psql -f /path/to/your/project/db/students.sql
     ```

     Make sure to replace `/path/to/your/project/` with the actual path to your project directory. OR copy and past the code in the file, AS INSTRUCTED in the comments provided there
