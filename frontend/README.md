# Books and Tasks Management App

- [Books and Tasks Management App](#books-and-tasks-management-app)
  - [Overview](#overview)
  - [Design](#design)
    - [Types and Interfaces (/src/types/index.ts):](#types-and-interfaces-srctypesindexts)
    - [State Management (/src/stores/):](#state-management-srcstores)
    - [Components:](#components)
    - [Views (/src/views/):](#views-srcviews)
    - [Router (/src/router/index.ts):](#router-srcrouterindexts)
    - [UI/UX:](#uiux)
  - [Usage](#usage)
    - [Navigation:](#navigation)
    - [Books Management:](#books-management)
    - [Tasks Management:](#tasks-management)
    - [Features Available:](#features-available)


## Overview

A Vue 3 application with TypeScript and Bootstrap for managing Books and Tasks. 

## Design
### Types and Interfaces (/src/types/index.ts):

Defines TypeScript interfaces for Book and Task entities
Ensures type safety throughout the application

### State Management (/src/stores/):

- Uses Pinia for state management
- Separate stores for books and tasks
- Implements CRUD operations and search functionality

### Components:

- Book components (/src/components/books/):
    - BookForm.vue: Form for creating/editing books
    - BookList.vue: Display and manage books

- Task components (/src/components/tasks/):
    - TaskForm.vue: Form for creating/editing tasks
    - TaskList.vue: Display and manage tasks
    -
### Views (/src/views/):

BooksView.vue: Books management page
TasksView.vue: Tasks management page

### Router (/src/router/index.ts):

Handles navigation between views
Default route redirects to books page
Features:

Full CRUD operations for both books and tasks
Form validation using Vuelidate
Search functionality for both books and tasks
Responsive design using Bootstrap
Date formatting using date-fns
Type-safe throughout with TypeScript

### UI/UX:

- Clean and intuitive interface
- Responsive design that works on all devices
- Form validation with error messages
- Confirmation dialogs for deletions
- Visual indicators for task priorities


The development server is now running. You can navigate between the Books and Tasks sections using the navigation bar at the top. Each section allows you to:

- View all items in a table format
- Search items using the search bar
- Add new items using the "Add New" button
- Edit existing items
- Delete items (with confirmation)
- Manage all properties of both Books and Tasks


The application is built with scalability in mind and can be easily extended with additional features as needed.


## Usage

The application is already set up with Vite as the development server.

To run the application, simply use the following command:

Start Development Server
```
npm run dev
```
The development server is now running. Here's what you can do:

### Navigation:

- Use the navigation bar at the top to switch between Books and Tasks sections
- The default route will take you to the Books page

### Books Management:

- Click "Add New Book" to create a book
- Fill in the required fields (ISBN, Title, Author, Price)
- Optionally add Borrow Time and Return Time
- Use the search bar to filter books
- Edit or delete existing books using the action buttons

### Tasks Management:

- Click "Add New Task" to create a task
- Fill in the required fields (Name, Description, Priority, Duration, Deadline)
- Optionally add Start Time, End Time, and Tags
- Use the search bar to filter tasks
- Edit or delete existing tasks using the action buttons

### Features Available:

- Full CRUD operations (Create, Read, Update, Delete)
- Search functionality
- Form validation
- Responsive design that works on all screen sizes

The application uses:

- Vue 3 for the frontend framework
- TypeScript for type safety
- Bootstrap for styling
- Pinia for state management
- Vue Router for navigation
- Vuelidate for form validation

The development server provides hot module replacement (HMR), so any changes you make to the code will be immediately reflected in the browser.