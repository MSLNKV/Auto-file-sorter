# Auto-file-sorter
A training project to get acquainted with the capabilities of the Go language.  The desired outcome is to study the main tools of the language, to immerse oneself in the specifics of working with the language.  

**"Technical task":**
To create an automatic sorter of files that end up in a directory (for example, /downloads), which will keep files with logging (recording file movements, errors), automatically clear old logs, remind the user to clear heavy files (or those not used for a long time)

**Additional difficulties and challenges**
1) Flexible settings for the user
2) Support for different criteria for sorting files
3) Individual reminder rules for each file category
4) Handling errors and unexpected user actions

## 📌 Current Status
**Last updated:** 2025-09-09 

### ✅ What is already implemented:
- 🔹Sorting files by:
- Extension (categories based on file types)
- Content type (images, documents, archives, etc.)
- 🔹 Automatic creation of folders for categories
- 🔹 Automatic logging with division by days
- 🔹 Cleaning old logs (older than 72 hours)
- 🔹 Deleting empty folders after sorting
- 🔹 `Undo/Redo` system for undoing/returning actions
  
### 🛠️In work:
- 🔸 Support for custom categories (folder structure customization)
- 🔸 Development of a module **reminders about cleaning** (Reminder)
- 🔸 Refactoring and optimization of algorithms
- 🔸 Preparation for integration of GUI or CLI interface
- 🔸 Increasing code stability and error handling

### 🎯 Future plans:
- 📂 Implementation of a database for storing the history of actions and settings
- 🔍 Extended metadata system for files (view flags, reminders)
- 🧩 Integration of a sorter into the future project **NoctisExplorer**
- 🌐 Cross-platform support (Linux/Windows/macOS)
- 🎨 Development of a user interface


## The project's philosophy is not to use neural networks in any form. Let it be "non-professional code", but it will give me experience

